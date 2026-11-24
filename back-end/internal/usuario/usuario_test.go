package usuario

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"discoveryrural/internal/falha"
)

// --- dublês de teste ---------------------------------------------------

type repoUsuariosFalso struct {
	porNome map[string]*Usuario
	hashes  map[string]string
	proximo int64
}

func novoRepoUsuariosFalso() *repoUsuariosFalso {
	return &repoUsuariosFalso{porNome: map[string]*Usuario{}, hashes: map[string]string{}}
}

func (r *repoUsuariosFalso) Criar(_ context.Context, n NovoUsuario) (*Usuario, error) {
	if _, existe := r.porNome[n.NomeUsuario]; existe {
		return nil, falha.Conflito("Nome de usuário já em uso")
	}
	r.proximo++
	u := &Usuario{ID: r.proximo, NomeUsuario: n.NomeUsuario, CPF: n.CPF, Email: n.Email, Reputacao: 5}
	r.porNome[n.NomeUsuario] = u
	r.hashes[n.NomeUsuario] = n.SenhaHash
	return u, nil
}

func (r *repoUsuariosFalso) BuscarCredencial(_ context.Context, ident string, porEmail bool) (*Usuario, string, error) {
	for nome, u := range r.porNome {
		if (porEmail && u.Email == ident) || (!porEmail && nome == ident) {
			copia := *u
			return &copia, r.hashes[nome], nil
		}
	}
	return nil, "", nil
}

func (r *repoUsuariosFalso) AjustarReputacao(context.Context, int64, float64) (float64, error) {
	return 0, errors.New("não usado")
}

type repoSessoesFalso struct {
	porHash  map[string]*Usuario
	usuarios *repoUsuariosFalso
}

func (r *repoSessoesFalso) Criar(_ context.Context, th string, usuarioID, _ int64) error {
	for _, u := range r.usuarios.porNome {
		if u.ID == usuarioID {
			copia := *u
			r.porHash[th] = &copia
		}
	}
	return nil
}
func (r *repoSessoesFalso) BuscarUsuario(_ context.Context, th string, _ int64) (*Usuario, error) {
	return r.porHash[th], nil
}
func (r *repoSessoesFalso) Remover(_ context.Context, th string) error {
	delete(r.porHash, th)
	return nil
}
func (r *repoSessoesFalso) RemoverExpiradas(context.Context, int64) error { return nil }

func novoServicoDeTeste(t *testing.T) *Servico {
	t.Helper()
	usuarios := novoRepoUsuariosFalso()
	sessoes := &repoSessoesFalso{porHash: map[string]*Usuario{}, usuarios: usuarios}
	s, err := NovoServico(usuarios, sessoes, HasherBcrypt{Custo: bcrypt.MinCost})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func cadastroValido() Cadastro {
	return Cadastro{NomeUsuario: "felipe_01", CPF: "529.982.247-25", Email: "Felipe@Exemplo.com", Senha: "senha-forte-123"}
}

// --- CPF ---------------------------------------------------------------

func TestCPF(t *testing.T) {
	if !NovoCPF("529.982.247-25").Valido() {
		t.Error("CPF válido foi recusado")
	}
	for _, ruim := range []string{"111.111.111-11", "529.982.247-24", "123", ""} {
		if NovoCPF(ruim).Valido() {
			t.Errorf("CPF %q deveria ser inválido", ruim)
		}
	}
	if got := NovoCPF("52998224725").Mascarado(); got != "***.982.247-**" {
		t.Errorf("máscara incorreta: %q", got)
	}
}

// --- Serviço -----------------------------------------------------------

func TestCadastrarValidacoes(t *testing.T) {
	casos := map[string]func(*Cadastro){
		"Nome de usuário deve ter 3 a 20 caracteres (letras, números e _)": func(c *Cadastro) { c.NomeUsuario = "a b" },
		"CPF inválido":   func(c *Cadastro) { c.CPF = "000.000.000-00" },
		"Email inválido": func(c *Cadastro) { c.Email = "sem-arroba" },
		"A senha deve ter entre 8 e 72 caracteres": func(c *Cadastro) { c.Senha = "curta" },
	}
	for msg, alterar := range casos {
		s := novoServicoDeTeste(t)
		c := cadastroValido()
		alterar(&c)
		_, err := s.Cadastrar(context.Background(), c)
		var f *falha.Falha
		if !errors.As(err, &f) || f.Tipo != falha.TipoValidacao || f.Mensagem != msg {
			t.Errorf("esperava validação %q, veio %v", msg, err)
		}
	}
}

func TestCadastrarMascaraCPFENormalizaEmail(t *testing.T) {
	s := novoServicoDeTeste(t)
	res, err := s.Cadastrar(context.Background(), cadastroValido())
	if err != nil {
		t.Fatal(err)
	}
	if res.Token == "" {
		t.Error("deveria devolver token")
	}
	if res.Usuario.CPF != "***.982.247-**" {
		t.Errorf("CPF deveria sair mascarado, veio %q", res.Usuario.CPF)
	}
	if res.Usuario.Email != "felipe@exemplo.com" {
		t.Errorf("e-mail deveria ser minúsculo, veio %q", res.Usuario.Email)
	}
}

func TestLoginPorUsuarioOuEmailESenhaErrada(t *testing.T) {
	s := novoServicoDeTeste(t)
	if _, err := s.Cadastrar(context.Background(), cadastroValido()); err != nil {
		t.Fatal(err)
	}

	for _, ident := range []string{"felipe_01", "felipe@exemplo.com"} {
		if _, err := s.Entrar(context.Background(), Credenciais{ident, "senha-forte-123"}); err != nil {
			t.Errorf("login com %q falhou: %v", ident, err)
		}
	}
	for _, c := range []Credenciais{{"felipe_01", "errada-errada"}, {"naoexiste", "senha-forte-123"}} {
		_, err := s.Entrar(context.Background(), c)
		if !falha.EhTipo(err, falha.TipoNaoAutorizado) {
			t.Errorf("credenciais %+v deveriam dar 401, veio %v", c, err)
		}
	}
}

func TestAutenticarELogout(t *testing.T) {
	s := novoServicoDeTeste(t)
	res, _ := s.Cadastrar(context.Background(), cadastroValido())

	u, th, err := s.Autenticar(context.Background(), "Bearer "+res.Token)
	if err != nil || u.NomeUsuario != "felipe_01" {
		t.Fatalf("token válido deveria autenticar: %v", err)
	}
	if err := s.Sair(context.Background(), th); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Autenticar(context.Background(), "Bearer "+res.Token); err == nil {
		t.Error("token deveria ser inválido após logout")
	}
	if _, _, err := s.Autenticar(context.Background(), res.Token); err == nil {
		t.Error("cabeçalho sem 'Bearer ' deveria ser recusado")
	}
}
