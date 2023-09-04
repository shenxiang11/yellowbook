package jwt_generator

import (
	"testing"
	"time"
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

func TestJWTGenerator_Generate(t *testing.T) {
	j := NewJWTGenerator("test", privateKey)

	j.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}

	token, err := j.Generate("1", time.Minute)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0Iiwic3ViIjoiMSIsImV4cCI6MTUxNjIzOTA4MiwiaWF0IjoxNTE2MjM5MDIyfQ.AnakHDzhdXJTbnMefAd3uR_l7lhz4Fzynsj4_PM4bXy7BOknczQDsYOBO3ZTeuNoEnwJZLdPjBIpT8bl381OAOhYmr1ywU-sRIpCK-Qq1WrQphHQnX-FRL9UXuhCeXCZTMgDb8mW_PYzNmCjMeKrV_y_SPE0v5bvmUKY2cFqo6c"

	if token != want {
		t.Fatalf("wrong token generated. want: %q, got: %q", want, token)
	}
}
