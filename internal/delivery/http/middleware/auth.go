package middleware

import (
	"cat-test/internal/domain"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

const (
	AuthIdentityCtxKey      = "identity"
	authTokenPattern        = "^Bearer\\s[A-z0-9-_]{1,256}$"
	authUnauthorizedMessage = "your request was unauthorized"
)

func Authenticate(usecase domain.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := parseTokenFromHeader(c.GetHeader("Authorization"))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": authUnauthorizedMessage})
			return
		}

		user, err := usecase.Authenticate(c, token)
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": authUnauthorizedMessage})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Set(AuthIdentityCtxKey, user)
	}
}

func parseTokenFromHeader(header string) string {
	if header == "" || !regexp.MustCompile(authTokenPattern).MatchString(header) {
		return ""
	}

	return strings.Split(header, "Bearer ")[1]
}
