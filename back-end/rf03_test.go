package main

import (
	"net/http"
	"testing"
)

func TestRF03_CoordenadasInvalidas(t *testing.T) {
	casos := []struct {
		nome  string
		corpo string
	}{
		{"latitude acima de 90", `{"animal_id":1,"latitude":91,"longitude":-34.9}`},
		{"latitude abaixo de -90", `{"animal_id":1,"latitude":-91,"longitude":-34.9}`},
		{"longitude acima de 180", `{"animal_id":1,"latitude":-8.0,"longitude":181}`},
		{"longitude abaixo de -180", `{"animal_id":1,"latitude":-8.0,"longitude":-181}`},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			r := postRelato(t, c.corpo, ipFicticio())
			if r.Status != http.StatusBadRequest {
				t.Errorf("status = %d, esperado 400 (corpo: %q)", r.Status, r.Texto())
			}
		})
	}
}

func TestRF03_CamposObrigatorios(t *testing.T) {
	casos := []struct {
		nome  string
		corpo string
	}{
		{"sem animal_id", `{"latitude":-8.0,"longitude":-34.9}`},
		{"sem latitude", `{"animal_id":1,"longitude":-34.9}`},
		{"sem longitude", `{"animal_id":1,"latitude":-8.0}`},
		{"corpo que não é JSON", `isto não é um json`},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			r := postRelato(t, c.corpo, ipFicticio())
			if r.Status != http.StatusBadRequest {
				t.Errorf("status = %d, esperado 400 (corpo: %q)", r.Status, r.Texto())
			}
		})
	}
}
