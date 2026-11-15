package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func apiURL() string {
	if v := os.Getenv("API_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}
	return "http://localhost:8080"
}

var httpTeste = &http.Client{Timeout: 90 * time.Second}

type respTeste struct {
	Status int
	Tipo   string // cabeçalho Content-Type
	Corpo  []byte
}

func (r respTeste) Texto() string { return string(r.Corpo) }

func chamarAPI(t *testing.T, metodo, caminho string, corpo []byte, cabecalhos map[string]string) respTeste {
	t.Helper()
	var leitor io.Reader
	if corpo != nil {
		leitor = bytes.NewReader(corpo)
	}
	req, err := http.NewRequest(metodo, apiURL()+caminho, leitor)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range cabecalhos {
		req.Header.Set(k, v)
	}
	res, err := httpTeste.Do(req)
	if err != nil {
		t.Skipf("servidor indisponível em %s: %v", apiURL(), err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return respTeste{Status: res.StatusCode, Tipo: res.Header.Get("Content-Type"), Corpo: b}
}

func apiGet(t *testing.T, caminho string) respTeste {
	t.Helper()
	return chamarAPI(t, http.MethodGet, caminho, nil, nil)
}

func ipFicticio() string {
	return fmt.Sprintf("198.51.%d.%d", rand.Intn(250)+1, rand.Intn(250)+1)
}

func postRelato(t *testing.T, corpo string, ip string) respTeste {
	t.Helper()
	return chamarAPI(t, http.MethodPost, "/api/relato", []byte(corpo), map[string]string{
		"Content-Type":    "application/json",
		"X-Forwarded-For": ip,
	})
}
