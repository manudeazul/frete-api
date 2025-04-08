# Superfrete API

Olá! Esta é minha resposta para o desafio de criação da API da **Superfrete**.

---

## Dependências

Para rodar a aplicação localmente, você vai precisar de:

- [Go](https://golang.org/dl/) na versão **1.24.2**
- [Docker](https://www.docker.com/)

---

## Como rodar

```bash
# Clone este repositório
git clone https://github.com/manudeazul/frete-api
cd frete-api

# Inicie o módulo Go
go mod init github.com/manudeazul/frete-api

# Baixe os módulos
go mod tidy

# Crie a imagem Docker da API
docker build -t superfrete-api .

# Suba os containers
docker compose up -d
```

Acesse a API em: [http://localhost:8000](http://localhost:8000)

---

## Rotas

### `POST /quote`

Realiza uma cotação fictícia com a API da Superfrete e salva no banco de dados PostgreSQL.

#### Exemplo de entrada:

```json
{
  "recipient": {
    "address": {
      "zipcode": "01311000"
    }
  },
  "volumes": [
    {
      "category": 7,
      "amount": 1,
      "unitary_weight": 5,
      "price": 349,
      "sku": "abc-teste-123",
      "height": 0.2,
      "width": 0.2,
      "length": 0.2
    },
    {
      "category": 7,
      "amount": 2,
      "unitary_weight": 4,
      "price": 556,
      "sku": "abc-teste-527",
      "height": 0.4,
      "width": 0.6,
      "length": 0.15
    }
  ]
}
```

#### Exemplo de saída:

```json
[
  {
    "name": "BOX DELIVERY",
    "service": "API_BOX",
    "deadline": 24,
    "price": 0
  },
  {
    "name": "BOX DELIVERY",
    "service": "API_BOX",
    "deadline": 24,
    "price": 0
  }
]
```

---

### `GET /metrics?last_quotes=<num>`

Consulta as últimas cotações salvas no banco e retorna:

- Quantidade de cotações por transportadora  
- Total do preço de frete (`final_price`) por transportadora  
- Média do preço de frete por transportadora  
- Frete mais barato geral  
- Frete mais caro geral  

> `last_quotes` é um parâmetro **opcional** que permite filtrar pelas últimas N cotações.

#### Exemplo de saída:

```json
{
  "carriers": [
    {
      "metrics": [
        {
          "name": "BOX DELIVERY",
          "count": 2,
          "total_price": 0,
          "average_price": 0
        }
      ],
      "cheapest_quote": {
        "name": "BOX DELIVERY",
        "service": "API_BOX",
        "deadline": 24,
        "price": 0
      },
      "higher_quote": {
        "name": "BOX DELIVERY",
        "service": "API_BOX",
        "deadline": 24,
        "price": 0
      }
    }
  ]
}
```
