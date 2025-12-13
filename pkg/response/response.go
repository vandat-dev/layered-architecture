package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ServiceResult - Unified structure for service responses
type ServiceResult struct {
	Data       interface{}
	Error      error
	StatusCode int
	ErrorCode  int
}

// NewServiceResult - Create a successful ServiceResult
func NewServiceResult(data interface{}) *ServiceResult {
	return &ServiceResult{
		Data:       data,
		Error:      nil,
		StatusCode: http.StatusOK,
		ErrorCode:  ErrCodeSuccess,
	}
}

// NewServiceError - Create a ServiceResult with error
func NewServiceError(err error, statusCode int, errorCode int) *ServiceResult {
	return &ServiceResult{
		Data:       nil,
		Error:      err,
		StatusCode: statusCode,
		ErrorCode:  errorCode,
	}
}

// NewServiceErrorWithCode - Create a ServiceResult error using error code (message comes from msg map)
func NewServiceErrorWithCode(statusCode int, errorCode int) *ServiceResult {
	return &ServiceResult{
		Data:       nil,
		Error:      fmt.Errorf(GetMessage(errorCode)),
		StatusCode: statusCode,
		ErrorCode:  errorCode,
	}
}

// DataDetailResponse - Return response with custom code, message, and data
func DataDetailResponse(c *gin.Context, statusCode int, code int, data interface{}) {
	c.JSON(statusCode, Response{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}

// SuccessResponse - Return success response (HTTP 200)
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// ErrorResponse - Return error response with given status code
func ErrorResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// HandleServiceResult - Handle ServiceResult automatically
func HandleServiceResult(c *gin.Context, result *ServiceResult) {
	if result.Error != nil {
		if result.ErrorCode != 0 {
			// Use DataDetailResponse for errors with error code
			DataDetailResponse(c, result.StatusCode, result.ErrorCode, nil)
		} else {
			// Use ErrorResponse for normal errors
			ErrorResponse(c, result.StatusCode, result.Error.Error())
		}
		return
	}

	// Return data if success
	SuccessResponse(c, result.Data)
}
