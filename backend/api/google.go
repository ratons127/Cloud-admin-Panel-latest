package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bertrandmartel/aws-admin/backend/auth"
	"github.com/bertrandmartel/aws-admin/backend/middleware"
	"github.com/bertrandmartel/aws-admin/backend/model"
	"github.com/bertrandmartel/aws-admin/backend/store"
	"github.com/bertrandmartel/aws-admin/backend/utils"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleUserInfo struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
}

func GoogleLogin(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	state, err := store.CreateOAuthState(db)
	if err != nil || state == "" {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create oauth state"))
	}
	cfg, err := googleOAuthConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", err.Error()))
	}
	url := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	state := c.QueryParam("state")
	code := c.QueryParam("code")
	if state == "" || code == "" {
		return redirectOAuthError(c, "missing_state_or_code")
	}
	ok, err := store.ConsumeOAuthState(db, state)
	if err != nil || !ok {
		return redirectOAuthError(c, "invalid_state")
	}
	cfg, err := googleOAuthConfig()
	if err != nil {
		return redirectOAuthError(c, "missing_google_config")
	}
	ctx := context.Background()
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return redirectOAuthError(c, "token_exchange_failed")
	}
	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode >= 400 {
		return redirectOAuthError(c, "userinfo_failed")
	}
	defer resp.Body.Close()

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return redirectOAuthError(c, "userinfo_decode_failed")
	}
	email := strings.ToLower(strings.TrimSpace(info.Email))
	if email == "" || !info.VerifiedEmail {
		return redirectOAuthError(c, "email_not_verified")
	}

	tenantID, userID, err := findOrCreateUserForGoogle(db, c, email)
	if err != nil {
		return redirectOAuthError(c, "user_lookup_failed")
	}
	if tenantID == "" || userID == "" {
		return redirectOAuthError(c, "no_tenant")
	}

	respTokens, err := issueTokens(db, userID, tenantID)
	if err != nil {
		return redirectOAuthError(c, "token_issue_failed")
	}

	appBase := strings.TrimRight(os.Getenv("APP_BASE_URL"), "/")
	if appBase == "" {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "APP_BASE_URL is required"))
	}
	redirectURL := appBase + "/oauth?accessToken=" + url.QueryEscape(respTokens.AccessToken) + "&refreshToken=" + url.QueryEscape(respTokens.RefreshToken)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func googleOAuthConfig() (*oauth2.Config, error) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET"))
	redirectURL := strings.TrimSpace(os.Getenv("GOOGLE_REDIRECT_URL"))
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REDIRECT_URL are required")
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}, nil
}

func findOrCreateUserForGoogle(db *store.DB, c echo.Context, email string) (string, string, error) {
	var userID string
	var disabled bool
	err := db.Conn.QueryRow(`SELECT id, is_disabled FROM users WHERE email = $1`, email).Scan(&userID, &disabled)
	if err == nil {
		if disabled {
			return "", "", errors.New("account_disabled")
		}
		_, _ = db.Conn.Exec(`UPDATE users SET email_verified = TRUE WHERE id = $1`, userID)
		return firstTenantForUser(db, userID)
	}

	if strings.ToLower(strings.TrimSpace(os.Getenv("GOOGLE_AUTO_CREATE"))) != "true" {
		return "", "", errors.New("user_not_found")
	}

	cfg := c.Get(middleware.ContextConfig).(*model.Config)
	tenantID := uuid.NewV4().String()
	newUserID := uuid.NewV4().String()
	password, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return "", "", err
	}
	tx, err := db.Conn.Begin()
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback()

	tenantName := strings.Split(email, "@")[0] + "'s Tenant"
	if _, err := tx.Exec(`INSERT INTO users (id, email, password_hash, email_verified) VALUES ($1,$2,$3,TRUE)`, newUserID, email, hash); err != nil {
		return "", "", err
	}
	if _, err := tx.Exec(`INSERT INTO tenants (id, name) VALUES ($1,$2)`, tenantID, tenantName); err != nil {
		return "", "", err
	}
	if _, err := tx.Exec(`INSERT INTO memberships (user_id, tenant_id, role) VALUES ($1,$2,$3)`, newUserID, tenantID, "owner"); err != nil {
		return "", "", err
	}
	if _, err := tx.Exec(`INSERT INTO user_settings (user_id, region) VALUES ($1,$2)`, newUserID, cfg.DefaultRegion); err != nil {
		return "", "", err
	}
	if err := tx.Commit(); err != nil {
		return "", "", err
	}
	_ = store.LogAudit(db, newUserID, email, "google_user_created", "user", newUserID, tenantID, map[string]interface{}{
		"tenantName": tenantName,
	})
	return tenantID, newUserID, nil
}

func firstTenantForUser(db *store.DB, userID string) (string, string, error) {
	var tenantID string
	if err := db.Conn.QueryRow(`SELECT tenant_id FROM memberships WHERE user_id = $1 ORDER BY created_at LIMIT 1`, userID).Scan(&tenantID); err != nil {
		return "", "", err
	}
	return tenantID, userID, nil
}

func redirectOAuthError(c echo.Context, reason string) error {
	appBase := strings.TrimRight(os.Getenv("APP_BASE_URL"), "/")
	if appBase == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("oauth_error", reason))
	}
	return c.Redirect(http.StatusTemporaryRedirect, appBase+"/oauth?error="+url.QueryEscape(reason))
}
