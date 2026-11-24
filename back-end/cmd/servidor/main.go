// Comando servidor: ponto de entrada e "raiz de composição" da aplicação.
// É o único lugar que conhece todas as implementações concretas e as liga
// umas às outras (injeção de dependências).
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"discoveryrural/internal/animal"
	"discoveryrural/internal/avistamento"
	"discoveryrural/internal/banco"
	"discoveryrural/internal/config"
	"discoveryrural/internal/httpx"
	"discoveryrural/internal/limitador"
	"discoveryrural/internal/usuario"
	"discoveryrural/internal/web"
)

func main() {
	_ = godotenv.Load()
	if err := executar(); err != nil {
		log.Fatal(err)
	}
}

func executar() error {
	cfg, err := config.Carregar()
	if err != nil {
		return err
	}
	if cfg.ChaveAdmin == "" {
		log.Println("aviso: ADMIN_KEY não definida; endpoints de administrador responderão 401")
	}

	ctx, parar := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer parar()

	db, err := banco.Abrir(cfg.DSN())
	if err != nil {
		return err
	}
	defer db.Close()
	if err := banco.Migrar(ctx, db); err != nil {
		return err
	}

	// Infraestrutura compartilhada
	admin := httpx.NovoAdmin(cfg.ChaveAdmin)
	ip := httpx.ExtratorIP{ConfiarProxy: cfg.ConfiarProxy}

	// Avistamentos
	repoAvistamentos := avistamento.NovoRepositorioSQL(db)

	// Animais
	servicoAnimais := animal.NovoServico(
		animal.NovoRepositorioSQL(db),
		repoAvistamentos,
		animal.NovoProcessadorJPEG(),
	)
	if err := servicoAnimais.SemearPadrao(ctx); err != nil {
		log.Printf("aviso: não foi possível cadastrar os animais padrão: %v", err)
	}

	servicoAvistamentos := avistamento.NovoServico(repoAvistamentos, servicoAnimais)

	// Usuários
	servicoUsuarios, err := usuario.NovoServico(
		usuario.NovoRepositorioUsuariosSQL(db),
		usuario.NovoRepositorioSessoesSQL(db),
		usuario.NovoHasherBcrypt(),
	)
	if err != nil {
		return err
	}
	if err := servicoUsuarios.LimparSessoesExpiradas(ctx); err != nil {
		log.Printf("aviso: não foi possível limpar sessões expiradas: %v", err)
	}

	// Rotas
	mux := http.NewServeMux()
	modulos := []httpx.Registravel{
		animal.NovoHandler(servicoAnimais, admin),
		avistamento.NovoHandler(servicoAvistamentos, limitador.Novo(1, 10*time.Second), ip),
		usuario.NovoHandler(servicoUsuarios, admin, limitador.Novo(10, 5*time.Minute), ip),
		web.NovoHandler(),
	}
	for _, m := range modulos {
		m.Registrar(mux)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Porta,
		Handler:           httpx.Observar(mux),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	erros := make(chan error, 1)
	go func() {
		log.Println("Servidor rodando na porta " + cfg.Porta)
		erros <- srv.ListenAndServe()
	}()

	select {
	case err := <-erros:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		log.Println("encerrando servidor...")
		fim, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelar()
		return srv.Shutdown(fim)
	}
	return nil
}
