
<h4 align="center">
    <p>
        <b>English</b> |
        <a href="./README.md">English</a>
    </p>
</h4>



<p align="center">
  <img src="https://i.imgur.com/orlCpQX.png" alt="Imagem logo" />
</p>

![GitHub license](https://img.shields.io/github/license/AlexWilliam12/silent-signal-backend)
![GitHub languages top](https://img.shields.io/github/languages/top/AlexWilliam12/silent-signal-backend)
![GitHub last commit](https://img.shields.io/github/last-commit/AlexWilliam12/silent-signal-backend)

# Silent Signal

Silent Signal é um aplicativo de mensagens seguro e privado, projetado para proteger as comunicações dos usuários. O servidor é construído com Go (Golang) e utiliza PostgreSQL para armazenamento de dados.

## Funcionalidades

- Criptografia de ponta a ponta
- Autenticação segura de usuários
- Armazenamento de dados criptografado
- Interface amigável

## Requisitos

- Go (Golang) 1.20 ou superior
- PostgreSQL 12 ou superior
- Make

## Instalação

1. Clone o repositório

```bash
git clone https://github.com/AlexWilliam12/silent-signal-backend.git
```
2. Mude para o diretório do projeto
    
    ```bash
    cd silent-signal-backend
    ```

## Configuração do PostgreSQL

Certifique-se de que o PostgreSQL está instalado e em execução. Crie um banco de dados para o Silent Signal:

```bash
psql -U postgres
CREATE DATABASE silentsignal;
\q
```

## Configurar Variáveis de Ambiente

Crie um arquivo .env no diretório raiz do projeto e adicione a seguinte configuração:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=suasenha
DB_NAME=silentsignal
```

## Instalar Dependências do Go

```bash
make install
```
## Executar Migrações

```bash
make migrate
```

## Construir e Executar o Servidor

```bash
make run
```
e 
```bash
make build
```

O servidor deve estar em execução em http://localhost:8080.

## Comandos do Makefile 

* make install: Instala as dependências do Go.
* make migrate: Executa as migrações do banco de dados.
* make build: Constrói a aplicação Go.
* make run: Executa a aplicação Go.

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma issue ou enviar um pull request.

## Licença

Silent Signal é distribuído sob a licença MIT.



