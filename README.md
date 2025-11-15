# Desafio - Stress test

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
