package main

import (
	"crypto/subtle"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func lerAnimais(db *sql.DB, idFiltro int) ([]map[string]any, error) {
	query := `SELECT id, nome_popular, nome_cientifico, breve_descricao, estado_conservacao, habito,
		vocalizacao_postura, peconhento, aviso_seguranca, protocolo_incidente, contatos_resgate,
		raridade, bitmap_id FROM animais`
	args := []any{}
	if idFiltro > 0 {
		query += " WHERE id = ?"
		args = append(args, idFiltro)
	}
	query += " ORDER BY id"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	animais := []map[string]any{}
	for rows.Next() {
		var id int
		var nomePopular, nomeCientifico, breve, estado, habito, vocal, aviso, protocolo, contatos sql.NullString
		var peconhento, raridade, bitmap sql.NullInt64
		if err := rows.Scan(&id, &nomePopular, &nomeCientifico, &breve, &estado, &habito,
			&vocal, &peconhento, &aviso, &protocolo, &contatos, &raridade, &bitmap); err != nil {
			continue
		}
		animais = append(animais, map[string]any{
			"id":                  id,
			"nome_popular":        nomePopular.String,
			"nome_cientifico":     nomeCientifico.String,
			"breve_descricao":     breve.String,
			"estado_conservacao":  estado.String,
			"habito":              habito.String,
			"vocalizacao_postura": vocal.String,
			"peconhento":          peconhento.Int64 == 1,
			"aviso_seguranca":     aviso.String,
			"protocolo_incidente": protocolo.String,
			"contatos_resgate":    contatos.String,
			"raridade":            int(raridade.Int64),
			"bitmap_id":           int(bitmap.Int64),
			"imagem_url":          "/api/animal/imagem?id=" + strconv.Itoa(id),
			"avistamentos":        []pontoAvistamento{},
		})
	}

	pontos := map[int][]pontoAvistamento{}
	pq := "SELECT animal_id, longitude, latitude FROM relatos"
	pargs := []any{}
	if idFiltro > 0 {
		pq += " WHERE animal_id = ?"
		pargs = append(pargs, idFiltro)
	}
	pq += " ORDER BY animal_id, id LIMIT ?"
	pargs = append(pargs, limiteAvistamentos)
	if prow, err := db.Query(pq, pargs...); err == nil {
		defer prow.Close()
		for prow.Next() {
			var animalID int
			var p pontoAvistamento
			if err := prow.Scan(&animalID, &p.Long, &p.Lat); err != nil {
				continue
			}
			pontos[animalID] = append(pontos[animalID], p)
		}
	}
	for _, a := range animais {
		if lista, ok := pontos[a["id"].(int)]; ok {
			a["avistamentos"] = lista
		}
	}
	return animais, nil
}

type edicaoAnimal struct {
	NomePopular        *string `json:"nome_popular"`
	NomeCientifico     *string `json:"nome_cientifico"`
	BreveDescricao     *string `json:"breve_descricao"`
	EstadoConservacao  *string `json:"estado_conservacao"`
	Habito             *string `json:"habito"`
	VocalizacaoPostura *string `json:"vocalizacao_postura"`
	Peconhento         *bool   `json:"peconhento"`
	AvisoSeguranca     *string `json:"aviso_seguranca"`
	ProtocoloIncidente *string `json:"protocolo_incidente"`
	ContatosResgate    *string `json:"contatos_resgate"`
	Raridade           *int    `json:"raridade"`
	BitmapID           *int    `json:"bitmap_id"`
}

func registrarAnimaisPadrao(adminKey string) {
	url := os.Getenv("TURSO_DATABASE_URL")
	token := os.Getenv("TURSO_AUTH_TOKEN")
	if url == "" || token == "" {
		log.Fatal("defina TURSO_DATABASE_URL e TURSO_AUTH_TOKEN (via .env ou variáveis de ambiente)")
	}

	db, err := sql.Open("libsql", url+"?authToken="+token)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/animais", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}
		animais, err := lerAnimais(db, 0)
		if err != nil {
			escreverErro(w, http.StatusInternalServerError, "Erro ao consultar animais")
			return
		}
		escreverJSON(w, http.StatusOK, animais)
	})

	http.HandleFunc("/api/animal/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch && r.Method != http.MethodPut {
			escreverErro(w, http.StatusMethodNotAllowed, "Método não permitido")
			return
		}
		chave := r.Header.Get("X-Admin-Key")
		if adminKey == "" || subtle.ConstantTimeCompare([]byte(chave), []byte(adminKey)) != 1 {
			escreverErro(w, http.StatusUnauthorized, "Não autorizado")
			return
		}

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id <= 0 {
			escreverErro(w, http.StatusBadRequest, "id inválido ou ausente")
			return
		}

		var in edicaoAnimal
		if !lerJSON(w, r, &in) {
			return
		}

		sets := []string{}
		args := []any{}

		textos := []struct {
			coluna      string
			valor       *string
			max         int
			obrigatorio bool
		}{
			{"nome_popular", in.NomePopular, 100, true},
			{"nome_cientifico", in.NomeCientifico, 100, false},
			{"breve_descricao", in.BreveDescricao, 500, false},
			{"estado_conservacao", in.EstadoConservacao, 100, false},
			{"habito", in.Habito, 100, false},
			{"vocalizacao_postura", in.VocalizacaoPostura, 500, false},
			{"aviso_seguranca", in.AvisoSeguranca, 500, false},
			{"protocolo_incidente", in.ProtocoloIncidente, 500, false},
			{"contatos_resgate", in.ContatosResgate, 300, false},
		}
		for _, c := range textos {
			if c.valor == nil {
				continue
			}
			t := strings.TrimSpace(*c.valor)
			if c.obrigatorio && t == "" {
				escreverErro(w, http.StatusBadRequest, c.coluna+" não pode ficar vazio")
				return
			}
			if utf8.RuneCountInString(t) > c.max {
				escreverErro(w, http.StatusBadRequest, c.coluna+" passa de "+strconv.Itoa(c.max)+" caracteres")
				return
			}
			sets = append(sets, c.coluna+" = ?")
			args = append(args, t)
		}
		if in.Peconhento != nil {
			v := 0
			if *in.Peconhento {
				v = 1
			}
			sets = append(sets, "peconhento = ?")
			args = append(args, v)
		}
		if in.Raridade != nil {
			if *in.Raridade < 1 || *in.Raridade > 5 {
				escreverErro(w, http.StatusBadRequest, "raridade deve ser de 1 a 5")
				return
			}
			sets = append(sets, "raridade = ?")
			args = append(args, *in.Raridade)
		}
		if in.BitmapID != nil {
			if *in.BitmapID < 0 {
				escreverErro(w, http.StatusBadRequest, "bitmap_id não pode ser negativo")
				return
			}
			sets = append(sets, "bitmap_id = ?")
			args = append(args, *in.BitmapID)
		}
		if len(sets) == 0 {
			escreverErro(w, http.StatusBadRequest, "nenhum campo para atualizar")
			return
		}

		var existe int
		if err := db.QueryRow("SELECT 1 FROM animais WHERE id = ?", id).Scan(&existe); err != nil {
			escreverErro(w, http.StatusNotFound, "Animal não encontrado")
			return
		}

		args = append(args, id)
		if _, err := db.Exec("UPDATE animais SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...); err != nil {
			escreverErro(w, http.StatusInternalServerError, "Erro ao atualizar animal")
			return
		}

		animais, err := lerAnimais(db, id)
		if err != nil || len(animais) == 0 {
			escreverErro(w, http.StatusInternalServerError, "Animal atualizado, mas não foi possível relê-lo")
			return
		}
		escreverJSON(w, http.StatusOK, animais[0])
	})
}
