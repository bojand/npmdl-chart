[![Go Report Card](https://goreportcard.com/badge/github.com/bojand/npmdl-chart)](https://goreportcard.com/report/github.com/bojand/npmdl-chart)

# npmdl

Plot NPM download counts over time 

## Usage

```
dep ensure
go run *.go
```

Then browse to http://localhost:8080 or http://localhost:8080/express

Example chart:

[![Express downloads over time](http://npmdl.bojan.codes/chart/express.svg)](https://npmdl.bojan.codes/chart/express)

## Why?

Yes this could be built with just frontend JavaScript.

Just experimenting with [Go](https://golang.org/).

Also this provides a mechanism for embedding the NPM download chart similar to [https://nodei.co](https://nodei.co/) which does not work any more. Plus provides some additional options.

## License

Apache-2.0
