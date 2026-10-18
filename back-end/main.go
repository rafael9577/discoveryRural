package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"strconv"

	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func processImage(img image.Image) []byte {
	bounds := img.Bounds()
	sw := bounds.Dx()
	sh := bounds.Dy()

	targetRatio := 3.0 / 4.0
	currentRatio := float64(sw) / float64(sh)

	var cw, ch, x0, y0 int
	if currentRatio > targetRatio {
		ch = sh
		cw = int(float64(sh) * targetRatio)
		x0 = (sw - cw) / 2
		y0 = 0
	} else {
		cw = sw
		ch = int(float64(sw) / targetRatio)
		x0 = 0
		y0 = (sh - ch) / 2
	}

	tw, th := 450, 600
	dst := image.NewRGBA(image.Rect(0, 0, tw, th))
	for y := 0; y < th; y++ {
		for x := 0; x < tw; x++ {
			srcX := x0 + (x * cw / tw)
			srcY := y0 + (y * ch / th)
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	var buf bytes.Buffer
	quality := 85
	for {
		buf.Reset()
		jpeg.Encode(&buf, dst, &jpeg.Options{Quality: quality})
		if buf.Len() <= 16*1024 || quality <= 10 {
			break
		}
		quality -= 5
	}

	return buf.Bytes()
}

func createPlaceholderImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 600, 800))
	for y := 0; y < 800; y++ {
		for x := 0; x < 600; x++ {
			img.Set(x, y, color.RGBA{40, 120, 80, 255})
		}
	}
	return processImage(img)
}

