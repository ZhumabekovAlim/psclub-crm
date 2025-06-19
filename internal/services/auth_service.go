package services

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"

	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

// AuthService handles user authentication and token management.
type AuthService struct {
	userRepo      *repositories.UserRepository
	tokenRepo     *repositories.TokenRepository
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

type jwtClaims struct {
	UserID int   `json:"user_id"`
	Exp    int64 `json:"exp"`
}

func NewAuthService(userRepo *repositories.UserRepository, tokenRepo *repositories.TokenRepository,
	accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

// Register creates a new user and returns token pair.
func (s *AuthService) Register(ctx context.Context, u *models.User) (string, string, error) {
	existing, err := s.userRepo.GetByPhone(ctx, u.Phone)
	if err != nil {
		return "", "", err
	}
	if existing != nil {
		return "", "", errors.New("user already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return "", "", err
	}
	u.Password = string(hashed)
	id, err := s.userRepo.Create(ctx, u)
	if err != nil {
		return "", "", err
	}
	u.ID = id
	return s.generateTokenPair(ctx, id)
}

// Login verifies credentials and returns new tokens.
func (s *AuthService) Login(ctx context.Context, phone, password string) (string, string, string, []string, string, error) {
	u, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return "", "", "", nil, "", err
	}
	if u == nil {
		return "", "", "", nil, "", errors.New("invalid credentials1")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", "", "", nil, "", errors.New("invalid credentials2")
	}

	token1, token2, err := s.generateTokenPair(ctx, u.ID)

	return token1, token2, u.Role, u.Permissions, u.Name, err
}

// Refresh validates refresh token and returns a new pair.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	hash := s.hashToken(refreshToken)
	userID, exp, err := s.tokenRepo.Get(ctx, hash)
	if err != nil {
		return "", "", errors.New("invalid token")
	}
	if time.Now().After(exp) {
		_ = s.tokenRepo.Delete(ctx, hash)
		return "", "", errors.New("token expired")
	}
	// rotate token
	_ = s.tokenRepo.Delete(ctx, hash)
	return s.generateTokenPair(ctx, userID)
}

func (s *AuthService) generateTokenPair(ctx context.Context, userID int) (string, string, error) {
	access, err := generateJWT(userID, s.accessSecret, s.accessTTL)
	if err != nil {
		return "", "", err
	}
	refreshRaw, err := randomString(32)
	if err != nil {
		return "", "", err
	}
	hash := s.hashToken(refreshRaw)
	exp := time.Now().Add(s.refreshTTL)
	if err := s.tokenRepo.Save(ctx, userID, hash, exp); err != nil {
		return "", "", err
	}
	return access, refreshRaw, nil
}

func (s *AuthService) hashToken(t string) string {
	h := hmac.New(sha256.New, []byte(s.refreshSecret))
	h.Write([]byte(t))
	return hex.EncodeToString(h.Sum(nil))
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateJWT(userID int, secret string, ttl time.Duration) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(jwtClaims{UserID: userID, Exp: time.Now().Add(ttl).Unix()})
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	unsigned := header + "." + payload
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsigned))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return unsigned + "." + signature, nil
}
