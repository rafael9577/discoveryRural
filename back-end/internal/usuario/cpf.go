package usuario

import "strings"

// CPF é um valor que encapsula as regras do documento.
type CPF string

// NovoCPF mantém apenas os dígitos do texto informado.
func NovoCPF(bruto string) CPF {
	var b strings.Builder
	for _, c := range bruto {
		if c >= '0' && c <= '9' {
			b.WriteRune(c)
		}
	}
	return CPF(b.String())
}

// Valido confere tamanho, sequências repetidas e os dois dígitos verificadores.
func (c CPF) Valido() bool {
	s := string(c)
	if len(s) != 11 {
		return false
	}
	iguais := true
	for i := 1; i < 11; i++ {
		if s[i] != s[0] {
			iguais = false
			break
		}
	}
	if iguais {
		return false
	}
	for t := 9; t < 11; t++ {
		soma := 0
		for i := 0; i < t; i++ {
			soma += int(s[i]-'0') * (t + 1 - i)
		}
		d := (soma * 10) % 11
		if d == 10 {
			d = 0
		}
		if d != int(s[t]-'0') {
			return false
		}
	}
	return true
}

// Mascarado devolve o CPF no formato ***.123.456-** (vazio se inválido).
func (c CPF) Mascarado() string {
	s := string(c)
	if len(s) != 11 {
		return ""
	}
	return "***." + s[3:6] + "." + s[6:9] + "-**"
}
