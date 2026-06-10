package service

import (
	"excel-comparer/internal/model"
	"excel-comparer/internal/utils"
	"time"
)

type CompareService struct{}

func NewCompareService() *CompareService {
	return &CompareService{}
}

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

	for _, rowA := range dadosA {
		encontrouMatchB := false
		indexMatchB := -1
		houveDivergencia := false

		for idxB, rowB := range dadosB {
			match := true

			for colA, colB := range config.Mapping {
				normA := utils.NormalizeValue(rowA[colA], config.Rules)
				normB := utils.NormalizeValue(rowB[colB], config.Rules)

				if normA != normB {
					match = false
					break
				}
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

	resumo.TempoProcessamento = time.Since(tempoInicio)
	return resumo, linhasResultado
}
