package sbdb

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	SetLogger(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	os.Exit(m.Run())
}

func TestDecode(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Payload
		wantErr bool
	}{
		{
			name: "valid",
			args: args{bytes.NewBufferString(`{"signature":{"version":"1","source":"test"},"fields":["spkid"],"data":[[123]],"count":1}`)},
			want: &Payload{
				Signature: struct {
					Version string `json:"version"`
					Source  string `json:"source"`
				}{Version: "1", Source: "test"},
				Fields: []string{"spkid"},
				Data:   [][]any{{json.Number("123")}},
				Count:  1,
			},
			wantErr: false,
		},
		{
			name:    "nil reader",
			args:    args{nil},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPayload_Bodies(t *testing.T) {
	wantBody := Body{Identity: Identity{SpkID: ptrTo(123), NEO: ptrTo(true)}}
	type fields struct {
		Signature struct {
			Version string `json:"version"`
			Source  string `json:"source"`
		}
		Fields []string
		Data   [][]any
		Count  int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Body
		wantErr bool
	}{

		{
			name: "ok",
			fields: struct {
				Signature struct {
					Version string `json:"version"`
					Source  string `json:"source"`
				}
				Fields []string
				Data   [][]any
				Count  int
			}{
				Fields: []string{"spkid", "neo"},
				Data:   [][]any{{json.Number("123"), "Y"}},
			},
			want:    []Body{wantBody},
			wantErr: false,
		},
		{
			name: "mismatched row",
			fields: struct {
				Signature struct {
					Version string `json:"version"`
					Source  string `json:"source"`
				}
				Fields []string
				Data   [][]any
				Count  int
			}{
				Fields: []string{"spkid", "neo"},
				Data:   [][]any{{json.Number("123")}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payload{
				Signature: tt.fields.Signature,
				Fields:    tt.fields.Fields,
				Data:      tt.fields.Data,
				Count:     tt.fields.Count,
			}
			got, err := p.Bodies()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bodies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bodies() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPayload_Records(t *testing.T) {
	type fields struct {
		Signature struct {
			Version string `json:"version"`
			Source  string `json:"source"`
		}
		Fields []string
		Data   [][]any
		Count  int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Record
		wantErr bool
	}{
		{
			name: "ok",
			fields: struct {
				Signature struct {
					Version string `json:"version"`
					Source  string `json:"source"`
				}
				Fields []string
				Data   [][]any
				Count  int
			}{
				Fields: []string{"a", "b"},
				Data:   [][]any{{json.Number("1"), json.Number("2")}},
			},
			want: []Record{
				{"a": json.Number("1"), "b": json.Number("2")},
			},
			wantErr: false,
		},
		{
			name: "mismatch",
			fields: struct {
				Signature struct {
					Version string `json:"version"`
					Source  string `json:"source"`
				}
				Fields []string
				Data   [][]any
				Count  int
			}{
				Fields: []string{"a", "b"},
				Data:   [][]any{{json.Number("1")}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payload{
				Signature: tt.fields.Signature,
				Fields:    tt.fields.Fields,
				Data:      tt.fields.Data,
				Count:     tt.fields.Count,
			}
			got, err := p.Records()
			if (err != nil) != tt.wantErr {
				t.Errorf("Records() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Records() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_getBool(t *testing.T) {
	type args struct {
		field Field
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
		field Field
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
		field Field
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
		field Field
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
