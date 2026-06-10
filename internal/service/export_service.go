package service

import (
	"excel-comparer/internal/model"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

func (es *ExportService) ExportarParaXLSX(
	caminhoDestino string,
	summary *model.ExecutiveSummary,
	resultados []model.ResultRow,
	config model.CompareConfig,
) error {

	f := excelize.NewFile()
	defer f.Close()

	abaDados := "Resultado Unificado"
	f.NewSheet(abaDados)

	colIdx := 1
	for _, cA := range config.OrdemColunas {
		eixo, _ := excelize.CoordinatesToCellName(colIdx, 1)
		f.SetCellValue(abaDados, eixo, strings.ToUpper(cA))
		colIdx++
	}

	eixoStatus, _ := excelize.CoordinatesToCellName(colIdx, 1)
	f.SetCellValue(abaDados, eixoStatus, "STATUS DO REGISTRO")

	for rIdx, r := range resultados {
		linhaPlanilha := rIdx + 2
		cIdx := 1

		for _, cA := range config.OrdemColunas {
			eixo, _ := excelize.CoordinatesToCellName(cIdx, linhaPlanilha)
			f.SetCellValue(abaDados, eixo, r.ValoresMapeados[cA])
			cIdx++
		}

		eixoVStatus, _ := excelize.CoordinatesToCellName(cIdx, linhaPlanilha)
		f.SetCellValue(abaDados, eixoVStatus, r.Status)
	}

	abaResumo := "Indicadores"
	f.NewSheet(abaResumo)
	f.SetCellValue(abaResumo, "A1", "INDICADOR")
	f.SetCellValue(abaResumo, "B1", "MAPEAMENTO")

	f.SetCellValue(abaResumo, "A2", "Linhas Processadas "+config.NomeCustomA)
	f.SetCellValue(abaResumo, "B2", summary.TotalA)
	f.SetCellValue(abaResumo, "A3", "Linhas Processadas "+config.NomeCustomB)
	f.SetCellValue(abaResumo, "B3", summary.TotalB)
	f.SetCellValue(abaResumo, "A4", "Registros Idênticos")
	f.SetCellValue(abaResumo, "B4", summary.Iguais)
	f.SetCellValue(abaResumo, "A5", "Registros com Divergência")
	f.SetCellValue(abaResumo, "B5", summary.Divergentes)
	f.SetCellValue(abaResumo, "A6", "Exclusivos em "+config.NomeCustomA)
	f.SetCellValue(abaResumo, "B6", summary.ApenasA)
	f.SetCellValue(abaResumo, "A7", "Exclusivos em "+config.NomeCustomB)
	f.SetCellValue(abaResumo, "B7", summary.ApenasB)
	f.SetCellValue(abaResumo, "A8", "Tempo Decorrido")
	f.SetCellValue(abaResumo, "B8", summary.TempoProcessamento.String())

	f.DeleteSheet("Sheet1")
	return f.SaveAs(caminhoDestino)
}
