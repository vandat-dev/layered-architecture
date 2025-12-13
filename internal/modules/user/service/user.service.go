package service

import (
	"app/global"
	"app/internal/modules/user/constants"
	"app/internal/modules/user/dto"
	"app/internal/modules/user/model"
	"app/internal/modules/user/repo"
	"app/internal/third_party/redis"
	"app/pkg/jwt"
	"app/pkg/response"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetUserByID(id uuid.UUID) *response.ServiceResult
	GetListUser(req dto.UserListRequestDto, userRole string) *response.ServiceResult
	CreateUser(userDto dto.UserRequestDto) *response.ServiceResult
	UpdateUser(id uuid.UUID, updateDto dto.UserUpdateRequestDto, userRole string, userID uuid.UUID) *response.ServiceResult
	Login(username string, password string) *response.ServiceResult
	Register(registerDto dto.RegisterRequestDto) *response.ServiceResult
	ReceiveMessages(msg []byte) error
}

type userService struct {
	userRepo      repo.IUserRepository
	redisProvider *redis.RedisProvider
}

func NewUserService(userRepo repo.IUserRepository, redisProvider *redis.RedisProvider) IUserService {
	return &userService{
		userRepo:      userRepo,
		redisProvider: redisProvider,
	}
}

func (us *userService) getUserFromCache(id uuid.UUID) (*dto.UserResponseDto, bool) {
	ctx := context.Background()
	key := fmt.Sprintf("user:%s", id.String())

	data, err := us.redisProvider.Get(ctx, key)
	if err != nil || data == "" {
		return nil, false
	}

	var user dto.UserResponseDto
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, false
	}

	return &user, true
}

func (us *userService) saveUserToCache(id uuid.UUID, user *dto.UserResponseDto) {
	ctx := context.Background()
	key := fmt.Sprintf("user:%s", id.String())

	if data, err := json.Marshal(user); err == nil {
		_ = us.redisProvider.Set(ctx, key, data, 10*time.Second)
	}
}

func (us *userService) GetUserByID(id uuid.UUID) *response.ServiceResult {
	// 1. Get from cache
	if user, ok := us.getUserFromCache(id); ok {
		global.Logger.Info("Cache hit for user: " + id.String())
		return response.NewServiceResult(user)
	}

	// 2. Get from DB
	user, errResult := us.getUserFromDB(id)
	if errResult != nil {
		return errResult
	}

	// 3. Save to cache
	us.saveUserToCache(id, user)

	return response.NewServiceResult(user)
}

func (us *userService) getUserFromDB(id uuid.UUID) (*dto.UserResponseDto, *response.ServiceResult) {
	result := us.userRepo.GetUserByID(id)
	if result == nil {
		return nil, response.NewServiceErrorWithCode(422, response.ErrCodeUserNotFound)
	}

	user := &dto.UserResponseDto{
		Id:          result.ID,
		Email:       result.Email,
		Username:    result.Username,
		FullName:    result.FullName,
		PhoneNumber: result.PhoneNumber,
		Gender:      result.Gender,
		Address:     result.Address,
		SystemRole:  result.SystemRole,
		IsActive:    *result.IsActive,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}

	return user, nil
}

