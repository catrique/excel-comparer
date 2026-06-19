// package service

// import (
// 	"excel-comparer/internal/model"
// 	"excel-comparer/internal/utils"
// 	"time"
// )

// type CompareService struct{}

// func NewCompareService() *CompareService {
// 	return &CompareService{}
// }

// func (cs *CompareService) ProcessarComparacao(
// 	dadosA []map[string]string,
// 	dadosB []map[string]string,
// 	config model.CompareConfig,
// ) (*model.ExecutiveSummary, []model.ResultRow) {

// 	tempoInicio := time.Now()
// 	var linhasResultado []model.ResultRow

// 	resumo := &model.ExecutiveSummary{
// 		TotalA: len(dadosA),
// 		TotalB: len(dadosB),
// 	}

// 	indicesUtilizadosB := make(map[int]bool)

// 	for _, rowA := range dadosA {
// 		encontrouMatchB := false
// 		indexMatchB := -1
// 		houveDivergencia := false

// 		for idxB, rowB := range dadosB {
// 			match := true

// 			for colA, colB := range config.Mapping {
// 				normA := utils.NormalizeValue(rowA[colA], config.Rules)
// 				normB := utils.NormalizeValue(rowB[colB], config.Rules)

// 				if normA != normB {
// 					match = false
// 					break
// 				}
// 			}

// 			if match && len(config.Mapping) > 0 {
// 				encontrouMatchB = true
// 				indexMatchB = idxB

// 				for colA, colB := range config.Mapping {
// 					nA := utils.NormalizeValue(rowA[colA], config.Rules)
// 					nB := utils.NormalizeValue(rowB[colB], config.Rules)
// 					if nA != nB {
// 						houveDivergencia = true
// 						break
// 					}
// 				}
// 				break
// 			}
// 		}

// 		valoresSaida := make(map[string]string)
// 		for colA, colB := range config.Mapping {
// 			valAOriginal := rowA[colA]
// 			var valBOriginal string
// 			if encontrouMatchB {
// 				valBOriginal = dadosB[indexMatchB][colB]
// 			}

// 			if valAOriginal != "" {
// 				valoresSaida[colA] = valAOriginal
// 			} else {
// 				valoresSaida[colA] = valBOriginal
// 			}
// 		}

// 		status := "Igual"
// 		if !encontrouMatchB {
// 			status = "Exclusivo na Planilha " + config.NomeCustomA
// 			resumo.ApenasA++
// 		} else if houveDivergencia {
// 			status = "Divergente"
// 			resumo.Divergentes++
// 			indicesUtilizadosB[indexMatchB] = true
// 		} else {
// 			resumo.Iguais++
// 			indicesUtilizadosB[indexMatchB] = true
// 		}

// 		linhasResultado = append(linhasResultado, model.ResultRow{
// 			ValoresMapeados: valoresSaida,
// 			Status:          status,
// 		})
// 	}

// 	for idxB, rowB := range dadosB {
// 		if indicesUtilizadosB[idxB] {
// 			continue
// 		}

// 		valoresSaida := make(map[string]string)
// 		for colA, colB := range config.Mapping {
// 			valoresSaida[colA] = rowB[colB]
// 		}

// 		resumo.ApenasB++
// 		linhasResultado = append(linhasResultado, model.ResultRow{
// 			ValoresMapeados: valoresSaida,
// 			Status:          "Exclusivo na Planilha " + config.NomeCustomB,
// 		})
// 	}

// 	resumo.TempoProcessamento = time.Since(tempoInicio)
// 	return resumo, linhasResultado
// }

package service

import (
	"excel-comparer/internal/model"
	"excel-comparer/internal/utils"
	"fmt"
	"os"
	"time"
)

type CompareService struct{}

func NewCompareService() *CompareService {
	return &CompareService{}
}

// debugAtivado: defina como true para imprimir no console os valores
// normalizados das comparações das primeiras linhas. Útil para descobrir
// por que duas linhas que "deveriam" bater não estão batendo.
var debugAtivado = true

// limiteDebug: quantas linhas de A serão impressas no debug
const limiteDebug = 5

