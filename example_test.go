package sbdb_test

import (
	"fmt"
	"strings"

	"github.com/alanmccallum/sbdb-go"
)

func ExampleDecode() {
	data := `{"fields":["spkid","full_name","neo"],"data":[[1234,"(1234) Example","Y"]]}`
	p, err := sbdb.Decode(strings.NewReader(data))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	bodies, err := p.Bodies()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(*bodies[0].Identity.FullName)
	// Output: (1234) Example
}

func ExamplePayload_Records() {
	data := `{"fields":["spkid","neo"],"data":[[555,"Y"]]}`
	p, _ := sbdb.Decode(strings.NewReader(data))
	records, _ := p.Records()
	fmt.Println(records[0][sbdb.SpkID])
	fmt.Println(records[0][sbdb.NEO])
	// Output:
	// 555
	// Y
}

func ExampleClient_GetURL() {
	c := &sbdb.Client{}
	f := sbdb.Filter{
		Fields: sbdb.NewFieldSet(sbdb.SpkID, sbdb.FullName),
		Limit:  1,
	}
	u, err := c.GetURL(f)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(u.String())
	// Output:
	// https://ssd-api.jpl.nasa.gov/sbdb_query.api?fields=full_name%2Cspkid&limit=1
}

func ExampleFilter_Values() {
	f := sbdb.Filter{
		Fields: sbdb.NewFieldSet(sbdb.SpkID, sbdb.FullName),
		Limit:  5,
		Kind:   sbdb.KindAsteroid,
	}
	v, err := f.Values()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(v.Encode())
	// Output:
	// fields=full_name%2Cspkid&limit=5&sb-kind=a
}
