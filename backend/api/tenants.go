package api

import (
	"net/http"
	"os"
	"strings"
	"time"

<<<<<<< HEAD
	"github.com/bertrandmartel/aws-admin/backend/auth"
	"github.com/bertrandmartel/aws-admin/backend/middleware"
	"github.com/bertrandmartel/aws-admin/backend/model"
=======
	"github.com/bertrandmartel/aws-admin/backend/middleware"
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
	"github.com/bertrandmartel/aws-admin/backend/store"
	"github.com/bertrandmartel/aws-admin/backend/utils"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type TenantRequest struct {
	Name string `json:"name"`
}

type InviteRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AcceptInviteRequest struct {
	Token string `json:"token"`
}

<<<<<<< HEAD
type AcceptInvitePublicRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
func ListTenants(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Get(middleware.ContextUserID).(string)
	rows, err := db.Conn.Query(`SELECT t.id, t.name, m.role FROM memberships m JOIN tenants t ON t.id = m.tenant_id WHERE m.user_id = $1 ORDER BY t.name`, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to list tenants"))
	}
	defer rows.Close()

	var tenants []map[string]string
	for rows.Next() {
		var id, name, role string
		if err := rows.Scan(&id, &name, &role); err == nil {
			tenants = append(tenants, map[string]string{
				"id":   id,
				"name": name,
				"role": role,
			})
		}
	}
	return c.JSON(http.StatusOK, tenants)
}

func CreateTenant(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Get(middleware.ContextUserID).(string)
	req := new(TenantRequest)
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
	if _, err := tx.Exec(`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1,$2,$3)`, userID, tenantID, "owner"); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create membership"))
	}
	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to commit tenant"))
	}
	return c.JSON(http.StatusCreated, map[string]string{"id": tenantID})
}

func InviteTenant(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Get(middleware.ContextUserID).(string)
	tenantID := c.Param("id")

	var role string
	if err := db.Conn.QueryRow(`SELECT role FROM memberships WHERE user_id = $1 AND tenant_id = $2`, userID, tenantID).Scan(&role); err != nil || (role != "owner" && role != "admin") {
		return c.JSON(http.StatusForbidden, utils.SendError("forbidden", "admin role required"))
	}

	req := new(InviteRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || strings.TrimSpace(req.Role) == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "email and role are required"))
	}
	if req.Role != "owner" && req.Role != "admin" && req.Role != "member" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "role must be owner, admin, or member"))
	}

	token, err := store.NewInviteToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create token"))
	}
	tokenHash := store.HashToken(token)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	if _, err := db.Conn.Exec(`INSERT INTO tenant_invites (token_hash, tenant_id, email, role, expires_at) VALUES ($1,$2,$3,$4,$5)`,
		tokenHash, tenantID, email, req.Role, expiresAt); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create invite"))
	}

	appBase := os.Getenv("APP_BASE_URL")
	inviteLink := appBase + "/invite?token=" + token
	emailCfg, err := utils.LoadEmailConfig()
	if err == nil {
		_ = utils.SendEmail(emailCfg, email, "You're invited to AWS Admin", "Accept invite: "+inviteLink)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func AcceptInvite(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Get(middleware.ContextUserID).(string)
	req := new(AcceptInviteRequest)
	if err := c.Bind(req); err != nil || req.Token == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "token is required"))
	}
	tokenHash := store.HashToken(req.Token)

	var tenantID, email, role string
	var expires time.Time
	if err := db.Conn.QueryRow(`SELECT tenant_id, email, role, expires_at FROM tenant_invites WHERE token_hash = $1`, tokenHash).Scan(&tenantID, &email, &role, &expires); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("invalid_token", "invite not found"))
	}
	if time.Now().After(expires) {
		return c.JSON(http.StatusBadRequest, utils.SendError("invalid_token", "invite expired"))
	}

	var userEmail string
	_ = db.Conn.QueryRow(`SELECT email FROM users WHERE id = $1`, userID).Scan(&userEmail)
	if strings.ToLower(userEmail) != strings.ToLower(email) {
		return c.JSON(http.StatusForbidden, utils.SendError("forbidden", "invite email mismatch"))
	}

	if _, err := db.Conn.Exec(`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`, userID, tenantID, role); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to add membership"))
	}
	_, _ = db.Conn.Exec(`DELETE FROM tenant_invites WHERE token_hash = $1`, tokenHash)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

<<<<<<< HEAD
func AcceptInvitePublic(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	cfg := c.Get(middleware.ContextConfig).(*model.Config)
	req := new(AcceptInvitePublicRequest)
	if err := c.Bind(req); err != nil || req.Token == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "token is required"))
	}

	tokenHash := store.HashToken(req.Token)
	var tenantID, email, role string
	var expires time.Time
	if err := db.Conn.QueryRow(`SELECT tenant_id, email, role, expires_at FROM tenant_invites WHERE token_hash = $1`, tokenHash).Scan(&tenantID, &email, &role, &expires); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("invalid_token", "invite not found"))
	}
	if time.Now().After(expires) {
		return c.JSON(http.StatusBadRequest, utils.SendError("invalid_token", "invite expired"))
	}

	email = strings.ToLower(strings.TrimSpace(email))
	var userID string
	var existing bool
	if err := db.Conn.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&userID); err == nil {
		existing = true
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "db transaction failed"))
	}
	defer tx.Rollback()

	if !existing {
		password := strings.TrimSpace(req.Password)
		if password == "" || len(password) < 8 {
			return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "password (min 8) is required"))
		}
		hash, err := auth.HashPassword(password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to hash password"))
		}
		userID = uuid.NewV4().String()
		if _, err := tx.Exec(`INSERT INTO users (id, email, password_hash, email_verified) VALUES ($1,$2,$3,TRUE)`, userID, email, hash); err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create user"))
		}
		if _, err := tx.Exec(`INSERT INTO user_settings (user_id, region) VALUES ($1,$2)`, userID, cfg.DefaultRegion); err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create settings"))
		}
	} else {
		_, _ = tx.Exec(`UPDATE users SET email_verified = TRUE WHERE id = $1`, userID)
	}

	if _, err := tx.Exec(`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`,
		userID, tenantID, role); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to add membership"))
	}
	_, _ = tx.Exec(`DELETE FROM tenant_invites WHERE token_hash = $1`, tokenHash)

	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to accept invite"))
	}

	_ = store.LogAudit(db, userID, email, "tenant_invite_accepted", "membership", userID, tenantID, map[string]interface{}{
		"role":     role,
		"existing": existing,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":       "ok",
		"requiresLogin": existing,
	})
}

=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
func SwitchTenant(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Get(middleware.ContextUserID).(string)
	tenantID := c.Param("id")

	var exists bool
	if err := db.Conn.QueryRow(`SELECT EXISTS(SELECT 1 FROM memberships WHERE user_id = $1 AND tenant_id = $2)`, userID, tenantID).Scan(&exists); err != nil || !exists {
		return c.JSON(http.StatusForbidden, utils.SendError("forbidden", "not a member of tenant"))
	}

	resp, err := issueTokens(db, userID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create tokens"))
	}
	return c.JSON(http.StatusOK, resp)
}
