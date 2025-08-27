# Stress Test CLI

Um utilitário simples e eficiente para realizar testes de carga em serviços web via linha de comando (CLI). Escrito em Go, este projeto permite enviar milhares de requisições HTTP com concorrência controlada e gerar um relatório detalhado com métricas de desempenho e status HTTP.

Ideal para testes de desempenho, simulação de carga e análise de resiliência de APIs.

---

## 📦 Funcionalidades

- Teste de carga com número configurável de requisições
- Controle de concorrência (requisições simultâneas)
- Relatório completo ao final da execução:
  - Tempo total de execução
  - Quantidade total de requisições
  - Número de respostas com status 200
  - Distribuição de outros códigos de status (404, 500, etc.)
- Barra de progresso em tempo real

- Execução via Docker para fácil distribuição

---

## 🚀 Como Usar

### 1. Clone o repositório

```bash
git clone https://github.com/oliveiracmorais/stress-test.git
cd stress-test
```

### 2. Construa a imagem Docker

```bash
docker build -t stress-test .
```

### 3. Execute o teste de carga (exemplo)

```bash
docker run --rm stress-test \
  --url=https://httpbin.org/status/200 \
  --requests=1000 \
  --concurrency=10
```
<small><em>Observação: utilizamos uma url de teste aqui. Pode ser substituída por http://google.com para cenário mais realístico<em><small>
