package animal

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png" // registra o decodificador PNG
	"sync"
)

const (
	larguraAlvo      = 450
	alturaAlvo       = 600
	tamanhoMaxBytes  = 16 * 1024
	qualidadeInicial = 85
	qualidadeMinima  = 10
)

// ProcessadorImagem normaliza as fotos enviadas.
type ProcessadorImagem interface {
	// Processar recebe os bytes enviados (JPEG/PNG) e devolve um JPEG 450x600
	// de até ~16 KB. Se os bytes forem vazios ou inválidos, devolve o placeholder.
	Processar(dados []byte) []byte
	// Placeholder devolve a imagem padrão (verde) usada quando não há foto.
	Placeholder() []byte
}

// ProcessadorJPEG recorta em 3:4, redimensiona e comprime até caber no limite.
type ProcessadorJPEG struct {
	uma         sync.Once
	placeholder []byte
}

func NovoProcessadorJPEG() *ProcessadorJPEG { return &ProcessadorJPEG{} }

func (p *ProcessadorJPEG) Processar(dados []byte) []byte {
	if len(dados) > 0 {
		if img, _, err := image.Decode(bytes.NewReader(dados)); err == nil {
			if out := p.codificar(img); len(out) > 0 {
				return out
			}
		}
	}
	return p.Placeholder()
}

func (p *ProcessadorJPEG) Placeholder() []byte {
	p.uma.Do(func() {
		img := image.NewRGBA(image.Rect(0, 0, 600, 800))
		draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: 40, G: 120, B: 80, A: 255}}, image.Point{}, draw.Src)
		p.placeholder = p.codificar(img)
	})
	return p.placeholder
}

func (p *ProcessadorJPEG) codificar(img image.Image) []byte {
	dst := recortarEEscalar(img)

	var buf bytes.Buffer
	qualidade := qualidadeInicial
	for {
		buf.Reset()
		if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: qualidade}); err != nil {
			return nil
		}
		if buf.Len() <= tamanhoMaxBytes || qualidade <= qualidadeMinima {
			break
		}
		qualidade -= 5
	}
	return buf.Bytes()
}

// recortarEEscalar faz um recorte central na proporção 3:4 e reduz para 450x600
// (vizinho mais próximo).
func recortarEEscalar(img image.Image) *image.RGBA {
	b := img.Bounds()
	sw, sh := b.Dx(), b.Dy()

	proporcaoAlvo := 3.0 / 4.0
	proporcaoAtual := float64(sw) / float64(sh)

	var cw, ch, x0, y0 int
	if proporcaoAtual > proporcaoAlvo {
		ch = sh
		cw = int(float64(sh) * proporcaoAlvo)
		x0 = (sw - cw) / 2
	} else {
		cw = sw
		ch = int(float64(sw) / proporcaoAlvo)
		y0 = (sh - ch) / 2
	}

	dst := image.NewRGBA(image.Rect(0, 0, larguraAlvo, alturaAlvo))
	for y := 0; y < alturaAlvo; y++ {
		for x := 0; x < larguraAlvo; x++ {
			srcX := b.Min.X + x0 + (x * cw / larguraAlvo)
			srcY := b.Min.Y + y0 + (y * ch / alturaAlvo)
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}
