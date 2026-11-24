package avistamento

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discoveryrural/internal/falha"
	"discoveryrural/internal/httpx"
	"discoveryrural/internal/limitador"
	"time"
)

type repoFalso struct{ inseridos []Novo }

func (r *repoFalso) BuscarPorID(context.Context, int) (*Avistamento, error) {
	return nil, falha.NaoEncontrado("Avistamento não encontrado")
}
func (r *repoFalso) Inserir(_ context.Context, n Novo) (int64, error) {
	r.inseridos = append(r.inseridos, n)
	return int64(len(r.inseridos)), nil
}
func (r *repoFalso) Linhas(context.Context, int, int) ([]Linha, error) { return nil, nil }
func (r *repoFalso) PontosPorAnimal(context.Context, int, int) (map[int][]Ponto, error) {
	return nil, nil
}

type animaisFalsos map[int]bool

func (a animaisFalsos) Existe(_ context.Context, id int) (bool, error) { return a[id], nil }

func f64(v float64) *float64 { return &v }

func TestValidar(t *testing.T) {
	casos := []struct {
		nome string
		in   Entrada
		erro string
	}{
		{"sem animal", Entrada{Latitude: f64(1), Longitude: f64(1)}, "animal_id obrigatório"},
		{"sem coordenadas", Entrada{AnimalID: 1}, "latitude e longitude obrigatórias"},
		{"latitude fora", Entrada{AnimalID: 1, Latitude: f64(91), Longitude: f64(0)}, "coordenadas fora do intervalo válido"},
		{"longitude fora", Entrada{AnimalID: 1, Latitude: f64(0), Longitude: f64(-181)}, "coordenadas fora do intervalo válido"},
		{"precisão negativa", Entrada{AnimalID: 1, Latitude: f64(0), Longitude: f64(0), Precisao: f64(-1)}, "precisão inválida"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			_, err := validar(c.in)
			f, _ := err.(*falha.Falha)
			if f == nil || f.Mensagem != c.erro {
				t.Fatalf("esperava %q, veio %v", c.erro, err)
			}
		})
	}
}

func TestLocalNomeEhAparadoETruncado(t *testing.T) {
	n, err := validar(Entrada{AnimalID: 1, Latitude: f64(0), Longitude: f64(0), LocalNome: "  " + strings.Repeat("é", 200) + "  "})
	if err != nil {
		t.Fatal(err)
	}
	if got := len([]rune(n.LocalNome)); got != 120 {
		t.Fatalf("local_nome deveria ter 120 caracteres, tem %d", got)
	}
}

func TestRegistrarAnimalInexistente(t *testing.T) {
	s := NovoServico(&repoFalso{}, animaisFalsos{})
	_, err := s.Registrar(context.Background(), Entrada{AnimalID: 9, Latitude: f64(0), Longitude: f64(0)})
	if !falha.EhTipo(err, falha.TipoNaoEncontrado) {
		t.Fatalf("esperava 404, veio %v", err)
	}
}

func TestAgruparMantemAnimalSemPontos(t *testing.T) {
	grupos := agrupar([]Linha{
		{AnimalID: 1, Ponto: &Ponto{Longitude: -35.1, Latitude: -8.2}},
		{AnimalID: 1, Ponto: &Ponto{Longitude: -35.2, Latitude: -8.3}},
		{AnimalID: 2},
	})
	if len(grupos) != 2 || len(grupos[0].Avistamentos) != 2 || len(grupos[1].Avistamentos) != 0 {
		t.Fatalf("agrupamento incorreto: %+v", grupos)
	}
	if grupos[0].Avistamentos[0] != [2]float64{-35.1, -8.2} {
		t.Errorf("ordem deveria ser [longitude, latitude]: %v", grupos[0].Avistamentos[0])
	}
}

func TestHandlerLimitaPostsSeguidos(t *testing.T) {
	servico := NovoServico(&repoFalso{}, animaisFalsos{1: true})
	h := NovoHandler(servico, limitador.Novo(1, 10*time.Second), httpx.ExtratorIP{})
	mux := http.NewServeMux()
	h.Registrar(mux)

	enviar := func() *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodPost, "/api/avistamento",
			strings.NewReader(`{"animal_id":1,"latitude":-8.0,"longitude":-35.0}`))
		r.RemoteAddr = "1.2.3.4:1000"
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}
	if w := enviar(); w.Code != http.StatusCreated {
		t.Fatalf("primeiro POST: %d (%s)", w.Code, w.Body.String())
	}
	if w := enviar(); w.Code != http.StatusTooManyRequests {
		t.Fatalf("segundo POST deveria ser 429, veio %d", w.Code)
	}
}
