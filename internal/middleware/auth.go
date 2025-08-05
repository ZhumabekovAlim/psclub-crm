package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type jwtClaims struct {
	UserID    int    `json:"user_id"`
	CompanyID int    `json:"company_id"`
	BranchID  int    `json:"branch_id"`
	Role      string `json:"role"`
	Exp       int64  `json:"exp"`
}

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		unsigned := parts[0] + "." + parts[1]
		h := hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(unsigned))
		expectedSig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
		if !hmac.Equal([]byte(expectedSig), []byte(parts[2])) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var claims jwtClaims
		if err := json.Unmarshal(payload, &claims); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if time.Now().Unix() > claims.Exp {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("company_id", claims.CompanyID)
		c.Set("branch_id", claims.BranchID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
