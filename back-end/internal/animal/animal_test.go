package animal

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discoveryrural/internal/avistamento"
	"discoveryrural/internal/falha"
	"discoveryrural/internal/httpx"
)

func ptr[T any](v T) *T { return &v }

// --- dublês ------------------------------------------------------------

type repoFalso struct {
	animais    map[int]Animal
	atualizado []Alteracao
	criados    []Animal
}

func (r *repoFalso) Listar(context.Context) ([]Animal, error) {
	var l []Animal
	for _, a := range r.animais {
		l = append(l, a)
	}
	return l, nil
}
func (r *repoFalso) BuscarPorID(_ context.Context, id int) (*Animal, error) {
	a, ok := r.animais[id]
	if !ok {
		return nil, falha.NaoEncontrado("Animal não encontrado")
	}
	return &a, nil
}
func (r *repoFalso) Resumos(context.Context) ([]Resumo, error) { return nil, nil }
func (r *repoFalso) IDs(context.Context) ([]int, error)        { return nil, nil }
func (r *repoFalso) Existe(_ context.Context, id int) (bool, error) {
	_, ok := r.animais[id]
	return ok, nil
}
func (r *repoFalso) Contar(context.Context) (int, error) { return len(r.animais), nil }
func (r *repoFalso) Criar(_ context.Context, a Animal, _ []byte) (int64, error) {
	r.criados = append(r.criados, a)
	return int64(len(r.criados)), nil
}
func (r *repoFalso) Atualizar(_ context.Context, _ int, alt []Alteracao) error {
	r.atualizado = alt
	return nil
}
func (r *repoFalso) Imagem(context.Context, int) ([]byte, error) { return nil, nil }

type pontosFalsos struct{}

func (pontosFalsos) PontosPorAnimal(context.Context, int, int) (map[int][]avistamento.Ponto, error) {
	return map[int][]avistamento.Ponto{1: {{Longitude: -35, Latitude: -8}}}, nil
}

func servicoDeTeste() (*Servico, *repoFalso) {
	repo := &repoFalso{animais: map[int]Animal{1: {ID: 1, NomePopular: "Coruja"}}}
	return NovoServico(repo, pontosFalsos{}, NovoProcessadorJPEG()), repo
}

// --- validação ---------------------------------------------------------

func TestNormalizarAplicaTrimELimites(t *testing.T) {
	e, err := EdicaoAnimal{NomePopular: ptr("  Sapo  ")}.Normalizar()
	if err != nil || *e.NomePopular != "Sapo" {
		t.Fatalf("trim falhou: %v %v", e.NomePopular, err)
	}

	casos := map[string]EdicaoAnimal{
		"nome_popular não pode ficar vazio": {NomePopular: ptr("   ")},
		"habito passa de 100 caracteres":    {Habito: ptr(strings.Repeat("x", 101))},
		"raridade deve ser de 1 a 5":        {Raridade: ptr(6)},
		"bitmap_id não pode ser negativo":   {BitmapID: ptr(-1)},
	}
	for msg, in := range casos {
		if _, err := in.Normalizar(); err == nil || err.(*falha.Falha).Mensagem != msg {
			t.Errorf("esperava %q, veio %v", msg, err)
		}
	}
}

func TestAlteracoesSoIncluemCamposEnviados(t *testing.T) {
	alt := EdicaoAnimal{Habito: ptr("Noturno"), Peconhento: ptr(true), Raridade: ptr(4)}.Alteracoes()
	if len(alt) != 3 || alt[0].Coluna != "habito" || alt[1].Coluna != "peconhento" || alt[1].Valor != 1 {
		t.Fatalf("alterações inesperadas: %+v", alt)
	}
}

// --- serviço -----------------------------------------------------------

func TestAtualizarSemCamposEhValidacao(t *testing.T) {
	s, _ := servicoDeTeste()
	_, err := s.Atualizar(context.Background(), 1, EdicaoAnimal{})
	if !falha.EhTipo(err, falha.TipoValidacao) {
		t.Fatalf("esperava 400, veio %v", err)
	}
}

func TestAtualizarAnimalInexistente(t *testing.T) {
	s, _ := servicoDeTeste()
	_, err := s.Atualizar(context.Background(), 99, EdicaoAnimal{Habito: ptr("x")})
	if !falha.EhTipo(err, falha.TipoNaoEncontrado) {
		t.Fatalf("esperava 404, veio %v", err)
	}
}