func (us *userService) GetListUser(req dto.UserListRequestDto, userRole string) *response.ServiceResult {
	// Check authorization - only ADMIN can get user list
	if userRole != constants.Admin {
		return response.NewServiceErrorWithCode(403, response.ErrCodeAccessDenied)
	}

	users, total, err := us.userRepo.GetListUser(req)
	if err != nil {
		global.Logger.Error("Failed to get users from repository: " + err.Error())
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	result := map[string]interface{}{
		"total": total,
		"data":  users,
	}
	return response.NewServiceResult(result)
}

func (us *userService) CreateUser(userDto dto.UserRequestDto) *response.ServiceResult {

	existingEmail := us.userRepo.GetUserByEmail(userDto.Email)
	if existingEmail != nil {
		return response.NewServiceErrorWithCode(409, response.ErrCodeUserHasExists)
	}

	existingUserName := us.userRepo.GetUserByUsername(userDto.Username)
	if existingUserName != nil {
		return response.NewServiceErrorWithCode(409, response.ErrCodeUserHasExists)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	active := true
	userID, err := uuid.NewV7()
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	user := &model.User{
		ID:          userID,
		Email:       userDto.Email,
		Username:    userDto.Username,
		FullName:    userDto.FullName,
		Password:    string(hashedPassword),
		PhoneNumber: userDto.PhoneNumber,
		Gender:      userDto.Gender,
		Address:     userDto.Address,
		SystemRole:  userDto.SystemRole,
		IsActive:    &active,
	}

	_, err = us.userRepo.CreateUser(user)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}
	return response.NewServiceResult(userID)
}

func (us *userService) UpdateUser(id uuid.UUID, updateDto dto.UserUpdateRequestDto, userRole string, userID uuid.UUID) *response.ServiceResult {

	existingUser := us.userRepo.GetUserByID(id)
	if existingUser == nil {
		return response.NewServiceErrorWithCode(404, response.ErrCodeUserNotFound)
	}

	if userRole != constants.Admin && userID != id {
		return response.NewServiceErrorWithCode(403, response.ErrCodeUserPermissionDenied)
	}

	updateUser := &model.User{}
	if updateDto.Email != "" {
		existingEmail := us.userRepo.GetUserByEmail(updateDto.Email)
		if existingEmail != nil {
			return response.NewServiceErrorWithCode(409, response.ErrCodeUserHasExists)
		}
		updateUser.Email = updateDto.Email
	}

	if updateDto.FullName != "" {
		updateUser.FullName = updateDto.FullName
	}
	if updateDto.PhoneNumber != "" {
		updateUser.PhoneNumber = updateDto.PhoneNumber
	}
	if updateDto.Gender != "" {
		updateUser.Gender = updateDto.Gender
	}
	if updateDto.Address != "" {
		updateUser.Address = updateDto.Address
	}
	if updateDto.SystemRole != "" {
		updateUser.SystemRole = updateDto.SystemRole
	}
	if updateDto.IsActive != nil {
		updateUser.IsActive = updateDto.IsActive
	}

	if updateDto.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateDto.Password), bcrypt.DefaultCost)
		if err != nil {
			return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
		}
		updateUser.Password = string(hashedPassword)
	}

	updatedUser, err := us.userRepo.UpdateUser(id, updateUser)
	if err != nil {
		return response.NewServiceErrorWithCode(400, response.ErrCodeUserHasExists)
	}

	userResponse := dto.UserResponseDto{
		Id:          updatedUser.ID,
		Email:       updatedUser.Email,
		Username:    updatedUser.Username,
		FullName:    updatedUser.FullName,
		PhoneNumber: updatedUser.PhoneNumber,
		Gender:      updatedUser.Gender,
		Address:     updatedUser.Address,
		SystemRole:  updatedUser.SystemRole,
		IsActive:    *updatedUser.IsActive,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
	}

	return response.NewServiceResult(&userResponse)
}

func (us *userService) Login(username string, password string) *response.ServiceResult {
	user := us.userRepo.GetUserByUsername(username)
	if user == nil {
		return response.NewServiceErrorWithCode(401, response.ErrCodeInvalidLogin)
	}

	// Compare password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return response.NewServiceErrorWithCode(401, response.ErrCodeInvalidLogin)
	}
	if user.IsActive != nil && !*user.IsActive {
		return response.NewServiceErrorWithCode(403, response.ErrCodeAccountLock)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.SystemRole, global.Config.JWT.SecretKey, global.Config.JWT.TokenExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Generate refresh token with longer expiry
	refreshToken, err := jwt.GenerateToken(user.ID, user.Email, user.SystemRole, global.Config.JWT.SecretKey, global.Config.JWT.RefreshExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	authResponse := &dto.AuthResponseDto{
		Token:        token,
		RefreshToken: refreshToken,
	}

	return response.NewServiceResult(authResponse)
}

func (us *userService) Register(registerDto dto.RegisterRequestDto) *response.ServiceResult {
	userDto := dto.UserRequestDto{
		Email:      registerDto.Email,
		Username:   registerDto.Username,
		Password:   registerDto.Password,
		SystemRole: registerDto.Role,
	}
	createResult := us.CreateUser(userDto)
	if createResult.Error != nil {
		return createResult // Return CreateUser
	}

	userID := createResult.Data.(uuid.UUID)
	user := us.userRepo.GetUserByID(userID)
	if user == nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.SystemRole, global.Config.JWT.SecretKey, global.Config.JWT.TokenExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	// Generate refresh token with longer expiry
	refreshToken, err := jwt.GenerateToken(user.ID, user.Email, user.SystemRole, global.Config.JWT.SecretKey, global.Config.JWT.RefreshExpiry)
	if err != nil {
		return response.NewServiceErrorWithCode(500, response.ErrCodeInternalError)
	}

	authResponse := &dto.AuthResponseDto{
		Token:        token,
		RefreshToken: refreshToken,
	}

	return response.NewServiceResult(authResponse)
}

func (us *userService) ReceiveMessages(msg []byte) error {
	global.Logger.Info("[USER-SERVICE] User Service received message: " + string(msg))
	// Add business logic here
	return nil
}
