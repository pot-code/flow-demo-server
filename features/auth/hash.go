package auth

import "golang.org/x/crypto/bcrypt"

type bcryptHash struct{}

func NewBcryptHash() *bcryptHash {
	return &bcryptHash{}
}

func (*bcryptHash) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (*bcryptHash) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
