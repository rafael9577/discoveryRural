// Package usuario trata de cadastro, login, sessões e reputação.
package usuario

// Usuario é a entidade exposta ao cliente (o CPF sai sempre mascarado).
type Usuario struct {
	ID          int64   `json:"id"`
	NomeUsuario string  `json:"nome_usuario"`
	CPF         string  `json:"cpf"`
	Email       string  `json:"email"`
	DataCriacao string  `json:"data_criacao"`
	Reputacao   float64 `json:"reputacao"`
}

// Cadastro é o corpo do registro.
type Cadastro struct {
	NomeUsuario string `json:"nome_usuario"`
	CPF         string `json:"cpf"`
	Email       string `json:"email"`
	Senha       string `json:"senha"`
}

// Credenciais é o corpo do login; Identificador é nome de usuário ou e-mail.
type Credenciais struct {
	Identificador string `json:"identificador"`
	Senha         string `json:"senha"`
}

// AjusteReputacao é o corpo (admin) para somar/subtrair reputação.
type AjusteReputacao struct {
	ID    int64   `json:"id"`
	Delta float64 `json:"delta"`
}

// Autenticado é a resposta de cadastro/login.
type Autenticado struct {
	Token   string   `json:"token"`
	Usuario *Usuario `json:"usuario"`
}

// NovoUsuario são os dados já validados e com senha em hash, prontos para gravar.
type NovoUsuario struct {
	NomeUsuario string
	CPF         string // somente dígitos
	Email       string
	SenhaHash   string
}
