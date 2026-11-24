// Package config carrega a configuração da aplicação a partir do ambiente.
package config

import (
	"errors"
	"os"
)

type Config struct {
	BancoURL     string
	BancoToken   string
	ChaveAdmin   string
	Porta        string
	ConfiarProxy bool // true quando RENDER=true (usa X-Forwarded-For)
}

// Carregar lê as variáveis de ambiente (o .env já deve ter sido carregado).
func Carregar() (Config, error) {
	c := Config{
		BancoURL:     os.Getenv("TURSO_DATABASE_URL"),
		BancoToken:   os.Getenv("TURSO_AUTH_TOKEN"),
		ChaveAdmin:   os.Getenv("ADMIN_KEY"),
		Porta:        os.Getenv("PORT"),
		ConfiarProxy: os.Getenv("RENDER") == "true",
	}
	if c.BancoURL == "" || c.BancoToken == "" {
		return Config{}, errors.New("defina TURSO_DATABASE_URL e TURSO_AUTH_TOKEN (via .env ou variáveis de ambiente)")
	}
	if c.Porta == "" {
		c.Porta = "8080"
	}
	return c, nil
}

// DSN monta a string de conexão do driver libsql.
func (c Config) DSN() string {
	return c.BancoURL + "?authToken=" + c.BancoToken
}
