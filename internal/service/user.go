package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepo
	jwtSecret string
	jwtExpire time.Duration
}

func NewUserService(repo *repository.UserRepo, jwtSecret string, jwtExpire time.Duration) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret, jwtExpire: jwtExpire}
}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (s *UserService) Register(req *RegisterReq) (*model.User, error) {
	if s.repo.ExistsByUsername(req.Username) {
		return nil, errors.New("用户名已存在")
	}
	if s.repo.ExistsByEmail(req.Email) {
		return nil, errors.New("邮箱已被注册")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         "user",
		Enabled:      true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("注册失败: %w", err)
	}

	return user, nil
}

func (s *UserService) Login(req *LoginReq) (*model.User, string, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	if !user.Enabled {
		return nil, "", errors.New("账号已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("生成令牌失败: %w", err)
	}

	return user, token, nil
}

func (s *UserService) GetByID(id int64) (*model.User, error) {
	return s.repo.FindByID(id)
}

// ---- Admin user management ----

type AdminListUsersReq struct {
	Keyword  string `form:"keyword"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

type AdminUpdateUserReq struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6,max=100"`
	Role     *string `json:"role" binding:"omitempty,oneof=user admin"`
	Enabled  *bool   `json:"enabled"`
}

type AdminCreateUserReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
}

func (s *UserService) AdminListUsers(req *AdminListUsersReq) ([]model.User, int64, error) {
	return s.repo.FindAll(req.Keyword, req.Page, req.PageSize)
}

func (s *UserService) AdminCreateUser(req *AdminCreateUserReq) (*model.User, error) {
	if s.repo.ExistsByUsername(req.Username) {
		return nil, errors.New("用户名已存在")
	}
	if s.repo.ExistsByEmail(req.Email) {
		return nil, errors.New("邮箱已被注册")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         role,
		Enabled:      true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}
	return user, nil
}

func (s *UserService) AdminUpdateUser(id int64, req *AdminUpdateUserReq) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.Username != nil && *req.Username != user.Username {
		if s.repo.ExistsByUsernameExcludeID(*req.Username, id) {
			return nil, errors.New("用户名已存在")
		}
		user.Username = *req.Username
	}

	if req.Email != nil && *req.Email != user.Email {
		if s.repo.ExistsByEmailExcludeID(*req.Email, id) {
			return nil, errors.New("邮箱已被注册")
		}
		user.Email = *req.Email
	}

	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("密码加密失败: %w", err)
		}
		user.PasswordHash = string(hashed)
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if req.Enabled != nil {
		user.Enabled = *req.Enabled
	}

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}
	return user, nil
}

func (s *UserService) AdminDeleteUser(id int64) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("用户不存在")
	}
	return s.repo.Delete(id)
}

func (s *UserService) generateToken(user *model.User) (string, error) {
	claims := &UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
