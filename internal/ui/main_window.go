package ui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"excel-comparer/internal/model"
	"excel-comparer/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	nativeDialog "github.com/sqweek/dialog"
)

type MainWindow struct {
	window             fyne.Window
	spreadsheetService *service.SpreadsheetService
	compareService     *service.CompareService
	exportService      *service.ExportService

	caminhoA        string
	caminhoB        string
	colunasA        []string
	colunasB        []string
	dadosA          []map[string]string
	dadosB          []map[string]string
	mapeamento      map[string]string
	ordemMapeamento []string
	regras          model.NormalizationRules

	apelidoA string
	apelidoB string

	progressA       *widget.ProgressBar
	progressB       *widget.ProgressBar
	progressProc    *widget.ProgressBarInfinite
	rectStatusA     *canvas.Rectangle
	rectStatusB     *canvas.Rectangle
	lblStatusA      *widget.Label
	lblStatusB      *widget.Label
	btnEditarNomeA  *widget.Button
	btnEditarNomeB  *widget.Button
	comboA          *widget.Select
	comboB          *widget.Select
	btnAdicionarVin *widget.Button
	containerLista  *fyne.Container
	btnComparar     *widget.Button
}

func NewMainWindow(w fyne.Window, ss *service.SpreadsheetService, cs *service.CompareService, es *service.ExportService) *MainWindow {
	return &MainWindow{
		window:             w,
		spreadsheetService: ss,
		compareService:     cs,
		exportService:      es,
		mapeamento:         make(map[string]string),
		ordemMapeamento:    []string{},
		apelidoA:           "Planilha A",
		apelidoB:           "Planilha B",
	}
}

