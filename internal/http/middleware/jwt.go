package middleware

import (
	"net/http"
	"strings"
	"time"

	"driving-authority-backend/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const authUserKey = "auth_user"

type AuthUser struct {
	ID    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Role  domain.Role        `json:"role"`
}

func GetAuthUser(c *gin.Context) AuthUser {
	v, _ := c.Get(authUserKey)
	return v.(AuthUser)
}

type JWT struct {
	secret    []byte
	issuer    string
	accessTTL time.Duration
}

func NewJWT(secret, issuer string, accessTTLMinutes int) *JWT {
	return &JWT{
		secret:    []byte(secret),
		issuer:    issuer,
		accessTTL: time.Duration(accessTTLMinutes) * time.Minute,
	}
}

type Claims struct {
	UserID string      `json:"uid"`
	Email  string      `json:"email"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWT) SignAccessToken(userID primitive.ObjectID, email string, role domain.Role) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID.Hex(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTTL)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.secret)
}

func (j *JWT) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		claims := &Claims{}
		tok, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return j.secret, nil
		})
		if err != nil || tok == nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		oid, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token user id"})
			return
		}

		c.Set(authUserKey, AuthUser{
			ID:    oid,
			Email: claims.Email,
			Role:  claims.Role,
		})
		c.Next()
	}
}
