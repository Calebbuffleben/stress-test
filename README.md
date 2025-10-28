# Stress Test CLI

Sistema CLI em Go para realizar testes de carga em serviços web, distribuído via Docker.

## 🚀 Funcionalidades

- **Testes de Carga**: Execute testes de carga em qualquer serviço web
- **Concorrência Configurável**: Controle o número de requests simultâneos
- **Relatórios Detalhados**: Obtenha métricas completas de performance
- **Docker Ready**: Execução simplificada via containers
- **Validação Robusta**: Validação completa de parâmetros de entrada

## 📋 Requisitos

- Docker instalado e funcionando

## 🐳 Instalação e Uso

### 1. Construir a Imagem Docker

```bash
# Clone o repositório
git clone <seu-repositorio>
cd stress-test

# Construa a imagem Docker
docker build -t stress-test .
```

### 2. Executar Testes de Carga

#### Sintaxe Básica
```bash
docker run --rm stress-test \
  --url=<URL_DO_SERVICO> \
  --requests=<NUMERO_TOTAL_REQUESTS> \
  --concurrency=<NUMERO_CHAMADAS_SIMULTANEAS>
```

#### Exemplos Práticos

**Teste básico com Google:**
```bash
docker run --rm stress-test \
  --url=http://google.com \
  --requests=100 \
  --concurrency=10
```

**Teste com serviço local:**
```bash
docker run --rm stress-test \
  --url=http://localhost:8080/api/health \
  --requests=500 \
  --concurrency=25
```

**Teste com endpoint específico:**
```bash
docker run --rm stress-test \
  --url=https://httpbin.org/status/200 \
  --requests=1000 \
  --concurrency=50
```

## 📊 Parâmetros de Entrada

| Parâmetro | Descrição | Obrigatório | Exemplo |
|-----------|-----------|-------------|---------|
| `--url` | URL do serviço a ser testado | ✅ Sim | `http://google.com` |
| `--requests` | Número total de requests | ✅ Sim | `1000` |
| `--concurrency` | Número de chamadas simultâneas | ✅ Sim | `10` |

### Validações Automáticas

- **URL**: Deve ser uma URL válida e não pode estar vazia
- **Requests**: Deve ser maior que 0
- **Concurrency**: Deve ser maior que 0
- **Ajuste Automático**: Se `concurrency > requests`, será ajustado automaticamente

## 📈 Relatório de Saída

Após a execução, o sistema gera um relatório completo com:

```
==== Relatório do Teste de Carga ====
Tempo total: 2.68048471s
Total de requests: 100
HTTP 200: 95
Outros status:
  404: 3
  500: 2
```

### Métricas Incluídas

- ⏱️ **Tempo Total**: Duração completa da execução
- 📊 **Total de Requests**: Quantidade total de requests realizados
- ✅ **HTTP 200**: Quantidade de requests bem-sucedidos
- ❌ **Outros Status**: Distribuição de códigos de erro (404, 500, etc.)
- 🔴 **Erros**: Requests que falharam por problemas de rede/conexão

## 🧪 Cenários de Teste

### Teste de Performance
```bash
# Teste pesado para avaliar performance
docker run --rm stress-test \
  --url=https://api.exemplo.com/endpoint \
  --requests=5000 \
  --concurrency=100
```

### Teste de Disponibilidade
```bash
# Teste rápido para verificar disponibilidade
docker run --rm stress-test \
  --url=https://meusite.com/health \
  --requests=50 \
  --concurrency=5
```

### Teste de Stress
```bash
# Teste de stress com alta concorrência
docker run --rm stress-test \
  --url=https://api.exemplo.com/heavy-endpoint \
  --requests=1000 \
  --concurrency=200
```

## 🔧 Configurações Avançadas

### Timeout de Requests
- **Timeout padrão**: 30 segundos por request
- **Reutilização de conexões**: Habilitada para melhor performance
- **Compressão**: Habilitada por padrão

### Otimizações de Performance
- Pool de conexões HTTP otimizado
- Gerenciamento eficiente de recursos
- Workers concorrentes para máxima throughput

## 🚨 Tratamento de Erros

### Tipos de Erro Detectados

1. **Erros de Validação**:
   ```bash
   erro: parâmetro --url é obrigatório
   erro: --requests deve ser maior que 0
   erro: --concurrency deve ser maior que 0
   erro: url inválida: parse error
   ```

2. **Erros de Rede**:
   - Timeouts de conexão
   - DNS resolution failures
   - Connection refused
   - Exibidos como "erro" no relatório

3. **Status HTTP de Erro**:
   - 4xx: Client errors (400, 401, 403, 404, etc.)
   - 5xx: Server errors (500, 502, 503, 504, etc.)

## 📝 Exemplos Completos

### Exemplo 1: Teste de API REST
```bash
docker run --rm stress-test \
  --url=https://jsonplaceholder.typicode.com/posts/1 \
  --requests=200 \
  --concurrency=20
```

**Saída esperada:**
```
==== Relatório do Teste de Carga ====
Tempo total: 3.245s
Total de requests: 200
HTTP 200: 200
Outros status: nenhum
```

### Exemplo 2: Teste com Falhas
```bash
docker run --rm stress-test \
  --url=https://httpbin.org/status/500 \
  --requests=50 \
  --concurrency=5
```

**Saída esperada:**
```
==== Relatório do Teste de Carga ====
Tempo total: 1.890s
Total de requests: 50
HTTP 200: 0
Outros status:
  500: 50
```

### Exemplo 3: Teste com URL Inválida
```bash
docker run --rm stress-test \
  --url=https://urlquenaoexiste.com.br \
  --requests=10 \
  --concurrency=2
```

**Saída esperada:**
```
==== Relatório do Teste de Carga ====
Tempo total: 2.110s
Total de requests: 10
HTTP 200: 0
Outros status:
  erro: 10
```

## 🏗️ Arquitetura

### Componentes Principais

- **CLI Load Test**: Executa os testes de carga
- **HTTP Client**: Cliente HTTP otimizado com pool de conexões
- **Worker Pool**: Goroutines para concorrência controlada
- **Result Collector**: Coleta e agrega resultados
- **Report Generator**: Gera relatórios detalhados

### Fluxo de Execução

1. **Validação**: Verifica parâmetros de entrada
2. **Configuração**: Configura cliente HTTP e workers
3. **Execução**: Distribui requests entre workers
4. **Coleta**: Agrega resultados de todos os workers
5. **Relatório**: Gera relatório final com métricas

## 🔍 Troubleshooting

### Problemas Comuns

**1. Erro de conexão com Docker:**
```bash
# Verifique se o Docker está rodando
docker --version
docker ps
```

**2. Timeout de requests:**
- Aumente o timeout se necessário (modificar código)
- Verifique conectividade de rede
- Teste com URLs menores primeiro

**3. Performance baixa:**
- Ajuste o nível de concorrência
- Verifique recursos do sistema
- Teste com menos requests primeiro

## 📚 Dependências

- **Go 1.22+**: Linguagem de programação
- **Docker**: Containerização
- **Alpine Linux**: Imagem base otimizada

---
