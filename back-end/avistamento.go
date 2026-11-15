package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

const limiteAvistamentos = 1000

type entradaAvistamento struct {
	AnimalID  int      `json:"animal_id"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Precisao  *float64 `json:"precisao"`
	LocalNome string   `json:"local_nome"`
}

type grupoAvistamentos struct {
	AnimalID     int          `json:"animal_id"`
	Avistamentos [][2]float64 `json:"avistamentos"`
}

type pontoAvistamento struct {
	Long float64 `json:"longitude"`
	Lat  float64 `json:"latitude"`
}

func registrarAvistamentos() {
	url := os.Getenv("TURSO_DATABASE_URL")
	token := os.Getenv("TURSO_AUTH_TOKEN")
	if url == "" || token == "" {
		log.Fatal("defina TURSO_DATABASE_URL e TURSO_AUTH_TOKEN (via .env ou variáveis de ambiente)")
	}

	db, err := sql.Open("libsql", url+"?authToken="+token)
	if err != nil {
		log.Fatal(err)
	}

	lim := novoLimitador(1, 10*time.Second)


	http.HandleFunc("/api/animais/ids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}
		rows, err := db.Query("SELECT id FROM animais ORDER BY id")
		if err != nil {
			escreverErro(w, http.StatusInternalServerError, "Erro ao consultar animais")
			return
		}
		defer rows.Close()

		ids := []int{}
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				continue
			}
			ids = append(ids, id)
		}
		escreverJSON(w, http.StatusOK, ids)
	})

	http.HandleFunc("/api/avistamento", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			buscarAvistamento(db, w, r)
		case http.MethodPost:
			inserirAvistamento(db, lim, w, r)
		default:
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
		}
	})

	http.HandleFunc("/api/avistamentos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}

		query := `
			SELECT a.id, r.longitude, r.latitude
			FROM animais a
			LEFT JOIN relatos r ON r.animal_id = a.id`
		args := []any{}
		filtrado := false

		if v := r.URL.Query().Get("animal_id"); v != "" {
			animalID, err := strconv.Atoi(v)
			if err != nil || animalID <= 0 {
				escreverErro(w, http.StatusBadRequest, "animal_id inválido")
				return
			}
			query += " WHERE a.id = ?"
			args = append(args, animalID)
			filtrado = true
		}
		query += " ORDER BY a.id, r.id LIMIT ?"
		args = append(args, limiteAvistamentos)

		rows, err := db.Query(query, args...)
		if err != nil {
			escreverErro(w, http.StatusInternalServerError, "Erro ao consultar avistamentos")
			return
		}
		defer rows.Close()

		grupos := []grupoAvistamentos{}
		indice := map[int]int{}
		for rows.Next() {
			var animalID int
			var lon, lat sql.NullFloat64
			if err := rows.Scan(&animalID, &lon, &lat); err != nil {
				continue
			}
			i, ok := indice[animalID]
			if !ok {
				grupos = append(grupos, grupoAvistamentos{AnimalID: animalID, Avistamentos: [][2]float64{}})
				i = len(grupos) - 1
				indice[animalID] = i
			}
			if lon.Valid && lat.Valid {
				grupos[i].Avistamentos = append(grupos[i].Avistamentos, [2]float64{lon.Float64, lat.Float64})
			}
		}

		if filtrado && len(grupos) == 0 {
			escreverErro(w, http.StatusNotFound, "Animal não encontrado")
			return
		}
		escreverJSON(w, http.StatusOK, grupos)
	})
}

func buscarAvistamento(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		escreverErro(w, http.StatusBadRequest, "id inválido ou ausente")
		return
	}

	var animalID int
	var lat, lon float64
	var precisao sql.NullFloat64
	var local, dataHora sql.NullString
	err = db.QueryRow(
		"SELECT animal_id, latitude, longitude, precisao, local_nome, data_hora FROM relatos WHERE id = ?",
		id,
	).Scan(&animalID, &lat, &lon, &precisao, &local, &dataHora)
	if err == sql.ErrNoRows {
		escreverErro(w, http.StatusNotFound, "Avistamento não encontrado")
		return
	}
	if err != nil {
		escreverErro(w, http.StatusInternalServerError, "Erro ao consultar avistamento")
		return
	}

	item := map[string]any{
		"id":         id,
		"animal_id":  animalID,
		"latitude":   lat,
		"longitude":  lon,
		"precisao":   nil,
		"local_nome": local.String,
		"data_hora":  dataHora.String,
	}
	if precisao.Valid {
		item["precisao"] = precisao.Float64
	}
	escreverJSON(w, http.StatusOK, item)
}

func inserirAvistamento(db *sql.DB, lim *limitador, w http.ResponseWriter, r *http.Request) {
	if !lim.permitir(ipDe(r)) {
		escreverErro(w, http.StatusTooManyRequests, "Muitos avistamentos seguidos, aguarde alguns segundos")
		return
	}

	var in entradaAvistamento
	if !lerJSON(w, r, &in) {
		return
	}

	if in.AnimalID <= 0 {
		escreverErro(w, http.StatusBadRequest, "animal_id obrigatório")
		return
	}
	if in.Latitude == nil || in.Longitude == nil {
		escreverErro(w, http.StatusBadRequest, "latitude e longitude obrigatórias")
		return
	}
	lat, lon := *in.Latitude, *in.Longitude
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		escreverErro(w, http.StatusBadRequest, "coordenadas fora do intervalo válido")
		return
	}
	if in.Precisao != nil && (*in.Precisao < 0 || *in.Precisao > 100000) {
		escreverErro(w, http.StatusBadRequest, "precisão inválida")
		return
	}

	local := strings.TrimSpace(in.LocalNome)
	if len([]rune(local)) > 120 {
		local = string([]rune(local)[:120])
	}

	var existe int
	if err := db.QueryRow("SELECT 1 FROM animais WHERE id = ?", in.AnimalID).Scan(&existe); err != nil {
		escreverErro(w, http.StatusNotFound, "Animal não encontrado")
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
		escreverErro(w, http.StatusInternalServerError, "Erro ao salvar avistamento")
		return
	}

	id, _ := res.LastInsertId()
	escreverJSON(w, http.StatusCreated, map[string]any{"id": id, "animal_id": in.AnimalID})
}
