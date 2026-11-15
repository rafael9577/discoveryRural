package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type itemResumo struct {
	ID          int    `json:"id"`
	NomePopular string `json:"nome_popular"`
	ImagemURL   string `json:"imagem_url"`
}

func primeiraEspecie(t *testing.T) itemResumo {
	t.Helper()
	r := apiGet(t, "/api/animais/resumo")
	var lista []itemResumo
	if err := json.Unmarshal(r.Corpo, &lista); err != nil {
		t.Fatalf("a lista deveria ser um JSON de array: %v (corpo: %q)", err, r.Texto())
	}
	if len(lista) == 0 {
		t.Skip("catálogo vazio: cadastre ao menos uma espécie para executar este teste")
	}
	return lista[0]
}

func TestRF01_ListaDeEspecies(t *testing.T) {
	r := apiGet(t, "/api/animais/resumo")
	if r.Status != http.StatusOK {
		t.Fatalf("status = %d, esperado 200", r.Status)
	}
	if !strings.HasPrefix(r.Tipo, "application/json") {
		t.Errorf("Content-Type = %q, esperado application/json", r.Tipo)
	}
	var lista []itemResumo
	if err := json.Unmarshal(r.Corpo, &lista); err != nil {
		t.Fatalf("a resposta deveria ser um JSON de array: %v", err)
	}
	if len(lista) == 0 {
		t.Skip("catálogo vazio: cadastre ao menos uma espécie para executar este teste")
	}
	for _, e := range lista {
		if e.ID <= 0 {
			t.Errorf("espécie com id inválido: %+v", e)
		}
		if strings.TrimSpace(e.NomePopular) == "" {
			t.Errorf("a espécie %d está sem nome", e.ID)
		}
		if !strings.HasPrefix(e.ImagemURL, "/api/animal/imagem?id=") {
			t.Errorf("a espécie %d tem endereço de imagem inesperado: %q", e.ID, e.ImagemURL)
		}
	}
}

func TestRF01_DetalheDeEspecie(t *testing.T) {
	esp := primeiraEspecie(t)
	r := apiGet(t, fmt.Sprintf("/api/animal?id=%d", esp.ID))
	if r.Status != http.StatusOK {
		t.Fatalf("status = %d, esperado 200", r.Status)
	}
	if !strings.HasPrefix(r.Tipo, "application/json") {
		t.Errorf("Content-Type = %q, esperado application/json", r.Tipo)
	}
	var d struct {
		ID          int    `json:"id"`
		NomePopular string `json:"nome_popular"`
		ImagemURL   string `json:"imagem_url"`
	}
	if err := json.Unmarshal(r.Corpo, &d); err != nil {
		t.Fatalf("a resposta deveria ser um JSON de objeto: %v", err)
	}
	if d.ID != esp.ID {
		t.Errorf("id = %d, esperado %d", d.ID, esp.ID)
	}
	if strings.TrimSpace(d.NomePopular) == "" {
		t.Error("o detalhe deveria trazer o nome da espécie")
	}
	if !strings.HasPrefix(d.ImagemURL, "/api/animal/imagem?id=") {
		t.Errorf("endereço de imagem inesperado: %q", d.ImagemURL)
	}
}
