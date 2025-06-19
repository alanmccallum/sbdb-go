package sbdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// NumStatusFilter limits results by numbered status.
// See the SBDB Query API documentation for details on each status code.
type NumStatusFilter uint

const (
	// NumStatusAny returns all records regardless of numbering.
	NumStatusAny NumStatusFilter = iota
	// NumStatusNumbered restricts results to numbered bodies.
	NumStatusNumbered
	// NumStatusUnnumbered restricts results to provisional designations.
	NumStatusUnnumbered
)

func (n NumStatusFilter) String() string {
	switch n {
	case NumStatusAny:
		return ""
	case NumStatusNumbered:
		return "n"
	case NumStatusUnnumbered:
		return "u"
	default:
		return fmt.Sprintf("Invalid NumStatusFilter(%d)", n)
	}
}

// KindFilter restricts results to asteroids or comets.
// The zero value includes both kinds.
type KindFilter uint

const (
	// KindAny does not filter by kind.
	KindAny KindFilter = iota
	// KindAsteroid returns asteroid records only.
	KindAsteroid
	// KindComet returns comet records only.
	KindComet
)

func (k KindFilter) String() string {
	switch k {
	case KindAny:
		return ""
	case KindAsteroid:
		return "a"
	case KindComet:
		return "c"
	default:
		return fmt.Sprintf("Invalid KindFilter(%d)", k)
	}
}

// GroupFilter narrows results to common NEO and PHA groups.
// The zero value does not apply group filtering.
type GroupFilter uint

const (
	// GroupAny performs no group filtering.
	GroupAny GroupFilter = iota
	// GroupNEO filters for near-Earth objects.
	GroupNEO
	// GroupPHA filters for potentially hazardous asteroids.
	GroupPHA
)

func (g GroupFilter) String() string {
	switch g {
	case GroupAny:
		return ""
	case GroupNEO:
		return "neo"
	case GroupPHA:
		return "pha"
	default:
		return fmt.Sprintf("Invalid GroupFilter(%d)", g)
	}
}

// ClassFilter restricts results by orbit class.
// Class codes correspond to those documented in the SBDB API.
type ClassFilter uint

// ClassFilters is a collection of ClassFilter values.
type ClassFilters []ClassFilter

const (
	IEO ClassFilter = iota + 1 // Atira
	ATE                        // Aten
	APO                        // Apollo
	AMO                        // Amor
	MCA                        // Mars-crossing Asteroid
	IMB                        // Inner Main-belt Asteroid
	MBA                        // Main-belt Asteroid
	OMB                        // Outer Main-belt Asteroid
	TJN                        // Jupiter Trojan
	AST                        // Asteroid
	CEN                        // Centaur
	TNO                        // TransNeptunian Object
	PAA                        // Parabolic “Asteroid”
	HYA                        // Hyperbolic “Asteroid”
	ETc                        // Encke-type Comet
	JFc                        // Jupiter-family Comet
	JFC                        // Jupiter-family Comet*
	CTc                        // Chiron-type Comet
	HTC                        // Halley-type Comet*
	PAR                        // Parabolic Comet
	HYP                        // Hyperbolic Comet
	COM                        // Comet
)

var classCodes = map[ClassFilter]string{
	IEO: "IEO", ATE: "ATE", APO: "APO", AMO: "AMO", MCA: "MCA", IMB: "IMB",
	MBA: "MBA", OMB: "OMB", TJN: "TJN", AST: "AST", CEN: "CEN", TNO: "TNO",
	PAA: "PAA", HYA: "HYA", ETc: "ETc", JFc: "JFc", JFC: "JFC", CTc: "CTc",
	HTC: "HTC", PAR: "PAR", HYP: "HYP", COM: "COM",
}

func (c ClassFilter) String() string {
	if s, ok := classCodes[c]; ok {
		return s
	}
	return fmt.Sprintf("Invalid ClassFilter(%d)", c)
}

func (c ClassFilters) String() string {
	parts := make([]string, len(c))
	for i, class := range c {
		parts[i] = class.String()
	}
	return strings.Join(parts, ",")
}

// Expr represents a filter expression that can be marshaled to JSON.
type Expr interface {
	MarshalJSON() ([]byte, error)
}

// And groups expressions that must all evaluate to true.
type And []Expr

// Or groups expressions where at least one must evaluate to true.
type Or []Expr

// ComparisonExpr encodes a single comparison operator expression.
type ComparisonExpr string