func (mw *MainWindow) BuildUI() {
	mw.progressA = widget.NewProgressBar()
	mw.progressA.Hide()
	mw.rectStatusA = canvas.NewRectangle(color.Transparent)
	mw.rectStatusA.SetMinSize(fyne.NewSize(15, 15))
	mw.lblStatusA = widget.NewLabel("Planilha A: Não selecionada")

	mw.btnEditarNomeA = widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		mw.abrirModalRenomear("A", func(novoNome string) {
			mw.apelidoA = novoNome
			if mw.caminhoA != "" {
				mw.lblStatusA.SetText(fmt.Sprintf("Planilha %s carregada!", mw.apelidoA))
			} else {
				mw.lblStatusA.SetText(fmt.Sprintf("%s: Não selecionada", mw.apelidoA))
			}
			mw.atualizarPlaceholdersCombos()
		})
	})
	mw.btnEditarNomeA.Disable()

	btnCarregarA := widget.NewButtonWithIcon("Buscar Arquivo A", theme.FileIcon(), func() {
		mw.abrirGerenciadorNativo(func(path string) {
			mw.lblStatusA.SetText("Carregando arquivo...")
			mw.progressA.Show()
			mw.progressA.SetValue(0.5)
			go func() {
				cols, dados, err := mw.spreadsheetService.LerCabecalhosEConteudo(path, ';')
				if err != nil {
					mw.progressA.Hide()
					dialog.ShowError(err, mw.window)
					return
				}
				mw.colunasA = cols
				mw.dadosA = dados
				mw.caminhoA = path
				mw.progressA.SetValue(1.0)
				mw.rectStatusA.FillColor = color.RGBA{R: 40, G: 167, B: 69, A: 255}
				mw.rectStatusA.Refresh()
				mw.lblStatusA.SetText(fmt.Sprintf("Planilha %s carregada!", mw.apelidoA))
				mw.btnEditarNomeA.Enable()
				mw.configurarMenusMapeamento()
			}()
		})
	})

	mw.progressB = widget.NewProgressBar()
	mw.progressB.Hide()
	mw.rectStatusB = canvas.NewRectangle(color.Transparent)
	mw.rectStatusB.SetMinSize(fyne.NewSize(15, 15))
	mw.lblStatusB = widget.NewLabel("Planilha B: Não selecionada")

	mw.btnEditarNomeB = widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		mw.abrirModalRenomear("B", func(novoNome string) {
			mw.apelidoB = novoNome
			if mw.caminhoB != "" {
				mw.lblStatusB.SetText(fmt.Sprintf("Planilha %s carregada!", mw.apelidoB))
			} else {
				mw.lblStatusB.SetText(fmt.Sprintf("%s: Não selecionada", mw.apelidoB))
			}
			mw.atualizarPlaceholdersCombos()
		})
	})
	mw.btnEditarNomeB.Disable()

	btnCarregarB := widget.NewButtonWithIcon("Buscar Arquivo B", theme.FileIcon(), func() {
		mw.abrirGerenciadorNativo(func(path string) {
			mw.lblStatusB.SetText("Carregando arquivo...")
			mw.progressB.Show()
			mw.progressB.SetValue(0.5)
			go func() {
				cols, dados, err := mw.spreadsheetService.LerCabecalhosEConteudo(path, ';')
				if err != nil {
					mw.progressB.Hide()
					dialog.ShowError(err, mw.window)
					return
				}
				mw.colunasB = cols
				mw.dadosB = dados
				mw.caminhoB = path
				mw.progressB.SetValue(1.0)
				mw.rectStatusB.FillColor = color.RGBA{R: 40, G: 167, B: 69, A: 255}
				mw.rectStatusB.Refresh()
				mw.lblStatusB.SetText(fmt.Sprintf("Planilha %s carregada!", mw.apelidoB))
				mw.btnEditarNomeB.Enable()
				mw.configurarMenusMapeamento()
			}()
		})
	})

	mw.comboA = widget.NewSelect([]string{}, func(s string) {})
	mw.comboB = widget.NewSelect([]string{}, func(s string) {})
	mw.atualizarPlaceholdersCombos()

	mw.containerLista = container.NewVBox()

	mw.btnAdicionarVin = widget.NewButtonWithIcon("Vincular Colunas", theme.ContentAddIcon(), func() {
		cA := mw.comboA.Selected
		cB := mw.comboB.Selected
		if cA == "" || cB == "" {
			return
		}

		if _, existe := mw.mapeamento[cA]; !existe {
			mw.ordemMapeamento = append(mw.ordemMapeamento, cA)
		}

		mw.mapeamento[cA] = cB
		mw.comboA.SetSelected("")
		mw.comboB.SetSelected("")
		mw.renderizarListaComLixeiras()
		mw.btnComparar.Enable()
	})
	mw.btnAdicionarVin.Disable()

	btnLimparMapeamento := widget.NewButtonWithIcon("Limpar Todos os Vínculos", theme.DeleteIcon(), func() {
		mw.mapeamento = make(map[string]string)
		mw.ordemMapeamento = []string{}
		mw.renderizarListaComLixeiras()
		mw.btnComparar.Disable()
	})

	checkCase := widget.NewCheck("Ignorar Case (Maiúsculas/Minúsculas)", func(b bool) { mw.regras.IgnoreCase = b })
	checkAcentos := widget.NewCheck("Ignorar Acentos", func(b bool) { mw.regras.IgnoreAccents = b })
	checkEspacos := widget.NewCheck("Forçar Remoção de Espaços Extras (Trim)", func(b bool) { mw.regras.TrimSpaces = b })
	checkDatas := widget.NewCheck("Auto-Padronizar Formato de Datas", func(b bool) { mw.regras.PadronizarDatas = b })
	checkCpf := widget.NewCheck("Limpar Símbolos de CPF/CNPJ", func(b bool) { mw.regras.RemoverPontuacaoCpf = b })
	checkZeros := widget.NewCheck("Suprimir Zeros à Esquerda", func(b bool) { mw.regras.RemoverZerosEsquerda = b })

	mw.progressProc = widget.NewProgressBarInfinite()
	mw.progressProc.Hide()

	mw.btnComparar = widget.NewButtonWithIcon("Executar Comparação e Gerar Excel", theme.ConfirmIcon(), func() {
		mw.executarFluxo()
	})
	mw.btnComparar.Importance = widget.HighImportance
	mw.btnComparar.Disable()

	mw.window.SetContent(container.NewVScroll(container.NewVBox(
		widget.NewLabelWithStyle("1. Arquivos de Origem (Use o lápis para dar apelidos)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(2,
			container.NewVBox(btnCarregarA, container.NewHBox(mw.rectStatusA, mw.lblStatusA, mw.btnEditarNomeA), mw.progressA),
			container.NewVBox(btnCarregarB, container.NewHBox(mw.rectStatusB, mw.lblStatusB, mw.btnEditarNomeB), mw.progressB),
		),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("2. Opções de Alinhamento e Tratamento", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(3, checkCase, checkAcentos, checkEspacos, checkDatas, checkCpf, checkZeros),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("3. Mapeamento dos Campos Cruzados", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(3, mw.comboA, mw.comboB, container.NewVBox(mw.btnAdicionarVin, btnLimparMapeamento)),
		widget.NewLabelWithStyle("Campos incluídos no relatório final (Clique na lixeira para remover individualmente):", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
		mw.containerLista,
		widget.NewSeparator(),
		mw.btnComparar,
		mw.progressProc,
	)))
}

func (mw *MainWindow) atualizarPlaceholdersCombos() {
	mw.comboA.PlaceHolder = "Coluna de " + mw.apelidoA
	mw.comboB.PlaceHolder = "Coluna de " + mw.apelidoB
	mw.comboA.Refresh()
	mw.comboB.Refresh()
}

func (mw *MainWindow) abrirModalRenomear(tipo string, onSalvar func(string)) {
	txtInput := widget.NewEntry()
	txtInput.SetPlaceHolder("Digite o apelido...")
	if tipo == "A" {
		txtInput.SetText(mw.apelidoA)
	} else {
		txtInput.SetText(mw.apelidoB)
	}

	dialog.ShowForm("Definir Apelido", "Salvar", "Cancelar", []*widget.FormItem{
		widget.NewFormItem("Nome:", txtInput),
	}, func(confirmado bool) {
		if confirmado && strings.TrimSpace(txtInput.Text) != "" {
			onSalvar(strings.TrimSpace(txtInput.Text))
		}
	}, mw.window)
}

func (mw *MainWindow) configurarMenusMapeamento() {
	if len(mw.colunasA) > 0 && len(mw.colunasB) > 0 {
		mw.comboA.Options = mw.colunasA
		mw.comboB.Options = mw.colunasB
		mw.comboA.Refresh()
		mw.comboB.Refresh()
		mw.btnAdicionarVin.Enable()
	}
}

func (mw *MainWindow) renderizarListaComLixeiras() {
	mw.containerLista.Objects = nil

	if len(mw.ordemMapeamento) == 0 {
		mw.containerLista.Add(widget.NewLabel("Nenhuma coluna vinculada ainda."))
		mw.containerLista.Refresh()
		return
	}

	for _, colA := range mw.ordemMapeamento {
		colunaAtualA := colA
		colunaCorrespondenteB := mw.mapeamento[colunaAtualA]

		textoVinculo := widget.NewLabel(fmt.Sprintf("• Saída: %s  ➔  (Alinhada com %s: %s)", colunaAtualA, mw.apelidoB, colunaCorrespondenteB))

		btnLixeira := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			delete(mw.mapeamento, colunaAtualA)
			for i, v := range mw.ordemMapeamento {
				if v == colunaAtualA {
					mw.ordemMapeamento = append(mw.ordemMapeamento[:i], mw.ordemMapeamento[i+1:]...)
					break
				}
			}
			mw.renderizarListaComLixeiras()
			if len(mw.mapeamento) == 0 {
				mw.btnComparar.Disable()
			}
		})

		linhaContainer := container.NewHBox(btnLixeira, textoVinculo)
		mw.containerLista.Add(linhaContainer)
	}

	mw.containerLista.Refresh()
}