func (cs *CompareService) ProcessarComparacao(
	dadosA []map[string]string,
	dadosB []map[string]string,
	config model.CompareConfig,
) (*model.ExecutiveSummary, []model.ResultRow) {

	tempoInicio := time.Now()
	var linhasResultado []model.ResultRow

	resumo := &model.ExecutiveSummary{
		TotalA: len(dadosA),
		TotalB: len(dadosB),
	}

	indicesUtilizadosB := make(map[int]bool)

	if debugAtivado {
		fmt.Fprintln(os.Stderr, "==================== DEBUG COMPARAÇÃO ====================")
		fmt.Fprintf(os.Stderr, "Mapeamento configurado (colunaA -> colunaB): %v\n", config.Mapping)
		fmt.Fprintf(os.Stderr, "Total de linhas A: %d | Total de linhas B: %d\n", len(dadosA), len(dadosB))

		// Verifica se as colunas mapeadas realmente existem nos dados
		if len(dadosA) > 0 {
			for colA := range config.Mapping {
				if _, existe := dadosA[0][colA]; !existe {
					fmt.Fprintf(os.Stderr, "[ALERTA] Coluna '%s' NÃO existe na Planilha A! "+
						"Verifique se carregou o arquivo correto como 'Arquivo A'.\n", colA)
				}
			}
		}
		if len(dadosB) > 0 {
			for _, colB := range config.Mapping {
				if _, existe := dadosB[0][colB]; !existe {
					fmt.Fprintf(os.Stderr, "[ALERTA] Coluna '%s' NÃO existe na Planilha B! "+
						"Verifique se carregou o arquivo correto como 'Arquivo B'.\n", colB)
				}
			}
		}
	}

	for i, rowA := range dadosA {
		encontrouMatchB := false
		indexMatchB := -1
		houveDivergencia := false

		imprimirEstaLinha := debugAtivado && i < limiteDebug

		if imprimirEstaLinha {
			fmt.Fprintf(os.Stderr, "\n--- Linha A #%d ---\n", i+1)
			for colA, colB := range config.Mapping {
				bruA := rowA[colA]
				normA := utils.NormalizeValue(bruA, config.Rules)
				fmt.Fprintf(os.Stderr, "  A[%q] = %q  -> normalizado: %q\n", colA, bruA, normA)
				_ = colB
			}
		}

		for idxB, rowB := range dadosB {
			match := true
			var detalhesComparacao []string

			for colA, colB := range config.Mapping {
				normA := utils.NormalizeValue(rowA[colA], config.Rules)
				normB := utils.NormalizeValue(rowB[colB], config.Rules)

				if imprimirEstaLinha && idxB < 3 {
					detalhesComparacao = append(detalhesComparacao, fmt.Sprintf(
						"%s(%q) vs %s(%q) => %q == %q ? %v",
						colA, rowA[colA], colB, rowB[colB], normA, normB, normA == normB,
					))
				}

				if normA != normB {
					match = false
				}
			}

			if imprimirEstaLinha && idxB < 3 {
				fmt.Fprintf(os.Stderr, "  Comparando com linha B #%d:\n", idxB+1)
				for _, d := range detalhesComparacao {
					fmt.Fprintf(os.Stderr, "    %s\n", d)
				}
				fmt.Fprintf(os.Stderr, "    => MATCH GERAL: %v\n", match)
			}

			if match && len(config.Mapping) > 0 {
				encontrouMatchB = true
				indexMatchB = idxB

				for colA, colB := range config.Mapping {
					nA := utils.NormalizeValue(rowA[colA], config.Rules)
					nB := utils.NormalizeValue(rowB[colB], config.Rules)
					if nA != nB {
						houveDivergencia = true
						break
					}
				}
				break
			}
		}

		valoresSaida := make(map[string]string)
		for colA, colB := range config.Mapping {
			valAOriginal := rowA[colA]
			var valBOriginal string
			if encontrouMatchB {
				valBOriginal = dadosB[indexMatchB][colB]
			}

			if valAOriginal != "" {
				valoresSaida[colA] = valAOriginal
			} else {
				valoresSaida[colA] = valBOriginal
			}
		}

		status := "Igual"
		if !encontrouMatchB {
			status = "Exclusivo na Planilha " + config.NomeCustomA
			resumo.ApenasA++
		} else if houveDivergencia {
			status = "Divergente"
			resumo.Divergentes++
			indicesUtilizadosB[indexMatchB] = true
		} else {
			resumo.Iguais++
			indicesUtilizadosB[indexMatchB] = true
		}

		linhasResultado = append(linhasResultado, model.ResultRow{
			ValoresMapeados: valoresSaida,
			Status:          status,
		})
	}

	for idxB, rowB := range dadosB {
		if indicesUtilizadosB[idxB] {
			continue
		}

		valoresSaida := make(map[string]string)
		for colA, colB := range config.Mapping {
			valoresSaida[colA] = rowB[colB]
		}

		resumo.ApenasB++
		linhasResultado = append(linhasResultado, model.ResultRow{
			ValoresMapeados: valoresSaida,
			Status:          "Exclusivo na Planilha " + config.NomeCustomB,
		})
	}

	if debugAtivado {
		fmt.Fprintln(os.Stderr, "===========================================================")
	}

	resumo.TempoProcessamento = time.Since(tempoInicio)
	return resumo, linhasResultado
}
