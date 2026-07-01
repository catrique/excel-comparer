package test

import (
	"testing"

	"excel-comparer/internal/model"
	"excel-comparer/internal/utils"
)

func TestNormalizeValue_IgnoreCase(t *testing.T) {
	rules := model.NormalizationRules{IgnoreCase: true}
	got := utils.NormalizeValue("joão da Silva", rules)
	want := "JOÃO DA SILVA"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_IgnoreAccents(t *testing.T) {
	rules := model.NormalizationRules{IgnoreAccents: true}
	got := utils.NormalizeValue("José Ávila", rules)
	want := "Jose Avila"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_TrimAndCollapseSpaces(t *testing.T) {
	rules := model.NormalizationRules{}
	got := utils.NormalizeValue("  texto   com   espaços   ", rules)
	want := "texto com espaços"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_SpacesAlwaysTrimmedRegardlessOfFlag(t *testing.T) {
	rules := model.NormalizationRules{TrimSpaces: false}
	got := utils.NormalizeValue("   valor   ", rules)
	want := "valor"
	if got != want {
		t.Errorf("got %q, want %q (TrimSpaces=false não deveria importar hoje)", got, want)
	}
}

func TestNormalizeValue_TabsAndNewlinesBecomeSpaces(t *testing.T) {
	rules := model.NormalizationRules{}
	got := utils.NormalizeValue("a\tb\nc\rd", rules)
	want := "a b c d"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_DecimalCommaToDot(t *testing.T) {
	rules := model.NormalizationRules{}
	got := utils.NormalizeValue("1234,56", rules)
	want := "1234.56"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_DecimalBrazilianThousandsAndComma(t *testing.T) {
	rules := model.NormalizationRules{}
	got := utils.NormalizeValue("1.234,56", rules)
	want := "1234.56"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_TrailingZerosAfterDecimalPointRemoved(t *testing.T) {
	rules := model.NormalizationRules{}
	got := utils.NormalizeValue("100.00", rules)
	want := "100"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_RemoverZerosEsquerda(t *testing.T) {
	rules := model.NormalizationRules{RemoverZerosEsquerda: true}
	got := utils.NormalizeValue("000123", rules)
	want := "123"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_RemoverZerosEsquerda_TodosZeros(t *testing.T) {
	rules := model.NormalizationRules{RemoverZerosEsquerda: true}
	got := utils.NormalizeValue("0000", rules)
	want := "0"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_PadronizarDatas_ISOParaBr(t *testing.T) {
	rules := model.NormalizationRules{PadronizarDatas: true}
	got := utils.NormalizeValue("2024-03-05", rules)
	want := "05/03/2024"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_PadronizarDatas_AnoCurto(t *testing.T) {
	rules := model.NormalizationRules{PadronizarDatas: true}
	got := utils.NormalizeValue("5/3/24", rules)
	want := "05/03/2024"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_RemoverPontuacaoCpf(t *testing.T) {
	rules := model.NormalizationRules{RemoverPontuacaoCpf: true}
	got := utils.NormalizeValue("123.456.789-00", rules)
	want := "12345678900"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeValue_RemoverPontuacaoCpf_NaoAfetaDatas(t *testing.T) {
	rules := model.NormalizationRules{RemoverPontuacaoCpf: true}
	got := utils.NormalizeValue("12/03/2024", rules)
	want := "12/03/2024"
	if got != want {
		t.Errorf("got %q, want %q (datas não devem ser tratadas como CPF/CNPJ)", got, want)
	}
}

func TestNormalizeValue_RegrasCombinadas(t *testing.T) {
	rules := model.NormalizationRules{
		IgnoreCase:    true,
		IgnoreAccents: true,
	}
	got := utils.NormalizeValue("  São Paulo  ", rules)
	want := "SAO PAULO"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTratarData_FormatoISO(t *testing.T) {
	got := utils.TratarData("2024-12-31")
	want := "31/12/2024"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTratarData_FormatoISOComBarra(t *testing.T) {
	got := utils.TratarData("2024/01/05")
	want := "05/01/2024"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTratarData_AnoCurtoBr(t *testing.T) {
	got := utils.TratarData("1/2/24")
	want := "01/02/2024"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTratarData_PreenchimentoDiaEMes(t *testing.T) {
	got := utils.TratarData("5/3/2024")
	want := "05/03/2024"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTratarData_RemoveSufixoDeHorario(t *testing.T) {
	got := utils.TratarData("2024-05-10 14:30:00")
	want := "10/05/2024"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRemoverAcentos(t *testing.T) {
	casos := map[string]string{
		"São Paulo":     "Sao Paulo",
		"Pão de Açúcar": "Pao de Acucar",
		"Crítica":       "Critica",
		"sem acento":    "sem acento",
	}
	for entrada, esperado := range casos {
		got := utils.RemoverAcentos(entrada)
		if got != esperado {
			t.Errorf("RemoverAcentos(%q) = %q, want %q", entrada, got, esperado)
		}
	}
}

func TestCorrigirEncoding(t *testing.T) {
	casos := map[string]string{
		"InformaÃ§Ã£o": "Informação",
		"AmanhÃ£":      "Amanhã",
		"VocÃª":        "Você",
	}
	for entrada, esperado := range casos {
		got := utils.CorrigirEncoding(entrada)
		if got != esperado {
			t.Errorf("CorrigirEncoding(%q) = %q, want %q", entrada, got, esperado)
		}
	}
}
