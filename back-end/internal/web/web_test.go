package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRaizServePaginaEOutrasRotasDao404JSON(t *testing.T) {
	mux := http.NewServeMux()
	NovoHandler().Registrar(mux)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Fauna Local") {
		t.Fatalf("raiz deveria servir a página: %d", w.Code)
	}

	w = httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/nao-existe", nil))
	if w.Code != http.StatusNotFound || !strings.Contains(w.Body.String(), `"erro"`) {
		t.Fatalf("rota inexistente: %d %s", w.Code, w.Body.String())
	}
}
