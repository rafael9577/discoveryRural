package avistamento

import (
	"net/http"

	"discoveryrural/internal/falha"
	"discoveryrural/internal/httpx"
	"discoveryrural/internal/limitador"
)

// Handler expõe os endpoints HTTP de avistamentos.
type Handler struct {
	servico   *Servico
	limitador *limitador.Limitador
	ip        httpx.ExtratorIP
}

func NovoHandler(servico *Servico, lim *limitador.Limitador, ip httpx.ExtratorIP) *Handler {
	return &Handler{servico: servico, limitador: lim, ip: ip}
}

func (h *Handler) Registrar(mux *http.ServeMux) {
	mux.HandleFunc("/api/avistamento", httpx.PermitirMetodos(h.avistamento, http.MethodGet, http.MethodPost))
	mux.HandleFunc("/api/avistamentos", httpx.PermitirMetodos(h.listar, http.MethodGet))
}

func (h *Handler) avistamento(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.registrar(w, r)
		return
	}
	h.buscar(w, r)
}

func (h *Handler) buscar(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.InteiroPositivo(r.URL.Query().Get("id"))
	if !ok {
		httpx.EscreverFalha(w, falha.Validacao("id inválido ou ausente"))
		return
	}
	a, err := h.servico.Buscar(r.Context(), id)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, a)
}

func (h *Handler) registrar(w http.ResponseWriter, r *http.Request) {
	if !h.limitador.Permitir(h.ip.De(r)) {
		httpx.EscreverFalha(w, falha.MuitasRequisicoes("Muitos avistamentos seguidos, aguarde alguns segundos"))
		return
	}
	var in Entrada
	if err := httpx.LerJSON(w, r, &in); err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	criado, err := h.servico.Registrar(r.Context(), in)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusCreated, criado)
}

func (h *Handler) listar(w http.ResponseWriter, r *http.Request) {
	animalID := 0
	if v := r.URL.Query().Get("animal_id"); v != "" {
		n, ok := httpx.InteiroPositivo(v)
		if !ok {
			httpx.EscreverFalha(w, falha.Validacao("animal_id inválido"))
			return
		}
		animalID = n
	}
	grupos, err := h.servico.Agrupar(r.Context(), animalID)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, grupos)
}
