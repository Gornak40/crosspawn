#!/usr/bin/bash
go build -o ./bin ./cmd/jwtsign && ./bin/jwtsign $@
