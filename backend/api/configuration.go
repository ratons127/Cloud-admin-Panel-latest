package api

import (
	"net/http"
	"strings"

	"github.com/bertrandmartel/aws-admin/backend/middleware"
	"github.com/bertrandmartel/aws-admin/backend/model"
	"github.com/bertrandmartel/aws-admin/backend/store"
	"github.com/bertrandmartel/aws-admin/backend/utils"
	"github.com/labstack/echo/v4"
)

func SetConfiguration(c echo.Context) (err error) {
	db := c.Get(middleware.ContextDB).(*store.DB)
	userID := c.Get(middleware.ContextUserID).(string)
	r := new(model.SetConfigurationInput)
	if err = c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "region is required"))
	}
	region := strings.TrimSpace(r.Region)
	if region == "" {
		return c.JSON(http.StatusBadRequest, utils.SendError("format_error", "region is required"))
	}
	if _, err := db.Conn.Exec(`UPDATE user_settings SET region = $1 WHERE user_id = $2`, region, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SendError("internal_error", "failed to update region"))
	}
	return c.JSON(http.StatusOK, &model.SetConfigurationOutput{
		Updated: true,
	})
}
