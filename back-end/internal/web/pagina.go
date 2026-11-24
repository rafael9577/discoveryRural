// Package web serve a página administrativa embutida e o tratamento de rotas
// inexistentes.
package web

import (
	_ "embed"
	"net/http"

	"discoveryrural/internal/httpx"
)

//go:embed pagina.html
var paginaHTML []byte

type Handler struct{}

func NovoHandler() *Handler { return &Handler{} }

func (h *Handler) Registrar(mux *http.ServeMux) {
	mux.HandleFunc("/{$}", h.pagina)
	mux.HandleFunc("/", h.naoEncontrado)
}

func (h *Handler) pagina(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(paginaHTML)
}

func (h *Handler) naoEncontrado(w http.ResponseWriter, r *http.Request) {
	httpx.EscreverErro(w, http.StatusNotFound, "Rota não encontrada")
}
