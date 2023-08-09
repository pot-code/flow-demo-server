package auth

import "golang.org/x/crypto/bcrypt"

type PasswordHash interface {
	Hash(password string) (string, error)
	VerifyPassword(password, hash string) error
}

func NewPasswordHash() PasswordHash {
	return &bcryptHash{}
}

type bcryptHash struct{}

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
