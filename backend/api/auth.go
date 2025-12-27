package api

import (
	"net/http"
	"os"
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

type SignupRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	TenantName string `json:"tenantName"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
	TenantID     string `json:"tenantId"`
}

func issueTokens(db *store.DB, userID, tenantID string) (*AuthResponse, error) {
	accessToken, err := auth.GenerateAccessToken(userID, tenantID, os.Getenv("JWT_SECRET"), 15*time.Minute)
	if err != nil {
		return nil, err
	}
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshID := uuid.NewV4().String()
	refreshHash := store.HashToken(refreshToken)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	if _, err := db.Conn.Exec(`INSERT INTO refresh_sessions (id, user_id, tenant_id, refresh_token_hash, expires_at) VALUES ($1,$2,$3,$4,$5)`,
		refreshID, userID, tenantID, refreshHash, expiresAt); err != nil {
		return nil, err
	}
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userID,
		TenantID:     tenantID,
	}, nil
}

func Signup(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	cfg := c.Get(middleware.ContextConfig).(*model.Config)

	req := new(SignupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || len(req.Password) < 8 || strings.TrimSpace(req.TenantName) == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "email, password (min 8), and tenantName are required"))
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to hash password"))
	}

	userID := uuid.NewV4().String()
	tenantID := uuid.NewV4().String()

	tx, err := db.Conn.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "db transaction failed"))
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`INSERT INTO users (id, email, password_hash) VALUES ($1,$2,$3)`, userID, email, hash); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "email already exists"))
	}
	if _, err := tx.Exec(`INSERT INTO tenants (id, name) VALUES ($1,$2)`, tenantID, req.TenantName); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create tenant"))
	}
	if _, err := tx.Exec(`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1,$2,$3)`, userID, tenantID, "owner"); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create membership"))
	}
	if _, err := tx.Exec(`INSERT INTO user_settings (user_id, region) VALUES ($1,$2)`, userID, cfg.DefaultRegion); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to initialize settings"))
	}

	token, err := auth.GenerateRefreshToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to generate verification token"))
	}
	tokenHash := store.HashToken(token)
	expiresAt := time.Now().Add(24 * time.Hour)
	if _, err := tx.Exec(`INSERT INTO email_verifications (token_hash, user_id, expires_at) VALUES ($1,$2,$3)`, tokenHash, userID, expiresAt); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create verification token"))
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to commit signup"))
	}

	appBase := os.Getenv("APP_BASE_URL")
	verifyLink := appBase + "/verify?token=" + token
	emailCfg, err := utils.LoadEmailConfig()
	if err == nil {
		_ = utils.SendEmail(emailCfg, email, "Verify your AWS Admin account", "Click to verify: "+verifyLink)
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
}

func VerifyEmail(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "token is required"))
	}
	tokenHash := store.HashToken(token)

	var userID string
	var expires time.Time
	err := db.Conn.QueryRow(`SELECT user_id, expires_at FROM email_verifications WHERE token_hash = $1`, tokenHash).Scan(&userID, &expires)
	if err != nil || time.Now().After(expires) {
		return c.JSON(http.StatusBadRequest, utils.SendError("invalid_token", "verification token is invalid or expired"))
	}

	if _, err := db.Conn.Exec(`UPDATE users SET email_verified = TRUE WHERE id = $1`, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to verify email"))
	}
	_, _ = db.Conn.Exec(`DELETE FROM email_verifications WHERE token_hash = $1`, tokenHash)
	return c.JSON(http.StatusOK, map[string]string{"status": "verified"})
}

func Login(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)

	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "email and password are required"))
	}

	var userID, passHash string
	var verified bool
	var disabled bool
	err := db.Conn.QueryRow(`SELECT id, password_hash, email_verified, is_disabled FROM users WHERE email = $1`, email).Scan(&userID, &passHash, &verified, &disabled)
	if err != nil || !auth.CheckPassword(passHash, req.Password) {
		return c.JSON(http.StatusUnauthorized, utils.SendError("unauthorized", "invalid credentials"))
	}
	if disabled {
		return c.JSON(http.StatusForbidden, utils.SendError("forbidden", "account disabled"))
	}
	if !verified {
		return c.JSON(http.StatusForbidden, utils.SendError("email_not_verified", "email verification required"))
	}

	var tenantID string
	if err := db.Conn.QueryRow(`SELECT tenant_id FROM memberships WHERE user_id = $1 LIMIT 1`, userID).Scan(&tenantID); err != nil {
		return c.JSON(http.StatusForbidden, utils.SendError("forbidden", "no tenant membership"))
	}

	resp, err := issueTokens(db, userID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create tokens"))
	}
	return c.JSON(http.StatusOK, resp)
}

func Refresh(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	req := new(RefreshRequest)
	if err := c.Bind(req); err != nil || req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "refreshToken is required"))
	}
	refreshHash := store.HashToken(req.RefreshToken)

	var sessionID, userID, tenantID string
	var expires time.Time
	err := db.Conn.QueryRow(`SELECT id, user_id, tenant_id, expires_at FROM refresh_sessions WHERE refresh_token_hash = $1`, refreshHash).Scan(&sessionID, &userID, &tenantID, &expires)
	if err != nil || time.Now().After(expires) {
		return c.JSON(http.StatusUnauthorized, utils.SendError("unauthorized", "refresh token invalid"))
	}

	accessToken, err := auth.GenerateAccessToken(userID, tenantID, os.Getenv("JWT_SECRET"), 15*time.Minute)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create access token"))
	}

	newRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to rotate refresh token"))
	}
	newHash := store.HashToken(newRefresh)
	newExpires := time.Now().Add(30 * 24 * time.Hour)

	if _, err := db.Conn.Exec(`UPDATE refresh_sessions SET refresh_token_hash = $1, expires_at = $2 WHERE id = $3`, newHash, newExpires, sessionID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to rotate refresh token"))
	}

	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefresh,
		UserID:       userID,
		TenantID:     tenantID,
	})
}

func Logout(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	req := new(RefreshRequest)
	if err := c.Bind(req); err != nil || req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "refreshToken is required"))
	}
	refreshHash := store.HashToken(req.RefreshToken)
	_, _ = db.Conn.Exec(`DELETE FROM refresh_sessions WHERE refresh_token_hash = $1`, refreshHash)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func Me(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Get(middleware.ContextUserID).(string)
	tenantID := c.Get(middleware.ContextTenantID).(string)

	var email, region string
	var isSuper bool
	_ = db.Conn.QueryRow(`SELECT email, is_super_admin FROM users WHERE id = $1`, userID).Scan(&email, &isSuper)
	_ = db.Conn.QueryRow(`SELECT region FROM user_settings WHERE user_id = $1`, userID).Scan(&region)
	if !isSuper && isEmailSuperAdmin(email) {
		isSuper = true
	}

	var tenantRole string
	_ = db.Conn.QueryRow(`SELECT role FROM memberships WHERE user_id = $1 AND tenant_id = $2`, userID, tenantID).Scan(&tenantRole)

	type AWSAccount struct {
		ID        string `json:"id"`
		AccountID string `json:"accountId"`
		Name      string `json:"name"`
		RoleArn   string `json:"roleArn"`
		Active    bool   `json:"active"`
	}
	var accounts []AWSAccount
	rows, _ := db.Conn.Query(`SELECT id, account_id, name, role_arn, active FROM aws_accounts WHERE tenant_id = $1 ORDER BY name`, tenantID)
	defer rows.Close()
	for rows.Next() {
		var a AWSAccount
		if err := rows.Scan(&a.ID, &a.AccountID, &a.Name, &a.RoleArn, &a.Active); err == nil {
			accounts = append(accounts, a)
		}
	}

	type Tenant struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	}
	var tenants []Tenant
	tRows, _ := db.Conn.Query(`SELECT t.id, t.name, m.role FROM memberships m JOIN tenants t ON t.id = m.tenant_id WHERE m.user_id = $1 ORDER BY t.name`, userID)
	defer tRows.Close()
	for tRows.Next() {
		var t Tenant
		if err := tRows.Scan(&t.ID, &t.Name, &t.Role); err == nil {
			tenants = append(tenants, t)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"userId":    userID,
		"tenantId":  tenantID,
		"email":     email,
		"region":    region,
		"role":      tenantRole,
		"isSuperAdmin": isSuper,
		"accounts":  accounts,
		"tenants":   tenants,
	})
}

func ForgotPassword(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	req := new(ForgotPasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "email is required"))
	}

	var userID string
	if err := db.Conn.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&userID); err != nil {
		// Do not leak account existence.
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}

	token, err := auth.GenerateRefreshToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to generate token"))
	}
	tokenHash := store.HashToken(token)
	expiresAt := time.Now().Add(2 * time.Hour)
	if _, err := db.Conn.Exec(`INSERT INTO password_resets (token_hash, user_id, expires_at) VALUES ($1,$2,$3)`,
		tokenHash, userID, expiresAt); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create reset token"))
	}

	appBase := os.Getenv("APP_BASE_URL")
	resetLink := appBase + "/reset?token=" + token
	emailCfg, err := utils.LoadEmailConfig()
	if err == nil && strings.TrimSpace(appBase) != "" {
		_ = utils.SendEmail(emailCfg, email, "Reset your AWS Admin password", "Reset your password: "+resetLink)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func ResetPassword(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	req := new(ResetPasswordRequest)
	if err := c.Bind(req); err != nil || req.Token == "" || strings.TrimSpace(req.Password) == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "token and password are required"))
	}
	if len(strings.TrimSpace(req.Password)) < 8 {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "password (min 8) is required"))
	}

	tokenHash := store.HashToken(req.Token)
	var userID string
	var expires time.Time
	err := db.Conn.QueryRow(`SELECT user_id, expires_at FROM password_resets WHERE token_hash = $1`, tokenHash).Scan(&userID, &expires)
	if err != nil || time.Now().After(expires) {
		return c.JSON(http.StatusBadRequest, utils.SendError("invalid_token", "reset token is invalid or expired"))
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to hash password"))
	}
	if _, err := db.Conn.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, hash, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to reset password"))
	}
	_, _ = db.Conn.Exec(`DELETE FROM password_resets WHERE token_hash = $1`, tokenHash)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func isEmailSuperAdmin(email string) bool {
	if email == "" {
		return false
	}
	list := strings.Split(os.Getenv("SUPER_ADMIN_EMAILS"), ",")
	for _, item := range list {
		if strings.ToLower(strings.TrimSpace(item)) == strings.ToLower(strings.TrimSpace(email)) {
			return true
		}
	}
	return false
}
