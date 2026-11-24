// Package animal trata do cadastro de espécies (dados, imagem e edição).
package animal

import "discoveryrural/internal/avistamento"

// Animal é a entidade de domínio.
type Animal struct {
	ID                 int    `json:"id"`
	NomePopular        string `json:"nome_popular"`
	NomeCientifico     string `json:"nome_cientifico"`
	BreveDescricao     string `json:"breve_descricao"`
	EstadoConservacao  string `json:"estado_conservacao"`
	Habito             string `json:"habito"`
	VocalizacaoPostura string `json:"vocalizacao_postura"`
	Peconhento         bool   `json:"peconhento"`
	AvisoSeguranca     string `json:"aviso_seguranca"`
	ProtocoloIncidente string `json:"protocolo_incidente"`
	ContatosResgate    string `json:"contatos_resgate"`
	Raridade           int    `json:"raridade"`
	BitmapID           int    `json:"bitmap_id"`
	ImagemURL          string `json:"imagem_url"`
}

// AnimalComAvistamentos é o Animal acrescido dos pontos [lon, lat] de avistamento.
type AnimalComAvistamentos struct {
	Animal
	Avistamentos []avistamento.Ponto `json:"avistamentos"`
}

// Resumo é a versão enxuta usada em listagens.
type Resumo struct {
	ID             int    `json:"id"`
	NomePopular    string `json:"nome_popular"`
	BreveDescricao string `json:"breve_descricao"`
	ImagemURL      string `json:"imagem_url"`
}

// EdicaoAnimal é o corpo da atualização parcial: só os campos enviados mudam.
type EdicaoAnimal struct {
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

// NovoAnimal reúne os campos do cadastro (vindos do formulário multipart).
type NovoAnimal struct {
	NomePopular        string
	NomeCientifico     string
	BreveDescricao     string
	EstadoConservacao  string
	Habito             string
	VocalizacaoPostura string
	Peconhento         bool
	AvisoSeguranca     string
	ProtocoloIncidente string
	ContatosResgate    string
	Raridade           int
	BitmapID           int
}

// Alteracao é um par coluna/valor a ser gravado.
type Alteracao struct {
	Coluna string
	Valor  any
}
