package animal

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"discoveryrural/internal/falha"
)

type campoTexto struct {
	coluna      string
	valor       *string
	max         int
	obrigatorio bool
}

// campos lista os campos de texto na ordem em que são validados.
func (e EdicaoAnimal) campos() []campoTexto {
	return []campoTexto{
		{"nome_popular", e.NomePopular, 100, true},
		{"nome_cientifico", e.NomeCientifico, 100, false},
		{"breve_descricao", e.BreveDescricao, 500, false},
		{"estado_conservacao", e.EstadoConservacao, 100, false},
		{"habito", e.Habito, 100, false},
		{"vocalizacao_postura", e.VocalizacaoPostura, 500, false},
		{"aviso_seguranca", e.AvisoSeguranca, 500, false},
		{"protocolo_incidente", e.ProtocoloIncidente, 500, false},
		{"contatos_resgate", e.ContatosResgate, 300, false},
	}
}

// Normalizar aplica trim nos textos e valida limites, devolvendo uma cópia
// pronta para ser gravada. A primeira regra violada interrompe a validação.
func (e EdicaoAnimal) Normalizar() (EdicaoAnimal, error) {
	campos := e.campos()
	destinos := []**string{
		&e.NomePopular, &e.NomeCientifico, &e.BreveDescricao, &e.EstadoConservacao, &e.Habito,
		&e.VocalizacaoPostura, &e.AvisoSeguranca, &e.ProtocoloIncidente, &e.ContatosResgate,
	}
	for i, c := range campos {
		if c.valor == nil {
			continue
		}
		t := strings.TrimSpace(*c.valor)
		if c.obrigatorio && t == "" {
			return e, falha.Validacao(c.coluna + " não pode ficar vazio")
		}
		if utf8.RuneCountInString(t) > c.max {
			return e, falha.Validacao(c.coluna + " passa de " + strconv.Itoa(c.max) + " caracteres")
		}
		*destinos[i] = &t
	}
	if e.Raridade != nil && (*e.Raridade < 1 || *e.Raridade > 5) {
		return e, falha.Validacao("raridade deve ser de 1 a 5")
	}
	if e.BitmapID != nil && *e.BitmapID < 0 {
		return e, falha.Validacao("bitmap_id não pode ser negativo")
	}
	return e, nil
}

// Alteracoes converte os campos presentes em pares coluna/valor.
func (e EdicaoAnimal) Alteracoes() []Alteracao {
	var alt []Alteracao
	for _, c := range e.campos() {
		if c.valor != nil {
			alt = append(alt, Alteracao{c.coluna, *c.valor})
		}
	}
	if e.Peconhento != nil {
		v := 0
		if *e.Peconhento {
			v = 1
		}
		alt = append(alt, Alteracao{"peconhento", v})
	}
	if e.Raridade != nil {
		alt = append(alt, Alteracao{"raridade", *e.Raridade})
	}
	if e.BitmapID != nil {
		alt = append(alt, Alteracao{"bitmap_id", *e.BitmapID})
	}
	return alt
}

// paraEdicao permite reutilizar a mesma validação no cadastro.
func (n NovoAnimal) paraEdicao() EdicaoAnimal {
	return EdicaoAnimal{
		NomePopular:        &n.NomePopular,
		NomeCientifico:     &n.NomeCientifico,
		BreveDescricao:     &n.BreveDescricao,
		EstadoConservacao:  &n.EstadoConservacao,
		Habito:             &n.Habito,
		VocalizacaoPostura: &n.VocalizacaoPostura,
		Peconhento:         &n.Peconhento,
		AvisoSeguranca:     &n.AvisoSeguranca,
		ProtocoloIncidente: &n.ProtocoloIncidente,
		ContatosResgate:    &n.ContatosResgate,
		Raridade:           &n.Raridade,
		BitmapID:           &n.BitmapID,
	}
}

// paraAnimal materializa a edição (campos ausentes ficam com valor zero).
func (e EdicaoAnimal) paraAnimal() Animal {
	s := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}
	a := Animal{
		NomePopular:        s(e.NomePopular),
		NomeCientifico:     s(e.NomeCientifico),
		BreveDescricao:     s(e.BreveDescricao),
		EstadoConservacao:  s(e.EstadoConservacao),
		Habito:             s(e.Habito),
		VocalizacaoPostura: s(e.VocalizacaoPostura),
		AvisoSeguranca:     s(e.AvisoSeguranca),
		ProtocoloIncidente: s(e.ProtocoloIncidente),
		ContatosResgate:    s(e.ContatosResgate),
	}
	if e.Peconhento != nil {
		a.Peconhento = *e.Peconhento
	}
	if e.Raridade != nil {
		a.Raridade = *e.Raridade
	}
	if e.BitmapID != nil {
		a.BitmapID = *e.BitmapID
	}
	return a
}
