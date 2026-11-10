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
    - [Infraestrutura & Autenticação (Público)](#infraestrutura--autenticação-público)
    - [Catálogo (Público)](#catálogo-público)
    - [Administração (Bearer Token)](#administração-bearer-token)
    - [Checkout (Público)](#checkout-público)
  - [Como executar](#como-executar)
    - [Requisitos](#requisitos)
    - [Passos](#passos)
  - [Decisões de projeto](#decisões-de-projeto)
    - [Arquitetura Hexagonal](#arquitetura-hexagonal)
    - [Spec-First com OpenAPI](#spec-first-com-openapi)
    - [Observabilidade de Nível de Produção](#observabilidade-de-nível-de-produção)
    - [Persistência Poliglota](#persistência-poliglota)
    - [Integração com WhatsApp](#integração-com-whatsapp)
  - [Performance e Resiliência](#performance-e-resiliência)
  - [Licença](#licença)
<!--toc:end-->

E-commerce Fullstack (Catálogo de Produtos) com arquitetura hexagonal, persistência poliglota e observabilidade distribuída.

---

## Visão geral

O **BazaarGo** é uma plataforma de e-commerce que permite o cadastro e a listagem de produtos, gerenciamento de carrinho de compras e finalização de pedidos via WhatsApp. O projeto foi construído para demonstrar, na prática, a aplicação do conceito de **Banco de Dados Poliglota** (utilizando a tecnologia certa para cada tipo de dado) aliado a uma arquitetura backend moderna e altamente observável.

O sistema atende a dois módulos principais:

1. **Módulo Administrativo:** Onde administradores autenticados cadastram e gerenciam o catálogo de produtos.
2. **Módulo Cliente:** Uma vitrine rápida que exibe os produtos, registra o histórico de "vistos recentemente" e permite fechar o pedido gerando uma mensagem formatada no WhatsApp.

> **Nota de transparência** : Este projeto utilizou ferramentas de Inteligência Artificial (IA) como assistente de desenvolvimento durante todo o processo. A IA foi utilizada para acelerar a construção das interfaces (SvelteKit e Streamlit), mas também atuou diretamente como um "copiloto" técnico: auxiliando na configuração da infraestrutura (Docker, Traefik, stack LGTM), instrumentação da observabilidade (OpenTelemetry), resolução de bugs complexos e na discussão de tomadas de decisão sobre padrões de projeto e arquitetura.

---

## Diferenciais do projeto

Além dos requisitos mínimos, o projeto foi construído seguindo práticas adotadas em ambientes de produção corporativos:

- **Persistência Poliglota:** Uso de PostgreSQL (Supabase) para autenticação (ACID), MongoDB para catálogo (flexibilidade de schema) e Valkey/Redis para cache (performance).
- **Observabilidade Completa (LGTM + Pyroscope):** Rastreamento distribuído (Tempo), agregação de logs estruturados (Loki), coleta de métricas (Prometheus/Mimir) e profiling contínuo de CPU/Memória (Pyroscope).
- **Arquitetura Hexagonal (Ports and Adapters):** Domínio puro e isolado de dependências externas, facilitando testes e troca de tecnologia.
- **Spec-First com OpenAPI (`oapi-codegen`):** A API é gerada a partir do contrato `openapi.yaml`, garantindo *Strict Server* (tipagem forte e validação automática de requests/responses).
- **Reverse Proxy com Traefik:** Roteamento inteligente e balanceamento de carga nativo para os containers Docker.
- **Autenticação com PASETO v4:** Substituindo o JWT por um padrão mais seguro contra ataques de análise de tempo e confusão de algoritmos.
- **Hash de senha com Argon2id:** Configurado com parâmetros acima das recomendações mínimas do OWASP.
- **Infraestrutura 100% Containerizada:** Orquestração completa via Docker Compose.

---

## Stack de tecnologias

| Camada | Tecnologia |
| --- | --- |
| **Backend** | |
| Linguagem | Go (Golang) |
| Roteador HTTP | chi |
| Design de API | OpenAPI 3.0 (Spec-First + Strict Server) |
| Autenticação | PASETO v4 (aidantwoods) |
| Hash de senha | Argon2id |
| **Bancos de Dados** | |
| Relacional (Auth) | PostgreSQL 18 (Supabase) |
| Documentos (Produtos) | MongoDB 9 |
| Cache/Listas (Histórico) | Valkey 9 (Redis) |
| **Observabilidade** | |
| Coleta e Visualização | Grafana, Tempo, Loki, Prometheus, Pyroscope |
| Instrumentação | OpenTelemetry (SDK Go + OTLP gRPC/HTTP) |
| **Frontend & UI** | |
| Vitrine (Cliente) | SvelteKit (Svelte 5) + Tailwind CSS |
| Painel Admin | Python + Streamlit |
| **Infraestrutura** | |
| Proxy Reverso | Traefik v3 |
| Containerização | Docker e Docker Compose |

---

## Arquitetura

O projeto adota uma arquitetura desacoplada, onde o backend segue os princípios da **Arquitetura Hexagonal** e os frontends atuam como consumidores puros da API REST.

```
                        ┌─────────────────────────────┐
                        │           Traefik           │
                        │ (Reverse Proxy / Roteamento)│
                        └──────┬───────┬──────┬───────┘
                               │       │      │
                        ┌──────▼──┐    │   ┌──▼──────┐
                        │SvelteKit│    │   │Streamlit│
                        │(Vitrine)│    │   │(Admin)  │
                        └────┬────┘    │   └────┬────┘
                             │         │        │
                             └─────┐   │   ┌────┘
                                   ▼  ▼   ▼
                               ┌─────────────┐
                               │   Go API    │
                               │ (Hexagonal) │
                               └──────┬──────┘
                   ┌──────────────────┼──────────────────┐
                   │                  │                  │
            ┌──────▼──────┐    ┌──────▼──────┐    ┌──────▼──────┐
            │ PostgreSQL  │    │   MongoDB   │    │   Valkey    │
            │(Auth/Admin) │    │ (Produtos)  │    │   (Cache)   │
            └─────────────┘    └─────────────┘    └─────────────┘
```

### Estrutura do projeto

```
BazaarGo/
├── cmd/
│   └── api/                 # Entrypoint da aplicação (Go)
├── internal/
│   ├── domain/              # Entidades (User, Product) e regras puras
│   ├── ports/               # Interfaces (Repositories, Cache, Auth)
│   ├── service/             # Casos de uso (UserService, ProductService)
│   ├── telemetry/           # Setup do OpenTelemetry (Traces, Metrics, Logs)
│   └── adapters/
│       ├── inbound/http/    # Handlers gerados pelo OpenAPI, middlewares e rotas
│       └── outbound/        # Implementações (Postgres, MongoDB, Valkey, PASETO, Argon2)
├── api/
│   └── openapi.yaml         # A "Fonte da Verdade" do contrato da API
├── database/
│   └── migrations/          # Scripts SQL para o PostgreSQL
├── frontend/                # Aplicação SvelteKit (Vitrine do Cliente)
├── admin_panel/             # Aplicação Streamlit (Painel Administrativo em Python)
├── docker-compose.yml       # Orquestração completa (API, Bancos, Observabilidade, Traefik)
└── README.md
```

---

## Banco de dados (Persistência Poliglota)

A grande estrela deste projeto é a escolha cirúrgica do banco de dados para cada problema:

### SGBDs Utilizados

- **PostgreSQL (Supabase):** Utilizado para o controle de **Autenticação e Usuários**. Sua natureza relacional e as garantias ACID garantem a integridade no cadastro de administradores.
- **MongoDB:** Utilizado para o **Catálogo de Produtos**. Seu modelo de documentos flexível permite que produtos tenham atributos dinâmicos (`Specs: map[string]any`), o que seria um pesadelo de normalização no SQL.
- **Valkey (Redis):** Utilizado para o **Cache** de produtos consultados e para o histórico de **"Vistos Recentemente"**. Utiliza estruturas de dados em memória para entregar respostas em menos de 1 milissegundo.

### Modelagem de Dados

- **Users (Postgres):** Tabela relacional contendo `id`, `email`, `password_hash`, etc.
- **Products (MongoDB):** Documento contendo `_id` (ObjectID), `name`, `price`, `stock` e `specs` (objeto flexível).
- **Recent Views (Valkey):** Utilização da estrutura de dados **List** (`LPUSH`, `LTRIM`), mantendo apenas os últimos 5 produtos vistos por cada visitante (identificado via Header `X-Client-ID`).

---

## Endpoints da API

A documentação completa e interativa está disponível via Swagger UI, mas o resumo das rotas é:

### Infraestrutura & Autenticação (Público)

| Método | Rota | Descrição |
| --- | --- | --- |
| POST | `/auth/token` | Login e geração do token PASETO |
| POST | `/auth/register` | Criação de um novo Administrador |

### Catálogo (Público)

| Método | Rota | Descrição |
| --- | --- | --- |
| GET | `/products/search/{name}` | Busca produtos por nome (Regex no MongoDB) |
| GET | `/products/{id}` | Retorna detalhes de um produto e registra no Redis |
| GET | `/products/recent` | Retorna os produtos vistos recentemente (via Redis) |
| POST | `/checkout` | Recebe o carrinho e gera a string formatada para o WhatsApp |

### Administração (Bearer Token)

| Método | Rota | Descrição |
| --- | --- | --- |
| POST | `/products` | Cadastra um novo produto (requer token) |
| PUT | `/products/{id}` | Atualiza um produto (requer token) |
| DELETE | `/products/{id}` | Deleta um produto (requer token) |

---

## Como executar

### Requisitos

- Docker e Docker Compose
- Make (Opcional)

### Passos

**1. Clone o repositório**

```bash
git clone https://github.com/seu-usuario/bazaargo.git
cd bazaargo
cp .env.example .env
```

**2. Suba a stack completa (Traefik, Go API, Bancos, Grafana, etc.)**

```bash
docker compose up -d --build
```

**3. Acesse as aplicações via Traefik:**
O Traefik gerencia as rotas baseadas em subdomínios locais. Certifique-se de que sua máquina resolve `*.localhost` automaticamente (padrão na maioria dos SOs modernos).

- **Vitrine (Cliente):** `http://app.localhost` (ou a porta mapeada para o Svelte)
- **API (Swagger):** `http://api.localhost/swagger/index.html`
- **Painel Admin (Streamlit):** `http://admin.localhost` (ou a porta mapeada)
- **Observabilidade (Grafana):** `http://localhost:3000` (admin/admin)
- **Dashboard Traefik:** `http://localhost:8080`

---

## Decisões de projeto

### Arquitetura Hexagonal

A escolha de separar o domínio de frameworks (como o `chi`) e de bibliotecas de banco de dados garante que o núcleo do sistema (preços, validações, regras) seja 100% testável e independente de infraestrutura.

### Spec-First com OpenAPI

Em vez de escrever handlers na mão, o arquivo `openapi.yaml` é a única fonte da verdade. A ferramenta `oapi-codegen` gerou todas as structs de requisição/resposta e as interfaces (Strict Server). Isso elimina a possibilidade de a documentação e o código ficarem dessincronizados.

### Observabilidade de Nível de Produção

Em vez de apenas "prints no terminal", a API foi instrumentada com **OpenTelemetry**.

- O middleware `otelchi` captura traces das requisições HTTP.
- O `otelslog` envia logs estruturados em JSON para o Loki com o `trace_id` atrelado.
- O **Pyroscope** roda em background capturando perfis de CPU e alocação de memória, permitindo identificar gargalos exatos em tempo real.

### Integração com WhatsApp

Para manter o escopo enxuto e sem dependência de gateways de pagamento complexos, o processo de finalização de compra (Checkout) monta uma string estruturada com os itens do carrinho (`strings.Builder`), aplica `url.QueryEscape` e gera um link `wa.me`. É a solução mais direta para e-commerces locais.

---

## Performance e Resiliência

O backend foi submetido a testes de carga massivos (utilizando `autocannon`) para validar a arquitetura.

- **Vazão (Throughput):** A API sustentou uma média de **~45.000 requisições por segundo** com 10.000 conexões simultâneas, retornando status 200.
- **Eficiência de CPU:** Durante o pico de testes, a aplicação Go consumiu **menos de 6% de um único núcleo de CPU**, provando a eficiência do modelo de concorrência do Go e do uso do cache no Redis.
- **Estabilidade de Memória:** Os gráficos de *Goroutines* e *Heap Allocation* no Grafana mantiveram-se estáveis (padrão dente de serra), provando a ausência de vazamentos de memória (Memory/Goroutine Leaks).

---

## Licença

Este projeto está licenciado sob a MIT License. Consulte o arquivo `LICENSE` para mais detalhes.
