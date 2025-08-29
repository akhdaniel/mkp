package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTConfig struct {
	Secret        []byte
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

type Claims struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	UserType    string `json:"user_type"`
	OperatorID  string `json:"operator_id,omitempty"`
	SessionID   string `json:"session_id"`
	TokenType   string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken       string
	RefreshToken      string
	AccessTokenHash   string
	RefreshTokenHash  string
	ExpiresAt        time.Time
}

// GenerateTokenPair creates both access and refresh tokens
func GenerateTokenPair(config *JWTConfig, userID, email, userType, operatorID, sessionID string) (*TokenPair, error) {
	now := time.Now()

	// Create access token claims
	accessClaims := Claims{
		UserID:     userID,
		Email:      email,
		UserType:   userType,
		OperatorID: operatorID,
		SessionID:  sessionID,
		TokenType:  "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(config.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(config.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create refresh token claims
	refreshClaims := Claims{
		UserID:     userID,
		Email:      email,
		UserType:   userType,
		OperatorID: operatorID,
		SessionID:  sessionID,
		TokenType:  "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Create refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(config.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:      accessTokenString,
		RefreshToken:     refreshTokenString,
		AccessTokenHash:  HashToken(accessTokenString),
		RefreshTokenHash: HashToken(refreshTokenString),
		ExpiresAt:       now.Add(config.RefreshExpiry),
	}, nil
}

// ValidateToken verifies and parses a JWT token
func ValidateToken(config *JWTConfig, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.Secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Additional validation
	if claims.Issuer != config.Issuer {
		return nil, fmt.Errorf("invalid token issuer")
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a valid refresh token
func RefreshAccessToken(config *JWTConfig, refreshTokenString string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := ValidateToken(config, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Ensure it's a refresh token
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	// Generate new token pair with same session ID
	return GenerateTokenPair(config, claims.UserID, claims.Email, claims.UserType, claims.OperatorID, claims.SessionID)
}

// HashToken creates a SHA256 hash of the token for storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ExtractTokenFromHeader extracts the token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return authHeader[7:], nil
}

// JWTUtil provides JWT operations
type JWTUtil struct {
	config *JWTConfig
}

// NewJWTUtil creates a new JWT utility
func NewJWTUtil(config *JWTConfig) *JWTUtil {
	return &JWTUtil{config: config}
}

// GenerateToken creates an access token
func (j *JWTUtil) GenerateToken(userID uuid.UUID, email, userType string) (string, error) {
	now := time.Now()
	sessionID := uuid.New().String()
	
	claims := Claims{
		UserID:    userID.String(),
		Email:     email,
		UserType:  userType,
		SessionID: sessionID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.config.Secret)
}

// GenerateRefreshToken creates a refresh token
func (j *JWTUtil) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	sessionID := uuid.New().String()
	
	claims := Claims{
		UserID:    userID.String(),
		SessionID: sessionID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.config.Secret)
}

// ValidateToken verifies and parses a JWT token
func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	return ValidateToken(j.config, tokenString)
}