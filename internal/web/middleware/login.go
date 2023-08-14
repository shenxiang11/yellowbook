package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoinMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		authorization := ctx.GetHeader("Authorization")
		tokenStr := strings.TrimPrefix(authorization, "Bearer ")

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			pk := `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDOJhgJhy+FX2ao6hwwxYucRPdm
mCRXhuS+xbHwvblYpVOcz7Y8xFteEsB36cqhKeb7kioFafOT6gvBRFYcCOKr1gcL
0SUK5g/0g2lf6Rg6MvufybhAlt9SAGSBgyS/jZVAHmuiJpLOWC1JCOyYlhfX6JPY
GcQaucz8O9XyFZkWJwIDAQAB
-----END PUBLIC KEY-----`
			return jwt.ParseRSAPublicKeyFromPEM([]byte(pk))
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sid, err := token.Claims.GetSubject()
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("UserId", sid)
	}
}
