# Cryptopals challenges

[![build & test](https://github.com/racsoraul/cryptopals/actions/workflows/main.yml/badge.svg)](https://github.com/racsoraul/cryptopals/actions/workflows/main.yml)

Solutions written in Go for https://cryptopals.com challenges.

## Run and validate challenges
`go test ./...`

## Generate files
The english distribution map is generated dynamically from the "Frankenstein" book. If you need to regenerate it run:  
`go generate ./...`

If you want to generate the distribution using a different book or set of letters check [gen.go](https://github.com/racsoraul/cryptopals/blob/master/set/one/gen.go) to understand the program used run by `go:generate` in [one.go, line 14](https://github.com/racsoraul/cryptopals/blob/master/set/one/one.go#L14).
