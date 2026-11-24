// Package falha define o tipo de erro de domínio usado em toda a aplicação.
//
// Camadas internas (serviços e repositórios) devolvem *Falha descrevendo o que
// deu errado; a camada HTTP traduz o tipo para o status adequado. Assim, as
// regras de negócio não conhecem detalhes de HTTP e os handlers não precisam
// decidir status em cada if.
package falha

import (
	"errors"
	"net/http"
)

// Tipo classifica a falha.
type Tipo int

const (
	TipoValidacao Tipo = iota + 1
	TipoNaoAutorizado
	TipoNaoEncontrado
	TipoConflito
	TipoMuitasRequisicoes
	TipoInterno
)

// Falha é um erro com tipo, mensagem segura para o cliente e causa interna
// (que nunca é enviada ao cliente, apenas registrada em log).
type Falha struct {
	Tipo     Tipo
	Mensagem string
	Causa    error
}

func (f *Falha) Error() string {
	if f.Causa != nil {
		return f.Mensagem + ": " + f.Causa.Error()
	}
	return f.Mensagem
}

func (f *Falha) Unwrap() error { return f.Causa }

// StatusHTTP devolve o status HTTP correspondente ao tipo.
func (f *Falha) StatusHTTP() int {
	switch f.Tipo {
	case TipoValidacao:
		return http.StatusBadRequest
	case TipoNaoAutorizado:
		return http.StatusUnauthorized
	case TipoNaoEncontrado:
		return http.StatusNotFound
	case TipoConflito:
		return http.StatusConflict
	case TipoMuitasRequisicoes:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

func Validacao(msg string) *Falha         { return &Falha{Tipo: TipoValidacao, Mensagem: msg} }
func NaoAutorizado(msg string) *Falha     { return &Falha{Tipo: TipoNaoAutorizado, Mensagem: msg} }
func NaoEncontrado(msg string) *Falha     { return &Falha{Tipo: TipoNaoEncontrado, Mensagem: msg} }
func Conflito(msg string) *Falha          { return &Falha{Tipo: TipoConflito, Mensagem: msg} }
func MuitasRequisicoes(msg string) *Falha { return &Falha{Tipo: TipoMuitasRequisicoes, Mensagem: msg} }

// Interno cria uma falha 500. A mensagem vai para o cliente; a causa só para o log.
func Interno(msg string, causa error) *Falha {
	return &Falha{Tipo: TipoInterno, Mensagem: msg, Causa: causa}
}

// Converter garante um *Falha: erros desconhecidos viram falha interna genérica.
func Converter(err error) *Falha {
	var f *Falha
	if errors.As(err, &f) {
		return f
	}
	return Interno("Erro interno", err)
}

// EhTipo informa se err é uma Falha do tipo t.
func EhTipo(err error, t Tipo) bool {
	var f *Falha
	return errors.As(err, &f) && f.Tipo == t
}
