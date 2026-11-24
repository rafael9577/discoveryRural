// Package httpx reúne utilitários HTTP reutilizáveis: respostas JSON, leitura
// de corpo, middlewares, autenticação de administrador e extração de IP.
package httpx

import (
	"encoding/json"
	"log"
	"net/http"

	"discoveryrural/internal/falha"
)

// EscreverJSON serializa v como JSON com o status informado.
func EscreverJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("erro ao escrever resposta JSON: %v", err)
	}
}

// EscreverErro escreve o formato padrão de erro: {"erro": "mensagem"}.
func EscreverErro(w http.ResponseWriter, status int, msg string) {
	EscreverJSON(w, status, map[string]string{"erro": msg})
}

// EscreverFalha traduz qualquer erro para a resposta HTTP correta. A causa
// interna de falhas 500 é registrada em log e nunca enviada ao cliente.
func EscreverFalha(w http.ResponseWriter, err error) {
	f := falha.Converter(err)
	if f.Tipo == falha.TipoInterno {
		log.Printf("erro interno: %v", f)
	}
	EscreverErro(w, f.StatusHTTP(), f.Mensagem)
}
