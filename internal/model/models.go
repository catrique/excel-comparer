package model

import "time"

type NormalizationRules struct {
	IgnoreCase           bool
	IgnoreAccents        bool
	TrimSpaces           bool
	PadronizarDatas      bool
	RemoverPontuacaoCpf  bool
	RemoverZerosEsquerda bool
}

type CompareConfig struct {
	Mapping      map[string]string
	OrdemColunas []string
	Rules        NormalizationRules
	NomeCustomA  string
	NomeCustomB  string
}

type ExecutiveSummary struct {
	TotalA             int
	TotalB             int
	Iguais             int
	Divergentes        int
	ApenasA            int
	ApenasB            int
	TempoProcessamento time.Duration
}

type ResultRow struct {
	ValoresMapeados map[string]string
	Status          string
}
