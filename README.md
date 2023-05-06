<p align="center" style="text-align:center;">
    <img alt="Correios CEP Admin Logo" src="docs/assets/logo.svg" width="200" />  
</p>

# Correios CEP Admin (API)

Within the [Open Transparency ecosystem](https://github.com/insighted4/), this repository
is responsible for proxying and caching addresses from Brazilian Correios [Busca por CEP](https://buscacepinter.correios.com.br/app/cep/index.php).

**Note**: The `master` branch may be in an *unstable or even broken state* during development. Please use [releases][github-release]
instead of the `master` branch in order to get stable binaries.

Quick Start
---

#### Requirements

* Go >= 1.19
* Docker

```bash
# Clone repository
$ git clone https://github.com/insighted4/correios-cep.git
$ cd correios-cep

# Make sure to edit the database address, log level, and others.
# For more options: ./cmd/admin/serve.go
$ cp .env.sample .env

# Build & run
$ docker-compose up -d
$ curl -s -XGET http://localhost:8080/api/v1/addresses/74001970 | jq .
{
    "cep": "74001970",
    "state": "GO",
    "city": "Goiânia",
    "neighborhood": "Setor Central",
    "location": "Praça Doutor Pedro Ludovico Teixeira, 11",
    "children": null,
    "created_at": "2023-05-06T20:05:03.340664587Z",
    "updated_at": "2023-05-06T20:05:03.340664587Z"
}
```

#### Running the live server

```bash
# Make sure to edit the database address, log level, and others.
# For more options: ./cmd/admin/serve.go
$ cp .env.sample .env

$ make build
$ ./bin/admin serve
```

#### Unit Tests

```bash
$ make test
```

### Integration with Postgres
```bash
# Please create a separate database for testing with Postgres.
$ psql -d postgresql://postgres:secret@localhost:5432/postgres -c "CREATE DATABASE cep_test;"

# Run the migration schema
$ psql -d postgresql://postgres:secret@localhost:5432/cep_test -f ./migration/schema.sql

# Make sure to export TEST_DATABASE_URL to enable testing the storage package. 
# You might want to include TEST_DATABASE_URL to your .env file for better convenience.
$ export TEST_DATABASE_URL=postgresql://postgres:secret@localhost:5432/cep_test?sslmode=disable
$ make test
```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

### License

This repository is under the AGPL 3.0 license. See the [LICENSE](LICENSE) file for details.
