package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
	"yellowbook/internal/domain"
)

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

func setJWTToken(ctx *gin.Context, user domain.User) error {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.RegisteredClaims{
		Issuer:    "yellowbook",
		Subject:   strconv.FormatUint(user.Id, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return err
	}

	ctx.Header("X-Jwt-Token", tokenStr)
	return nil
}
