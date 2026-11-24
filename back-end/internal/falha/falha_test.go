package falha

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestStatusHTTPPorTipo(t *testing.T) {
	casos := map[*Falha]int{
		Validacao("x"):                 http.StatusBadRequest,
		NaoAutorizado("x"):             http.StatusUnauthorized,
		NaoEncontrado("x"):             http.StatusNotFound,
		Conflito("x"):                  http.StatusConflict,
		MuitasRequisicoes("x"):         http.StatusTooManyRequests,
		Interno("x", errors.New("db")): http.StatusInternalServerError,
	}
	for f, esperado := range casos {
		if got := f.StatusHTTP(); got != esperado {
			t.Errorf("tipo %d: status %d, esperado %d", f.Tipo, got, esperado)
		}
	}
}

func TestConverterErroDesconhecidoViraInternoGenerico(t *testing.T) {
	f := Converter(errors.New("senha do banco vazou aqui"))
	if f.Tipo != TipoInterno || f.Mensagem != "Erro interno" {
		t.Fatalf("esperava falha interna genérica, veio %+v", f)
	}
}

func TestConverterPreservaFalhaEmbrulhada(t *testing.T) {
	embrulhado := fmt.Errorf("contexto: %w", NaoEncontrado("Animal não encontrado"))
	if !EhTipo(embrulhado, TipoNaoEncontrado) {
		t.Fatal("deveria reconhecer a falha embrulhada")
	}
}
