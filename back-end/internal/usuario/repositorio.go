package usuario

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"discoveryrural/internal/falha"
)

// RepositorioUsuarios abstrai a persistência de usuários.
type RepositorioUsuarios interface {
	// Criar devolve o usuário gravado (CPF ainda sem máscara) ou Conflito.
	Criar(ctx context.Context, n NovoUsuario) (*Usuario, error)
	// BuscarCredencial procura por e-mail (porEmail) ou nome de usuário.
	// Devolve também o hash da senha; ErrNoRows vira (nil, "", nil).
	BuscarCredencial(ctx context.Context, identificador string, porEmail bool) (*Usuario, string, error)
	AjustarReputacao(ctx context.Context, id int64, delta float64) (float64, error)
}

// RepositorioSessoes abstrai a persistência de sessões.
type RepositorioSessoes interface {
	Criar(ctx context.Context, tokenHash string, usuarioID, expiraEm int64) error
	// BuscarUsuario devolve o dono de uma sessão válida; (nil, nil) se não houver.
	BuscarUsuario(ctx context.Context, tokenHash string, agora int64) (*Usuario, error)
	Remover(ctx context.Context, tokenHash string) error
	RemoverExpiradas(ctx context.Context, agora int64) error
}

// RepositorioUsuariosSQL implementa RepositorioUsuarios.
type RepositorioUsuariosSQL struct{ db *sql.DB }

func NovoRepositorioUsuariosSQL(db *sql.DB) *RepositorioUsuariosSQL {
	return &RepositorioUsuariosSQL{db: db}
}

func (r *RepositorioUsuariosSQL) Criar(ctx context.Context, n NovoUsuario) (*Usuario, error) {
	u := Usuario{NomeUsuario: n.NomeUsuario, Email: n.Email, CPF: n.CPF}
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO usuarios (nome_usuario, cpf, email, senha_hash) VALUES (?, ?, ?, ?) RETURNING id, data_criacao, reputacao",
		n.NomeUsuario, n.CPF, n.Email, n.SenhaHash,
	).Scan(&u.ID, &u.DataCriacao, &u.Reputacao)
	if err != nil {
		return nil, traduzirErroCriacao(err)
	}
	return &u, nil
}

// traduzirErroCriacao converte violações de UNIQUE em Conflito específico.
func traduzirErroCriacao(err error) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "UNIQUE") && strings.Contains(msg, "email"):
		return falha.Conflito("Email já cadastrado")
	case strings.Contains(msg, "UNIQUE") && strings.Contains(msg, "cpf"):
		return falha.Conflito("CPF já cadastrado")
	case strings.Contains(msg, "UNIQUE") && strings.Contains(msg, "nome_usuario"):
		return falha.Conflito("Nome de usuário já em uso")
	default:
		return falha.Interno("Erro ao criar usuário", err)
	}
}

func (r *RepositorioUsuariosSQL) BuscarCredencial(ctx context.Context, identificador string, porEmail bool) (*Usuario, string, error) {
	coluna := "nome_usuario"
	if porEmail {
		coluna = "email"
	}
	var u Usuario
	var hash string
	err := r.db.QueryRowContext(ctx,
		"SELECT id, nome_usuario, cpf, email, data_criacao, reputacao, senha_hash FROM usuarios WHERE "+coluna+" = ?",
		identificador,
	).Scan(&u.ID, &u.NomeUsuario, &u.CPF, &u.Email, &u.DataCriacao, &u.Reputacao, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}
	return &u, hash, nil
}

func (r *RepositorioUsuariosSQL) AjustarReputacao(ctx context.Context, id int64, delta float64) (float64, error) {
	var nova float64
	err := r.db.QueryRowContext(ctx,
		"UPDATE usuarios SET reputacao = MAX(0, MIN(10, reputacao + ?)) WHERE id = ? RETURNING reputacao",
		delta, id,
	).Scan(&nova)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, falha.NaoEncontrado("Usuário não encontrado")
	}
	if err != nil {
		return 0, falha.Interno("Erro ao atualizar reputação", err)
	}
	return nova, nil
}

// RepositorioSessoesSQL implementa RepositorioSessoes.
type RepositorioSessoesSQL struct{ db *sql.DB }

func NovoRepositorioSessoesSQL(db *sql.DB) *RepositorioSessoesSQL {
	return &RepositorioSessoesSQL{db: db}
}

func (r *RepositorioSessoesSQL) Criar(ctx context.Context, tokenHash string, usuarioID, expiraEm int64) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO sessoes (token_hash, usuario_id, expira_em) VALUES (?, ?, ?)",
		tokenHash, usuarioID, expiraEm)
	return err
}

func (r *RepositorioSessoesSQL) BuscarUsuario(ctx context.Context, tokenHash string, agora int64) (*Usuario, error) {
	var u Usuario
	err := r.db.QueryRowContext(ctx, `
		SELECT u.id, u.nome_usuario, u.cpf, u.email, u.data_criacao, u.reputacao
		FROM sessoes s JOIN usuarios u ON u.id = s.usuario_id
		WHERE s.token_hash = ? AND s.expira_em > ?`,
		tokenHash, agora,
	).Scan(&u.ID, &u.NomeUsuario, &u.CPF, &u.Email, &u.DataCriacao, &u.Reputacao)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *RepositorioSessoesSQL) Remover(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessoes WHERE token_hash = ?", tokenHash)
	return err
}

func (r *RepositorioSessoesSQL) RemoverExpiradas(ctx context.Context, agora int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessoes WHERE expira_em < ?", agora)
	return err
}
