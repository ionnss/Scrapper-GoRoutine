# Scrapper Concorrente

[![N|Solid](https://cldup.com/dTxpPi9lDf.thumb.png)](https://nodesource.com/products/nsolid)

[![Concurrency](https://th.bing.com/th/id/OIP.izLHuE6uuuFn5yByFPmrIgHaFj?rs=1&pid=ImgDetMain)](https://th.bing.com/th/id/OIP.izLHuE6uuuFn5yByFPmrIgHaFj?rs=1&pid=ImgDetMain)


## Objetivo do programa
Acessar páginas web ao mesmo tempo para extrair o título de cada página e exibir esses títulos no terminal. Isso é feito utilizando concorrência em Go, que permite acessar várias páginas simultâneamente, economizando tempo. 

## Explicação do Código

1. Pacotes utilizados
```Golang
import (
	"fmt"
	"net/http"
	"sync"
	"github.com/PuerkitoBio/goquery"
)
```

2. Função `fetchTitle`

Essa função é reponsável por:
- Acessar uma página web (url)
- Extrair o título da página
- Evniar o resultado para um canal

```
func fetchTitle(url string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done() // Marca a goroutine como concluída no WaitGroup
```
Parâmetros da função:
- `url string`: Representa o endereço da página web (url) que vamos acessar para obter o título
- `wg *sync.WaitGroup`: Pointeiro para um `WaitGroup, que usamos para sincornizar o término de todas as tarefas (goroutines) que estão rodando ao mesmo tempo. O `*` indica que estamos passando um "enderço" para o `WaitGroup` e não uma cópia dela.
- `results chan<- string`: Este é um canal unidirecional que permite enviar strings para outra parte do programa. Ele é usado para passar resultados (títulos ou mensagens de erro) para a função `main`

A linha `defer wg.Done()` diz ao programa para marcar esta tarefa (goroutine) como concluída quando a função fetchTitle terminar. Isso é importante para que o `main` saiba quando todas as tarefas foram concluídas.



---

##### Requisição HTTP

```go
	req, err := http.Get(url)
	if err != nil {
		results <- fmt.Sprintf("Erro ao acessar %s: %v", url, err)
		return
	}
	defer req.Body.Close()
```

- **`http.Get(url)`**: Esta linha faz um **pedido HTTP GET** para a URL. Isso significa que estamos acessando a página e pedindo ao servidor o conteúdo dela.
- **`err != nil`**: Aqui verificamos se houve algum erro ao acessar a página (por exemplo, se a página não existe ou o servidor não está respondendo). Se houver erro, enviamos uma mensagem para o canal `results` e encerramos a função com `return`.
- **`defer req.Body.Close()`**: Isso garante que, depois que terminarmos de usar o conteúdo da página, liberaremos a memória alocada para armazená-lo.

---

##### Verificação de Status

```go
	if req.StatusCode != 200 {
		results <- fmt.Sprintf("Erro ao acessar %s: status %d %s", url, req.StatusCode, req.Status)
		return
	}
```

- **`req.StatusCode != 200`**: Verificamos se o servidor respondeu com o código **200 OK** (indica sucesso). Se não for 200, isso significa que a página não foi carregada corretamente. Então, enviamos uma mensagem de erro para o canal `results` e encerramos a função.

---

##### Carregamento e Busca do Título

```go
	doc, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		results <- fmt.Sprintf("Erro ao carregar documento de %s: %v", url, err)
		return
	}
	title := doc.Find("title").Text()
	results <- fmt.Sprintf("Título de %s: %s", url, title)
}
```

- **`goquery.NewDocumentFromReader(req.Body)`**: Carregamos o conteúdo HTML da página (fornecido por `req.Body`) no `goquery`, que permite navegar e buscar partes específicas do HTML.
- **`doc.Find("title").Text()`**: Procuramos a tag `<title>` no HTML da página e pegamos o texto dentro dela (ou seja, o título).
- **`results <- fmt.Sprintf("Título de %s: %s", url, title)`**: Enviamos o título extraído para o canal `results`, onde ele será lido mais tarde.
---

#### 3. Função `main`

A função `main` é a função principal que configura e controla o programa.

```go
func main() {
	urls := []string{
		"http://olos.novagne.com.br/Olos/login.aspx?logout=true",
		"http://sistema.novagne.com.br/novagne/",
	}
```

- **`urls := []string{...}`**: Definimos uma lista de URLs que queremos processar. Cada URL será passada para uma goroutine que extrairá o título da página.

---

##### Configuração do WaitGroup e do Canal

```go
	var wg sync.WaitGroup
	results := make(chan string, len(urls)) // Canal para armazenar os resultados
```

- **`var wg sync.WaitGroup`**: Criamos uma nova instância de `WaitGroup`, que controlará o número de goroutines e garantirá que todas terminem antes que o programa finalize.
- **`results := make(chan string, len(urls))`**: Criamos um canal `results` com capacidade igual ao número de URLs. Esse canal armazenará as mensagens com os títulos ou erros.

---

##### Início das Goroutines

```go
	for _, url := range urls {
		wg.Add(1)
		go fetchTitle(url, &wg, results)
	}
```

- **`for _, url := range urls`**: Aqui, percorremos cada URL da lista.
- **`wg.Add(1)`**: Para cada URL, incrementamos o contador do `WaitGroup` para indicar que uma nova tarefa (goroutine) será iniciada.
- **`go fetchTitle(url, &wg, results)`**: Chamamos `fetchTitle` como uma **goroutine** para cada URL, ou seja, fazemos com que ela rode em paralelo com as outras.

---

##### Espera e Exibição dos Resultados

```go
	wg.Wait()
	close(results)
	for result := range results {
		fmt.Println(result)
	}
}
```

- **`wg.Wait()`**: Aqui o programa principal espera que todas as goroutines sejam concluídas. Ele só avança para a próxima linha quando todas as tarefas terminarem.
- **`close(results)`**: Fechamos o canal `results` para indicar que não serão enviados mais dados para ele. Isso ajuda a evitar erros ao tentar ler um canal ainda aberto.
- **`for result := range results`**: Finalmente, percorremos todos os valores no canal `results`, exibindo o título ou a mensagem de erro de cada URL.

---

### Resumo Final

- O programa acessa várias páginas web em paralelo para extrair seus títulos.
- Cada URL é processada por uma goroutine, que faz o scraping de forma independente.
- Um `WaitGroup` controla a sincronização, garantindo que o `main` espere até que todas as goroutines terminem.
- Um canal (`results`) coleta e exibe todos os resultados (títulos ou mensagens de erro) de forma organizada.

Esse programa demonstra o uso de **concorrência em Go** de maneira eficiente, permitindo processar múltiplas páginas ao mesmo tempo e otimizar o tempo de execução.
