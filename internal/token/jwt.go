package token

import "github.com/golang-jwt/jwt/v5"

type JwtIssuer struct {
	key string
}

func NewJwtIssuer(key string) *JwtIssuer {
	return &JwtIssuer{key: key}
}

func (j *JwtIssuer) Sign(claims jwt.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(j.key))
}

func (j *JwtIssuer) Verify(token string) (jwt.Claims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.key), nil
	})
	if err != nil {
		return nil, err
	}
	return t.Claims, nil
}
