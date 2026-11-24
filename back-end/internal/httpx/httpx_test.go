package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func ok(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }

func TestPermitirMetodos(t *testing.T) {
	h := PermitirMetodos(ok, http.MethodGet)

	w := httptest.NewRecorder()
	h(w, httptest.NewRequest(http.MethodPost, "/", nil))
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST: status %d, esperado 405", w.Code)
	}
	w = httptest.NewRecorder()
	h(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET: status %d, esperado 200", w.Code)
	}
}

func TestAdminExigir(t *testing.T) {
	casos := []struct {
		nome, chaveServidor, chaveEnviada string
		esperado                          int
	}{
		{"chave correta", "segredo", "segredo", http.StatusOK},
		{"chave errada", "segredo", "outra", http.StatusUnauthorized},
		{"sem cabeçalho", "segredo", "", http.StatusUnauthorized},
		{"servidor sem chave configurada", "", "", http.StatusUnauthorized},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			if c.chaveEnviada != "" {
				r.Header.Set(CabecalhoAdmin, c.chaveEnviada)
			}
			w := httptest.NewRecorder()
			NovoAdmin(c.chaveServidor).Exigir(ok)(w, r)
			if w.Code != c.esperado {
				t.Fatalf("status %d, esperado %d", w.Code, c.esperado)
			}
		})
	}
}

func TestExtratorIP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "10.0.0.9:5555"
	r.Header.Set("X-Forwarded-For", "203.0.113.7, 10.0.0.1")

	if got := (ExtratorIP{}).De(r); got != "10.0.0.9" {
		t.Errorf("sem proxy confiável: %q", got)
	}
	if got := (ExtratorIP{ConfiarProxy: true}).De(r); got != "203.0.113.7" {
		t.Errorf("com proxy confiável: %q", got)
	}
}
