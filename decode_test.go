package sbdb

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	os.Exit(m.Run())
}

//func TestDecode(t *testing.T) {
//	type args struct {
//		r io.Reader
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *Payload
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := Decode(tt.args.r)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Decode() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPayload_Bodies(t *testing.T) {
//	type fields struct {
//		Signature struct {
//			Version string `json:"version"`
//			Source  string `json:"source"`
//		}
//		Fields []string
//		Data   [][]any
//		Count  int
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    []Body
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := &Payload{
//				Signature: tt.fields.Signature,
//				Fields:    tt.fields.Fields,
//				Data:      tt.fields.Data,
//				Count:     tt.fields.Count,
//			}
//			got, err := p.Bodies()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Bodies() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Bodies() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPayload_Records(t *testing.T) {
//	type fields struct {
//		Signature struct {
//			Version string `json:"version"`
//			Source  string `json:"source"`
//		}
//		Fields []string
//		Data   [][]any
//		Count  int
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    []Record
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := &Payload{
//				Signature: tt.fields.Signature,
//				Fields:    tt.fields.Fields,
//				Data:      tt.fields.Data,
//				Count:     tt.fields.Count,
//			}
//			got, err := p.Records()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Records() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Records() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestRecord_getBool(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		r    Record
		args args
		want *bool
	}{
		{
			name: "String",
			r: Record{
				"string": "Hello World",
			},
			args: args{"string"},
			want: nil,
		},
		{
			name: "Float",
			r: Record{
				"float": 3.14159,
			},
			args: args{"float"},
			want: nil,
		},
		{
			name: "Int",
			r: Record{
				"int": 3,
			},
			args: args{"int"},
			want: nil,
		},
		{
			name: "Bool True",
			r: Record{
				"bool": true,
			},
			args: args{"bool"},
			want: ptrTo(true),
		},
		{
			name: "Bool Y",
			r: Record{
				"bool": "Y",
			},
			args: args{"bool"},
			want: ptrTo(true),
		},
		{
			name: "Bool T",
			r: Record{
				"bool": "T",
			},
			args: args{"bool"},
			want: ptrTo(true),
		},
		{
			name: "Bool N",
			r: Record{
				"bool": "N",
			},
			args: args{"bool"},
			want: ptrTo(false),
		},
		{
			name: "Bool F",
			r: Record{
				"bool": "F",
			},
			args: args{"bool"},
			want: ptrTo(false),
		},
		{
			name: "Nil",
			r: Record{
				"bool": nil,
			},
			args: args{"bool"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.getBool(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_getFloat(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		r    Record
		args args
		want *float64
	}{
		{
			name: "String",
			r: Record{
				"string": "Hello World",
			},
			args: args{"string"},
			want: nil,
		},
		{
			name: "Float String",
			r: Record{
				"string": "3.14159",
			},
			args: args{"string"},
			want: ptrTo(3.14159),
		},
		{
			name: "JSON Number (float)",
			r: Record{
				"num": json.Number("3.14159"),
			},
			args: args{"num"},
			want: ptrTo(3.14159),
		},
		{
			name: "JSON Number (int)",
			r: Record{
				"num": json.Number("3"),
			},
			args: args{"num"},
			want: ptrTo(3.0),
		},
		{
			name: "Bool",
			r: Record{
				"bool": true,
			},
			args: args{"bool"},
			want: nil,
		},
		{
			name: "Nil",
			r: Record{
				"nil": nil,
			},
			args: args{"nil"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.getFloat(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_getInt(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		r    Record
		args args
		want *int
	}{
		{
			name: "String",
			r: Record{
				"string": "Hello World",
			},
			args: args{"string"},
			want: nil,
		},
		{
			name: "Int String",
			r: Record{
				"string": "3",
			},
			args: args{"string"},
			want: ptrTo(3),
		},
		{
			name: "JSON Number (int)",
			r: Record{
				"num": json.Number("3"),
			},
			args: args{"num"},
			want: ptrTo(3),
		},
		{
			name: "JSON Number (float)",
			r: Record{
				"num": json.Number("3.14159"),
			},
			args: args{"num"},
			want: ptrTo(3),
		},
		{
			name: "Bool",
			r: Record{
				"bool": true,
			},
			args: args{"bool"},
			want: nil,
		},
		{
			name: "Nil",
			r: Record{
				"nil": nil,
			},
			args: args{"nil"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.getInt(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_getString(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		r    Record
		args args
		want *string
	}{
		{
			name: "Hello World",
			r: Record{
				"string": "Hello World",
			},
			args: args{"string"},
			want: ptrTo("Hello World"),
		},
		{
			name: "JSON Number (int)",
			r: Record{
				"num": json.Number("3"),
			},
			args: args{"num"},
			want: ptrTo("3"),
		},
		{
			name: "JSON Number (float)",
			r: Record{
				"num": json.Number("3.14159"),
			},
			args: args{"num"},
			want: ptrTo("3.14159"),
		},
		{
			name: "Float",
			r: Record{
				"float": 3.14159,
			},
			args: args{"float"},
			want: ptrTo("3.14159"),
		},
		{
			name: "Int",
			r: Record{
				"int": 3,
			},
			args: args{"int"},
			want: ptrTo("3"),
		},
		{
			name: "Bool",
			r: Record{
				"bool": true,
			},
			args: args{"bool"},
			want: ptrTo("true"),
		},
		{
			name: "Nil",
			r: Record{
				"nil": nil,
			},
			args: args{"nil"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.getString(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getString() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestRecord_identity(t *testing.T) {
//	tests := []struct {
//		name string
//		r    Record
//		want Identity
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.r.identity(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("identity() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestRecord_nonGrav(t *testing.T) {
	tests := []struct {
		name string
		r    Record
		want NonGrav
	}{
		{
			name: "mixed types",
			r: Record{
				A1:      1.23,
				A2:      "2.34",
				A3:      nil,
				DT:      4,
				S0:      "5.6",
				A1Sigma: 0.1,
				A2Sigma: "0.2",
				A3Sigma: nil,
				DTSigma: "0.4",
				S0Sigma: 0.5,
			},
			want: NonGrav{
				A1:      ptrTo(1.23),
				A2:      ptrTo(2.34),
				A3:      nil,
				DT:      ptrTo(4.0),
				S0:      ptrTo(5.6),
				A1Sigma: ptrTo(0.1),
				A2Sigma: ptrTo(0.2),
				A3Sigma: nil,
				DTSigma: ptrTo(0.4),
				S0Sigma: ptrTo(0.5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.nonGrav(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nonGrav() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestRecord_orbit(t *testing.T) {
//	tests := []struct {
//		name string
//		r    Record
//		want Orbit
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.r.orbit(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("orbit() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestRecord_physical(t *testing.T) {
	tests := []struct {
		name string
		r    Record
		want Physical
	}{

		{
			name: "mixed types",
			r: Record{
				H:             1.1,
				G:             "0.15",
				M1:            nil,
				K1:            0,
				M2:            2.2,
				K2:            "1.1",
				PC:            0.5,
				HSigma:        "0.1",
				Diameter:      100,
				Extent:        "10x20",
				GM:            "123.4",
				Density:       3.0,
				RotPer:        7,
				Pole:          "90,0",
				Albedo:        0.1,
				BV:            "0.2",
				UB:            0.3,
				IR:            "0.4",
				SpecT:         "S",
				SpecB:         "B",
				DiameterSigma: "0.7",
			},
			want: Physical{
				H:             ptrTo(1.1),
				G:             ptrTo(0.15),
				M1:            nil,
				K1:            ptrTo(0.0),
				M2:            ptrTo(2.2),
				K2:            ptrTo(1.1),
				PC:            ptrTo(0.5),
				HSigma:        ptrTo(0.1),
				Diameter:      ptrTo(100.0),
				Extent:        ptrTo("10x20"),
				GM:            ptrTo(123.4),
				Density:       ptrTo(3.0),
				RotPer:        ptrTo(7.0),
				Pole:          ptrTo("90,0"),
				Albedo:        ptrTo(0.1),
				BV:            ptrTo(0.2),
				UB:            ptrTo(0.3),
				IR:            ptrTo(0.4),
				SpecT:         ptrTo("S"),
				SpecB:         ptrTo("B"),
				DiameterSigma: ptrTo(0.7),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.physical(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("physical() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestRecord_quality(t *testing.T) {
//	tests := []struct {
//		name string
//		r    Record
//		want Quality
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.r.quality(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("quality() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRecord_solution(t *testing.T) {
//	tests := []struct {
//		name string
//		r    Record
//		want Solution
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.r.solution(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("solution() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRecord_uncertainty(t *testing.T) {
//	tests := []struct {
//		name string
//		r    Record
//		want Uncertainty
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.r.uncertainty(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("uncertainty() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func FuzzBodies(f *testing.F) {
	f.Add([]byte(`{"fields":["spkid","full_name","neo","t_jup"],"data":[[1234,"name","Y","3.14"]]}`)) // realistic, full
	f.Add([]byte(`{"fields":["spkid"],"data":[[1234]]}`))                                             // minimal int
	f.Add([]byte(`{"fields":["neo"],"data":[["T"]]}`))                                                // minimal bool
	f.Add([]byte(`{"fields":["t_jup"],"data":[["NaN"]]}`))
	f.Add([]byte(`{"fields":["t_jup"],"data":[[1.0]]}`)) // edge-case float
	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic for input %q: %v", data, r)
			}
		}()

		r := bytes.NewReader(data)
		p, err := Decode(r)
		if err == nil && p != nil {
			_, _ = p.Bodies()
		}
	})
}

func ptrTo[T any](v T) *T {
	return &v
}
