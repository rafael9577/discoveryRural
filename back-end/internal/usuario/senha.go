package usuario

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HasherSenha abstrai o algoritmo de hash de senhas.
type HasherSenha interface {
	Hash(senha string) (string, error)
	Comparar(hash, senha string) error
}

// HasherBcrypt implementa HasherSenha com bcrypt.
type HasherBcrypt struct {
	Custo int
}

func NovoHasherBcrypt() HasherBcrypt { return HasherBcrypt{Custo: bcrypt.DefaultCost} }

func (h HasherBcrypt) Hash(senha string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(senha), h.Custo)
	return string(b), err
}

func (h HasherBcrypt) Comparar(hash, senha string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha))
}

// novoToken gera um token opaco aleatório (256 bits, em hexadecimal).
func novoToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashToken é o que é gravado no banco: o token em si nunca é armazenado.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
