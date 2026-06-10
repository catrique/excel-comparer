package test

import (
	"strings"
	"testing"
)

func TestAmbiente(t *testing.T) {
	resultado := strings.ToLower("JOÃO")
	esperado := "joão"

	if resultado != esperado {
		t.Errorf("Esperado %s, mas obteve %s", esperado, resultado)
	} else {
		t.Log("Parabéns! O motor Go está funcionando perfeitamente sem CGO!")
	}
}
