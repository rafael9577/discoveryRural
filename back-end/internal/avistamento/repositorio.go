package avistamento

import (
	"context"
	"database/sql"
	"errors"

	"discoveryrural/internal/falha"
)

// Repositorio abstrai a persistência de avistamentos.
type Repositorio interface {
	BuscarPorID(ctx context.Context, id int) (*Avistamento, error)
	Inserir(ctx context.Context, n Novo) (int64, error)
	// Linhas devolve animais com seus pontos; animalID 0 = todos.
	Linhas(ctx context.Context, animalID, limite int) ([]Linha, error)
	// PontosPorAnimal devolve os pontos agrupados por animal; animalID 0 = todos.
	PontosPorAnimal(ctx context.Context, animalID, limite int) (map[int][]Ponto, error)
}

// RepositorioSQL implementa Repositorio sobre database/sql.
type RepositorioSQL struct {
	db *sql.DB
}

func NovoRepositorioSQL(db *sql.DB) *RepositorioSQL { return &RepositorioSQL{db: db} }

func (r *RepositorioSQL) BuscarPorID(ctx context.Context, id int) (*Avistamento, error) {
	var a Avistamento
	var precisao sql.NullFloat64
	var local, dataHora sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT animal_id, latitude, longitude, precisao, local_nome, data_hora FROM relatos WHERE id = ?", id,
	).Scan(&a.AnimalID, &a.Latitude, &a.Longitude, &precisao, &local, &dataHora)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, falha.NaoEncontrado("Avistamento não encontrado")
	}
	if err != nil {
		return nil, falha.Interno("Erro ao consultar avistamento", err)
	}
	a.ID = id
	a.LocalNome = local.String
	a.DataHora = dataHora.String
	if precisao.Valid {
		a.Precisao = &precisao.Float64
	}
	return &a, nil
}

func (r *RepositorioSQL) Inserir(ctx context.Context, n Novo) (int64, error) {
	var precisao any
	if n.Precisao != nil {
		precisao = *n.Precisao
	}
	var local any
	if n.LocalNome != "" {
		local = n.LocalNome
	}
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO relatos (animal_id, latitude, longitude, precisao, local_nome) VALUES (?, ?, ?, ?, ?)",
		n.AnimalID, n.Latitude, n.Longitude, precisao, local,
	)
	if err != nil {
		return 0, falha.Interno("Erro ao salvar avistamento", err)
	}
	id, _ := res.LastInsertId()
	return id, nil
}

// Linhas faz o LEFT JOIN para que animais sem relatos também apareçam.
func (r *RepositorioSQL) Linhas(ctx context.Context, animalID, limite int) ([]Linha, error) {
	query := `
		SELECT a.id, r.longitude, r.latitude
		FROM animais a
		LEFT JOIN relatos r ON r.animal_id = a.id`
	args := []any{}
	if animalID > 0 {
		query += " WHERE a.id = ?"
		args = append(args, animalID)
	}
	query += " ORDER BY a.id, r.id LIMIT ?"
	args = append(args, limite)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, falha.Interno("Erro ao consultar avistamentos", err)
	}
	defer rows.Close()

	linhas := []Linha{}
	for rows.Next() {
		var l Linha
		var lon, lat sql.NullFloat64
		if err := rows.Scan(&l.AnimalID, &lon, &lat); err != nil {
			return nil, falha.Interno("Erro ao consultar avistamentos", err)
		}
		if lon.Valid && lat.Valid {
			l.Ponto = &Ponto{Longitude: lon.Float64, Latitude: lat.Float64}
		}
		linhas = append(linhas, l)
	}
	if err := rows.Err(); err != nil {
		return nil, falha.Interno("Erro ao consultar avistamentos", err)
	}
	return linhas, nil
}

func (r *RepositorioSQL) PontosPorAnimal(ctx context.Context, animalID, limite int) (map[int][]Ponto, error) {
	query := "SELECT animal_id, longitude, latitude FROM relatos"
	args := []any{}
	if animalID > 0 {
		query += " WHERE animal_id = ?"
		args = append(args, animalID)
	}
	query += " ORDER BY animal_id, id LIMIT ?"
	args = append(args, limite)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, falha.Interno("Erro ao consultar avistamentos", err)
	}
	defer rows.Close()

	pontos := map[int][]Ponto{}
	for rows.Next() {
		var id int
		var p Ponto
		if err := rows.Scan(&id, &p.Longitude, &p.Latitude); err != nil {
			return nil, falha.Interno("Erro ao consultar avistamentos", err)
		}
		pontos[id] = append(pontos[id], p)
	}
	if err := rows.Err(); err != nil {
		return nil, falha.Interno("Erro ao consultar avistamentos", err)
	}
	return pontos, nil
}
