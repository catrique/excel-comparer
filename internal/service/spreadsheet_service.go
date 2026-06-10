package service

import (
	"encoding/csv"
	"errors"
	"excel-comparer/internal/utils"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

type SpreadsheetService struct{}

func NewSpreadsheetService() *SpreadsheetService {
	return &SpreadsheetService{}
}

func (ss *SpreadsheetService) LerCabecalhosEConteudo(path string, separadorCsv rune) ([]string, []map[string]string, error) {
	ext := filepath.Ext(path)

	if ext == ".csv" {
		return ss.lerCSV(path, separadorCsv)
	} else if ext == ".xlsx" || ext == ".xls" {
		return ss.lerExcel(path)
	}

	return nil, nil, errors.New("formato de arquivo não suportado (use .csv, .xlsx ou .xls)")
}

func (ss *SpreadsheetService) lerExcel(path string) ([]string, []map[string]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	listaAbas := f.GetSheetList()
	if len(listaAbas) > 1 {
		return nil, nil, errors.New("apenas arquivos com uma única aba (worksheet) são suportados")
	}
	if len(listaAbas) == 0 {
		return nil, nil, errors.New("nenhuma aba encontrada no arquivo excel")
	}

	nomeAba := listaAbas[0]
	linhas, err := f.GetRows(nomeAba)
	if err != nil {
		return nil, nil, err
	}

	if len(linhas) == 0 {
		return nil, nil, errors.New("o arquivo excel está vazio")
	}

	var cabecalhos []string
	for _, h := range linhas[0] {
		cabecalhos = append(cabecalhos, utils.CorrigirEncoding(h))
	}
	var dados []map[string]string

	for i := 1; i < len(linhas); i++ {
		linha := linhas[i]
		mapaLinha := make(map[string]string)

		for idxCol, nomeCol := range cabecalhos {
			val := ""
			if idxCol < len(linha) {
				val = utils.CorrigirEncoding(linha[idxCol])
			}
			mapaLinha[nomeCol] = val
		}
		dados = append(dados, mapaLinha)
	}

	return cabecalhos, dados, nil
}

func (ss *SpreadsheetService) lerCSV(path string, separador rune) ([]string, []map[string]string, error) {
	arquivo, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer arquivo.Close()

	leitor := csv.NewReader(arquivo)
	leitor.Comma = separador
	leitor.LazyQuotes = true

	registros, err := leitor.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if len(registros) == 0 {
		return nil, nil, errors.New("o arquivo csv está vazio")
	}

	var cabecalhos []string
	for _, h := range registros[0] {
		cabecalhos = append(cabecalhos, utils.CorrigirEncoding(h))
	}

	var dados []map[string]string

	for i := 1; i < len(registros); i++ {
		linha := registros[i]
		mapaLinha := make(map[string]string)

		for idxCol, nomeCol := range cabecalhos {
			val := ""
			if idxCol < len(linha) {
				val = utils.CorrigirEncoding(linha[idxCol])
			}
			mapaLinha[nomeCol] = val
		}
		dados = append(dados, mapaLinha)
	}

	return cabecalhos, dados, nil
}
