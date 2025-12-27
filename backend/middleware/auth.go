package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/bertrandmartel/aws-admin/backend/auth"
	"github.com/bertrandmartel/aws-admin/backend/model"
	"github.com/bertrandmartel/aws-admin/backend/session"
	"github.com/bertrandmartel/aws-admin/backend/store"
	"github.com/labstack/echo/v4"
)

const (
	ContextDB        = "db"
	ContextConfig    = "config"
	ContextUserID    = "user_id"
	ContextUserEmail = "user_email"
	ContextTenantID  = "tenant_id"
	ContextRegion    = "region"
	ContextAWSAcctID = "aws_account_id"
	ContextIsSuperAdmin = "is_super_admin"
)

func AttachDB(db *store.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(ContextDB, db)
			return next(c)
		}
	}
}

func AttachConfig(cfg *model.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(ContextConfig, cfg)
			return next(c)
		}
	}
}

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized", ErrorDescription: "missing bearer token"})
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ParseAccessToken(tokenStr, os.Getenv("JWT_SECRET"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized", ErrorDescription: "invalid token"})
		}
		if dbVal := c.Get(ContextDB); dbVal != nil {
			db := dbVal.(*store.DB)
			var email string
			var disabled bool
			var isSuper bool
			if err := db.Conn.QueryRow(`SELECT email, is_disabled, is_super_admin FROM users WHERE id = $1`, claims.UserID).Scan(&email, &disabled, &isSuper); err != nil {
				return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized", ErrorDescription: "invalid user"})
			}
			if disabled {
				return c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "forbidden", ErrorDescription: "account disabled"})
			}
			c.Set(ContextUserEmail, email)
			if isSuper || isEmailSuperAdmin(email) {
				c.Set(ContextIsSuperAdmin, true)
			} else {
				c.Set(ContextIsSuperAdmin, false)
			}
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextTenantID, claims.TenantID)
		return next(c)
	}
}

func RequireSuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := c.Get(ContextDB).(*store.DB)
		userID := c.Get(ContextUserID).(string)
		email := resolveUserEmail(c)
		isSuper, err := userIsSuperAdmin(db, userID, email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal_error", ErrorDescription: "failed to validate admin"})
		}
		if !isSuper {
			return c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "forbidden", ErrorDescription: "super admin required"})
		}
		c.Set(ContextIsSuperAdmin, true)
		return next(c)
	}
}

func RequireTenantAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := c.Get(ContextDB).(*store.DB)
		userID := c.Get(ContextUserID).(string)
		tenantID := c.Get(ContextTenantID).(string)
		email := resolveUserEmail(c)
		isSuper, err := userIsSuperAdmin(db, userID, email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal_error", ErrorDescription: "failed to validate admin"})
		}
		if isSuper {
			return next(c)
		}
		var role string
		if err := db.Conn.QueryRow(`SELECT role FROM memberships WHERE user_id = $1 AND tenant_id = $2`, userID, tenantID).Scan(&role); err != nil {
			return c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "forbidden", ErrorDescription: "tenant admin required"})
		}
		if !isTenantAdminRole(role) {
			return c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "forbidden", ErrorDescription: "tenant admin required"})
		}
		return next(c)
	}
}

func RequireAWSAccount(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := c.Get(ContextDB).(*store.DB)
		tenantID := c.Get(ContextTenantID).(string)
		accountID := c.Request().Header.Get("X-AWS-ACCOUNT-ID")
		if accountID == "" {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "format_error", ErrorDescription: "X-AWS-ACCOUNT-ID header is required"})
		}

		var roleArn, externalID, name string
		var active bool
		err := db.Conn.QueryRow(
			`SELECT role_arn, external_id, name, active FROM aws_accounts WHERE id = $1 AND tenant_id = $2`,
			accountID, tenantID,
		).Scan(&roleArn, &externalID, &name, &active)
		if err != nil {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "not_found", ErrorDescription: "aws account not found"})
		}
		if !active {
			return c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "forbidden", ErrorDescription: "aws account is disabled"})
		}

		region := resolveRegion(c, db)
		sess, err := assumeRoleSession(region, roleArn, externalID)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "assume_role_failed", ErrorDescription: err.Error()})
		}

		c.Set("session", sess)
		c.Set(ContextRegion, region)
		c.Set(ContextAWSAcctID, accountID)
		return next(c)
	}
}

func resolveRegion(c echo.Context, db *store.DB) string {
	if val := c.Request().Header.Get("X-AWS-REGION"); val != "" {
		return val
	}
	userID := c.Get(ContextUserID).(string)
	var region string
	_ = db.Conn.QueryRow(`SELECT region FROM user_settings WHERE user_id = $1`, userID).Scan(&region)
	if region != "" {
		return region
	}
	cfg := c.Get(ContextConfig).(*model.Config)
	return cfg.DefaultRegion
}

func assumeRoleSession(region, roleArn, externalID string) (*session.Session, error) {
	if roleArn == "" || externalID == "" {
		return nil, errors.New("role_arn or external_id missing")
	}
	baseSess := awsSession.Must(awsSession.NewSession(&aws.Config{Region: aws.String(region)}))
	stsClient := sts.New(baseSess)

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String("aws-admin-session"),
		ExternalId:      aws.String(externalID),
		DurationSeconds: aws.Int64(3600),
	}
	resp, err := stsClient.AssumeRole(input)
	if err != nil {
		return nil, err
	}

	creds := &session.StaticCredentials{
		AccessID:        aws.StringValue(resp.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(resp.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(resp.Credentials.SessionToken),
	}

	s := &session.Session{
		ID:      "",
		Region:  region,
		Credentials: session.Credentials{
			CredentialType:    session.StaticCred,
			StaticCredentials: *creds,
		},
	}
	ConfigureStaticCredSession(s, creds)
	return s, nil
}

func resolveUserEmail(c echo.Context) string {
	if val := c.Get(ContextUserEmail); val != nil {
		if email, ok := val.(string); ok {
			return email
		}
	}
	return ""
}

func userIsSuperAdmin(db *store.DB, userID, email string) (bool, error) {
	if email != "" && isEmailSuperAdmin(email) {
		return true, nil
	}
	var isSuper bool
	if err := db.Conn.QueryRow(`SELECT is_super_admin FROM users WHERE id = $1`, userID).Scan(&isSuper); err != nil {
		return false, err
	}
	return isSuper, nil
}

func isTenantAdminRole(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "owner", "admin":
		return true
	default:
		return false
	}
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
