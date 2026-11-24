package usuario

import (
	"net/http"

	"discoveryrural/internal/falha"
	"discoveryrural/internal/httpx"
	"discoveryrural/internal/limitador"
)

// Handler expõe os endpoints HTTP de usuários.
type Handler struct {
	servico   *Servico
	admin     *httpx.Admin
	limitador *limitador.Limitador // compartilhado entre cadastro e login
	ip        httpx.ExtratorIP
}

func NovoHandler(servico *Servico, admin *httpx.Admin, lim *limitador.Limitador, ip httpx.ExtratorIP) *Handler {
	return &Handler{servico: servico, admin: admin, limitador: lim, ip: ip}
}

func (h *Handler) Registrar(mux *http.ServeMux) {
	post := func(f http.HandlerFunc) http.HandlerFunc { return httpx.PermitirMetodos(f, http.MethodPost) }

	mux.HandleFunc("/api/users/register", post(h.cadastrar))
	mux.HandleFunc("/api/users/login", post(h.entrar))
	mux.HandleFunc("/api/users/me", httpx.PermitirMetodos(h.eu, http.MethodGet))
	mux.HandleFunc("/api/users/logout", post(h.sair))
	mux.HandleFunc("/api/users/reputacao", post(h.admin.Exigir(h.reputacao)))
}

// limitar devolve false (e já responde 429) quando o IP excedeu as tentativas.
func (h *Handler) limitar(w http.ResponseWriter, r *http.Request) bool {
	if h.limitador.Permitir(h.ip.De(r)) {
		return true
	}
	httpx.EscreverFalha(w, falha.MuitasRequisicoes("Muitas tentativas, aguarde alguns minutos"))
	return false
}

func (h *Handler) cadastrar(w http.ResponseWriter, r *http.Request) {
	if !h.limitar(w, r) {
		return
	}
	var in Cadastro
	if err := httpx.LerJSON(w, r, &in); err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	res, err := h.servico.Cadastrar(r.Context(), in)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusCreated, res)
}

func (h *Handler) entrar(w http.ResponseWriter, r *http.Request) {
	if !h.limitar(w, r) {
		return
	}
	var in Credenciais
	if err := httpx.LerJSON(w, r, &in); err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	res, err := h.servico.Entrar(r.Context(), in)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, res)
}

func (h *Handler) eu(w http.ResponseWriter, r *http.Request) {
	u, _, err := h.servico.Autenticar(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, u)
}

func (h *Handler) sair(w http.ResponseWriter, r *http.Request) {
	_, th, err := h.servico.Autenticar(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	if err := h.servico.Sair(r.Context(), th); err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) reputacao(w http.ResponseWriter, r *http.Request) {
	var in AjusteReputacao
	if err := httpx.LerJSON(w, r, &in); err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	nova, err := h.servico.AjustarReputacao(r.Context(), in)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, map[string]any{"id": in.ID, "reputacao": nova})
}