func (a And) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]Expr{"AND": a})
}
func (o Or) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]Expr{"OR": o})
}
func (c ComparisonExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

type operator uint

const (
	OpEQ operator = iota + 1 //	equal
	OpNE                     //	not equal
	OpLT                     //	less than
	OpGT                     //	greater than
	OpLE                     //	less than or equal
	OpGE                     //	greater than or equal
	OpRG                     //	inclusive range: matches values greater than or equal to the minimum and less than or equal to the maximum
	OpRE                     //	regular expression
	OpDF                     //	the value is defined (not NULL)
	OpND                     //	the value is not defined (is NULL)
)

var opCode = map[operator]string{
	OpEQ: "EQ", OpNE: "NE", OpLT: "LT", OpGT: "GT", OpLE: "LE",
	OpGE: "GE", OpRG: "RG", OpRE: "RE", OpDF: "DF", OpND: "ND",
}

func (op operator) String() string {
	if s, ok := opCode[op]; ok {
		return s
	}
	return fmt.Sprintf("InvalidOp(%d)", op)
}

func c(field string, args ...string) ComparisonExpr {
	return ComparisonExpr(strings.Join(append([]string{field}, args...), "|"))
}

// EQ creates an equality comparison expression.
func EQ(field, value string) ComparisonExpr { return c(field, OpEQ.String(), value) }

// NE creates a not-equal comparison expression.
func NE(field, value string) ComparisonExpr { return c(field, OpNE.String(), value) }

// LT creates a less-than comparison expression.
func LT(field, value string) ComparisonExpr { return c(field, OpLT.String(), value) }

// GT creates a greater-than comparison expression.
func GT(field, value string) ComparisonExpr { return c(field, OpGT.String(), value) }

// LE creates a less-than-or-equal comparison expression.
func LE(field, value string) ComparisonExpr { return c(field, OpLE.String(), value) }

// GE creates a greater-than-or-equal comparison expression.
func GE(field, value string) ComparisonExpr { return c(field, OpGE.String(), value) }

// RG creates an inclusive range comparison expression.
func RG(field, min, max string) ComparisonExpr { return c(field, OpRG.String(), min, max) }

// RE creates a regular-expression comparison expression.
func RE(field, value string) ComparisonExpr { return c(field, OpRE.String(), value) }

// DF checks that a field is defined.
func DF(field string) ComparisonExpr { return c(field, OpDF.String()) }

// ND checks that a field is not defined.
func ND(field string) ComparisonExpr { return c(field, OpND.String()) }

// FieldSet represents the list of fields requested from the API.
// It ensures stable ordering when encoded into a query string.
type FieldSet map[string]struct{}

// NewFieldSet returns a FieldSet initialized with the provided fields.
func NewFieldSet(fields ...Field) FieldSet {
	fs := FieldSet{}
	for _, f := range fields {
		fs.Add(f)
	}
	return fs
}

// Add inserts a field into the set.
func (fs FieldSet) Add(field Field) {
	fs[field.String()] = struct{}{}
}

// AddFields inserts multiple fields.
func (fs FieldSet) AddFields(fields ...Field) {
	for _, f := range fields {
		fs.Add(f)
	}
}

// Remove deletes a field from the set.
func (fs FieldSet) Remove(field Field) {
	delete(fs, field.String())
}

// List returns the fields in sorted order.
func (fs FieldSet) List() []string {
	out := make([]string, 0, len(fs))
	for f := range fs {
		out = append(out, f)
	}
	sort.Strings(out) // Optional: stable ordering
	return out
}

// String implements fmt.Stringer and returns a comma separated field list.
func (fs FieldSet) String() string {
	return strings.Join(fs.List(), ",")
}

// Filter defines the search parameters for a query. It mirrors the
// parameters described in the SBDB Query API documentation.
type Filter struct {
	Fields         FieldSet
	Limit          uint
	LimitFrom      uint
	NumberedStatus NumStatusFilter
	Kind           KindFilter
	Group          GroupFilter
	// Classes limits results by up to 3 orbital classes
	// Refer to orbit class table at https://ssd-api.jpl.nasa.gov/doc/sbdb_filter.html
	Classes ClassFilters
	// MustHaveSatellite, when true, will filter for bodies with at least one know satellite.
	MustHaveSatellite bool
	// ExcludeFragments, when true, will exclude all comet fragments from results.
	ExcludeFragments bool
	// FieldConstraints applies advanced field-level filters encoded as
	// AND/OR expressions. See the SBDB filter documentation for syntax.
	FieldConstraints Expr
}

// Values converts the Filter into URL query parameters.
func (f Filter) Values() (url.Values, error) {
	if len(f.Fields) <= 0 {
		return nil, errors.New("must provide at least one field")
	}
	if len(f.Classes) > 3 {
		return nil, fmt.Errorf("len(ClassFilters) = %d, max = 3", len(f.Classes))
	}

	v := url.Values{}
	v.Set("fields", f.Fields.String())
	if f.Limit > 0 {
		v.Set("limit", strconv.FormatUint(uint64(f.Limit), 10))
	}
	if f.LimitFrom > 0 {
		v.Set("limit-from", strconv.FormatUint(uint64(f.LimitFrom), 10))
	}
	if f.NumberedStatus > NumStatusAny {
		v.Set("sb-ns", f.NumberedStatus.String())
	}
	if f.Kind > KindAny {
		v.Set("sb-kind", f.Kind.String())
	}
	if f.Group > GroupAny {
		v.Set("sb-group", f.Group.String())
	}
	if f.MustHaveSatellite {
		v.Set("sb-sat", strconv.FormatBool(f.MustHaveSatellite))
	}
	if f.ExcludeFragments {
		v.Set("sb-xfrag", strconv.FormatBool(f.ExcludeFragments))
	}
	if len(f.Classes) > 0 {
		v.Set("sb-class", f.Classes.String())
	}
	if f.FieldConstraints != nil {
		marshal, err := json.Marshal(f.FieldConstraints)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal fieldConstraints: %w", err)
		}
		v.Set("sb-cf", string(marshal))
	}

	return v, nil
}
