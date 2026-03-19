# pismo-account

Simple Transaction Routine service, which manages accounts and transactions associated with those accounts.

## What you need

- Docker
- Docker Compose
- Make

## How to run it

1. Make sure you have a `config/config.json` file (a `config` folder with a `config.json` file inside it).
2. In a terminal, inside this folder, run:

```bash
make up
```

This will:
- start MySQL
- build the Go app


To see logs:

```bash
make logs
```

To stop everything:

```bash
make down
```

@TODO More info to be added