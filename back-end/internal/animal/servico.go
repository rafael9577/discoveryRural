package animal

import (
	"context"
	"strconv"

	"discoveryrural/internal/avistamento"
	"discoveryrural/internal/falha"
)

// ProvedorPontos fornece os pontos de avistamento por animal (definido aqui,
// por quem consome; implementado pelo módulo de avistamentos).
type ProvedorPontos interface {
	PontosPorAnimal(ctx context.Context, animalID, limite int) (map[int][]avistamento.Ponto, error)
}

// Servico concentra as regras de negócio de animais.
type Servico struct {
	repo    Repositorio
	pontos  ProvedorPontos
	imagens ProcessadorImagem
}

func NovoServico(repo Repositorio, pontos ProvedorPontos, imagens ProcessadorImagem) *Servico {
	return &Servico{repo: repo, pontos: pontos, imagens: imagens}
}

func urlImagem(id int) string { return "/api/animal/imagem?id=" + strconv.Itoa(id) }

// Listar devolve todos os animais com seus avistamentos.
func (s *Servico) Listar(ctx context.Context) ([]AnimalComAvistamentos, error) {
	animais, err := s.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	pontos, err := s.pontos.PontosPorAnimal(ctx, 0, avistamento.LimitePontos)
	if err != nil {
		return nil, err
	}
	return montar(animais, pontos), nil
}

// Buscar devolve os dados de um animal (sem a lista de avistamentos).
func (s *Servico) Buscar(ctx context.Context, id int) (*Animal, error) {
	a, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	a.ImagemURL = urlImagem(a.ID)
	return a, nil
}

func (s *Servico) Resumos(ctx context.Context) ([]Resumo, error) {
	lista, err := s.repo.Resumos(ctx)
	if err != nil {
		return nil, err
	}
	for i := range lista {
		lista[i].ImagemURL = urlImagem(lista[i].ID)
	}
	return lista, nil
}

func (s *Servico) IDs(ctx context.Context) ([]int, error) { return s.repo.IDs(ctx) }

func (s *Servico) Imagem(ctx context.Context, id int) ([]byte, error) {
	return s.repo.Imagem(ctx, id)
}

// Existe permite que outros módulos verifiquem a existência de um animal.
func (s *Servico) Existe(ctx context.Context, id int) (bool, error) {
	return s.repo.Existe(ctx, id)
}

// Criar valida os dados, processa a foto (ou usa o placeholder) e grava.
func (s *Servico) Criar(ctx context.Context, in NovoAnimal, foto []byte) (int64, error) {
	normalizado, err := in.paraEdicao().Normalizar()
	if err != nil {
		return 0, err
	}
	return s.repo.Criar(ctx, normalizado.paraAnimal(), s.imagens.Processar(foto))
}

// Atualizar aplica uma edição parcial e devolve o animal já atualizado.
func (s *Servico) Atualizar(ctx context.Context, id int, in EdicaoAnimal) (*AnimalComAvistamentos, error) {
	normalizado, err := in.Normalizar()
	if err != nil {
		return nil, err
	}
	alteracoes := normalizado.Alteracoes()
	if len(alteracoes) == 0 {
		return nil, falha.Validacao("nenhum campo para atualizar")
	}

	existe, err := s.repo.Existe(ctx, id)
	if err != nil {
		return nil, falha.Interno("Erro ao atualizar animal", err)
	}
	if !existe {
		return nil, falha.NaoEncontrado("Animal não encontrado")
	}

	if err := s.repo.Atualizar(ctx, id, alteracoes); err != nil {
		return nil, err
	}

	a, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, falha.Interno("Animal atualizado, mas não foi possível relê-lo", err)
	}
	pontos, err := s.pontos.PontosPorAnimal(ctx, id, avistamento.LimitePontos)
	if err != nil {
		return nil, falha.Interno("Animal atualizado, mas não foi possível relê-lo", err)
	}
	resultado := montar([]Animal{*a}, pontos)
	return &resultado[0], nil
}

func montar(animais []Animal, pontos map[int][]avistamento.Ponto) []AnimalComAvistamentos {
	saida := make([]AnimalComAvistamentos, 0, len(animais))
	for _, a := range animais {
		a.ImagemURL = urlImagem(a.ID)
		lista, ok := pontos[a.ID]
		if !ok {
			lista = []avistamento.Ponto{}
		}
		saida = append(saida, AnimalComAvistamentos{Animal: a, Avistamentos: lista})
	}
	return saida
}
