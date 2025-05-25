package auth

import (
	"os"
	"time"

	"github.com/DrD-media/common/errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req *RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	user := &User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     "user",
	}

	return s.repo.Create(user)
}

func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate JWT")
	}

	return &LoginResponse{Token: token}, nil
}

func (s *UserService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		secret := []byte(os.Getenv("JWT_SECRET"))
		return secret, nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

func (s *UserService) generateJWT(user *User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign JWT")
	}

	return tokenString, nil
}

func (s *UserService) GetUserByID(id int) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by ID")
	}
	return user, nil
}
