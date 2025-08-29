package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
	t.Run("Hash and verify password", func(t *testing.T) {
		password := "TestPassword123!"
		
		// Hash password
		hash, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.True(t, strings.HasPrefix(hash, "$argon2id$"))
		
		// Verify correct password
		valid, err := VerifyPassword(password, hash)
		assert.NoError(t, err)
		assert.True(t, valid)
		
		// Verify incorrect password
		valid, err = VerifyPassword("WrongPassword123!", hash)
		assert.NoError(t, err)
		assert.False(t, valid)
	})
	
	t.Run("Different hashes for same password", func(t *testing.T) {
		password := "TestPassword123!"
		
		hash1, err := HashPassword(password)
		require.NoError(t, err)
		
		hash2, err := HashPassword(password)
		require.NoError(t, err)
		
		// Hashes should be different due to random salt
		assert.NotEqual(t, hash1, hash2)
		
		// Both should verify correctly
		valid1, err := VerifyPassword(password, hash1)
		assert.NoError(t, err)
		assert.True(t, valid1)
		
		valid2, err := VerifyPassword(password, hash2)
		assert.NoError(t, err)
		assert.True(t, valid2)
	})
	
	t.Run("Invalid hash format", func(t *testing.T) {
		password := "TestPassword123!"
		
		// Test various invalid formats
		invalidHashes := []string{
			"",
			"invalid",
			"$bcrypt$2a$10$invalid",
			"$argon2id$invalid",
			"$argon2id$v=19$invalid",
		}
		
		for _, hash := range invalidHashes {
			valid, err := VerifyPassword(password, hash)
			assert.Error(t, err)
			assert.False(t, valid)
		}
	})
}

