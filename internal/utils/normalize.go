package utils

import (
	"excel-comparer/internal/model"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var reEspacosMultiplos = regexp.MustCompile(`\s+`)
var reDataISO = regexp.MustCompile(`^(\d{4})[-/](\d{2})[-/](\d{2})`)
var reDataBrCurta = regexp.MustCompile(`^(\d{2})/(\d{2})/(\d{2})$`)

var reElesParecemDocumento = regexp.MustCompile(`^\d{2,3}[\.\-/\s]\d{3}`)

var reElePareceData = regexp.MustCompile(`\d{2,4}[-/]\d{2}[-/]\d{2,4}`)

func NormalizeValue(val string, rules model.NormalizationRules) string {
	val = strings.NewReplacer("\t", " ", "\n", " ", "\r", " ").Replace(val)
	val = strings.TrimSpace(val)
	val = reEspacosMultiplos.ReplaceAllString(val, " ")

	if rules.IgnoreCase {
		val = strings.ToUpper(val)
	}

	if rules.IgnoreAccents {
		val = RemoverAcentos(val)
	}

	if rules.RemoverZerosEsquerda {
		val = strings.TrimLeft(val, "0")
		if val == "" {
			val = "0"
		}
	}

	if rules.PadronizarDatas && reElePareceData.MatchString(val) {
		val = TratarData(val)
	}

	if rules.RemoverPontuacaoCpf && reElesParecemDocumento.MatchString(val) {
		if !reElePareceData.MatchString(val) {
			val = strings.NewReplacer(".", "", "-", "", "/", "").Replace(val)
		}
	}

	return val
}

func TratarData(val string) string {
	valLimpa := strings.TrimSpace(val)

	if len(valLimpa) > 10 && (strings.Contains(valLimpa, "/") || strings.Contains(valLimpa, "-")) {
		partes := strings.Fields(valLimpa)
		if len(partes) > 0 {
			valLimpa = partes[0]
		}
	}

	if reDataISO.MatchString(valLimpa) {
		matches := reDataISO.FindStringSubmatch(valLimpa)
		if len(matches) == 4 {
			valLimpa = matches[3] + "/" + matches[2] + "/" + matches[1]
		}
	}

	if reDataBrCurta.MatchString(valLimpa) {
		matches := reDataBrCurta.FindStringSubmatch(valLimpa)
		if len(matches) == 4 {
			ano := matches[3]
			if len(ano) == 2 {
				valLimpa = matches[1] + "/" + matches[2] + "/20" + ano
			}
		}
	}

	return valLimpa
}

func RemoverAcentos(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	resultado, _, _ := transform.String(t, s)
	return resultado
}

func CorrigirEncoding(s string) string {
	m := map[string]string{
		"Ã§Ã£": "çã", "Ã£": "ã", "Ã©": "é", "Ã³": "ó", "Ã­": "í",
		"Ãº": "ú", "Ã¢": "â", "Ãª": "ê", "Ã´": "ô", "Ã‡": "Ç",
		"Ã": "Á", "Ã‰": "É", "Ã“": "Ó", "Âº": "º", "Âª": "ª",
		"Ã§": "ç",
	}
	for errado, correto := range m {
		s = strings.ReplaceAll(s, errado, correto)
	}
	return s
}
