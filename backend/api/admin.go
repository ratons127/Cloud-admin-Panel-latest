package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bertrandmartel/aws-admin/backend/auth"
	"github.com/bertrandmartel/aws-admin/backend/middleware"
	"github.com/bertrandmartel/aws-admin/backend/model"
	"github.com/bertrandmartel/aws-admin/backend/store"
	"github.com/bertrandmartel/aws-admin/backend/utils"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type AdminUserCreateRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	IsSuperAdmin bool   `json:"isSuperAdmin"`
	SendEmail    bool   `json:"sendEmail"`
	TenantID     string `json:"tenantId"`
	Role         string `json:"role"`
}

type AdminUserUpdateRequest struct {
	Email         string `json:"email"`
	IsSuperAdmin  bool   `json:"isSuperAdmin"`
	IsDisabled    bool   `json:"isDisabled"`
	EmailVerified bool   `json:"emailVerified"`
}

type AdminResetPasswordRequest struct {
	Password  string `json:"password"`
	SendEmail bool   `json:"sendEmail"`
}

type AdminTenantRequest struct {
	Name        string `json:"name"`
	OwnerUserID string `json:"ownerUserId"`
}

type AdminTenantMemberRequest struct {
	UserID         string `json:"userId"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	CreateIfMissing bool  `json:"createIfMissing"`
	SendEmail      bool   `json:"sendEmail"`
	Password       string `json:"password"`
}

type AdminTenantMemberUpdateRequest struct {
	Role string `json:"role"`
}

type AdminAWSAccountRequest struct {
	AccountID string `json:"accountId"`
	RoleArn   string `json:"roleArn"`
	ExternalID string `json:"externalId"`
	Name      string `json:"name"`
	Active    *bool  `json:"active"`
}

func ListUsers(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	limit := parseLimit(c.QueryParam("limit"), 200)
	offset := parseLimit(c.QueryParam("offset"), 0)
	query := strings.ToLower(strings.TrimSpace(c.QueryParam("q")))

	var rows *sql.Rows
	var err error
	if query != "" {
		rows, err = db.Conn.Query(
			`SELECT id, email, email_verified, is_super_admin, is_disabled, created_at
			 FROM users
			 WHERE LOWER(email) LIKE $1
			 ORDER BY created_at DESC
			 LIMIT $2 OFFSET $3`,
			"%"+query+"%", limit, offset,
		)
	} else {
		rows, err = db.Conn.Query(
			`SELECT id, email, email_verified, is_super_admin, is_disabled, created_at
			 FROM users
			 ORDER BY created_at DESC
			 LIMIT $1 OFFSET $2`,
			limit, offset,
		)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to list users"))
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id, email string
		var emailVerified, isSuper, isDisabled bool
		var createdAt time.Time
		if err := rows.Scan(&id, &email, &emailVerified, &isSuper, &isDisabled, &createdAt); err == nil {
			users = append(users, map[string]interface{}{
				"id":            id,
				"email":         email,
				"emailVerified": emailVerified,
				"isSuperAdmin":  isSuper,
				"isDisabled":    isDisabled,
				"createdAt":     createdAt,
			})
		}
	}
	return c.JSON(http.StatusOK, users)
}

func CreateUser(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	cfg := c.Get(middleware.ContextConfig).(*model.Config)
	req := new(AdminUserCreateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "email is required"))
	}
	password := strings.TrimSpace(req.Password)
	if password == "" && req.SendEmail {
		var err error
		password, err = auth.GenerateRefreshToken()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to generate password"))
		}
	}
	if password == "" || len(password) < 8 {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "password (min 8) is required"))
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to hash password"))
	}

	userID := uuid.NewV4().String()
	tx, err := db.Conn.Begin()
	if err != nil {
		log.Printf("[admin] create user: begin tx failed: %v", err)
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "db transaction failed"))
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`INSERT INTO users (id, email, password_hash, is_super_admin, email_verified) VALUES ($1,$2,$3,$4,TRUE)`,
		userID, email, hash, req.IsSuperAdmin); err != nil {
		log.Printf("[admin] create user: insert users failed for %s: %v", email, err)
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "email already exists"))
	}
	if _, err := tx.Exec(`INSERT INTO user_settings (user_id, region) VALUES ($1,$2)`, userID, cfg.DefaultRegion); err != nil {
		log.Printf("[admin] create user: insert user_settings failed for %s (region=%s): %v", email, cfg.DefaultRegion, err)
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create settings"))
	}
	tenantID := strings.TrimSpace(req.TenantID)
	if tenantID == "" {
		if val := c.Get(middleware.ContextTenantID); val != nil {
			if ctxTenantID, ok := val.(string); ok {
				tenantID = strings.TrimSpace(ctxTenantID)
			}
		}
	}
	if tenantID != "" {
		role := strings.TrimSpace(req.Role)
		if role == "" {
			role = "member"
		}
		if !isValidRole(role) {
			return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "role must be owner, admin, or member"))
		}
		if _, err := tx.Exec(`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1,$2,$3) ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = EXCLUDED.role`,
			userID, tenantID, role); err != nil {
			log.Printf("[admin] create user: insert membership failed for %s (tenant=%s role=%s): %v", email, tenantID, role, err)
			return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to add membership"))
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[admin] create user: commit failed for %s: %v", email, err)
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to commit user"))
	}

	if req.SendEmail {
		if err := sendPasswordEmail(email, password, true); err != nil {
			log.Printf("[admin] failed to send welcome email to %s: %v", email, err)
		}
	}

	logAdminAction(c, "user_created", "user", userID, tenantID, map[string]interface{}{"email": email})
	return c.JSON(http.StatusCreated, map[string]string{"id": userID})
}

func UpdateUser(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Param("id")
	req := new(AdminUserUpdateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	if strings.TrimSpace(req.Email) == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "email is required"))
	}
	if _, err := db.Conn.Exec(
		`UPDATE users SET email = $1, is_super_admin = $2, is_disabled = $3, email_verified = $4 WHERE id = $5`,
		strings.ToLower(strings.TrimSpace(req.Email)), req.IsSuperAdmin, req.IsDisabled, req.EmailVerified, userID,
	); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to update user"))
	}
	logAdminAction(c, "user_updated", "user", userID, "", map[string]interface{}{
		"email":         req.Email,
		"isSuperAdmin":  req.IsSuperAdmin,
		"isDisabled":    req.IsDisabled,
		"emailVerified": req.EmailVerified,
	})
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func DeleteUser(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Param("id")
	if _, err := db.Conn.Exec(`DELETE FROM users WHERE id = $1`, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to delete user"))
	}
	logAdminAction(c, "user_deleted", "user", userID, "", nil)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func ResetUserPassword(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Param("id")
	req := new(AdminResetPasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	password := strings.TrimSpace(req.Password)
	if password == "" && req.SendEmail {
		var err error
		password, err = auth.GenerateRefreshToken()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to generate password"))
		}
	}
	if password == "" || len(password) < 8 {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "password (min 8) is required"))
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to hash password"))
	}
	if _, err := db.Conn.Exec(`UPDATE users SET password_hash = $1, email_verified = TRUE WHERE id = $2`, hash, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to update password"))
	}

	if req.SendEmail {
		var email string
		if err := db.Conn.QueryRow(`SELECT email FROM users WHERE id = $1`, userID).Scan(&email); err == nil {
			if err := sendPasswordEmail(email, password, false); err != nil {
				log.Printf("[admin] failed to send reset email to %s: %v", email, err)
			}
		}
	}

	logAdminAction(c, "user_password_reset", "user", userID, "", map[string]interface{}{"sentEmail": req.SendEmail})
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func VerifyUserEmail(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Param("id")
	if _, err := db.Conn.Exec(`UPDATE users SET email_verified = TRUE WHERE id = $1`, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to verify user"))
	}
	logAdminAction(c, "user_email_verified", "user", userID, "", nil)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func ListTenantsAdmin(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	rows, err := db.Conn.Query(
		`SELECT t.id, t.name, t.created_at,
		 (SELECT COUNT(*) FROM memberships m WHERE m.tenant_id = t.id) AS members
		 FROM tenants t
		 ORDER BY t.created_at DESC`,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to list tenants"))
	}
	defer rows.Close()

	var tenants []map[string]interface{}
	for rows.Next() {
		var id, name string
		var createdAt time.Time
		var members int
		if err := rows.Scan(&id, &name, &createdAt, &members); err == nil {
			tenants = append(tenants, map[string]interface{}{
				"id":        id,
				"name":      name,
				"createdAt": createdAt,
				"members":   members,
			})
		}
	}
	return c.JSON(http.StatusOK, tenants)
}

func CreateTenantAdmin(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	req := new(AdminTenantRequest)
	if err := c.Bind(req); err != nil || strings.TrimSpace(req.Name) == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "name is required"))
	}
	tenantID := uuid.NewV4().String()
	tx, err := db.Conn.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "db transaction failed"))
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`INSERT INTO tenants (id, name) VALUES ($1,$2)`, tenantID, req.Name); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create tenant"))
	}

	ownerID := strings.TrimSpace(req.OwnerUserID)
	if ownerID == "" {
		ownerID = c.Get(middleware.ContextUserID).(string)
	}
	if _, err := tx.Exec(`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`,
		ownerID, tenantID, "owner"); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create membership"))
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to commit tenant"))
	}
	logAdminAction(c, "tenant_created", "tenant", tenantID, tenantID, map[string]interface{}{"name": req.Name})
	return c.JSON(http.StatusCreated, map[string]string{"id": tenantID})
}

func UpdateTenantAdmin(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	tenantID := c.Param("id")
	req := new(AdminTenantRequest)
	if err := c.Bind(req); err != nil || strings.TrimSpace(req.Name) == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "name is required"))
	}
	if _, err := db.Conn.Exec(`UPDATE tenants SET name = $1 WHERE id = $2`, req.Name, tenantID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to update tenant"))
	}
	logAdminAction(c, "tenant_updated", "tenant", tenantID, tenantID, map[string]interface{}{"name": req.Name})
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func DeleteTenantAdmin(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	tenantID := c.Param("id")
	if _, err := db.Conn.Exec(`DELETE FROM tenants WHERE id = $1`, tenantID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to delete tenant"))
	}
	logAdminAction(c, "tenant_deleted", "tenant", tenantID, tenantID, nil)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func ListTenantMembersAdmin(c echo.Context) error {
	tenantID := c.Param("id")
	return listMembersForTenant(c, tenantID)
}

func ListTenantMembersSelf(c echo.Context) error {
	tenantID := c.Get(middleware.ContextTenantID).(string)
	return listMembersForTenant(c, tenantID)
}

func AddTenantMemberAdmin(c echo.Context) error {
	tenantID := c.Param("id")
	return addTenantMember(c, tenantID)
}

func AddTenantMemberSelf(c echo.Context) error {
	tenantID := c.Get(middleware.ContextTenantID).(string)
	return addTenantMember(c, tenantID)
}

func UpdateTenantMemberAdmin(c echo.Context) error {
	tenantID := c.Param("id")
	userID := c.Param("userId")
	return updateTenantMemberRole(c, tenantID, userID)
}

func UpdateTenantMemberSelf(c echo.Context) error {
	tenantID := c.Get(middleware.ContextTenantID).(string)
	userID := c.Param("userId")
	return updateTenantMemberRole(c, tenantID, userID)
}

func RemoveTenantMemberAdmin(c echo.Context) error {
	tenantID := c.Param("id")
	userID := c.Param("userId")
	return removeTenantMember(c, tenantID, userID)
}

func RemoveTenantMemberSelf(c echo.Context) error {
	tenantID := c.Get(middleware.ContextTenantID).(string)
	userID := c.Param("userId")
	return removeTenantMember(c, tenantID, userID)
}

func ListAWSAccountsAdmin(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	rows, err := db.Conn.Query(
		`SELECT a.id, a.tenant_id, t.name, a.account_id, a.role_arn, a.external_id, a.name, a.active
		 FROM aws_accounts a
		 JOIN tenants t ON t.id = a.tenant_id
		 ORDER BY t.name, a.name`,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to list accounts"))
	}
	defer rows.Close()

	var accounts []map[string]interface{}
	for rows.Next() {
		var id, tenantID, tenantName, accountID, roleArn, externalID, name string
		var active bool
		if err := rows.Scan(&id, &tenantID, &tenantName, &accountID, &roleArn, &externalID, &name, &active); err == nil {
			accounts = append(accounts, map[string]interface{}{
				"id":         id,
				"tenantId":   tenantID,
				"tenantName": tenantName,
				"accountId":  accountID,
				"roleArn":    roleArn,
				"externalId": externalID,
				"name":       name,
				"active":     active,
			})
		}
	}
	return c.JSON(http.StatusOK, accounts)
}

func ListAWSAccountsSelf(c echo.Context) error {
	tenantID := c.Get(middleware.ContextTenantID).(string)
	return listAccountsForTenant(c, tenantID)
}

func ListAWSAccountsForTenantAdmin(c echo.Context) error {
	tenantID := c.Param("id")
	return listAccountsForTenant(c, tenantID)
}

func CreateAWSAccountAdmin(c echo.Context) error {
	tenantID := c.Param("id")
	return createAWSAccountForTenant(c, tenantID)
}

func CreateAWSAccountSelf(c echo.Context) error {
	tenantID := c.Get(middleware.ContextTenantID).(string)
	return createAWSAccountForTenant(c, tenantID)
}

func UpdateAWSAccountAdmin(c echo.Context) error {
	tenantID := c.Param("id")
	accountID := c.Param("accountId")
	return updateAWSAccountForTenant(c, tenantID, accountID)
}

func UpdateAWSAccountSelf(c echo.Context) error {
	tenantID := c.Get(middleware.ContextTenantID).(string)
	accountID := c.Param("accountId")
	return updateAWSAccountForTenant(c, tenantID, accountID)
}

func DeleteAWSAccountAdmin(c echo.Context) error {
	tenantID := c.Param("id")
	accountID := c.Param("accountId")
	return deleteAWSAccountForTenant(c, tenantID, accountID)
}

func DeleteAWSAccountSelf(c echo.Context) error {
	tenantID := c.Get(middleware.ContextTenantID).(string)
	accountID := c.Param("accountId")
	return deleteAWSAccountForTenant(c, tenantID, accountID)
}

func ListAuditLogs(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	limit := parseLimit(c.QueryParam("limit"), 200)
	tenantID := strings.TrimSpace(c.QueryParam("tenant_id"))
	isSuper := false
	if val := c.Get(middleware.ContextIsSuperAdmin); val != nil {
		if isAdmin, ok := val.(bool); ok {
			isSuper = isAdmin
		}
	}
	if !isSuper {
		tenantID = c.Get(middleware.ContextTenantID).(string)
	}

	var rows *sql.Rows
	var err error
	if tenantID != "" {
		rows, err = db.Conn.Query(
			`SELECT id, actor_user_id, actor_email, tenant_id, action, entity_type, entity_id, details, created_at
			 FROM audit_logs
			 WHERE tenant_id = $1
			 ORDER BY created_at DESC
			 LIMIT $2`,
			tenantID, limit,
		)
	} else {
		rows, err = db.Conn.Query(
			`SELECT id, actor_user_id, actor_email, tenant_id, action, entity_type, entity_id, details, created_at
			 FROM audit_logs
			 ORDER BY created_at DESC
			 LIMIT $1`,
			limit,
		)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to list audit logs"))
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id, actorEmail string
		var actorID, tenant sql.NullString
		var action, entityType string
		var entityID sql.NullString
		var detailsBytes []byte
		var createdAt time.Time
		if err := rows.Scan(&id, &actorID, &actorEmail, &tenant, &action, &entityType, &entityID, &detailsBytes, &createdAt); err == nil {
			var details interface{}
			if len(detailsBytes) > 0 {
				_ = json.Unmarshal(detailsBytes, &details)
			}
			payload := map[string]interface{}{
				"id":         id,
				"actorUserId": actorID.String,
				"actorEmail": actorEmail,
				"tenantId":   tenant.String,
				"action":     action,
				"entityType": entityType,
				"entityId":   entityID.String,
				"details":    details,
				"createdAt":  createdAt,
			}
			logs = append(logs, payload)
		}
	}
	return c.JSON(http.StatusOK, logs)
}

func listMembersForTenant(c echo.Context, tenantID string) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	rows, err := db.Conn.Query(
		`SELECT u.id, u.email, u.email_verified, u.is_disabled, m.role, m.created_at
		 FROM memberships m
		 JOIN users u ON u.id = m.user_id
		 WHERE m.tenant_id = $1
		 ORDER BY u.email`,
		tenantID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to list members"))
	}
	defer rows.Close()

	var members []map[string]interface{}
	for rows.Next() {
		var userID, email, role string
		var emailVerified, isDisabled bool
		var createdAt time.Time
		if err := rows.Scan(&userID, &email, &emailVerified, &isDisabled, &role, &createdAt); err == nil {
			members = append(members, map[string]interface{}{
				"userId":       userID,
				"email":        email,
				"role":         role,
				"emailVerified": emailVerified,
				"isDisabled":   isDisabled,
				"createdAt":    createdAt,
			})
		}
	}
	return c.JSON(http.StatusOK, members)
}

func addTenantMember(c echo.Context, tenantID string) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	cfg := c.Get(middleware.ContextConfig).(*model.Config)
	req := new(AdminTenantMemberRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	role := strings.TrimSpace(req.Role)
	if role == "" {
		role = "member"
	}
	if !isValidRole(role) {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "role must be owner, admin, or member"))
	}

	var userID, email string
	if strings.TrimSpace(req.UserID) != "" {
		if err := db.Conn.QueryRow(`SELECT id, email FROM users WHERE id = $1`, strings.TrimSpace(req.UserID)).Scan(&userID, &email); err != nil {
			return c.JSON(http.StatusNotFound, utils.SendError("not_found", "user not found"))
		}
	} else {
		email = strings.ToLower(strings.TrimSpace(req.Email))
		if email == "" {
			return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "email is required"))
		}
		err := db.Conn.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&userID)
		if err == sql.ErrNoRows && req.CreateIfMissing {
			password := strings.TrimSpace(req.Password)
			if password == "" && req.SendEmail {
				var genErr error
				password, genErr = auth.GenerateRefreshToken()
				if genErr != nil {
					return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to generate password"))
				}
			}
			if password == "" || len(password) < 8 {
				return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "password (min 8) is required"))
			}
			hash, err := auth.HashPassword(password)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to hash password"))
			}
			userID = uuid.NewV4().String()
			tx, err := db.Conn.Begin()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "db transaction failed"))
			}
			defer tx.Rollback()
			if _, err := tx.Exec(`INSERT INTO users (id, email, password_hash, email_verified) VALUES ($1,$2,$3,TRUE)`, userID, email, hash); err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create user"))
			}
			if _, err := tx.Exec(`INSERT INTO user_settings (user_id, region) VALUES ($1,$2)`, userID, cfg.DefaultRegion); err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create settings"))
			}
			if _, err := tx.Exec(`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1,$2,$3) ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = EXCLUDED.role`,
				userID, tenantID, role); err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to add membership"))
			}
			if err := tx.Commit(); err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to commit member"))
			}
			if req.SendEmail {
				_ = sendPasswordEmail(email, password, true)
			}
			logAdminAction(c, "tenant_member_created", "membership", userID, tenantID, map[string]interface{}{
				"email": email,
				"role":  role,
			})
			return c.JSON(http.StatusCreated, map[string]string{"userId": userID})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to lookup user"))
		}
	}

	if _, err := db.Conn.Exec(`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1,$2,$3) ON CONFLICT (user_id, tenant_id) DO UPDATE SET role = EXCLUDED.role`,
		userID, tenantID, role); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to add membership"))
	}
	logAdminAction(c, "tenant_member_added", "membership", userID, tenantID, map[string]interface{}{"role": role})
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func updateTenantMemberRole(c echo.Context, tenantID, userID string) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	req := new(AdminTenantMemberUpdateRequest)
	if err := c.Bind(req); err != nil || strings.TrimSpace(req.Role) == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "role is required"))
	}
	if !isValidRole(req.Role) {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "role must be owner, admin, or member"))
	}
	if _, err := db.Conn.Exec(`UPDATE memberships SET role = $1 WHERE tenant_id = $2 AND user_id = $3`, req.Role, tenantID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to update role"))
	}
	logAdminAction(c, "tenant_member_role_updated", "membership", userID, tenantID, map[string]interface{}{"role": req.Role})
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func removeTenantMember(c echo.Context, tenantID, userID string) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	if _, err := db.Conn.Exec(`DELETE FROM memberships WHERE tenant_id = $1 AND user_id = $2`, tenantID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to remove member"))
	}
	logAdminAction(c, "tenant_member_removed", "membership", userID, tenantID, nil)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func listAccountsForTenant(c echo.Context, tenantID string) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	rows, err := db.Conn.Query(
		`SELECT id, account_id, role_arn, external_id, name, active
		 FROM aws_accounts
		 WHERE tenant_id = $1
		 ORDER BY name`,
		tenantID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to load accounts"))
	}
	defer rows.Close()

	var accounts []map[string]interface{}
	for rows.Next() {
		var id, accountID, roleArn, externalID, name string
		var active bool
		if err := rows.Scan(&id, &accountID, &roleArn, &externalID, &name, &active); err == nil {
			accounts = append(accounts, map[string]interface{}{
				"id":         id,
				"accountId":  accountID,
				"roleArn":    roleArn,
				"externalId": externalID,
				"name":       name,
				"active":     active,
			})
		}
	}
	return c.JSON(http.StatusOK, accounts)
}

func createAWSAccountForTenant(c echo.Context, tenantID string) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	req := new(AdminAWSAccountRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	if strings.TrimSpace(req.AccountID) == "" || strings.TrimSpace(req.RoleArn) == "" || strings.TrimSpace(req.ExternalID) == "" || strings.TrimSpace(req.Name) == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "accountId, roleArn, externalId, name are required"))
	}
	id := uuid.NewV4().String()
	if _, err := db.Conn.Exec(
		`INSERT INTO aws_accounts (id, tenant_id, account_id, role_arn, external_id, name) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, tenantID, req.AccountID, req.RoleArn, req.ExternalID, req.Name,
	); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create account"))
	}
	logAdminAction(c, "aws_account_created", "aws_account", id, tenantID, map[string]interface{}{"name": req.Name})
	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}

