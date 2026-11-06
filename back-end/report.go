package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type novoRelato struct {
	AnimalID  int      `json:"animal_id"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Precisao  *float64 `json:"precisao"`
	LocalNome string   `json:"local_nome"`
}

type limitadorRelato struct {
	mu      sync.Mutex
	ultimos map[string]time.Time
	janela  time.Duration
}

func novoLimitadorRelato(janela time.Duration) *limitadorRelato {
	l := &limitadorRelato{ultimos: make(map[string]time.Time), janela: janela}
	go func() {
		for range time.Tick(10 * time.Minute) {
			l.mu.Lock()
			agora := time.Now()
			for k, t := range l.ultimos {
				if agora.Sub(t) > l.janela {
					delete(l.ultimos, k)
				}
			}
			l.mu.Unlock()
		}
	}()
	return l
}

func (l *limitadorRelato) permitir(chave string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	agora := time.Now()
	if t, ok := l.ultimos[chave]; ok && agora.Sub(t) < l.janela {
		return false
	}
	l.ultimos[chave] = agora
	return true
}

func ipCliente(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func comCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

func responderJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func registrarRelatos() {
	url := os.Getenv("TURSO_DATABASE_URL")
	token := os.Getenv("TURSO_AUTH_TOKEN")

	if url == "" || token == "" {
		log.Fatal("defina TURSO_DATABASE_URL e TURSO_AUTH_TOKEN (via .env ou variáveis de ambiente)")
	}

	db, err := sql.Open("libsql", url+"?authToken="+token)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS relatos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		animal_id INTEGER NOT NULL,
		latitude REAL NOT NULL,
		longitude REAL NOT NULL,
		precisao REAL,
		local_nome TEXT,
		data_hora DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (animal_id) REFERENCES animais(id)
	);
	CREATE INDEX IF NOT EXISTS idx_relatos_animal ON relatos(animal_id);
	`)
	if err != nil {
		log.Fatal(err)
	}

	lim := novoLimitadorRelato(10 * time.Second)

	http.HandleFunc("/api/relato", comCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		if !lim.permitir(ipCliente(r)) {
			http.Error(w, "Muitos relatos seguidos, aguarde alguns segundos", http.StatusTooManyRequests)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)

		var in novoRelato
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&in); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if in.AnimalID <= 0 {
			http.Error(w, "animal_id obrigatório", http.StatusBadRequest)
			return
		}
		if in.Latitude == nil || in.Longitude == nil {
			http.Error(w, "latitude e longitude obrigatórias", http.StatusBadRequest)
			return
		}
		lat, lon := *in.Latitude, *in.Longitude
		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			http.Error(w, "coordenadas fora do intervalo válido", http.StatusBadRequest)
			return
		}
		if in.Precisao != nil && (*in.Precisao < 0 || *in.Precisao > 100000) {
			http.Error(w, "precisão inválida", http.StatusBadRequest)
			return
		}

		local := strings.TrimSpace(in.LocalNome)
		if len([]rune(local)) > 120 {
			local = string([]rune(local)[:120])
		}

		var existe int
		if err := db.QueryRow("SELECT 1 FROM animais WHERE id = ?", in.AnimalID).Scan(&existe); err != nil {
			http.Error(w, "Animal não encontrado", http.StatusNotFound)
			return
		}

		var precisao any
		if in.Precisao != nil {
			precisao = *in.Precisao
		}
		var localNome any
		if local != "" {
			localNome = local
		}

		res, err := db.Exec(
			"INSERT INTO relatos (animal_id, latitude, longitude, precisao, local_nome) VALUES (?, ?, ?, ?, ?)",
			in.AnimalID, lat, lon, precisao, localNome,
		)
		if err != nil {
			http.Error(w, "Erro ao salvar relato", http.StatusInternalServerError)
			return
		}

		id, _ := res.LastInsertId()
		responderJSON(w, http.StatusCreated, map[string]any{"id": id, "ok": true})
	}))

	http.HandleFunc("/api/relatos", comCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		query := `
			SELECT r.id, r.animal_id, a.nome_popular, r.latitude, r.longitude, r.precisao, r.local_nome, r.data_hora
			FROM relatos r
			JOIN animais a ON a.id = r.animal_id`
		args := []any{}

		if v := r.URL.Query().Get("animal_id"); v != "" {
			animalID, err := strconv.Atoi(v)
			if err != nil || animalID <= 0 {
				http.Error(w, "animal_id inválido", http.StatusBadRequest)
				return
			}
			query += " WHERE r.animal_id = ?"
			args = append(args, animalID)
		}

		limite := 200
		if v := r.URL.Query().Get("limite"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil || n <= 0 || n > 1000 {
				http.Error(w, "limite inválido", http.StatusBadRequest)
				return
			}
			limite = n
		}
		query += " ORDER BY r.data_hora DESC, r.id DESC LIMIT ?"
		args = append(args, limite)

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Erro ao consultar relatos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		lista := []map[string]any{}
		for rows.Next() {
			var id, animalID int
			var nome string
			var lat, lon float64
			var precisao sql.NullFloat64
			var localNome, dataHora sql.NullString
			if err := rows.Scan(&id, &animalID, &nome, &lat, &lon, &precisao, &localNome, &dataHora); err != nil {
				continue
			}
			item := map[string]any{
				"id":           id,
				"animal_id":    animalID,
				"nome_popular": nome,
				"latitude":     lat,
				"longitude":    lon,
				"local_nome":   localNome.String,
				"data_hora":    dataHora.String,
				"precisao":     nil,
			}
			if precisao.Valid {
				item["precisao"] = precisao.Float64
			}
			lista = append(lista, item)
		}

		responderJSON(w, http.StatusOK, lista)
	}))
}
