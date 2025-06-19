sbdb-go is a Go library for interacting with NASA JPL's Small-Body Database (SBDB) Query API. It provides a query builder and decoders so you can inspect results programmatically.

### Features
- Build SBDB queries using a typed `Filter` API
- Decode responses into rich Go structs
- Works with the [SBDB Query API](https://ssd-api.jpl.nasa.gov/doc/sbdb_query.html)
- Supports advanced field filtering as described in the [SBDB filter documentation](https://ssd-api.jpl.nasa.gov/doc/sbdb_filter.html)

## Installation

```
go get github.com/alanmccallum/sbdb-go
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alanmccallum/sbdb-go"
)

func main() {
	c := &sbdb.Client{}
	f := sbdb.Filter{
		Fields: sbdb.NewFieldSet(sbdb.SpkID, sbdb.FullName),
		Limit:  1,
	}
	resp, err := c.Get(context.Background(), f)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	p, err := sbdb.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodies, err := p.Bodies()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*bodies[0].Identity.FullName)
}

```

## Usage

`sbdb.Decode` reads a JSON payload and returns a `Payload` containing the raw data. Use `Payload.Records` to get a slice of generic map-based records or `Payload.Bodies` to populate the strongly typed `Body` struct.

The `Filter` type and helper functions allow you to build complex queries in Go. Field names mirror those documented by the [SBDB Query API](https://ssd-api.jpl.nasa.gov/doc/sbdb_query.html) and [filter syntax](https://ssd-api.jpl.nasa.gov/doc/sbdb_filter.html).

Constants such as `sbdb.SpkID`, `sbdb.NEO`, and others mirror the field names used by the SBDB API. These can be helpful when constructing queries or inspecting `Record` values.

For additional examples see the package documentation on [pkg.go.dev](https://pkg.go.dev/github.com/alanmccallum/sbdb-go).

## Debug Logging

`sbdb-go` uses Go's `slog` package for optional debug output of type conversions. Replace the
package logger using `SetLogger` to enable these messages. The tests show an
example of configuring a text logger at the debug level in
[`decode_test.go`](decode_test.go):

```go
func TestMain(m *testing.M) {
    SetLogger(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    })))
    os.Exit(m.Run())
}
```


## Data Source

This project uses publicly available data from NASA JPL's Small-Body Database API. Data is provided by the [Jet Propulsion Laboratory](https://ssd-api.jpl.nasa.gov/) under U.S. Government public domain.

## License

This library is released under the MIT License. See [LICENSE](LICENSE) for details.