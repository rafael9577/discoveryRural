package usuario

import (
	"context"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"discoveryrural/internal/falha"
)

// DuracaoSessao é por quanto tempo um token permanece válido.
const DuracaoSessao = 30 * 24 * time.Hour

const prefixoBearer = "Bearer "

var reNomeUsuario = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

// Servico concentra as regras de negócio de usuários e sessões.
type Servico struct {
	usuarios  RepositorioUsuarios
	sessoes   RepositorioSessoes
	hasher    HasherSenha
	hashFalso string // usado para igualar o tempo de resposta quando o usuário não existe
	agora     func() time.Time
}

func NovoServico(usuarios RepositorioUsuarios, sessoes RepositorioSessoes, hasher HasherSenha) (*Servico, error) {
	falso, err := hasher.Hash("senha-falsa-para-tempo-constante")
	if err != nil {
		return nil, err
	}
	return &Servico{usuarios: usuarios, sessoes: sessoes, hasher: hasher, hashFalso: falso, agora: time.Now}, nil
}

// LimparSessoesExpiradas remove sessões vencidas (chamado na inicialização).
func (s *Servico) LimparSessoesExpiradas(ctx context.Context) error {
	return s.sessoes.RemoverExpiradas(ctx, s.agora().Unix())
}

// Cadastrar valida os dados, cria o usuário e já abre uma sessão.
func (s *Servico) Cadastrar(ctx context.Context, in Cadastro) (*Autenticado, error) {
	nome := strings.TrimSpace(in.NomeUsuario)
	cpf := NovoCPF(in.CPF)
	email := strings.ToLower(strings.TrimSpace(in.Email))

	if !reNomeUsuario.MatchString(nome) {
		return nil, falha.Validacao("Nome de usuário deve ter 3 a 20 caracteres (letras, números e _)")
	}
	if !cpf.Valido() {
		return nil, falha.Validacao("CPF inválido")
	}
	if addr, err := mail.ParseAddress(email); err != nil || addr.Address != email || len(email) > 254 {
		return nil, falha.Validacao("Email inválido")
	}
	if len(in.Senha) < 8 || len(in.Senha) > 72 {
		return nil, falha.Validacao("A senha deve ter entre 8 e 72 caracteres")
	}

	hash, err := s.hasher.Hash(in.Senha)
	if err != nil {
		return nil, falha.Interno("Erro interno", err)
	}

	u, err := s.usuarios.Criar(ctx, NovoUsuario{NomeUsuario: nome, CPF: string(cpf), Email: email, SenhaHash: hash})
	if err != nil {
		return nil, err
	}
	u.CPF = cpf.Mascarado()
	return s.abrirSessao(ctx, u)
}

// Entrar autentica por nome de usuário ou e-mail e abre uma sessão.
func (s *Servico) Entrar(ctx context.Context, in Credenciais) (*Autenticado, error) {
	ident := strings.TrimSpace(in.Identificador)
	porEmail := strings.Contains(ident, "@")
	if porEmail {
		ident = strings.ToLower(ident)
	}

	u, hash, err := s.usuarios.BuscarCredencial(ctx, ident, porEmail)
	if err != nil {
		return nil, falha.Interno("Erro interno", err)
	}
	if u == nil {
		s.hasher.Comparar(s.hashFalso, in.Senha) // mesmo custo de tempo de um login real
		return nil, falha.NaoAutorizado("Credenciais inválidas")
	}
	if s.hasher.Comparar(hash, in.Senha) != nil {
		return nil, falha.NaoAutorizado("Credenciais inválidas")
	}
	u.CPF = CPF(u.CPF).Mascarado()
	return s.abrirSessao(ctx, u)
}

// Autenticar resolve o cabeçalho "Authorization: Bearer <token>" no usuário
// dono da sessão e devolve também o hash do token (útil para o logout).
func (s *Servico) Autenticar(ctx context.Context, cabecalho string) (*Usuario, string, error) {
	if !strings.HasPrefix(cabecalho, prefixoBearer) {
		return nil, "", falha.NaoAutorizado("Não autorizado")
	}
	th := hashToken(strings.TrimSpace(cabecalho[len(prefixoBearer):]))
	u, err := s.sessoes.BuscarUsuario(ctx, th, s.agora().Unix())
	if err != nil || u == nil {
		return nil, "", falha.NaoAutorizado("Não autorizado")
	}
	u.CPF = CPF(u.CPF).Mascarado()
	return u, th, nil
}

// Sair encerra a sessão identificada pelo hash do token.
func (s *Servico) Sair(ctx context.Context, tokenHash string) error {
	if err := s.sessoes.Remover(ctx, tokenHash); err != nil {
		return falha.Interno("Erro ao encerrar sessão", err)
	}
	return nil
}

// AjustarReputacao soma delta à reputação, mantendo-a entre 0 e 10.
func (s *Servico) AjustarReputacao(ctx context.Context, in AjusteReputacao) (float64, error) {
	return s.usuarios.AjustarReputacao(ctx, in.ID, in.Delta)
}

func (s *Servico) abrirSessao(ctx context.Context, u *Usuario) (*Autenticado, error) {
	token, err := novoToken()
	if err != nil {
		return nil, falha.Interno("Erro ao criar sessão", err)
	}
	expira := s.agora().Add(DuracaoSessao).Unix()
	if err := s.sessoes.Criar(ctx, hashToken(token), u.ID, expira); err != nil {
		return nil, falha.Interno("Erro ao criar sessão", err)
	}
	return &Autenticado{Token: token, Usuario: u}, nil
}
