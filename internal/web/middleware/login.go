package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
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

const privateKey = `-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAM4mGAmHL4VfZqjq
HDDFi5xE92aYJFeG5L7FsfC9uVilU5zPtjzEW14SwHfpyqEp5vuSKgVp85PqC8FE
VhwI4qvWBwvRJQrmD/SDaV/pGDoy+5/JuECW31IAZIGDJL+NlUAea6Imks5YLUkI
7JiWF9fok9gZxBq5zPw71fIVmRYnAgMBAAECgYEAp1bW5k0dbyeM/wrjDVgeRyDY
ryhLP92ZK57xHZn0rZeusrkNlnBSNqAEKpLWUFLiVE5G3BQwjF5NYnolaCZyUFOE
kZ26aSVJ4CJCKIvEY32Vfxkis6ajxU7PnBorwLHaloNrXk/KIgSya80nmC+ibLRq
WEBVBP2rq1bwa5yjj1ECQQDto0Jo7JPopG6q5ingW1zmY3PYs5PZyupHtKwrm5Up
SKvjrMNB0sEvMUG7Wj/h2xotvxkwMqIfPCnNc3QNcodpAkEA3hPyveZ2Se7AjFSY
1QbqvBnXL/dxRM20q1QsKcwbjtPJyJVaXfNw4yYc6VaN5C1v3GBHlAHnnbbnsqVU
e9QPDwJAReUbB1luN6MFmeaQspisvmbKEBbhidGRDv4pFbpxKO9i/1g1JgsjHwpR
1xU4bOnQzVvDwNVjseQ0N2WZ4Mqq4QJBAMS2AMeLU24LqMzkxnez58r0TLL1SIS8
fXNhXLktTZ/HI66j9ObRk0XxZZyeiZL7WGFpex20TjhaYoPQhLQm06sCQADNQz5H
S/t+zcNA/uwSyGOP+zXIL2+WBC1tsKvuNyM5YX5yWU6hiGKFmd6LYgA9yWkyBtcl
4L3+mbK6rVrMOdA=
-----END PRIVATE KEY-----
`

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

		// 续期
		if claims.ExpiresAt.Sub(time.Now()) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 10))
			key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
			if err != nil {
				ctx.String(http.StatusInternalServerError, "系统错误")
				return
			}
			tokenStr, err := token.SignedString(key)
			if err != nil {
				log.Panicln("jwt 续约失败", err)
			}
			ctx.Header("X-Jwt-Token", tokenStr)
		}

		ctx.Set("UserId", sid)
	}
}
