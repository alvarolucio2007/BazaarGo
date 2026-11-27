# BazaarGo

<!--toc:start-->
- [BazaarGo](#bazaargo)
  - [Visão geral](#visão-geral)
  - [Diferenciais do projeto](#diferenciais-do-projeto)
  - [Stack de tecnologias](#stack-de-tecnologias)
  - [Arquitetura](#arquitetura)
    - [Estrutura do projeto](#estrutura-do-projeto)
  - [Banco de dados (Persistência Poliglota)](#banco-de-dados-persistência-poliglota)
    - [SGBDs Utilizados](#sgbds-utilizados)
    - [Modelagem de Dados](#modelagem-de-dados)
  - [Endpoints da API](#endpoints-da-api)
    - [Infraestrutura & Autenticação](#infraestrutura-autenticação)
    - [Catálogo de Produtos](#catálogo-de-produtos)
    - [Checkout](#checkout)
    - [Assistente de IA (Chatbot & RAG)](#assistente-de-ia-chatbot-rag)
  - [Como executar](#como-executar)
    - [Requisitos](#requisitos)
    - [Passos](#passos)
  - [Decisões de projeto](#decisões-de-projeto)
    - [Arquitetura Hexagonal](#arquitetura-hexagonal)
    - [Spec-First com OpenAPI](#spec-first-com-openapi)
    - [Assistente de IA com RAG e PGVector](#assistente-de-ia-com-rag-e-pgvector)
    - [Observabilidade de Nível de Produção](#observabilidade-de-nível-de-produção)
    - [Roteamento com Traefik](#roteamento-com-traefik)
    - [Integração com WhatsApp](#integração-com-whatsapp)
    - [Cache Semântico (Otimização de Custos e Latência)](#cache-semântico-otimização-de-custos-e-latência)
  - [Decisões de Arquitetura (ADRs)](#decisões-de-arquitetura-adrs)
  - [Performance e Resiliência](#performance-e-resiliência)
    - [Escalonamento Horizontal (Load Balancing)](#escalonamento-horizontal-load-balancing)
    - [Resultados do Teste de Estresse](#resultados-do-teste-de-estresse)
    - [Análise de Gargalos (Profiling)](#análise-de-gargalos-profiling)
  - [Licença](#licença)
<!--toc:end-->

E-commerce Fullstack (Catálogo de Produtos e Assistente de IA com RAG) com arquitetura hexagonal, persistência poliglota, roteamento inteligente e observabilidade distribuída.

---

## Visão geral

O **BazaarGo** é uma plataforma de e-commerce que permite o cadastro e a listagem de produtos, gerenciamento de carrinho de compras, atendimento automatizado via **Chatbot de Inteligência Artificial (RAG)** e finalização de pedidos via WhatsApp. O projeto foi construído para demonstrar, na prática, a aplicação do conceito de **Banco de Dados Poliglota** (utilizando a tecnologia certa para cada tipo de dado, incluindo busca vetorial) aliado a uma arquitetura hexagonal monolítica modular que gera módulos leves, desacoplados e altamente observáveis.

O sistema atende a três frentes principais:

1. **Módulo Administrativo:** Onde administradores autenticados (via PostgreSQL) cadastram e gerenciam o catálogo de produtos (MongoDB).
2. **Módulo Cliente:** Uma vitrine que exibe os produtos, registra o histórico de "vistos recentemente" (Valkey) e permite fechar o pedido gerando uma mensagem formatada para o WhatsApp.
3. **Assistente de IA (Chatbot com RAG):** Microsserviço dedicado que utiliza **Google Gemini** (`google.golang.org/genai`) e busca semântica com **PGVector** no PostgreSQL para responder dúvidas dos clientes sobre os produtos e políticas da loja com base em contexto real recuperado.

> **Nota de transparência:** Este projeto utilizou ferramentas de Inteligência Artificial (IA) como assistente de desenvolvimento durante todo o processo. A IA foi utilizada para acelerar a construção das interfaces visuais, mas também atuou diretamente como um "copiloto" técnico: auxiliando na configuração da infraestrutura (Docker, Traefik, stack LGTM), instrumentação da observabilidade (OpenTelemetry), integração do pipeline de RAG/Gemini, resolução de bugs complexos e na discussão de tomadas de decisão sobre padrões de projeto e arquitetura. Apesar do forte auxílio da IA, a implementação, os testes e a compreensão de cada decisão técnica foram de responsabilidade do desenvolvedor.

---

## Diferenciais do projeto

Além dos requisitos mínimos, o projeto foi construído seguindo práticas adotadas em ambientes de produção corporativos:

- **Persistência Poliglota + Busca Vetorial:** Uso de PostgreSQL 18 com **PGVector** (Autenticação/ACID + Embeddings para RAG), MongoDB (Catálogo/Flexibilidade de Schema) e Valkey/Redis (Cache e Histórico de visualizações).
- **Chatbot de IA com RAG (Google Gemini):** Serviço dedicado em Go integrado ao SDK oficial da Google (`genai`) para geração de embeddings e respostas contextualizadas via Retrieval-Augmented Generation.
- **Cache Semântico para LLM (Otimização de Custos e Latência):** Implementação de uma camada de cache inteligente usando PGVector. Perguntas com significado similar são respondidas em milissegundos sem acionar a API do Google Gemini, reduzindo custos e protegendo o sistema contra indisponibilidades (erros 503) do provedor de IA.
- **Observabilidade Completa (LGTM + Pyroscope):** Rastreamento distribuído (Tempo), agregação de logs estruturados (Loki), coleta de métricas (Prometheus/Mimir) e profiling contínuo de CPU/Memória (Pyroscope) cobrindo tanto a API principal quanto o serviço de Chatbot.
- **Arquitetura Hexagonal (Ports and Adapters):** Domínio puro e isolado de frameworks, bancos de dados, provedores de LLM e detalhes de HTTP.
- **Spec-First com OpenAPI (`oapi-codegen`):** A API é gerada a partir do contrato `openapi.yaml`, garantindo _Strict Server_ (tipagem forte e validação automática de requests/responses).
- **Reverse Proxy com Traefik:** Roteamento inteligente e isolamento de rede com subdomínios locais (`api.localhost`, `chat.localhost`, `admin.localhost`, `front.localhost`, `telemetry.localhost`).
- **Autenticação com PASETO v4:** Substituindo o JWT por um padrão mais seguro contra ataques de análise de tempo e confusão de algoritmos.
- **Hash de senha com Argon2id:** Configurado com parâmetros acima das recomendações mínimas do OWASP.
- **Infraestrutura 100% Containerizada:** Orquestração completa via Docker Compose.

---

## Stack de tecnologias

| Camada                     | Tecnologia                                        |
| -------------------------- | ------------------------------------------------- |
| **Backend & IA**           |                                                   |
| Linguagem                  | Go (Golang)                                       |
| Roteador HTTP              | chi                                               |
| Design de API              | OpenAPI 3.0 (Spec-First + Strict Server)          |
| IA Generativa & Embeddings | Google Gemini (`google.golang.org/genai`)         |
| Pipeline RAG               | Busca Vetorial via PGVector + Prompt Augmentation |
| Autenticação               | PASETO v4                                         |
| Hash de senha              | Argon2id                                          |
| **Bancos de Dados**        |                                                   |
| Relacional & Vetorial      | PostgreSQL 18 + Extensão **PGVector**             |
| Documentos (Produtos)      | MongoDB 9                                         |
| Cache/Listas (Histórico)   | Valkey 9 (Redis)                                  |
| **Observabilidade**        |                                                   |
| Coleta e Visualização      | Grafana, Tempo, Loki, Prometheus, Pyroscope       |
| Instrumentação             | OpenTelemetry (SDK Go + OTLP gRPC/HTTP)           |
| **Frontend & UI**          |                                                   |
| Vitrine e Painel Admin     | Python + Streamlit                                |
| **Infraestrutura**         |                                                   |
| Proxy Reverso              | Traefik v3                                        |
| Containerização            | Docker e Docker Compose                           |

---

## Arquitetura

O projeto adota uma arquitetura de microsserviços desacoplados, orquestrada pelo Traefik. Os serviços backend (**API Core** e **Chatbot RAG**) seguem os princípios da **Arquitetura Hexagonal**, isolando regras de negócio, bancos de dados e a integração com a API do Google Gemini.

```mermaid
flowchart TD
    %% Atores
    User([Cliente])
    Admin([Administrador])

    %% Proxy e Roteamento
    User <-->|http://front.localhost| Traefik[Traefik Reverse Proxy]
    Admin <-->|http://admin.localhost| Traefik

    %% Frontends
    Traefik <--> Vitrine[Vitrine Streamlit]
    Traefik <--> AdminPanel[Painel Admin Streamlit]

    %% Backends
    Traefik <-->|http://api.localhost| API[Go API - Core Hexagonal]
    Traefik <-->|http://chat.localhost| ChatSvc[Go Chatbot - RAG Service]
    Vitrine <--> API
    Vitrine <-.-> ChatSvc
    AdminPanel <--> API

    %% Bancos de Dados (Persistência Poliglota + Vetorial)
    API <--> PG[(PostgreSQL 18 + PGVector)]
    API <--> Mongo[(MongoDB)]
    API <--> Valkey[(Valkey / Redis)]
    ChatSvc <-->|Busca Vetorial / Embeddings| PG

    %% Serviços Externos (IA)
    ChatSvc <-->|SDK GenAI / LLM + Embeddings| Gemini[[Google Gemini API]]

    %% Observabilidade
    API -.->|OTLP / Traces, Metrics, Logs| Telemetry[Stack LGTM + Pyroscope]
    ChatSvc -.->|OTLP / Traces, Metrics, Logs| Telemetry
    Traefik -.->|http://telemetry.localhost| Telemetry

    %% Estilização
    classDef proxy fill:#f9f,stroke:#333,stroke-width:2px,color:#000;
    classDef frontend fill:#bbf,stroke:#333,stroke-width:2px,color:#000;
    classDef backend fill:#fbb,stroke:#333,stroke-width:2px,color:#000;
    classDef database fill:#bfb,stroke:#333,stroke-width:2px,color:#000;
    classDef obs fill:#ffb,stroke:#333,stroke-width:2px,color:#000;
    classDef external fill:#e1d5e7,stroke:#333,stroke-width:2px,color:#000;

    class Traefik proxy;
    class Vitrine,AdminPanel frontend;
    class API,ChatSvc backend;
    class PG,Mongo,Valkey database;
    class Telemetry obs;
    class Gemini external;
```

### Estrutura do projeto

```
BazaarGo/
├── cmd/
│   ├── api/                 # Entrypoint da API principal (Go)
│   ├── chatbot/             # Entrypoint do microsserviço de Chatbot IA (Go)
│   └── seed/                # Script para popular o banco vetorial com FAQs
├── internal/
│   ├── domain/              # Entidades (User, Product, Chat) e regras puras
│   ├── ports/               # Interfaces (Repositories, Cache, Auth, VectorStore, LLM)
│   ├── service/             # Casos de uso (UserService, ProductService, ChatService)
│   ├── telemetry/           # Setup do OpenTelemetry (Traces, Metrics, Logs)
│   ├── chatbot/             # Domínio e lógica específica do serviço de IA
│   └── adapters/
│       ├── inbound/http/    # Handlers (OpenAPI e Chat), middlewares e rotas
│       └── outbound/        # Implementações (Postgres, PGVector, MongoDB, Valkey, Gemini, PASETO, Argon2)
├── api/
│   ├── openapi.yaml         # A "Fonte da Verdade" do contrato da API
│   └── chatbot.yaml         # A "Fonte da Verdade" do contrato do chatbot
├── database/
│   └── migrations/          # Scripts SQL para o PostgreSQL e extensão PGVector
├── front/                   # Aplicações Python (Vitrine e Admin em Streamlit)
│   └── admin/               # Painel Administrativo
├── docker-compose.yml       # Orquestração (Traefik, API, Chatbot, Bancos, LGTM)
└── README.md
```

---

## Banco de dados (Persistência Poliglota)

A grande estrela deste projeto é a escolha cirúrgica do banco de dados para cada problema, agora estendida para busca semântica:

### SGBDs Utilizados

- **PostgreSQL 18 + PGVector:** Além de controlar **Autenticação e Usuários** com garantias ACID, utiliza a extensão **PGVector** para armazenar vetores de alta dimensão (_embeddings_), viabilizando buscas por similaridade semântica no pipeline de RAG sem a necessidade de manter um banco vetorial separado.
- **MongoDB:** Utilizado para o **Catálogo de Produtos**. Seu modelo de documentos flexível permite que produtos tenham atributos dinâmicos (`Specs: map[string]any`), o que seria um pesadelo de normalização no SQL.
- **Valkey (Redis):** Utilizado para o **Cache** de produtos consultados e para o histórico de **"Vistos Recentemente"**. Utiliza estruturas de dados em memória para entregar respostas em menos de 1 milissegundo.

### Modelagem de Dados

- **Users (Postgres):** Tabela relacional contendo `id`, `email`, `password_hash`, etc.
- **Embeddings / Knowledge Base (Postgres + PGVector):** Tabela vetorial armazenando os trechos de contexto/produtos e suas representações vetoriais (`vector`), consultadas via distância de cosseno/L2.
- **Products (MongoDB):** Documento contendo `_id` (ObjectID hex de 24 caracteres), `name`, `price`, `stock` e `specs` (objeto flexível).
- **Recent Views (Valkey):** Utilização da estrutura de dados **List** (`LPUSH`, `LTRIM`), mantendo apenas os últimos 5 produtos vistos por cada visitante (identificado via Header `X-Client-ID`).

---

## Endpoints da API

A documentação completa e interativa está disponível via Swagger UI em `[http://api.localhost/swagger/index.html](http://api.localhost/swagger/index.html)` (API Core) e `[http://chat.localhost/swagger/index.html](http://chat.localhost/swagger/index.html)` (Chatbot IA).

### Infraestrutura & Autenticação

| Método | Rota             | Descrição                                          |
| ------ | ---------------- | -------------------------------------------------- |
| GET    | `/health`        | Health check concorrente (Postgres, Mongo, Valkey) |
| POST   | `/auth/register` | Criação de um novo Administrador                   |
| POST   | `/auth/token`    | Login e geração do token PASETO                    |

### Catálogo de Produtos

| Método | Rota                      | Descrição                                                             |
| ------ | ------------------------- | --------------------------------------------------------------------- |
| GET    | `/products`               | Lista todos os produtos do catálogo                                   |
| POST   | `/products`               | Cadastra um novo produto (Requer Bearer Token)                        |
| GET    | `/products/{id}`          | Detalhes do produto (Requer Header `X-Client-ID`)                     |
| PUT    | `/products/{id}`          | Atualiza um produto (Requer Bearer Token e ID Mongo válido)           |
| DELETE | `/products/{id}`          | Deleta um produto (Requer Bearer Token)                               |
| GET    | `/products/search/{name}` | Busca produtos por nome (Case-insensitive)                            |
| GET    | `/products/recent`        | Retorna os produtos vistos recentemente (Requer Header `X-Client-ID`) |

### Checkout

| Método | Rota        | Descrição                                            |
| ------ | ----------- | ---------------------------------------------------- |
| POST   | `/checkout` | Recebe o carrinho e gera a URL formatada do WhatsApp |

### Assistente de IA (Chatbot & RAG)

Servido pelo microsserviço `bazaar-chat` ([http://chat.localhost](http://chat.localhost)):

| Método | Rota      | Descrição                                                                                             |
| ------ | --------- | ----------------------------------------------------------------------------------------------------- |
| GET    | `/health` | Health check do serviço de Chatbot                                                                    |
| POST   | `/chat`   | Envia uma pergunta do usuário, busca contexto relevante no PGVector e gera resposta via Google Gemini |

---

## Como executar

### Requisitos

- Docker e Docker Compose
- Chave de API do Google Gemini (`GEMINI_API_KEY`)
- Make (Opcional)

### Passos

**1. Clone o repositório e configure as variáveis de ambiente**

```bash
git clone https://github.com/seu-usuario/bazaargo.git
cd bazaargo
cp .env.example .env
```

Configure o seu `.env` com as seguintes variáveis:

```env
DATABASE_URL=postgres://admin:root@db:5432/bazaargo?sslmode=disable
MONGO_URL=mongodb://admin:admin123@mongo:27017/bazaargo?authSource=admin
CACHE_URL=cache:6379
SERVER_ADDR=0.0.0.0:8080
GEMINI_API_KEY=sua_chave_aqui
```

**2. Suba a stack completa (Traefik, Go API, Go Chatbot, Bancos, LGTM)**

```bash
docker compose up -d --build
```

**3. Popule o banco de dados vetorial (PGVector)**
Para que o assistente de IA funcione, é necessário ingerir a base de conhecimento (FAQs) e gerar os embeddings via Gemini. Execute o script de seed:

```bash
# Se estiver rodando os containers
docker compose exec api go run cmd/seed/main.go

# Ou se estiver rodando localmente
go run cmd/seed/main.go
```

**4. Acesse as aplicações via Traefik:**
O Traefik gerencia as rotas baseadas em subdomínios locais (padrão na maioria dos SOs modernos).

- **Vitrine (Cliente):** <http://front.localhost>
- **Painel Admin (Streamlit):** <http://admin.localhost>
- **API Core (Swagger):** <http://api.localhost/swagger/index.html>
- **Chatbot IA (Swagger):** <http://chat.localhost/swagger/index.html>
- **Observabilidade (Grafana LGTM):** <http://telemetry.localhost>
- **Dashboard Traefik:** <http://localhost:8080>

---

## Decisões de projeto

### Arquitetura Hexagonal

A escolha de separar o domínio de frameworks (como o `chi`), bancos de dados e SDKs externos (como o `genai` da Google) garante que o núcleo do sistema (preços, validações, regras de negócio e orquestração de chat) seja 100% testável e independente de infraestrutura.

### Spec-First com OpenAPI

Em vez de escrever handlers na mão, o arquivo `openapi.yaml` é a única fonte da verdade. A ferramenta `oapi-codegen` gerou todas as structs de requisição/resposta e as interfaces (Strict Server). Isso elimina a possibilidade de a documentação e o código ficarem dessincronizados. O contrato também inclui validações como tamanho máximo do `X-Client-ID` e regex para validar IDs do MongoDB.

### Assistente de IA com RAG e PGVector

Para evitar alucinações do modelo de linguagem e garantir respostas alinhadas ao catálogo real do BazaarGo, o microsserviço de chat implementa **Retrieval-Augmented Generation (RAG)**:

1. A pergunta do cliente é convertida em vetor (_embedding_) usando o SDK oficial do **Google Gemini** (`google.golang.org/genai`).
2. O adaptador `pgvector` realiza uma busca por similaridade semântica no PostgreSQL 18 para recuperar os documentos/produtos mais relevantes.
3. Os trechos recuperados são injetados como contexto no prompt enviado ao Gemini, gerando uma resposta precisa e fundamentada.

### Observabilidade de Nível de Produção

Em vez de apenas "prints no terminal", tanto a API principal quanto o serviço de Chatbot foram instrumentados com **OpenTelemetry**.

- O middleware `otelchi` captura traces das requisições HTTP.
- O `otelslog` envia logs estruturados em JSON para o Loki com o `trace_id` atrelado.
- Métricas e traces são exportados via OTLP (gRPC/HTTP) para a stack LGTM, permitindo acompanhar latência de inferência do Gemini, buscas no PGVector e throughput das APIs.
- O **Pyroscope** roda em background capturando perfis de CPU e alocação de memória, permitindo identificar gargalos exatos em tempo real.

### Roteamento com Traefik

Em vez de expor portas soltas (`localhost:8080`, `localhost:9090`, `localhost:8501`), toda a comunicação externa passa pelo Traefik. Isso simula um ambiente de produção real com API Gateway, permitindo escalar containers, rotear múltiplos serviços por subdomínio (`api.localhost`, `chat.localhost`, `front.localhost`, `admin.localhost`) e isolar os bancos de dados na rede interna do Docker.

### Integração com WhatsApp

Para manter o escopo enxuto e sem dependência de gateways de pagamento complexos, o processo de finalização de compra (Checkout) monta uma string estruturada com os itens do carrinho, aplica `url.QueryEscape` e gera um link `wa.me`. É a solução mais direta para e-commerces locais.

### Cache Semântico (Otimização de Custos e Latência)

Em sistemas RAG, a latência e o custo estão diretamente atrelados à chamada do Large Language Model (LLM). Para mitigar isso, não utilizamos um cache tradicional (baseado em strings exatas), mas sim um **Cache Semântico** construído sobre o próprio PGVector.

**Como funciona:**

1. Quando uma pergunta chega, seu _embedding_ é gerado.
2. Antes de buscar o contexto ou chamar o LLM, o sistema realiza uma busca por similaridade de cosseno na tabela `semantic_cache`.
3. Se a pergunta for **≥ 95% similar** a uma pergunta feita anteriormente, a resposta salva é retornada instantaneamente.
4. Se não houver _cache hit_, o fluxo normal do RAG é executado e a nova pergunta/resposta é salva na tabela para consultas futuras.

**Impacto Mensurado (via OpenTelemetry):**

- **Redução de Latência:** O tempo de resposta caiu de **1.55s** (Cache Miss) para **~357ms** (Cache Hit), um ganho de performance de **~4.5x**. O tempo de processamento do LLM (1.22s) é totalmente eliminado.
- **Economia de Custos:** Redução drástica no consumo de tokens da API do Google Gemini.
- **Resiliência:** Perguntas repetidas continuam sendo respondidas mesmo se a API do Gemini estiver fora do ar.

---

## Decisões de Arquitetura (ADRs)

Para garantir a manutenibilidade a longo prazo e a transparência das escolhas técnicas, as decisões arquiteturais mais críticas deste projeto estão documentadas no formato **ADR (Architecture Decision Records)**.

Cada documento detalha o contexto da época, as alternativas avaliadas, a decisão tomada e as suas consequências (trade-offs). O índice completo com o histórico de decisões está disponível em [`docs/adr/`](./docs/adr/README.md).

- [001 - Hexagonal Architecture and Modular Monolith](./docs/adr/001-hexagonal-architecture.md)
- [002 - Contract-First API Design using OpenAPI](./docs/adr/002-openapi-spec-first.md)
- [003 - Polyglot Persistence](./docs/adr/003-polyglot-persistence.md)
- [004 - PASETO v4 instead of JWT](./docs/adr/004-paseto-over-jwt.md)
- [005 - OpenTelemetry and LGTM Stack](./docs/adr/005-observability-lgtm.md)
- [006 - Modular Monolith over Microservices](./docs/adr/006-modular-monolith-over-microservices.md)
- [007 - Semantic Cache for RAG](./docs/adr/007-semantic-cache.md)

---

## Performance e Resiliência

O sistema foi submetido a testes de carga massivos para validar a arquitetura e a capacidade de escalonamento horizontal.

### Escalonamento Horizontal (Load Balancing)

A API foi escalada para **3 instâncias** simultâneas (`docker compose up -d --scale api=3`), com o Traefik atuando como Load Balancer para distribuir o tráfego entre elas.

### Resultados do Teste de Estresse

Utilizando a ferramenta `autocannon` com **10.000 conexões simultâneas**, o sistema apresentou os seguintes resultados:

- **Vazão (Throughput):** Média de **~43.000 requisições por segundo**.
- **Confiabilidade:** **347.693 requisições** processadas com **100% de sucesso** (Status 200), sem nenhuma queda de conexão ou erro 500.
- **Latência (p99):** 1.72s.

### Análise de Gargalos (Profiling)

Os testes foram executados em um ambiente de desenvolvimento local (todas as instâncias e o gerador de carga na mesma máquina física). Através da stack de Observabilidade, foi possível identificar os seguintes gargalos:

- **MongoDB (233% de CPU):** A rota testada utiliza _Regex Scan_ (Case-Insensitive) no catálogo. Como o teste exige leitura extrema, o banco varre a coleção inteira. Em produção, um índice de texto (`Text Index`) resolveria esse gargalo.
- **OpenTelemetry Stack (1.3GB RAM):** A stack LGTM processou a telemetria de 3 instâncias da API simultaneamente sob carga máxima, provando a robustez da coleta de traces e métricas.
- **Go API (Estável):** Mesmo sob estresse, as instâncias da API mantiveram a estabilidade de memória (sem Goroutine Leaks) e uso de CPU focado no processamento concorrente.

---

## Licença

Este projeto está licenciado sob a MIT License. Consulte o arquivo `LICENSE` para mais detalhes.
