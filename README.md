# Load Tester CLI

Uma ferramenta CLI em Go para realizar testes de carga em serviços web com execução concorrente e relatórios detalhados.

## Características

- ✅ Testes de carga HTTP com concorrência configurável
- ✅ Relatórios detalhados com métricas de performance
- ✅ Distribuição de códigos de status HTTP
- ✅ Containerização com Docker
- ✅ Interface CLI simples e intuitiva

## Instalação

### Opção 1: Usando Docker

```bash
# Build da imagem Docker
docker build -t load-tester .

# Executar teste de carga
docker run load-tester --url=http://google.com --requests=1000 --concurrency=10
```

### Opção 2: Build Local

```bash
# Build da aplicação
go build -o load-tester main.go

# Executar teste de carga
./load-tester --url=http://google.com --requests=1000 --concurrency=10
```

## Uso

### Parâmetros Obrigatórios

- `--url`: URL do serviço a ser testado
- `--requests`: Número total de requests a serem realizados
- `--concurrency`: Número de chamadas simultâneas

### Exemplos de Uso

```bash
# Teste básico
docker run load-tester --url=http://httpbin.org/get --requests=100 --concurrency=5

# Teste de alta concorrência
docker run load-tester --url=http://google.com --requests=1000 --concurrency=50

# Teste local
docker run load-tester --url=http://localhost:8080/api/health --requests=500 --concurrency=20
```

### Exemplo de Saída

```
Starting load test...
URL: http://google.com
Total Requests: 1000
Concurrency: 10

=== LOAD TEST REPORT ===
Total execution time: 15.234s
Total requests made: 1000
Successful requests (HTTP 200): 950

HTTP Status Code Distribution:
  200: 950 requests (95.0%)
  301: 30 requests (3.0%)
  404: 20 requests (2.0%)

Performance Metrics:
  Average request time: 152ms
  Requests per second: 65.67
```

## Estrutura do Projeto

```
desafio-cli/
├── main.go        # Código principal da aplicação
├── go.mod         # Módulo Go
├── go.sum         # Checksum das dependências
├── Dockerfile     # Containerização
└── README.md      # Esta documentação
```

## Funcionalidades Implementadas

### 1. Parser CLI
- Validação de parâmetros obrigatórios
- Tratamento de erros de entrada
- Validação de consistência (ex: concorrência não pode ser maior que total de requests)

### 2. Motor de Testes de Carga
- Execução concorrente usando goroutines
- Pool de workers configurável
- Timeout de 30 segundos por request
- Coleta de métricas em tempo real

### 3. Sistema de Relatórios
- Tempo total de execução
- Contagem total de requests
- Requests bem-sucedidos (HTTP 200)
- Distribuição completa de códigos de status
- Métricas de performance (requests/segundo, tempo médio)

### 4. Containerização
- Multi-stage build para otimização de tamanho
- Imagem final baseada em Alpine Linux
- Suporte a HTTPS com certificados CA

## Arquitetura

A aplicação utiliza o padrão **Producer-Consumer** com:

- **Channel `requestChan`**: Queue de requests a serem processados
- **Channel `resultChan`**: Coleta dos resultados de cada request
- **Workers (goroutines)**: Executam os requests HTTP de forma concorrente
- **WaitGroup**: Sincronização para aguardar conclusão de todos os workers

### Fluxo de Execução

1. **Parsing e Validação**: Processa e valida os argumentos CLI
2. **Inicialização**: Cria channels e pool de workers
3. **Distribuição**: Envia requests para o channel de trabalho
4. **Execução Concorrente**: Workers processam requests em paralelo
5. **Coleta de Resultados**: Agrega métricas de todos os requests
6. **Relatório**: Apresenta estatísticas consolidadas

## Tecnologias Utilizadas

- **Go 1.21**: Linguagem de programação
- **Goroutines**: Concorrência nativa do Go
- **Channels**: Comunicação entre goroutines
- **net/http**: Cliente HTTP padrão do Go
- **Docker**: Containerização
- **Alpine Linux**: Imagem base otimizada

## Limitações e Considerações

- Timeout fixo de 30 segundos por request
- Suporte apenas para método HTTP GET
- Não implementa autenticação personalizada
- Não salva logs detalhados em arquivos

## Possíveis Melhorias Futuras

- [ ] Suporte a outros métodos HTTP (POST, PUT, etc.)
- [ ] Headers personalizados
- [ ] Autenticação (Bearer token, Basic Auth)
- [ ] Exportação de relatórios em formatos diversos (JSON, CSV)
- [ ] Configuração de timeout personalizável
- [ ] Histograma de latências
- [ ] Modo de aquecimento (warm-up)
- [ ] Monitoramento em tempo real

## Contribuição

Para contribuir com o projeto:

1. Faça um fork do repositório
2. Crie uma branch para sua feature (`git checkout -b feature/amazing-feature`)
3. Faça commit das suas mudanças (`git commit -m 'Add some amazing feature'`)
4. Push para a branch (`git push origin feature/amazing-feature`)
5. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
