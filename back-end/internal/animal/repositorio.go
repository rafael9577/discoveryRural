package animal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"discoveryrural/internal/falha"
)

// Repositorio abstrai a persistência de animais.
type Repositorio interface {
	Listar(ctx context.Context) ([]Animal, error)
	BuscarPorID(ctx context.Context, id int) (*Animal, error)
	Resumos(ctx context.Context) ([]Resumo, error)
	IDs(ctx context.Context) ([]int, error)
	Existe(ctx context.Context, id int) (bool, error)
	Contar(ctx context.Context) (int, error)
	Criar(ctx context.Context, a Animal, imagem []byte) (int64, error)
	Atualizar(ctx context.Context, id int, alteracoes []Alteracao) error
	Imagem(ctx context.Context, id int) ([]byte, error)
}

const colunasAnimal = `id, nome_popular, nome_cientifico, breve_descricao, estado_conservacao, habito,
	vocalizacao_postura, peconhento, aviso_seguranca, protocolo_incidente, contatos_resgate,
	raridade, bitmap_id`

// colunasEditaveis é a lista branca de colunas aceitas em Atualizar.
var colunasEditaveis = map[string]struct{}{
	"nome_popular": {}, "nome_cientifico": {}, "breve_descricao": {}, "estado_conservacao": {},
	"habito": {}, "vocalizacao_postura": {}, "aviso_seguranca": {}, "protocolo_incidente": {},
	"contatos_resgate": {}, "peconhento": {}, "raridade": {}, "bitmap_id": {},
}

// RepositorioSQL implementa Repositorio sobre database/sql.
type RepositorioSQL struct {
	db *sql.DB
}

func NovoRepositorioSQL(db *sql.DB) *RepositorioSQL { return &RepositorioSQL{db: db} }

type escaneavel interface{ Scan(dest ...any) error }

func escanear(s escaneavel) (Animal, error) {
	var a Animal
	var nomeCientifico, breve, estado, habito, vocal, aviso, protocolo, contatos sql.NullString
	var peconhento, raridade, bitmap sql.NullInt64
	if err := s.Scan(&a.ID, &a.NomePopular, &nomeCientifico, &breve, &estado, &habito, &vocal,
		&peconhento, &aviso, &protocolo, &contatos, &raridade, &bitmap); err != nil {
		return Animal{}, err
	}
	a.NomeCientifico = nomeCientifico.String
	a.BreveDescricao = breve.String
	a.EstadoConservacao = estado.String
	a.Habito = habito.String
	a.VocalizacaoPostura = vocal.String
	a.Peconhento = peconhento.Int64 == 1
	a.AvisoSeguranca = aviso.String
	a.ProtocoloIncidente = protocolo.String
	a.ContatosResgate = contatos.String
	a.Raridade = int(raridade.Int64)
	a.BitmapID = int(bitmap.Int64)
	return a, nil
}

func (r *RepositorioSQL) Listar(ctx context.Context) ([]Animal, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT "+colunasAnimal+" FROM animais ORDER BY id")
	if err != nil {
		return nil, falha.Interno("Erro ao consultar animais", err)
	}
	defer rows.Close()

	animais := []Animal{}
	for rows.Next() {
		a, err := escanear(rows)
		if err != nil {
			return nil, falha.Interno("Erro ao consultar animais", err)
		}
		animais = append(animais, a)
	}
	if err := rows.Err(); err != nil {
		return nil, falha.Interno("Erro ao consultar animais", err)
	}
	return animais, nil
}

func (r *RepositorioSQL) BuscarPorID(ctx context.Context, id int) (*Animal, error) {
	a, err := escanear(r.db.QueryRowContext(ctx, "SELECT "+colunasAnimal+" FROM animais WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, falha.NaoEncontrado("Animal não encontrado")
	}
	if err != nil {
		return nil, falha.Interno("Erro ao consultar animais", err)
	}
	return &a, nil
}

func (r *RepositorioSQL) Resumos(ctx context.Context) ([]Resumo, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, nome_popular, breve_descricao FROM animais ORDER BY id")
	if err != nil {
		return nil, falha.Interno("Erro ao consultar animais", err)
	}
	defer rows.Close()

	lista := []Resumo{}
	for rows.Next() {
		var x Resumo
		var breve sql.NullString
		if err := rows.Scan(&x.ID, &x.NomePopular, &breve); err != nil {
			return nil, falha.Interno("Erro ao consultar animais", err)
		}
		x.BreveDescricao = breve.String
		lista = append(lista, x)
	}
	if err := rows.Err(); err != nil {
		return nil, falha.Interno("Erro ao consultar animais", err)
	}
	return lista, nil
}

func (r *RepositorioSQL) IDs(ctx context.Context) ([]int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id FROM animais ORDER BY id")
	if err != nil {
		return nil, falha.Interno("Erro ao consultar animais", err)
	}
	defer rows.Close()

	ids := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, falha.Interno("Erro ao consultar animais", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, falha.Interno("Erro ao consultar animais", err)
	}
	return ids, nil
}

func (r *RepositorioSQL) Existe(ctx context.Context, id int) (bool, error) {
	var um int
	err := r.db.QueryRowContext(ctx, "SELECT 1 FROM animais WHERE id = ?", id).Scan(&um)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RepositorioSQL) Contar(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM animais").Scan(&n)
	return n, err
}

func (r *RepositorioSQL) Criar(ctx context.Context, a Animal, imagem []byte) (int64, error) {
	peconhento := 0
	if a.Peconhento {
		peconhento = 1
	}
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO animais (
			nome_popular, nome_cientifico, imagem, breve_descricao,
			estado_conservacao, habito, vocalizacao_postura, peconhento,
			aviso_seguranca, protocolo_incidente, contatos_resgate, raridade, bitmap_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.NomePopular, a.NomeCientifico, imagem, a.BreveDescricao,
		a.EstadoConservacao, a.Habito, a.VocalizacaoPostura, peconhento,
		a.AvisoSeguranca, a.ProtocoloIncidente, a.ContatosResgate, a.Raridade, a.BitmapID,
	)
	if err != nil {
		return 0, falha.Interno("Erro ao cadastrar animal", err)
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func (r *RepositorioSQL) Atualizar(ctx context.Context, id int, alteracoes []Alteracao) error {
	sets := make([]string, 0, len(alteracoes))
	args := make([]any, 0, len(alteracoes)+1)
	for _, a := range alteracoes {
		if _, ok := colunasEditaveis[a.Coluna]; !ok {
			return falha.Interno("Erro ao atualizar animal", fmt.Errorf("coluna não editável: %q", a.Coluna))
		}
		sets = append(sets, a.Coluna+" = ?")
		args = append(args, a.Valor)
	}
	args = append(args, id)
	if _, err := r.db.ExecContext(ctx, "UPDATE animais SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...); err != nil {
		return falha.Interno("Erro ao atualizar animal", err)
	}
	return nil
}

func (r *RepositorioSQL) Imagem(ctx context.Context, id int) ([]byte, error) {
	var img []byte
	err := r.db.QueryRowContext(ctx, "SELECT imagem FROM animais WHERE id = ?", id).Scan(&img)
	if (err != nil && errors.Is(err, sql.ErrNoRows)) || (err == nil && len(img) == 0) {
		return nil, falha.NaoEncontrado("Imagem não encontrada")
	}
	if err != nil {
		return nil, falha.Interno("Erro ao consultar imagem", err)
	}
	return img, nil
}