func (mw *MainWindow) executarFluxo() {
	mw.btnComparar.Disable()
	mw.btnComparar.SetText("Processando batimento...")
	mw.progressProc.Show()
	mw.progressProc.Start()

	go func() {
		time.Sleep(150 * time.Millisecond)

		config := model.CompareConfig{
			Mapping:      mw.mapeamento,
			OrdemColunas: mw.ordemMapeamento,
			Rules:        mw.regras,
			NomeCustomA:  mw.apelidoA,
			NomeCustomB:  mw.apelidoB,
		}

		resumo, resultados := mw.compareService.ProcessarComparacao(mw.dadosA, mw.dadosB, config)

		caminhoSaida, err := nativeDialog.File().Title("Salvar Relatório Combinado").Filter("Planilha Excel (*.xlsx)", "xlsx").Save()
		if err != nil || caminhoSaida == "" {
			mw.restaurarBotaoComparar()
			return
		}
		if !strings.HasSuffix(caminhoSaida, ".xlsx") {
			caminhoSaida += ".xlsx"
		}

		err = mw.exportService.ExportarParaXLSX(caminhoSaida, resumo, resultados, config)
		if err != nil {
			mw.restaurarBotaoComparar()
			dialog.ShowError(err, mw.window)
			return
		}

		mw.restaurarBotaoComparar()

		msg := fmt.Sprintf("Conciliação Concluída!\n\nIguais: %d\nDivergentes: %d\nExclusivos [%s]: %d\nExclusivos [%s]: %d",
			resumo.Iguais, resumo.Divergentes, config.NomeCustomA, resumo.ApenasA, config.NomeCustomB, resumo.ApenasB)
		dialog.ShowInformation("Sucesso", msg, mw.window)
	}()
}

func (mw *MainWindow) restaurarBotaoComparar() {
	mw.progressProc.Stop()
	mw.progressProc.Hide()
	mw.btnComparar.SetText("Executar Comparação e Gerar Excel")
	mw.btnComparar.Enable()
}

func (mw *MainWindow) abrirGerenciadorNativo(onSelecionado func(string)) {
	path, err := nativeDialog.File().Title("Escolher Planilha").Filter("Arquivos suportados (*.xlsx, *.csv)", "xlsx", "csv").Load()
	if err != nil || path == "" {
		return
	}
	onSelecionado(path)
}
