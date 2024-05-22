
<h4 align="center">
    <p>
        <b>English</b> |
        <a href="./README_PT-br.md">Рortuguês</a>
    </p>
</h4>



<p align="center">
  <img src="https://i.imgur.com/orlCpQX.png" alt="Imagem logo" />
</p>


# Silent Signal 

Silent Signal is a secure and private messaging app designed to protect user communications. The server is built with Go (Golang) and uses PostgreSQL for data storage.

## Features

- End-to-end encryption
- Secure user authentication
- Encrypted data storage
- User-friendly interface

## Requirements

- Go (Golang) 1.20 or later
- PostgreSQL 12 or later
- Make

## Installation

1. Clone the repository

```bash
git clone https://github.com/AlexWilliam12/silent-signal-backend.git
```

2. Change to the project directory

```bash
cd silent-signal-backend
```

### Set Up PostgreSQL

Ensure PostgreSQL is installed and running. Create a database for Silent Signal:

```bash
psql -U postgres
CREATE DATABASE silentsignal;
\q
```

### Configure Environment Variables

Create a .env file in the root directory of the project and add the following configuration:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=silentsignal
```

### Install Go Dependencies

```bash
make install
```
### Run Migrations

```bash
make migrate
```

### Build and Run the Server

```bash
make run
```
and 
```bash
make build
```

The server should now be running on http://localhost:8080.


## Makefile Commands

* make install: Installs Go dependencies.
* make migrate: Runs database migrations.
* make build: Builds the Go application.
* make run: Runs the Go application.

## Contributing

We welcome contributions from the community.

## License
Silent Signal is licensed under the MIT License. 
