package api

import (
	"net/http"
	"strings"

	"github.com/bertrandmartel/aws-admin/backend/middleware"
	"github.com/bertrandmartel/aws-admin/backend/store"
	"github.com/bertrandmartel/aws-admin/backend/utils"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type AWSAccountRequest struct {
	AccountID string `json:"accountId"`
	RoleArn   string `json:"roleArn"`
	ExternalID string `json:"externalId"`
	Name      string `json:"name"`
	Active    *bool  `json:"active"`
}

func ListAWSAccounts(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	tenantID := c.Get(middleware.ContextTenantID).(string)

	rows, err := db.Conn.Query(`SELECT id, account_id, role_arn, external_id, name, active FROM aws_accounts WHERE tenant_id = $1 ORDER BY name`, tenantID)
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

func CreateAWSAccount(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	tenantID := c.Get(middleware.ContextTenantID).(string)
	req := new(AWSAccountRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "invalid payload"))
	}
	if strings.TrimSpace(req.AccountID) == "" || strings.TrimSpace(req.RoleArn) == "" || strings.TrimSpace(req.ExternalID) == "" || strings.TrimSpace(req.Name) == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "accountId, roleArn, externalId, name are required"))
	}
	id := uuid.NewV4().String()
	if _, err := db.Conn.Exec(`INSERT INTO aws_accounts (id, tenant_id, account_id, role_arn, external_id, name) VALUES ($1,$2,$3,$4,$5,$6)`,
		id, tenantID, req.AccountID, req.RoleArn, req.ExternalID, req.Name); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to create account"))
	}
	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}

func UpdateAWSAccount(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	tenantID := c.Get(middleware.ContextTenantID).(string)
	accountID := c.Param("id")
	req := new(AWSAccountRequest)
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
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func DeleteAWSAccount(c echo.Context) error {
	db := c.Get(middleware.ContextDB).(*store.DB)
	tenantID := c.Get(middleware.ContextTenantID).(string)
	accountID := c.Param("id")
	if _, err := db.Conn.Exec(`DELETE FROM aws_accounts WHERE id = $1 AND tenant_id = $2`, accountID, tenantID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to delete account"))
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
