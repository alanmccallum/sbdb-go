sbdb-go is a Go library for working with NASA JPL's Small-Body Database (SBDB) Query API. It decodes JSON responses into convenient Go structures so you can inspect results programmatically.

## Installation

```
go get github.com/alanmccallum/sbdb-go
```

## Quickstart

```go
resp, err := http.Get("https://ssd-api.jpl.nasa.gov/sbdb_query.api?limit=1&fields=spkid,full_name")
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
```

## Usage

`sbdb.Decode` reads a JSON payload and returns a `Payload` containing the raw data. Use `Payload.Records` to get a slice of generic map-based records or `Payload.Bodies` to populate the strongly typed `Body` struct.

Constants such as `sbdb.SpkID`, `sbdb.NEO`, and others mirror the field names used by the SBDB API. These can be helpful when constructing queries or inspecting `Record` values.

For additional examples see the package documentation on [pkg.go.dev](https://pkg.go.dev/github.com/alanmccallum/sbdb-go).

## Data Source

This project uses publicly available data from NASA JPL's Small-Body Database API. Data is provided by the [Jet Propulsion Laboratory](https://ssd-api.jpl.nasa.gov/) under U.S. Government public domain.

## License

This library is released under the MIT License. See [LICENSE](LICENSE) for details.