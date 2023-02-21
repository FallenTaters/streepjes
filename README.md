# Streepjes

Streepjes is a custom POS app for two specific sports clubs and a hobby project written in pure Go.

The frontend is WASM powered by Vugu, and the backend relies on a sqlite database.

## Development

### Requirements

* go >= 1.18
* go tooling:
    * [enumer](https://github.com/alvaroloes/enumer)
    * [vugugen](https://www.vugu.org/doc/start)

### Run locally

1. `make generate`
2. `make run`
    * re-run after changes

### Build for production

1. `make`

## Settings


The listed values are the default values.

```
STREEPJES_PORT=80
STREEPJES_DB_PATH=streepjes.db
```

### Requirements
* Browser must support webassembly
