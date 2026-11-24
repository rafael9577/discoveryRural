// Package limitador implementa um limitador de taxa por chave (ex.: por IP)
// usando janela deslizante em memória.
package limitador

import (
	"sync"
	"time"
)

type Limitador struct {
	mu       sync.Mutex
	acessos  map[string][]time.Time
	max      int
	janela   time.Duration
	agora    func() time.Time
	ultLimpa time.Time
}

// Novo cria um limitador que aceita até max eventos por janela para cada chave.
func Novo(max int, janela time.Duration) *Limitador {
	return &Limitador{acessos: map[string][]time.Time{}, max: max, janela: janela, agora: time.Now}
}

// Permitir registra o evento e informa se ele está dentro do limite.
func (l *Limitador) Permitir(chave string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	agora := l.agora()
	l.limparSeNecessario(agora)

	recentes := l.acessos[chave][:0]
	for _, t := range l.acessos[chave] {
		if agora.Sub(t) < l.janela {
			recentes = append(recentes, t)
		}
	}
	if len(recentes) >= l.max {
		l.acessos[chave] = recentes
		return false
	}
	l.acessos[chave] = append(recentes, agora)
	return true
}

// limparSeNecessario remove chaves ociosas no máximo uma vez por janela, para
// que o mapa não cresça indefinidamente.
func (l *Limitador) limparSeNecessario(agora time.Time) {
	if agora.Sub(l.ultLimpa) < l.janela {
		return
	}
	l.ultLimpa = agora
	for chave, marcas := range l.acessos {
		if len(marcas) == 0 || agora.Sub(marcas[len(marcas)-1]) >= l.janela {
			delete(l.acessos, chave)
		}
	}
}
