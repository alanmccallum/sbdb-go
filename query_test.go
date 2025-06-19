package sbdb

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAnd_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		a       And
		want    []byte
		wantErr bool
	}{
		{
			name:    "Happy Path",
			a:       And{ComparisonExpr("pole|DF"), ComparisonExpr("condition_code|EQ|0"), ComparisonExpr("albedo|DF"), ComparisonExpr("H|DF")},
			want:    []byte(`{"AND":["pole|DF","condition_code|EQ|0","albedo|DF","H|DF"]}`),
			wantErr: false,
		},
		{
			name:    "Empty Or",
			a:       And{},
			want:    []byte(`{"AND":[]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestOr_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		o       Or
		want    []byte
		wantErr bool
	}{
		{
			name:    "Happy Path",
			o:       Or{ComparisonExpr("pole|DF"), ComparisonExpr("condition_code|EQ|0"), ComparisonExpr("albedo|DF"), ComparisonExpr("H|DF")},
			want:    []byte(`{"OR":["pole|DF","condition_code|EQ|0","albedo|DF","H|DF"]}`),
			wantErr: false,
		},
		{
			name:    "Empty Or",
			o:       Or{},
			want:    []byte(`{"OR":[]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.o.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestComparisonExpr_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		c       ComparisonExpr
		want    []byte
		wantErr bool
	}{
		{
			name:    "Simple String",
			c:       "Hello World",
			want:    []byte(`"Hello World"`),
			wantErr: false,
		},
		{
			name:    "Empty String",
			c:       "",
			want:    []byte(`""`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestClass_String(t *testing.T) {
	tests := []struct {
		name string
		c    SBClass
		want string
	}{
		{
			name: "Out of Range 0",
			c:    0,
			want: "Invalid SBClass(0)",
		},
		{
			name: "Out of Range (High)",
			c:    999,
			want: "Invalid SBClass(999)",
		},
		{
			name: "IEO",
			c:    IEO,
			want: "IEO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClasses_String(t *testing.T) {
	tests := []struct {
		name string
		c    SBClasses
		want string
	}{
		{
			name: "No SBClass",
			c:    SBClasses{},
			want: "",
		},
		{
			name: "One SBClass",
			c:    SBClasses{IEO},
			want: "IEO",
		},
		{
			name: "Two SBClasses",
			c:    SBClasses{IEO, ATE},
			want: "IEO,ATE",
		},
		{
			name: "Three SBClasses",
			c:    SBClasses{IEO, ATE, APO},
			want: "IEO,ATE,APO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_String(t *testing.T) {
	tests := []struct {
		name string
		g    SGGroup
		want string
	}{
		{
			name: "Any",
			g:    GroupAny,
			want: "",
		},
		{
			name: "GroupNEO",
			g:    GroupNEO,
			want: "neo",
		},
		{
			name: "Out of Range (High)",
			g:    999,
			want: "Invalid SGGroup(999)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKind_String(t *testing.T) {
	tests := []struct {
		name string
		k    SBKind
		want string
	}{
		{
			name: "Any",
			k:    KindAny,
			want: "",
		},
		{
			name: "KindAsteroid",
			k:    KindAsteroid,
			want: "a",
		},
		{
			name: "Out of Range (High)",
			k:    999,
			want: "Invalid SBKind(999)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberedStatus_String(t *testing.T) {
	tests := []struct {
		name string
		n    SBNS
		want string
	}{
		{
			name: "Any",
			n:    NumStatusAny,
			want: "",
		},
		{
			name: "NumStatusNumbered",
			n:    NumStatusNumbered,
			want: "n",
		},
		{
			name: "Out of Range (High)",
			n:    999,
			want: "Invalid SBNS(999)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_operator_String(t *testing.T) {
	tests := []struct {
		name string
		op   operator
		want string
	}{
		{
			name: "Out of Range (Low)",
			op:   0,
			want: "InvalidOp(0)",
		},
		{
			name: "EQ",
			op:   OpEQ,
			want: "EQ",
		},
		{
			name: "Out of Range (High)",
			op:   999,
			want: "InvalidOp(999)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Values(t *testing.T) {
	q := And{ComparisonExpr("pole|DF"), ComparisonExpr("condition_code|EQ|0"), ComparisonExpr("albedo|DF"), ComparisonExpr("H|DF")}
	json, err := q.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		fields  Filter
		want    url.Values
		wantErr bool
	}{
		{
			name:    "Nil",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Need at least one field",
			fields:  Filter{Fields: FieldSet{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Minimum Filter",
			fields:  Filter{Fields: NewFieldSet("field")},
			want:    url.Values{"fields": []string{"field"}},
			wantErr: false,
		},
		{
			name:    "Several Fields",
			fields:  Filter{Fields: NewFieldSet("field1", "field2", "field3", "field4", "field5", "field6")},
			want:    url.Values{"fields": []string{"field1,field2,field3,field4,field5,field6"}},
			wantErr: false,
		},
		{
			name: "Sort - Too many",
			fields: Filter{
				Fields: NewFieldSet("field"),
				Sort:   Fields{"field1", "field2", "field3", "field4", "field5", "field6"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sort - Valid",
			fields: Filter{
				Fields: NewFieldSet("field"),
				Sort:   Fields{"field1", "field2", "field3"}},
			want:    url.Values{"fields": []string{"field"}, "sort": []string{"field1,field2,field3"}},
			wantErr: false,
		},
		{
			name: "SBClasses - Too Many",
			fields: Filter{
				Fields:  NewFieldSet("field"),
				Classes: SBClasses{IEO, ATE, APO, AMO, MCA, IMB, MBA, OMB, TJN},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "SBClasses - Valid",
			fields: Filter{
				Fields:  NewFieldSet("field"),
				Classes: SBClasses{IEO, ATE, APO},
			},
			want:    url.Values{"fields": []string{"field"}, "sb-class": []string{"IEO,ATE,APO"}},
			wantErr: false,
		},
		{
			name: "Limit, LimitFrom, SBNS, SBKind, SGGroup, MustHaveSatellite, & ExcludeFragments - Valid",
			fields: Filter{
				Fields:            NewFieldSet("field"),
				Limit:             10,
				LimitFrom:         10,
				NumberedStatus:    NumStatusNumbered,
				Kind:              KindAsteroid,
				Group:             GroupNEO,
				MustHaveSatellite: true,
				ExcludeFragments:  true,
			},
			want: url.Values{
				"fields":     []string{"field"},
				"limit":      []string{"10"},
				"limit-from": []string{"10"},
				"sb-ns":      []string{"n"},
				"sb-kind":    []string{"a"},
				"sb-group":   []string{"neo"},
				"sb-sat":     []string{"true"},
				"sb-xfrag":   []string{"true"},
			},
			wantErr: false,
		},
		{
			name: "Field Constraints - Valid",
			fields: Filter{
				Fields:           NewFieldSet("field"),
				FieldConstraints: And{ComparisonExpr("pole|DF"), ComparisonExpr("condition_code|EQ|0"), ComparisonExpr("albedo|DF"), ComparisonExpr("H|DF")}},
			want: url.Values{
				"fields": []string{"field"},
				"sb-cf":  []string{string(json)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filter{
				Fields:            tt.fields.Fields,
				Sort:              tt.fields.Sort,
				Limit:             tt.fields.Limit,
				LimitFrom:         tt.fields.LimitFrom,
				NumberedStatus:    tt.fields.NumberedStatus,
				Kind:              tt.fields.Kind,
				Group:             tt.fields.Group,
				Classes:           tt.fields.Classes,
				MustHaveSatellite: tt.fields.MustHaveSatellite,
				ExcludeFragments:  tt.fields.ExcludeFragments,
				FieldConstraints:  tt.fields.FieldConstraints,
			}
			got, err := f.Values()
			if (err != nil) != tt.wantErr {
				t.Errorf("Values() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFilter_Values_Query(t *testing.T) {
	got, err := Filter{
		Fields:         NewFieldSet("spkid", "full_name", "kind", "pdes", "name", "prefix", "class", "neo", "pha", "sats"),
		Limit:          1,
		NumberedStatus: NumStatusUnnumbered,
	}.Values()
	if err != nil {
		t.Fatal(err)
	}

	want, err := url.ParseQuery("limit=1&sb-ns=u&fields=class,full_name,kind,name,neo,pdes,pha,prefix,sats,spkid")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func Test_Ops(t *testing.T) {
	type testFunc func() ComparisonExpr
	var tests = []struct {
		name string
		test testFunc
		want ComparisonExpr
	}{
		{
			name: "EQ",
			test: func() ComparisonExpr { return EQ("field", "value") },
			want: "field|EQ|value",
		},
		{
			name: "NE",
			test: func() ComparisonExpr { return NE("field", "value") },
			want: "field|NE|value",
		},
		{
			name: "LT",
			test: func() ComparisonExpr { return LT("field", "value") },
			want: "field|LT|value",
		},
		{
			name: "GT",
			test: func() ComparisonExpr { return GT("field", "value") },
			want: "field|GT|value",
		},
		{
			name: "LE",
			test: func() ComparisonExpr { return LE("field", "value") },
			want: "field|LE|value",
		},
		{
			name: "GE",
			test: func() ComparisonExpr { return GE("field", "value") },
			want: "field|GE|value",
		},
		{
			name: "RG",
			test: func() ComparisonExpr { return RG("field", "min", "max") },
			want: "field|RG|min|max",
		},
		{
			name: "RE",
			test: func() ComparisonExpr { return RE("field", "value") },
			want: "field|RE|value",
		},
		{
			name: "DF",
			test: func() ComparisonExpr { return DF("field") },
			want: "field|DF",
		},
		{
			name: "ND",
			test: func() ComparisonExpr { return ND("field") },
			want: "field|ND",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.test(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
