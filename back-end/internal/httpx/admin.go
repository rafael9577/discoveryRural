package httpx

import (
	"crypto/subtle"
	"net/http"

	"discoveryrural/internal/falha"
)

// CabecalhoAdmin é o cabeçalho com a chave de administrador.
const CabecalhoAdmin = "X-Admin-Key"

// Admin valida a chave de administrador. Se nenhuma chave foi configurada,
// nenhuma requisição é autorizada.
type Admin struct {
	chave string
}

func NovoAdmin(chave string) *Admin { return &Admin{chave: chave} }

// Validar compara a chave candidata em tempo constante.
func (a *Admin) Validar(candidata string) error {
	if a.chave == "" || subtle.ConstantTimeCompare([]byte(candidata), []byte(a.chave)) != 1 {
		return falha.NaoAutorizado("Não autorizado")
	}
	return nil
}

// Exigir é um middleware que só deixa passar quem enviou X-Admin-Key válida.
func (a *Admin) Exigir(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a.Validar(r.Header.Get(CabecalhoAdmin)); err != nil {
			EscreverFalha(w, err)
			return
		}
		next(w, r)
	}
}
