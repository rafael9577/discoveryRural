// Package banco cuida da conexão com o Turso (libsql) e do esquema das tabelas.
package banco

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Abrir cria o pool de conexões (compartilhado por toda a aplicação).
func Abrir(dsn string) (*sql.DB, error) {
	return sql.Open("libsql", dsn)
}

var esquema = []string{
	`CREATE TABLE IF NOT EXISTS animais (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome_popular TEXT NOT NULL,
		nome_cientifico TEXT,
		imagem BLOB,
		breve_descricao TEXT,
		estado_conservacao TEXT,
		habito TEXT,
		vocalizacao_postura TEXT,
		peconhento INTEGER NOT NULL,
		aviso_seguranca TEXT,
		protocolo_incidente TEXT,
		contatos_resgate TEXT,
		raridade INTEGER CHECK(raridade BETWEEN 1 AND 5),
		bitmap_id INTEGER DEFAULT 0
	)`,
	`CREATE TABLE IF NOT EXISTS relatos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		animal_id INTEGER NOT NULL,
		latitude REAL NOT NULL,
		longitude REAL NOT NULL,
		precisao REAL,
		local_nome TEXT,
		data_hora DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (animal_id) REFERENCES animais(id)
	)`,
	`CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome_usuario TEXT NOT NULL UNIQUE COLLATE NOCASE,
		cpf TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE COLLATE NOCASE,
		senha_hash TEXT NOT NULL,
		data_criacao DATETIME DEFAULT CURRENT_TIMESTAMP,
		reputacao REAL NOT NULL DEFAULT 5 CHECK(reputacao BETWEEN 0 AND 10)
	)`,
	`CREATE TABLE IF NOT EXISTS sessoes (
		token_hash TEXT PRIMARY KEY,
		usuario_id INTEGER NOT NULL,
		expira_em INTEGER NOT NULL,
		FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
	)`,
}

// Migrar cria as tabelas que ainda não existem.
func Migrar(ctx context.Context, db *sql.DB) error {
	for i, ddl := range esquema {
		if _, err := db.ExecContext(ctx, ddl); err != nil {
			return fmt.Errorf("migração %d: %w", i+1, err)
		}
	}
	return nil
}
