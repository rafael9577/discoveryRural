// Package avistamento trata dos relatos de avistamento de animais.
package avistamento

// LimitePontos é o máximo de linhas retornadas nas consultas de pontos.
const LimitePontos = 1000

// Ponto é uma coordenada de avistamento.
type Ponto struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// Avistamento é um relato completo.
type Avistamento struct {
	ID        int      `json:"id"`
	AnimalID  int      `json:"animal_id"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Precisao  *float64 `json:"precisao"`
	LocalNome string   `json:"local_nome"`
	DataHora  string   `json:"data_hora"`
}

// Entrada é o corpo aceito no registro de um novo avistamento.
type Entrada struct {
	AnimalID  int      `json:"animal_id"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Precisao  *float64 `json:"precisao"`
	LocalNome string   `json:"local_nome"`
}

// Novo é um avistamento já validado, pronto para ser persistido.
type Novo struct {
	AnimalID  int
	Latitude  float64
	Longitude float64
	Precisao  *float64
	LocalNome string // vazio => NULL
}

// Criado é a resposta do registro.
type Criado struct {
	ID       int64 `json:"id"`
	AnimalID int   `json:"animal_id"`
}

// Grupo reúne os pontos [longitude, latitude] de um animal.
type Grupo struct {
	AnimalID     int          `json:"animal_id"`
	Avistamentos [][2]float64 `json:"avistamentos"`
}

// Linha é uma linha do LEFT JOIN animais x relatos (Ponto nil = sem relatos).
type Linha struct {
	AnimalID int
	Ponto    *Ponto
}
