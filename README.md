# Comparador de planilhas

O **Comparador de planilhas** é uma ferramenta de interface gráfica desenvolvida em Go (Fyne) para consolidação, cruzamento e auditoria de dados entre duas planilhas distintas (Excel ou CSV). O sistema permite o mapeamento dinâmico de chaves e colunas de dados, aplicando regras inteligentes de higienização sem alterar a integridade visual do relatório de saída.

---

## 🚀 Funcionalidades

* **Interface Gráfica Nativa (GUI):** Operação simples e intuitiva via menus visuais, eliminando o uso de janelas de terminal.
* **Apelidos Dinâmicos:** Permite renomear os blocos das planilhas em tempo real, atualizando os menus suspensos e os status do relatório final automaticamente.
* **Mapeamento de Colunas Personalizado:** Respeita estritamente a ordem com que os campos foram vinculados para estruturar a sequência do arquivo de saída.
* **Gestão de Vínculos por Lixeira:** Botões de exclusão individual para cada relacionamento de coluna criado, permitindo ajustes rápidos sem perder as outras configurações.
* **Higienização Avançada de Dados:** Opções opcionais para ignorar maiúsculas/minúsculas, remover acentos, suprimir espaços extras, limpar pontuações de CPF/CNPJ, padronizar datas e zeros à esquerda durante o cruzamento.
* **Preservação de Caixa (Case):** O motor de comparação preserva a grafia original do texto (letras maiúsculas e minúsculas intactas) no documento gerado.

---

## 🛠️ Pré-requisitos de Desenvolvimento

Para rodar ou compilar o projeto, você precisará ter instalado em sua máquina:

1. **Go (Golang):** Versão 1.18 ou superior. [Download Oficial](https://go.dev/dl/)
2. **Compilador C (GCC):** Necessário para a biblioteca gráfica (Fyne) compilar nativamente no ecossistema Windows.
   * *Recomendação para Windows:* Instalar o [TDM-GCC](https://jasonwilliams2006.github.io/tdm-gcc/) ou configurar via MSYS2.

---

## 📦 Como Compilar o Código

Abra o terminal de sua preferência (como o Git Bash ou Prompt de Comando) na raiz do projeto onde está o arquivo `go.mod` e execute o comando abaixo:

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o Comparar_Planilhas.exe cmd/main.go

```

### O que este comando faz?

* **`GOOS=windows GOARCH=amd64`**: Força a compilação de um binário executável para Windows 64-bits.
* **`-ldflags="-H=windowsgui"`**: Remove a janela preta do prompt de comando em segundo plano, fazendo com que apenas a tela do aplicativo apareça ao ser iniciado.
* **`-o Comparar_Planilhas.exe`**: Define o nome do arquivo executável final que será gerado na raiz.

---

## 💡 Como Usar o Aplicativo

1. **Carregar os Arquivos:** Clique em "Buscar Arquivo A" e "Buscar Arquivo B" para selecionar os documentos que deseja analisar.
2. **Definir Apelidos:** Clique no botão de lápis (📝) ao lado do status de carregamento para dar um nome personalizado a cada planilha.
3. **Configurar Regras:** Marque as caixas de seleção desejadas na seção de tratamento (remover acentos, ignorar maiúsculas, etc.) para calibrar o nível de precisão do batimento.
4. **Vincular as Colunas:** Use os menus suspensos para selecionar uma coluna de cada planilha correspondente e clique em **"Vincular Colunas"**. Repita o processo para todas as colunas que devem constar no relatório final.
5. **Gerar Relatório:** Clique no botão destacado **"Executar Comparação e Gerar Excel"**, escolha o local de destino e salve o novo arquivo `.xlsx`.

---

## 👥 Estrutura do Relatório Gerado

O arquivo Excel de saída é gerado contendo duas abas principais:

1. **Resultado Unificado:** Apresenta todas as linhas processadas com as colunas na ordem exata definida pelos seus cliques na tela. A última coluna (`STATUS DO REGISTRO`) classifica a linha como `Igual`, `Divergente` ou `Exclusivo na Planilha [Nome Customizado]`.
2. **Indicadores:** Fornece um sumário gerencial completo com estatísticas de linhas processadas, totalizadores de igualdades, divergências capturadas e o tempo exato decorrido na operação.
