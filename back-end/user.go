package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"golang.org/x/crypto/bcrypt"
)

const duracaoSessao = 30 * 24 * time.Hour

var reNomeUsuario = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

type usuario struct {
	ID          int64   `json:"id"`
	NomeUsuario string  `json:"nome_usuario"`
	CPF         string  `json:"cpf"`
	Email       string  `json:"email"`
	DataCriacao string  `json:"data_criacao"`
	Reputacao   float64 `json:"reputacao"`
}

type limitador struct {
	mu      sync.Mutex
	acessos map[string][]time.Time
	max     int
	janela  time.Duration
}

func novoLimitador(max int, janela time.Duration) *limitador {
	return &limitador{acessos: map[string][]time.Time{}, max: max, janela: janela}
}

func (l *limitador) permitir(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	agora := time.Now()
	recentes := l.acessos[ip][:0]
	for _, t := range l.acessos[ip] {
		if agora.Sub(t) < l.janela {
			recentes = append(recentes, t)
		}
	}
	if len(recentes) >= l.max {
		l.acessos[ip] = recentes
		return false
	}
	l.acessos[ip] = append(recentes, agora)
	return true
}

func ipDe(r *http.Request) string {
	if os.Getenv("RENDER") == "true" {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			return strings.TrimSpace(strings.Split(xff, ",")[0])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func escreverJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func escreverErro(w http.ResponseWriter, status int, msg string) {
	escreverJSON(w, status, map[string]string{"erro": msg})
}

func lerJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<16)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		escreverErro(w, http.StatusBadRequest, "JSON inválido")
		return false
	}
	return true
}