func updateAWSAccountForTenant(c echo.Context, tenantID, accountID string) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	req := new(AdminAWSAccountRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	_, err := db.Conn.Exec(
		`UPDATE aws_accounts SET account_id = COALESCE(NULLIF($1,''), account_id),
		 role_arn = COALESCE(NULLIF($2,''), role_arn),
		 external_id = COALESCE(NULLIF($3,''), external_id),
		 name = COALESCE(NULLIF($4,''), name),
		 active = COALESCE($5, active)
		 WHERE id = $6 AND tenant_id = $7`,
		req.AccountID, req.RoleArn, req.ExternalID, req.Name, req.Active, accountID, tenantID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to update account"))
	}
	logAdminAction(c, "aws_account_updated", "aws_account", accountID, tenantID, nil)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func deleteAWSAccountForTenant(c echo.Context, tenantID, accountID string) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	if _, err := db.Conn.Exec(`DELETE FROM aws_accounts WHERE id = $1 AND tenant_id = $2`, accountID, tenantID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to delete account"))
	}
	logAdminAction(c, "aws_account_deleted", "aws_account", accountID, tenantID, nil)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func logAdminAction(c echo.Context, action, entityType, entityID, tenantID string, details map[string]interface{}) {
	db := c.Get(middleware.ContextDB).(*store.DB)
	actorID := c.Get(middleware.ContextUserID).(string)
	actorEmail := ""
	if val := c.Get(middleware.ContextUserEmail); val != nil {
		if email, ok := val.(string); ok {
			actorEmail = email
		}
	}
	_ = store.LogAudit(db, actorID, actorEmail, action, entityType, entityID, tenantID, details)
}

func sendPasswordEmail(email, password string, isNewUser bool) error {
	cfg, err := utils.LoadEmailConfig()
	if err != nil {
		return err
	}
	loginURL := buildLoginURL()
	var subject string
	var bodyLines []string
	if isNewUser {
		subject = "Welcome to AWS Admin"
		bodyLines = []string{
			"Welcome to AWS Admin!",
			"",
			"Your account has been created by an administrator.",
			"Your temporary password is: " + password,
		}
	} else {
		subject = "Your AWS Admin password reset"
		bodyLines = []string{
			"Your AWS Admin password was reset by an administrator.",
			"Your temporary password is: " + password,
		}
	}
	if loginURL != "" {
		bodyLines = append(bodyLines, "Login: "+loginURL)
	}
	bodyLines = append(bodyLines, "Please log in and change it.")
	body := strings.Join(bodyLines, "\n")
	return utils.SendEmail(cfg, email, subject, body)
}

func buildLoginURL() string {
	appBase := strings.TrimSpace(os.Getenv("APP_BASE_URL"))
	if appBase == "" {
		return ""
	}
	return strings.TrimRight(appBase, "/") + "/"
}

func parseLimit(value string, def int) int {
	if value == "" {
		return def
	}
	if v, err := strconv.Atoi(value); err == nil && v >= 0 {
		return v
	}
	return def
}

func isValidRole(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "owner", "admin", "member":
		return true
	default:
		return false
	}
}
