package response

const (
	ErrCodeSuccess              = 2001  //Success
	ErrCodeInvalidParams        = 2002  //Email invalid
	ErrInvalidToken             = 3001  //Token invalid
	ErrCodeUserHasExists        = 50001 // User already exist
	ErrCodeUserNotFound         = 4000  // User not found
	ErrCodeInvalidLogin         = 4001  // Invalid login credentials
	ErrCodeAccessDenied         = 4003  // Access denied
	ErrCodeAccountLock          = 4004  // Your account has been locked
	ErrCodeUserPermissionDenied = 4005  // You do not have permission to interact with this user
	ErrCodeInternalError        = 5000  // Internal server error
	ErrCodeInvalidData          = 4221  // Invalid request data
	ErrCodeUnauthorized         = 4010  // Unauthorized
	ErrCodeDataNotFound         = 4009  // Data not found
)

var (
	msg = map[int]string{
		//	common
		ErrCodeSuccess:       "SUCCESS",
		ErrInvalidToken:      "TOKEN_INVALID",
		ErrCodeInvalidLogin:  "LOGIN_FAILED",
		ErrCodeAccessDenied:  "ACCESS_DENIED",
		ErrCodeInternalError: "INTERNAL_SERVER_ERROR",
		ErrCodeInvalidData:   "INVALID_DATA",
		ErrCodeUnauthorized:  "UNAUTHORIZED",

		//	user
		ErrCodeInvalidParams:        "EMAIL_INVALID",
		ErrCodeUserHasExists:        "USER_ALREADY_EXISTS",
		ErrCodeUserNotFound:         "USER_NOT_FOUND",
		ErrCodeAccountLock:          "USER_ACCOUNT_LOCKED",
		ErrCodeUserPermissionDenied: "YOU_DO_NOT_HAVE_PERMISSION_TO_INTERACT_WITH_THIS_USER",
		ErrCodeDataNotFound:         "DATA_NOT_FOUND",
	}
)

// GetMessage - Get message from error code
func GetMessage(errorCode int) string {
	if message, exists := msg[errorCode]; exists {
		return message
	}
	return "Unknown error"
}