func TestPasswordStrength(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Valid strong password",
			password:    "StrongPass123!",
			shouldError: false,
		},
		{
			name:        "Too short",
			password:    "Short1!",
			shouldError: true,
			errorMsg:    "at least 8 characters",
		},
		{
			name:        "No uppercase",
			password:    "weakpass123!",
			shouldError: true,
			errorMsg:    "uppercase letter",
		},
		{
			name:        "No lowercase",
			password:    "WEAKPASS123!",
			shouldError: true,
			errorMsg:    "lowercase letter",
		},
		{
			name:        "No number",
			password:    "WeakPassword!",
			shouldError: true,
			errorMsg:    "number",
		},
		{
			name:        "No special character",
			password:    "WeakPassword123",
			shouldError: true,
			errorMsg:    "special character",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if tt.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestJWTTokens(t *testing.T) {
	config := &JWTConfig{
		Secret:        []byte("test-secret-key-for-testing-only"),
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
		Issuer:        "ferryflow-test",
	}
	
	t.Run("Generate and validate token pair", func(t *testing.T) {
		userID := "user-123"
		email := "test@example.com"
		userType := "customer"
		operatorID := ""
		sessionID := "session-456"
		
		// Generate token pair
		pair, err := GenerateTokenPair(config, userID, email, userType, operatorID, sessionID)
		assert.NoError(t, err)
		assert.NotEmpty(t, pair.AccessToken)
		assert.NotEmpty(t, pair.RefreshToken)
		assert.NotEmpty(t, pair.AccessTokenHash)
		assert.NotEmpty(t, pair.RefreshTokenHash)
		assert.NotEqual(t, pair.AccessToken, pair.RefreshToken)
		
		// Validate access token
		claims, err := ValidateToken(config, pair.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, userType, claims.UserType)
		assert.Equal(t, sessionID, claims.SessionID)
		assert.Equal(t, "access", claims.TokenType)
		
		// Validate refresh token
		refreshClaims, err := ValidateToken(config, pair.RefreshToken)
		assert.NoError(t, err)
		assert.Equal(t, userID, refreshClaims.UserID)
		assert.Equal(t, "refresh", refreshClaims.TokenType)
	})
	
	t.Run("Token with operator ID", func(t *testing.T) {
		userID := "agent-123"
		email := "agent@ferry.com"
		userType := "agent"
		operatorID := "operator-789"
		sessionID := "session-abc"
		
		pair, err := GenerateTokenPair(config, userID, email, userType, operatorID, sessionID)
		assert.NoError(t, err)
		
		claims, err := ValidateToken(config, pair.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, operatorID, claims.OperatorID)
		assert.Equal(t, userType, claims.UserType)
	})
	
	t.Run("Refresh access token", func(t *testing.T) {
		// Generate initial token pair
		pair1, err := GenerateTokenPair(config, "user-1", "user@test.com", "customer", "", "session-1")
		require.NoError(t, err)
		
		// Refresh using refresh token
		pair2, err := RefreshAccessToken(config, pair1.RefreshToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, pair2.AccessToken)
		assert.NotEqual(t, pair1.AccessToken, pair2.AccessToken)
		
		// Validate new access token
		claims, err := ValidateToken(config, pair2.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, "user-1", claims.UserID)
		assert.Equal(t, "session-1", claims.SessionID)
	})
	
	t.Run("Cannot refresh with access token", func(t *testing.T) {
		pair, err := GenerateTokenPair(config, "user-1", "user@test.com", "customer", "", "session-1")
		require.NoError(t, err)
		
		// Try to refresh using access token instead of refresh token
		_, err = RefreshAccessToken(config, pair.AccessToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a refresh token")
	})
	
	t.Run("Invalid token validation", func(t *testing.T) {
		// Invalid token format
		_, err := ValidateToken(config, "invalid.token.format")
		assert.Error(t, err)
		
		// Token with wrong secret
		wrongConfig := &JWTConfig{
			Secret:        []byte("wrong-secret"),
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
			Issuer:        "ferryflow-test",
		}
		
		pair, err := GenerateTokenPair(config, "user-1", "user@test.com", "customer", "", "session-1")
		require.NoError(t, err)
		
		_, err = ValidateToken(wrongConfig, pair.AccessToken)
		assert.Error(t, err)
	})
	
	t.Run("Expired token", func(t *testing.T) {
		// Create config with very short expiry
		shortConfig := &JWTConfig{
			Secret:        []byte("test-secret"),
			AccessExpiry:  1 * time.Nanosecond,
			RefreshExpiry: 1 * time.Nanosecond,
			Issuer:        "ferryflow-test",
		}
		
		pair, err := GenerateTokenPair(shortConfig, "user-1", "user@test.com", "customer", "", "session-1")
		require.NoError(t, err)
		
		// Wait for token to expire
		time.Sleep(10 * time.Millisecond)
		
		_, err = ValidateToken(shortConfig, pair.AccessToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})
}

func TestTokenHashing(t *testing.T) {
	token1 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test1"
	token2 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test2"
	
	hash1 := HashToken(token1)
	hash2 := HashToken(token2)
	
	// Hashes should be deterministic
	assert.Equal(t, hash1, HashToken(token1))
	assert.Equal(t, hash2, HashToken(token2))
	
	// Different tokens should have different hashes
	assert.NotEqual(t, hash1, hash2)
	
	// Hash should be hex encoded SHA256 (64 characters)
	assert.Len(t, hash1, 64)
	assert.Len(t, hash2, 64)
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantError bool
	}{
		{
			name:      "Valid Bearer token",
			header:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			wantError: false,
		},
		{
			name:      "Missing Bearer prefix",
			header:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
			wantToken: "",
			wantError: true,
		},
		{
			name:      "Wrong prefix",
			header:    "Basic dXNlcjpwYXNz",
			wantToken: "",
			wantError: true,
		},
		{
			name:      "Empty header",
			header:    "",
			wantToken: "",
			wantError: true,
		},
		{
			name:      "Just Bearer",
			header:    "Bearer",
			wantToken: "",
			wantError: true,
		},
		{
			name:      "Bearer with space",
			header:    "Bearer ",
			wantToken: "",
			wantError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractTokenFromHeader(tt.header)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}