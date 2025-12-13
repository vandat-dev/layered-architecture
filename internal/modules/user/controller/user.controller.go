package controller

import (
	"app/global"
	"app/internal/modules/user/dto"
	"app/internal/modules/user/service"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequestDto true "User Registration Data"
// @Success 200 {object} response.Response{data=dto.AuthResponseDto} "Registration successful"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 422 {object} response.Response "User already exists"
// @Router /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var registerRequest dto.RegisterRequestDto
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := uc.userService.Register(registerRequest)
	response.HandleServiceResult(c, result)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequestDto true "User Login Data"
// @Success 200 {object} response.Response{data=dto.AuthResponseDto} "Login successful"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 401 {object} response.Response "Invalid credentials"
// @Router /user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var loginRequest dto.LoginRequestDto
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := uc.userService.Login(loginRequest.Username, loginRequest.Password)
	response.HandleServiceResult(c, result)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieves a user by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=dto.UserResponseDto} "User details"
// @Failure 422 {object} response.Response "Invalid user ID"
// @Router /user/get_user/{id} [get]
func (uc *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	result := uc.userService.GetUserByID(id)
	response.HandleServiceResult(c, result)
}

// GetListUser godoc
// @Summary Get all users (Admin only)
// @Description Returns a paginated list of users with filtering options. Only admin users can access this endpoint.
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(10)
// @Param email query string false "Email"
// @Param username query string false "Username"
// @Param system_role query string false "Status filter" Enums(ADMIN, USER, SUPER_ADMIN)
// @Success 200 {object} response.Response{data=dto.UserListResponseDto} "Paginated list of users"
// @Failure 400 {object} response.Response "Invalid query parameters"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Access denied: Only admin can view user list"
// @Router /user/list_user [get]
func (uc *UserController) GetListUser(c *gin.Context) {
	var req dto.UserListRequestDto

	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Error("Failed to bind query parameters: " + err.Error())
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidParams, nil)
		return
	}

	role, _ := c.Get("system_role")
	result := uc.userService.GetListUser(req, role.(string))
	response.HandleServiceResult(c, result)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user with the provided information
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body dto.CreateUserDto true "User Information"
// @Success 200 {object} response.Response{data=map[string]interface{}} "User created successfully"
// @Failure 400 {object} response.Response "Invalid request payload"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 405 {object} response.Response "User already exists"
// @Router /user/create_user [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	userRequest := dto.CreateUserDto{}

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	dataUser := dto.UserRequestDto{
		Email:      userRequest.Email,
		Username:   userRequest.Username,
		SystemRole: userRequest.SystemRole,
		Password:   global.Config.System.DefaultPassWord,
	}
	if userRequest.FullName != "" {
		dataUser.FullName = userRequest.FullName
	}

	result := uc.userService.CreateUser(dataUser)
	response.HandleServiceResult(c, result)
}

// UpdateUser godoc
// @Summary Update user by ID
// @Description Updates user information by user ID (email cannot be updated)
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param user body dto.UserUpdateRequestDto true "User Update Data (username, password, role only)"
// @Success 200 {object} response.Response{data=dto.UserResponseDto} "User updated successfully"
// @Failure 400 {object} response.Response "Invalid request data"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "User not found"
// @Router /user/update_user/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidParams, nil)
		return
	}

	var updateRequest dto.UserUpdateRequestDto
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		response.DataDetailResponse(c, 422, response.ErrCodeInvalidData, nil)
		return
	}

	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("system_role")
	result := uc.userService.UpdateUser(id, updateRequest, userRole.(string), userID.(uuid.UUID))
	response.HandleServiceResult(c, result)
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get the currently log in user's information
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=dto.UserResponseDto} "Current user"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /user/me [get]
func (uc *UserController) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	result := uc.userService.GetUserByID(userID.(uuid.UUID))
	response.HandleServiceResult(c, result)
}
