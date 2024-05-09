# CrossPawn

*Cross review microservice developed for the Ejudge ecosystem.*

## Config

Create `config.ini` in project directory.

```ini
[ejudge]
API_KEY = <KEY>
API_SECRET = <SECRET>
URL = https://ejudge.algocourses.ru

[database]
SQLITE_PATH = <PATH>

[server]
GIN_SECRET = <SECRET>
JWT_SECRET = <SECRET>
POLL_DELAY_SECONDS = 10
REVIEW_LIMIT = 3
```

## Build

```bash
make
```

## Generate JWT

You can generate personal token for admin user.

```bash
./jwtsign.sh --user gorilla --duration 24h
```

## Usage

Start poller.

```bash
make run-poller
```

Start server.

```bash
make run-crossspawn
```
