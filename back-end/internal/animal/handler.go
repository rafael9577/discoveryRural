package animal

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"discoveryrural/internal/falha"
	"discoveryrural/internal/httpx"
)

const (
	limiteCorpoCadastro = 12 << 20 // 12 MiB (foto + campos)
	memoriaFormulario   = 10 << 20
)

// Handler expõe os endpoints HTTP de animais.
type Handler struct {
	servico *Servico
	admin   *httpx.Admin
}

func NovoHandler(servico *Servico, admin *httpx.Admin) *Handler {
	return &Handler{servico: servico, admin: admin}
}

func (h *Handler) Registrar(mux *http.ServeMux) {
	get := func(f http.HandlerFunc) http.HandlerFunc { return httpx.PermitirMetodos(f, http.MethodGet) }

	mux.HandleFunc("/api/animais", get(h.listar))
	mux.HandleFunc("/api/animais/ids", get(h.ids))
	mux.HandleFunc("/api/animais/resumo", get(h.resumo))
	mux.HandleFunc("/api/animal", get(h.detalhe))
	mux.HandleFunc("/api/animal/imagem", get(h.imagem))
	mux.HandleFunc("/api/animal/create", httpx.PermitirMetodos(h.criar, http.MethodPost))
	mux.HandleFunc("/api/animal/update",
		httpx.PermitirMetodos(h.admin.Exigir(h.atualizar), http.MethodPatch, http.MethodPut))
}

func (h *Handler) listar(w http.ResponseWriter, r *http.Request) {
	animais, err := h.servico.Listar(r.Context())
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, animais)
}

func (h *Handler) ids(w http.ResponseWriter, r *http.Request) {
	ids, err := h.servico.IDs(r.Context())
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, ids)
}

func (h *Handler) resumo(w http.ResponseWriter, r *http.Request) {
	lista, err := h.servico.Resumos(r.Context())
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, lista)
}

func (h *Handler) detalhe(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) imagem(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.InteiroPositivo(r.URL.Query().Get("id"))
	if !ok {
		httpx.EscreverFalha(w, falha.Validacao("id inválido ou ausente"))
		return
	}
	dados, err := h.servico.Imagem(r.Context(), id)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(dados)
}

func (h *Handler) atualizar(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.InteiroPositivo(r.URL.Query().Get("id"))
	if !ok {
		httpx.EscreverFalha(w, falha.Validacao("id inválido ou ausente"))
		return
	}
	var in EdicaoAnimal
	if err := httpx.LerJSON(w, r, &in); err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	a, err := h.servico.Atualizar(r.Context(), id, in)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}
	httpx.EscreverJSON(w, http.StatusOK, a)
}

// criar recebe multipart/form-data (ou urlencoded). A chave de administrador
// vai em X-Admin-Key ou no campo admin_key (usado pela página administrativa).
func (h *Handler) criar(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, limiteCorpoCadastro)
	if err := r.ParseMultipartForm(memoriaFormulario); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		var grande *http.MaxBytesError
		if errors.As(err, &grande) {
			httpx.EscreverErro(w, http.StatusRequestEntityTooLarge, "Requisição grande demais")
			return
		}
		httpx.EscreverFalha(w, falha.Validacao("Formulário inválido"))
		return
	}

	chave := r.Header.Get(httpx.CabecalhoAdmin)
	if chave == "" {
		chave = r.FormValue("admin_key")
	}
	if err := h.admin.Validar(chave); err != nil {
		httpx.EscreverFalha(w, err)
		return
	}

	entrada, err := lerFormulario(r)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}

	var foto []byte
	if arquivo, _, err := r.FormFile("imagem"); err == nil {
		defer arquivo.Close()
		foto, _ = io.ReadAll(arquivo)
	}

	id, err := h.servico.Criar(r.Context(), entrada, foto)
	if err != nil {
		httpx.EscreverFalha(w, err)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		httpx.EscreverJSON(w, http.StatusCreated, map[string]int64{"id": id})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func lerFormulario(r *http.Request) (NovoAnimal, error) {
	raridade, err := strconv.Atoi(strings.TrimSpace(r.FormValue("raridade")))
	if err != nil {
		return NovoAnimal{}, falha.Validacao("raridade deve ser de 1 a 5")
	}
	bitmap := 0
	if v := strings.TrimSpace(r.FormValue("bitmap_id")); v != "" {
		if bitmap, err = strconv.Atoi(v); err != nil {
			return NovoAnimal{}, falha.Validacao("bitmap_id inválido")
		}
	}
	return NovoAnimal{
		NomePopular:        r.FormValue("nome_popular"),
		NomeCientifico:     r.FormValue("nome_cientifico"),
		BreveDescricao:     r.FormValue("breve_descricao"),
		EstadoConservacao:  r.FormValue("estado_conservacao"),
		Habito:             r.FormValue("habito"),
		VocalizacaoPostura: r.FormValue("vocalizacao_postura"),
		Peconhento:         marcado(r.FormValue("peconhento")),
		AvisoSeguranca:     r.FormValue("aviso_seguranca"),
		ProtocoloIncidente: r.FormValue("protocolo_incidente"),
		ContatosResgate:    r.FormValue("contatos_resgate"),
		Raridade:           raridade,
		BitmapID:           bitmap,
	}, nil
}

// marcado aceita o "on" do checkbox HTML e também "true"/"1".
func marcado(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "on", "true", "1":
		return true
	}
	return false
}
