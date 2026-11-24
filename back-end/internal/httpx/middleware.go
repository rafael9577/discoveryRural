package httpx

import (
	"log"
	"net/http"
	"slices"
	"time"
)

// Registravel é implementado por tudo que sabe registrar suas rotas num mux.
type Registravel interface {
	Registrar(mux *http.ServeMux)
}

// PermitirMetodos responde 405 para métodos fora da lista.
func PermitirMetodos(next http.HandlerFunc, metodos ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(metodos, r.Method) {
			EscreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}
		next(w, r)
	}
}

type gravadorStatus struct {
	http.ResponseWriter
	status int
}

func (g *gravadorStatus) WriteHeader(status int) {
	g.status = status
	g.ResponseWriter.WriteHeader(status)
}

// Observar registra cada requisição em log e converte panics em resposta 500.
func Observar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		g := &gravadorStatus{ResponseWriter: w, status: http.StatusOK}
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic em %s %s: %v", r.Method, r.URL.Path, rec)
				EscreverErro(g, http.StatusInternalServerError, "Erro interno")
			}
			log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, g.status, time.Since(inicio).Round(time.Millisecond))
		}()
		next.ServeHTTP(g, r)
	})
}
