package auth

import "golang.org/x/crypto/bcrypt"

type PasswordHash interface {
	Hash(password string) (string, error)
	VerifyPassword(password, hash string) error
}

type passwordHash struct{}

func NewPasswordHash() PasswordHash {
	return &passwordHash{}
}

func (*passwordHash) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (*passwordHash) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
