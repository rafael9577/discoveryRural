package avistamento

import (
	"context"
	"strings"

	"discoveryrural/internal/falha"
)

const maxLocalNome = 120

// VerificadorAnimal é a dependência (definida aqui, por quem a consome) que
// informa se um animal existe. Quem implementa é o módulo de animais.
type VerificadorAnimal interface {
	Existe(ctx context.Context, id int) (bool, error)
}

// Servico concentra as regras de negócio de avistamentos.
type Servico struct {
	repo    Repositorio
	animais VerificadorAnimal
}

func NovoServico(repo Repositorio, animais VerificadorAnimal) *Servico {
	return &Servico{repo: repo, animais: animais}
}

func (s *Servico) Buscar(ctx context.Context, id int) (*Avistamento, error) {
	return s.repo.BuscarPorID(ctx, id)
}

// Registrar valida a entrada, confere o animal e persiste o avistamento.
func (s *Servico) Registrar(ctx context.Context, in Entrada) (*Criado, error) {
	novo, err := validar(in)
	if err != nil {
		return nil, err
	}

	existe, err := s.animais.Existe(ctx, novo.AnimalID)
	if err != nil {
		return nil, falha.Interno("Erro ao salvar avistamento", err)
	}
	if !existe {
		return nil, falha.NaoEncontrado("Animal não encontrado")
	}

	id, err := s.repo.Inserir(ctx, novo)
	if err != nil {
		return nil, err
	}
	return &Criado{ID: id, AnimalID: novo.AnimalID}, nil
}

// Agrupar devolve os pontos [longitude, latitude] por animal. Se animalID > 0
// e o animal não existir, devolve "não encontrado".
func (s *Servico) Agrupar(ctx context.Context, animalID int) ([]Grupo, error) {
	linhas, err := s.repo.Linhas(ctx, animalID, LimitePontos)
	if err != nil {
		return nil, err
	}
	grupos := agrupar(linhas)
	if animalID > 0 && len(grupos) == 0 {
		return nil, falha.NaoEncontrado("Animal não encontrado")
	}
	return grupos, nil
}

// PontosPorAnimal é usado pelo módulo de animais para anexar os pontos.
func (s *Servico) PontosPorAnimal(ctx context.Context, animalID, limite int) (map[int][]Ponto, error) {
	return s.repo.PontosPorAnimal(ctx, animalID, limite)
}

func agrupar(linhas []Linha) []Grupo {
	grupos := []Grupo{}
	indice := map[int]int{}
	for _, l := range linhas {
		i, ok := indice[l.AnimalID]
		if !ok {
			grupos = append(grupos, Grupo{AnimalID: l.AnimalID, Avistamentos: [][2]float64{}})
			i = len(grupos) - 1
			indice[l.AnimalID] = i
		}
		if l.Ponto != nil {
			grupos[i].Avistamentos = append(grupos[i].Avistamentos, [2]float64{l.Ponto.Longitude, l.Ponto.Latitude})
		}
	}
	return grupos
}

func validar(in Entrada) (Novo, error) {
	if in.AnimalID <= 0 {
		return Novo{}, falha.Validacao("animal_id obrigatório")
	}
	if in.Latitude == nil || in.Longitude == nil {
		return Novo{}, falha.Validacao("latitude e longitude obrigatórias")
	}
	lat, lon := *in.Latitude, *in.Longitude
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return Novo{}, falha.Validacao("coordenadas fora do intervalo válido")
	}
	if in.Precisao != nil && (*in.Precisao < 0 || *in.Precisao > 100000) {
		return Novo{}, falha.Validacao("precisão inválida")
	}

	local := strings.TrimSpace(in.LocalNome)
	if r := []rune(local); len(r) > maxLocalNome {
		local = string(r[:maxLocalNome])
	}
	return Novo{AnimalID: in.AnimalID, Latitude: lat, Longitude: lon, Precisao: in.Precisao, LocalNome: local}, nil
}