func main() {
	url := os.Getenv("TURSO_DATABASE_URL")
	token := os.Getenv("TURSO_AUTH_TOKEN")
	if url == "" || token == "" {
		log.Fatal("defina TURSO_DATABASE_URL e TURSO_AUTH_TOKEN")
	}

	db, err := sql.Open("libsql", url+"?authToken="+token)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS animais (
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
	);

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
	`)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM animais").Scan(&count)
	if count == 0 {
		imgBytes := createPlaceholderImage()
		stmt, err := db.Prepare(`INSERT INTO animais (nome_popular, nome_cientifico, imagem, breve_descricao, estado_conservacao, habito, vocalizacao_postura, peconhento, aviso_seguranca, protocolo_incidente, contatos_resgate, raridade, bitmap_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
		if err == nil {
			stmt.Exec("Coruja-buraqueira", "Athene cunicularia", imgBytes, "Pequena coruja de hábitos terrestres.", "Pouco preocupante", "Diurno/Crepuscular", "Emite chamados agudos e balança a cabeça.", 0, "Nunca tente manipular, alimentar ou se aproximar excessivamente de ninhos no chão.", "Manter distância segura e isolar a área.", "193 / Setor Ambiental", 2, 18)
			stmt.Exec("Cágado-de-barbicha", "Phrynops geoffroanus", imgBytes, "Réptil aquático encontrado em lagos e açudes.", "Pouco preocupante", "Diurno", "Silencioso na água; emite som soprado quando acuado.", 0, "Nunca tente manipular, alimentar ou remover o animal.", "Evitar contato direto e direcioná-lo com cuidado para a água se estiver em risco.", "Setor de Meio Ambiente / Biologia", 3, 384)
			stmt.Close()
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlPage))
	})

	http.HandleFunc("/api/animais/resumo", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, nome_popular FROM animais")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		lista := []map[string]any{}
		for rows.Next() {
			var id int
			var nomePopular string
			rows.Scan(&id, &nomePopular)
			lista = append(lista, map[string]any{
				"id":           id,
				"nome_popular": nomePopular,
				"imagem_url":   "/api/animal/imagem?id=" + strconv.Itoa(id),
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lista)
	})

	http.HandleFunc("/api/animal", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID obrigatório", http.StatusBadRequest)
			return
		}

		a := make(map[string]any)
		var idVal, raridade, peconhento, bitmapID int
		var nomePopular, nomeCientifico, breveDescricao, estadoConservacao, habito, vocalizacaoPostura, avisoSeguranca, protocoloIncidente, contatosResgate sql.NullString

		err := db.QueryRow("SELECT id, nome_popular, nome_cientifico, breve_descricao, estado_conservacao, habito, vocalizacao_postura, peconhento, aviso_seguranca, protocolo_incidente, contatos_resgate, raridade, bitmap_id FROM animais WHERE id = ?", id).
			Scan(&idVal, &nomePopular, &nomeCientifico, &breveDescricao, &estadoConservacao, &habito, &vocalizacaoPostura, &peconhento, &avisoSeguranca, &protocoloIncidente, &contatosResgate, &raridade, &bitmapID)

		if err != nil {
			http.Error(w, "Animal não encontrado", http.StatusNotFound)
			return
		}

		a["id"] = idVal
		a["nome_popular"] = nomePopular.String
		a["nome_cientifico"] = nomeCientifico.String
		a["breve_descricao"] = breveDescricao.String
		a["estado_conservacao"] = estadoConservacao.String
		a["habito"] = habito.String
		a["vocalizacao_postura"] = vocalizacaoPostura.String
		a["peconhento"] = peconhento == 1
		a["aviso_seguranca"] = avisoSeguranca.String
		a["protocolo_incidente"] = protocoloIncidente.String
		a["contatos_resgate"] = contatosResgate.String
		a["raridade"] = raridade
		a["bitmap_id"] = bitmapID
		a["imagem_url"] = "/api/animal/imagem?id=" + id

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(a)
	})

	http.HandleFunc("/api/animal/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		r.ParseMultipartForm(10 << 20)

		nomePopular := r.FormValue("nome_popular")
		nomeCientifico := r.FormValue("nome_cientifico")
		breveDescricao := r.FormValue("breve_descricao")
		estadoConservacao := r.FormValue("estado_conservacao")
		habito := r.FormValue("habito")
		vocalizacaoPostura := r.FormValue("vocalizacao_postura")

		peconhento := 0
		if r.FormValue("peconhento") == "on" {
			peconhento = 1
		}

		avisoSeguranca := r.FormValue("aviso_seguranca")
		protocoloIncidente := r.FormValue("protocolo_incidente")
		contatosResgate := r.FormValue("contatos_resgate")
		raridade := r.FormValue("raridade")
		bitmapID := r.FormValue("bitmap_id")

		var imgBytes []byte
		file, _, err := r.FormFile("imagem")
		if err == nil {
			defer file.Close()
			fileBytes, err := io.ReadAll(file)
			if err == nil {
				img, _, err := image.Decode(bytes.NewReader(fileBytes))
				if err == nil {
					imgBytes = processImage(img)
				}
			}
		}

		if len(imgBytes) == 0 {
			imgBytes = createPlaceholderImage()
		}

		_, err = db.Exec(`
			INSERT INTO animais (
				nome_popular, nome_cientifico, imagem, breve_descricao, 
				estado_conservacao, habito, vocalizacao_postura, peconhento, 
				aviso_seguranca, protocolo_incidente, contatos_resgate, raridade, bitmap_id
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, nomePopular, nomeCientifico, imgBytes, breveDescricao, estadoConservacao, habito, vocalizacaoPostura, peconhento, avisoSeguranca, protocoloIncidente, contatosResgate, raridade, bitmapID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/api/animal/imagem", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		var imgBytes []byte
		err := db.QueryRow("SELECT imagem FROM animais WHERE id = ?", id).Scan(&imgBytes)
		if err != nil || len(imgBytes) == 0 {
			http.Error(w, "Imagem não encontrada", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(imgBytes)
	})

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

const htmlPage = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <title>Fauna Local - UFRPE</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 800px; margin: 40px auto; padding: 0 20px; background: #f8f9fa; color: #212529; }
        .card { background: white; padding: 20px; margin-bottom: 15px; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border: 1px solid #dee2e6; display: flex; gap: 20px; align-items: flex-start; cursor: pointer; }
        .card img { width: 120px; height: 160px; object-fit: cover; border-radius: 6px; border: 1px solid #dee2e6; flex-shrink: 0; }
        .info { flex-grow: 1; }
        h2 { margin-top: 0; color: #198754; }
        .cientifico { font-style: italic; color: #6c757d; font-weight: normal; font-size: 0.9em; }
        .badge { display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 12px; font-weight: 600; background: #e9ecef; color: #495057; margin-right: 5px; margin-bottom: 8px; }
        .warning { background: #fff3cd; border: 1px solid #ffeeba; color: #856404; padding: 12px; border-radius: 6px; margin-top: 12px; font-size: 14px; }
        .form-box { background: white; padding: 20px; margin-bottom: 25px; border-radius: 8px; border: 1px solid #dee2e6; }
        .form-group { margin-bottom: 12px; display: flex; flex-direction: column; }
        .form-group label { font-weight: 600; margin-bottom: 4px; font-size: 14px; display: flex; justify-content: space-between; align-items: center; }
        .form-group input, .form-group textarea, .form-group select { padding: 8px; border: 1px solid #ced4da; border-radius: 4px; font-size: 14px; }
        .checkbox-group { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; }
        p { line-height: 1.5; margin: 8px 0; }
        button { padding: 10px 16px; background: #198754; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: 600; }
        button:hover { background: #157347; }
        .help-link { color: #0d6efd; cursor: pointer; font-size: 12px; font-weight: normal; text-decoration: underline; }
        
        .modal { display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.5); justify-content: center; align-items: center; }
        .modal-content { background: white; padding: 25px; border-radius: 8px; width: 450px; max-width: 90%; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px; }
        .modal-header h3 { margin: 0; color: #198754; }
        .close-btn { background: none; border: none; font-size: 18px; cursor: pointer; color: #6c757d; padding: 0; }
        .modal-body table { width: 100%; border-collapse: collapse; font-size: 13px; }
        .modal-body th, .modal-body td { border: 1px solid #dee2e6; padding: 6px 8px; text-align: left; }
        .modal-body th { background: #f8f9fa; }
    </style>
</head>
<body>
    <h1>Fauna Local - UFRPE</h1>

    <div class="form-box">
        <h3>Cadastrar Nova Espécie</h3>
        <form action="/api/animal/create" method="POST" enctype="multipart/form-data">
            <div class="form-group">
                <label>Nome Popular:</label>
                <input type="text" name="nome_popular" required>
            </div>
            <div class="form-group">
                <label>Nome Científico:</label>
                <input type="text" name="nome_cientifico">
            </div>
            <div class="form-group">
                <label>Breve Descrição:</label>
                <textarea name="breve_descricao" rows="2" required></textarea>
            </div>
            <div class="form-group">
                <label>Estado de Conservação:</label>
                <input type="text" name="estado_conservacao" placeholder="Ex: Pouco preocupante">
            </div>
            <div class="form-group">
                <label>Hábito:</label>
                <input type="text" name="habito" placeholder="Ex: Diurno / Noturno">
            </div>
            <div class="form-group">
                <label>Vocalização e Postura:</label>
                <textarea name="vocalizacao_postura" rows="2"></textarea>
            </div>
            <div class="checkbox-group">
                <input type="checkbox" id="peconhento" name="peconhento">
                <label for="peconhento" style="margin: 0; font-weight: 600;">Animal Peçonhento / Venenoso</label>
            </div>
            <div class="form-group">
                <label>Aviso de Segurança:</label>
                <textarea name="aviso_seguranca" rows="2"></textarea>
            </div>
            <div class="form-group">
                <label>Protocolo de Incidente:</label>
                <textarea name="protocolo_incidente" rows="2"></textarea>
            </div>
            <div class="form-group">
                <label>Contatos para Resgate:</label>
                <input type="text" name="contatos_resgate">
            </div>
            <div class="form-group">
                <label>Raridade (1 a 5):</label>
                <select name="raridade">
                    <option value="1">1 - Muito Comum</option>
                    <option value="2">2 - Comum</option>
                    <option value="3" selected>3 - Moderado</option>
                    <option value="4">4 - Raro</option>
                    <option value="5">5 - Muito Raro</option>
                </select>
            </div>
            <div class="form-group">
                <label>
                    Bitmap ID (Soma dos bits):
                    <span class="help-link" onclick="abrirModal()">O que é isso? Como calcular?</span>
                </label>
                <input type="number" name="bitmap_id" value="0" required>
            </div>
            <div class="form-group">
                <label>Foto do Animal:</label>
                <input type="file" name="imagem" accept="image/*" required>
            </div>
            <button type="submit">Salvar Nova Espécie</button>
        </form>
    </div>

    <div id="bitmapModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h3>Como calcular o Bitmap ID</h3>
                <button class="close-btn" onclick="fecharModal()">&times;</button>
            </div>
            <div class="modal-body">
                <p>O Bitmap ID agrupa as características do animal em uma única soma matemática baseada em potências de 2.</p>
                <p><strong>Somente some os valores correspondentes ao animal:</strong></p>
                <table>
                    <tr><th>Característica</th><th>Valor (Bit)</th></tr>
                    <tr><td>No alto (Árvores/Telhados)</td><td>1</td></tr>
                    <tr><td>No chão (Grama/Asfalto)</td><td>2</td></tr>
                    <tr><td>Na água ou margem</td><td>4</td></tr>
                    <tr><td>É um Pássaro / Ave</td><td>8</td></tr>
                    <tr><td>É uma Coruja</td><td>16</td></tr>
                    <tr><td>É um Cágado / Tartaruga</td><td>32</td></tr>
                    <tr><td>É um Lagarto</td><td>64</td></tr>
                    <tr><td>É uma Cobra / Serpente</td><td>128</td></tr>
                    <tr><td>É um Mamífero</td><td>256</td></tr>
                </table>
                <p style="margin-top: 10px; font-size: 12px; color: #6c757d;"><em>Exemplo: Se for um pássaro (8) que vive no chão (2), o Bitmap ID será 8 + 2 = 10.</em></p>
            </div>
        </div>
    </div>

    <div id="conteudo">Carregando...</div>
    <script>
        function abrirModal() {
            document.getElementById('bitmapModal').style.display = 'flex';
        }

        function fecharModal() {
            document.getElementById('bitmapModal').style.display = 'none';
        }

        window.onclick = function(event) {
            let modal = document.getElementById('bitmapModal');
            if (event.target == modal) {
                modal.style.display = 'none';
            }
        }

        function carregarLista() {
            fetch('/api/animais/resumo')
                .then(res => res.json())
                .then(data => {
                    const div = document.getElementById('conteudo');
                    if (!data || data.length === 0) {
                        div.innerHTML = '<p>Nenhum animal cadastrado.</p>';
                        return;
                    }
                    div.innerHTML = '<h3>Lista de Espécies Cadastradas</h3>' + data.map(a => ` + "`" + `
                        <div class="card" onclick="carregarDetalhe(${a.id})">
                            <img src="${a.imagem_url}" alt="${a.nome_popular}">
                            <div class="info">
                                <h2>ID ${a.id} -${a.nome_popular}</h2>
                                <button>Ver Detalhes</button>
                            </div>
                        </div>
                    ` + "`" + `).join('');
                });
        }

        function carregarDetalhe(id) {
            fetch('/api/animal?id=' + id)
                .then(res => res.json())
                .then(a => {
                    const div = document.getElementById('conteudo');
                    div.innerHTML = ` + "`" + `
                        <button onclick="carregarLista()" style="margin-bottom: 15px;">← Voltar</button>
                        <div class="card" style="cursor: default;">
                            <img src="${a.imagem_url}" alt="${a.nome_popular}">
                            <div class="info">
                                <h2>${a.nome_popular} <span class="cientifico">(${a.nome_cientifico})</span></h2>
                                <p>${a.breve_descricao}</p>
                                <div>
                                    <span class="badge">Hábito: ${a.habito}</span>
                                    <span class="badge">Conservação: ${a.estado_conservacao}</span>
                                    <span class="badge">Raridade: ${a.raridade}/5</span>
                                    <span class="badge">Bitmap ID: ${a.bitmap_id}</span>${a.peconhento ? '<span class="badge" style="background:#dc3545; color:white;">Peçonhento</span>' : ''}
                                </div>
                                <p><strong>Vocalização / Postura:</strong> ${a.vocalizacao_postura}</p>
                                <div class="warning">
                                    <strong>⚠️ Aviso de Segurança:</strong> ${a.aviso_seguranca}<br>
                                    <strong>Protocolo:</strong> ${a.protocolo_incidente}<br>
                                    <strong>Resgate:</strong> ${a.contatos_resgate}
                                </div>
                            </div>
                        </div>
                    ` + "`" + `;
                });
        }

        carregarLista();
    </script>
</body>
</html>`
