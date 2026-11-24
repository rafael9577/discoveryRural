package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"

	"discoveryrural/internal/falha"
)

// LimiteCorpoJSON é o tamanho máximo aceito para corpos JSON (64 KiB).
const LimiteCorpoJSON = 1 << 16

// LerJSON decodifica o corpo da requisição em v, limitando o tamanho.
func LerJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, LimiteCorpoJSON)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return falha.Validacao("JSON inválido")
	}
	return nil
}

// InteiroPositivo converte s em inteiro > 0.
func InteiroPositivo(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}
