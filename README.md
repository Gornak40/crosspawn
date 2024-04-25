# crosspawn

*Cross review microservice developed for the Ejudge ecosystem.*

## Config

Create `config.ini` file in project directory.

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
POLL_BATCH_SIZE = 50
POLL_DELAY_SECONDS = 10
```

## Usage

Generate personal admin `jwt` with `./jwtsign.sh` script.

Run poller with `make run-poller` command.

Run server with `make run-crosspawn` command.