func TestAtualizarDevolveAnimalComAvistamentos(t *testing.T) {
	s, repo := servicoDeTeste()
	res, err := s.Atualizar(context.Background(), 1, EdicaoAnimal{Habito: ptr("Noturno")})
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.atualizado) != 1 || len(res.Avistamentos) != 1 || res.ImagemURL != "/api/animal/imagem?id=1" {
		t.Fatalf("resultado inesperado: %+v / %+v", res, repo.atualizado)
	}
}

func TestCriarValidaEUsaMesmasRegrasDaEdicao(t *testing.T) {
	s, repo := servicoDeTeste()
	_, err := s.Criar(context.Background(), NovoAnimal{NomePopular: "", Raridade: 3}, nil)
	if !falha.EhTipo(err, falha.TipoValidacao) {
		t.Fatalf("nome vazio deveria ser 400, veio %v", err)
	}
	if _, err := s.Criar(context.Background(), NovoAnimal{NomePopular: " Tatu ", Raridade: 3}, nil); err != nil {
		t.Fatal(err)
	}
	if repo.criados[0].NomePopular != "Tatu" {
		t.Errorf("nome deveria ser aparado: %q", repo.criados[0].NomePopular)
	}
}

func TestSemearPadraoSoQuandoVazio(t *testing.T) {
	s, repo := servicoDeTeste()
	s.SemearPadrao(context.Background())
	if len(repo.criados) != 0 {
		t.Fatal("não deveria semear com a tabela populada")
	}
	vazio := &repoFalso{animais: map[int]Animal{}}
	s = NovoServico(vazio, pontosFalsos{}, NovoProcessadorJPEG())
	s.SemearPadrao(context.Background())
	if len(vazio.criados) != 2 {
		t.Fatalf("deveria semear 2 animais, semeou %d", len(vazio.criados))
	}
}

// --- imagem ------------------------------------------------------------

func TestProcessadorDevolveJPEG450x600AteLimite(t *testing.T) {
	var png bytes.Buffer
	origem := image.NewRGBA(image.Rect(0, 0, 1600, 900))
	if err := jpeg.Encode(&png, origem, nil); err != nil {
		t.Fatal(err)
	}
	p := NovoProcessadorJPEG()
	saida := p.Processar(png.Bytes())
	cfg, formato, err := image.DecodeConfig(bytes.NewReader(saida))
	if err != nil || formato != "jpeg" || cfg.Width != 450 || cfg.Height != 600 {
		t.Fatalf("saída inesperada: %v %q %+v", err, formato, cfg)
	}
	if len(saida) > tamanhoMaxBytes {
		t.Errorf("imagem com %d bytes passa do limite", len(saida))
	}
}

func TestProcessadorUsaPlaceholderParaEntradaInvalida(t *testing.T) {
	p := NovoProcessadorJPEG()
	if !bytes.Equal(p.Processar([]byte("não é imagem")), p.Placeholder()) || !bytes.Equal(p.Processar(nil), p.Placeholder()) {
		t.Fatal("entrada inválida/vazia deveria devolver o placeholder")
	}
}

// --- HTTP --------------------------------------------------------------

func TestUpdateExigeChaveDeAdmin(t *testing.T) {
	s, _ := servicoDeTeste()
	mux := http.NewServeMux()
	NovoHandler(s, httpx.NovoAdmin("segredo")).Registrar(mux)

	req := func(chave string) int {
		r := httptest.NewRequest(http.MethodPatch, "/api/animal/update?id=1", strings.NewReader(`{"habito":"Noturno"}`))
		if chave != "" {
			r.Header.Set("X-Admin-Key", chave)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w.Code
	}
	if c := req(""); c != http.StatusUnauthorized {
		t.Errorf("sem chave: %d", c)
	}
	if c := req("errada"); c != http.StatusUnauthorized {
		t.Errorf("chave errada: %d", c)
	}
	if c := req("segredo"); c != http.StatusOK {
		t.Errorf("chave certa: %d", c)
	}
}

func TestDetalheRejeitaIDInvalido(t *testing.T) {
	s, _ := servicoDeTeste()
	mux := http.NewServeMux()
	NovoHandler(s, httpx.NovoAdmin("x")).Registrar(mux)

	for _, url := range []string{"/api/animal", "/api/animal?id=abc", "/api/animal?id=-3"} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, url, nil))
		if w.Code != http.StatusBadRequest {
			t.Errorf("%s: status %d, esperado 400", url, w.Code)
		}
	}
}
