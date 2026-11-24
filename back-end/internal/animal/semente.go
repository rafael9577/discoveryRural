package animal

import "context"

// SemearPadrao cadastra as duas espécies iniciais quando a tabela está vazia.
func (s *Servico) SemearPadrao(ctx context.Context) error {
	total, err := s.repo.Contar(ctx)
	if err != nil || total > 0 {
		return err
	}
	imagem := s.imagens.Placeholder()
	for _, a := range animaisPadrao() {
		if _, err := s.repo.Criar(ctx, a, imagem); err != nil {
			return err
		}
	}
	return nil
}

func animaisPadrao() []Animal {
	return []Animal{
		{
			NomePopular:        "Coruja-buraqueira",
			NomeCientifico:     "Athene cunicularia",
			BreveDescricao:     "Pequena coruja de hábitos terrestres.",
			EstadoConservacao:  "Pouco preocupante",
			Habito:             "Diurno/Crepuscular",
			VocalizacaoPostura: "Emite chamados agudos e balança a cabeça.",
			AvisoSeguranca:     "Nunca tente manipular, alimentar ou se aproximar excessivamente de ninhos no chão.",
			ProtocoloIncidente: "Manter distância segura e isolar a área.",
			ContatosResgate:    "193 / Setor Ambiental",
			Raridade:           2,
			BitmapID:           18,
		},
		{
			NomePopular:        "Cágado-de-barbicha",
			NomeCientifico:     "Phrynops geoffroanus",
			BreveDescricao:     "Réptil aquático encontrado em lagos e açudes.",
			EstadoConservacao:  "Pouco preocupante",
			Habito:             "Diurno",
			VocalizacaoPostura: "Silencioso na água; emite som soprado quando acuado.",
			AvisoSeguranca:     "Nunca tente manipular, alimentar ou remover o animal.",
			ProtocoloIncidente: "Evitar contato direto e direcioná-lo com cuidado para a água se estiver em risco.",
			ContatosResgate:    "Setor de Meio Ambiente / Biologia",
			Raridade:           3,
			BitmapID:           384,
		},
	}
}
