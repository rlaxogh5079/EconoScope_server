// pkg/response/response.go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rlaxogh5079/EconoScope/pkg/errs"
	"github.com/rlaxogh5079/EconoScope/pkg/logger"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, status int, data interface{}) {
	if status == 0 {
		status = http.StatusOK
	}

	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, appErr *errs.AppError) {
	if appErr == nil {
		appErr = errs.ErrInternal
	}

	status := appErr.HTTPStatus
	if status == 0 {
		status = http.StatusInternalServerError
	}

	// 로그 기록
	logger.Log.WithFields(map[string]interface{}{
		"code":    appErr.Code,
		"message": appErr.Message,
		"status":  status,
		"error":   appErr.Err,
	}).Error("request error")

	c.JSON(status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}

func ErrorFromStd(c *gin.Context, err error) {
	appErr := errs.FromError(err)
	Error(c, appErr)
}