func apenasDigitos(s string) string {
	var b strings.Builder
	for _, c := range s {
		if c >= '0' && c <= '9' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func cpfValido(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}
	iguais := true
	for i := 1; i < 11; i++ {
		if cpf[i] != cpf[0] {
			iguais = false
			break
		}
	}
	if iguais {
		return false
	}
	for t := 9; t < 11; t++ {
		soma := 0
		for i := 0; i < t; i++ {
			soma += int(cpf[i]-'0') * (t + 1 - i)
		}
		d := (soma * 10) % 11
		if d == 10 {
			d = 0
		}
		if d != int(cpf[t]-'0') {
			return false
		}
	}
	return true
}

func mascararCPF(cpf string) string {
	if len(cpf) != 11 {
		return ""
	}
	return "***." + cpf[3:6] + "." + cpf[6:9] + "-**"
}

func hashToken(t string) string {
	h := sha256.Sum256([]byte(t))
	return hex.EncodeToString(h[:])
}

func criarSessao(db *sql.DB, usuarioID int64) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	_, err := db.Exec(
		"INSERT INTO sessoes (token_hash, usuario_id, expira_em) VALUES (?, ?, ?)",
		hashToken(token), usuarioID, time.Now().Add(duracaoSessao).Unix(),
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func autenticar(db *sql.DB, r *http.Request) (*usuario, string, bool) {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return nil, "", false
	}
	th := hashToken(strings.TrimSpace(h[7:]))
	var u usuario
	err := db.QueryRow(`
		SELECT u.id, u.nome_usuario, u.cpf, u.email, u.data_criacao, u.reputacao
		FROM sessoes s JOIN usuarios u ON u.id = s.usuario_id
		WHERE s.token_hash = ? AND s.expira_em > ?`,
		th, time.Now().Unix(),
	).Scan(&u.ID, &u.NomeUsuario, &u.CPF, &u.Email, &u.DataCriacao, &u.Reputacao)
	if err != nil {
		return nil, "", false
	}
	u.CPF = mascararCPF(u.CPF)
	return &u, th, true
}

func registrarUsuarios(adminKey string) {
	url := os.Getenv("USERS_DATABASE_URL")
	token := os.Getenv("USERS_AUTH_TOKEN")

	if url == "" || token == "" {
		log.Fatal("defina USERS_DATABASE_URL e USERS_AUTH_TOKEN (via .env ou variáveis de ambiente)")
	}

	db, err := sql.Open("libsql", url+"?authToken="+token)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome_usuario TEXT NOT NULL UNIQUE COLLATE NOCASE,
		cpf TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE COLLATE NOCASE,
		senha_hash TEXT NOT NULL,
		data_criacao DATETIME DEFAULT CURRENT_TIMESTAMP,
		reputacao REAL NOT NULL DEFAULT 5 CHECK(reputacao BETWEEN 0 AND 10)
	);

	CREATE TABLE IF NOT EXISTS sessoes (
		token_hash TEXT PRIMARY KEY,
		usuario_id INTEGER NOT NULL,
		expira_em INTEGER NOT NULL,
		FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
	);
	`)
	if err != nil {
		log.Fatal(err)
	}

	db.Exec("DELETE FROM sessoes WHERE expira_em < ?", time.Now().Unix())

	hashFalso, _ := bcrypt.GenerateFromPassword([]byte("senha-falsa-para-tempo-constante"), bcrypt.DefaultCost)
	limite := novoLimitador(10, 5*time.Minute)

	http.HandleFunc("/api/users/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}
		if !limite.permitir(ipDe(r)) {
			escreverErro(w, http.StatusTooManyRequests, "Muitas tentativas, aguarde alguns minutos")
			return
		}

		var in struct {
			NomeUsuario string `json:"nome_usuario"`
			CPF         string `json:"cpf"`
			Email       string `json:"email"`
			Senha       string `json:"senha"`
		}
		if !lerJSON(w, r, &in) {
			return
		}

		nome := strings.TrimSpace(in.NomeUsuario)
		cpf := apenasDigitos(in.CPF)
		email := strings.ToLower(strings.TrimSpace(in.Email))

		if !reNomeUsuario.MatchString(nome) {
			escreverErro(w, http.StatusBadRequest, "Nome de usuário deve ter 3 a 20 caracteres (letras, números e _)")
			return
		}
		if !cpfValido(cpf) {
			escreverErro(w, http.StatusBadRequest, "CPF inválido")
			return
		}
		addr, err := mail.ParseAddress(email)
		if err != nil || addr.Address != email || len(email) > 254 {
			escreverErro(w, http.StatusBadRequest, "Email inválido")
			return
		}
		if len(in.Senha) < 8 || len(in.Senha) > 72 {
			escreverErro(w, http.StatusBadRequest, "A senha deve ter entre 8 e 72 caracteres")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(in.Senha), bcrypt.DefaultCost)
		if err != nil {
			escreverErro(w, http.StatusInternalServerError, "Erro interno")
			return
		}

		u := usuario{NomeUsuario: nome, Email: email}
		err = db.QueryRow(
			"INSERT INTO usuarios (nome_usuario, cpf, email, senha_hash) VALUES (?, ?, ?, ?) RETURNING id, data_criacao, reputacao",
			nome, cpf, email, string(hash),
		).Scan(&u.ID, &u.DataCriacao, &u.Reputacao)
		if err != nil {
			msg := err.Error()
			switch {
			case strings.Contains(msg, "UNIQUE") && strings.Contains(msg, "email"):
				escreverErro(w, http.StatusConflict, "Email já cadastrado")
			case strings.Contains(msg, "UNIQUE") && strings.Contains(msg, "cpf"):
				escreverErro(w, http.StatusConflict, "CPF já cadastrado")
			case strings.Contains(msg, "UNIQUE") && strings.Contains(msg, "nome_usuario"):
				escreverErro(w, http.StatusConflict, "Nome de usuário já em uso")
			default:
				escreverErro(w, http.StatusInternalServerError, "Erro ao criar usuário")
			}
			return
		}
		u.CPF = mascararCPF(cpf)

		tk, err := criarSessao(db, u.ID)
		if err != nil {
			escreverErro(w, http.StatusInternalServerError, "Erro ao criar sessão")
			return
		}
		escreverJSON(w, http.StatusCreated, map[string]any{"token": tk, "usuario": u})
	})

	http.HandleFunc("/api/users/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}
		if !limite.permitir(ipDe(r)) {
			escreverErro(w, http.StatusTooManyRequests, "Muitas tentativas, aguarde alguns minutos")
			return
		}

		var in struct {
			Identificador string `json:"identificador"`
			Senha         string `json:"senha"`
		}
		if !lerJSON(w, r, &in) {
			return
		}

		ident := strings.TrimSpace(in.Identificador)
		coluna := "nome_usuario"
		if strings.Contains(ident, "@") {
			coluna = "email"
			ident = strings.ToLower(ident)
		}

		var u usuario
		var senhaHash string
		err := db.QueryRow(
			"SELECT id, nome_usuario, cpf, email, data_criacao, reputacao, senha_hash FROM usuarios WHERE "+coluna+" = ?",
			ident,
		).Scan(&u.ID, &u.NomeUsuario, &u.CPF, &u.Email, &u.DataCriacao, &u.Reputacao, &senhaHash)

		if err != nil {
			bcrypt.CompareHashAndPassword(hashFalso, []byte(in.Senha))
			escreverErro(w, http.StatusUnauthorized, "Credenciais inválidas")
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(senhaHash), []byte(in.Senha)) != nil {
			escreverErro(w, http.StatusUnauthorized, "Credenciais inválidas")
			return
		}
		u.CPF = mascararCPF(u.CPF)

		tk, err := criarSessao(db, u.ID)
		if err != nil {
			escreverErro(w, http.StatusInternalServerError, "Erro ao criar sessão")
			return
		}
		escreverJSON(w, http.StatusOK, map[string]any{"token": tk, "usuario": u})
	})

	http.HandleFunc("/api/users/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}
		u, _, ok := autenticar(db, r)
		if !ok {
			escreverErro(w, http.StatusUnauthorized, "Não autorizado")
			return
		}
		escreverJSON(w, http.StatusOK, u)
	})

	http.HandleFunc("/api/users/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}
		_, th, ok := autenticar(db, r)
		if !ok {
			escreverErro(w, http.StatusUnauthorized, "Não autorizado")
			return
		}
		db.Exec("DELETE FROM sessoes WHERE token_hash = ?", th)
		escreverJSON(w, http.StatusOK, map[string]bool{"ok": true})
	})

	http.HandleFunc("/api/users/reputacao", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}
		key := r.Header.Get("X-Admin-Key")
		if adminKey == "" || subtle.ConstantTimeCompare([]byte(key), []byte(adminKey)) != 1 {
			escreverErro(w, http.StatusUnauthorized, "Não autorizado")
			return
		}

		var in struct {
			ID    int64   `json:"id"`
			Delta float64 `json:"delta"`
		}
		if !lerJSON(w, r, &in) {
			return
		}

		var nova float64
		err := db.QueryRow(
			"UPDATE usuarios SET reputacao = MAX(0, MIN(10, reputacao + ?)) WHERE id = ? RETURNING reputacao",
			in.Delta, in.ID,
		).Scan(&nova)
		if err != nil {
			escreverErro(w, http.StatusNotFound, "Usuário não encontrado")
			return
		}
		escreverJSON(w, http.StatusOK, map[string]any{"id": in.ID, "reputacao": nova})
	})
}
