package limitador

import (
	"testing"
	"time"
)

func TestBloqueiaAposMaximoEDesbloqueiaAposJanela(t *testing.T) {
	agora := time.Unix(1000, 0)
	l := Novo(2, 10*time.Second)
	l.agora = func() time.Time { return agora }

	if !l.Permitir("ip") || !l.Permitir("ip") {
		t.Fatal("as duas primeiras tentativas deveriam passar")
	}
	if l.Permitir("ip") {
		t.Fatal("a terceira tentativa deveria ser bloqueada")
	}
	if !l.Permitir("outro-ip") {
		t.Fatal("outro IP não deve ser afetado")
	}

	agora = agora.Add(11 * time.Second)
	if !l.Permitir("ip") {
		t.Fatal("após a janela a tentativa deveria passar")
	}
}

func TestLimpaChavesOciosas(t *testing.T) {
	agora := time.Unix(1000, 0)
	l := Novo(1, time.Second)
	l.agora = func() time.Time { return agora }
	l.Permitir("a")
	l.Permitir("b")

	agora = agora.Add(5 * time.Second)
	l.Permitir("c")
	if len(l.acessos) != 1 {
		t.Fatalf("esperava só a chave ativa, há %d chaves", len(l.acessos))
	}
}
