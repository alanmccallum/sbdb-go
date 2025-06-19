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

type SBNS uint

const (
	NumStatusAny SBNS = iota
	NumStatusNumbered
	NumStatusUnnumbered
)

func (n SBNS) String() string {
	switch n {
	case NumStatusAny:
		return ""
	case NumStatusNumbered:
		return "n"
	case NumStatusUnnumbered:
		return "u"
	default:
		return fmt.Sprintf("Invalid SBNS(%d)", n)
	}
}

type SBKind uint

const (
	KindAny SBKind = iota
	KindAsteroid
	KindComet
)

func (k SBKind) String() string {
	switch k {
	case KindAny:
		return ""
	case KindAsteroid:
		return "a"
	case KindComet:
		return "c"
	default:
		return fmt.Sprintf("Invalid SBKind(%d)", k)
	}
}

type SGGroup uint

const (
	GroupAny SGGroup = iota
	GroupNEO
	GroupPHA
)

func (g SGGroup) String() string {
	switch g {
	case GroupAny:
		return ""
	case GroupNEO:
		return "neo"
	case GroupPHA:
		return "pha"
	default:
		return fmt.Sprintf("Invalid SGGroup(%d)", g)
	}
}

type SBClass uint
type SBClasses []SBClass

const (
	IEO SBClass = iota + 1 // Atira
	ATE                    // Aten
	APO                    // Apollo
	AMO                    // Amor
	MCA                    // Mars-crossing Asteroid
	IMB                    // Inner Main-belt Asteroid
	MBA                    // Main-belt Asteroid
	OMB                    // Outer Main-belt Asteroid
	TJN                    // Jupiter Trojan
	AST                    // Asteroid
	CEN                    // Centaur
	TNO                    // TransNeptunian Object
	PAA                    // Parabolic “Asteroid”
	HYA                    // Hyperbolic “Asteroid”
	ETc                    // Encke-type Comet
	JFc                    // Jupiter-family Comet
	JFC                    // Jupiter-family Comet*
	CTc                    // Chiron-type Comet
	HTC                    // Halley-type Comet*
	PAR                    // Parabolic Comet
	HYP                    // Hyperbolic Comet
	COM                    // Comet
)

var classCodes = map[SBClass]string{
	IEO: "IEO", ATE: "ATE", APO: "APO", AMO: "AMO", MCA: "MCA", IMB: "IMB",
	MBA: "MBA", OMB: "OMB", TJN: "TJN", AST: "AST", CEN: "CEN", TNO: "TNO",
	PAA: "PAA", HYA: "HYA", ETc: "ETc", JFc: "JFc", JFC: "JFC", CTc: "CTc",
	HTC: "HTC", PAR: "PAR", HYP: "HYP", COM: "COM",
}

func (c SBClass) String() string {
	if s, ok := classCodes[c]; ok {
		return s
	}
	return fmt.Sprintf("Invalid SBClass(%d)", c)
}

func (c SBClasses) String() string {
	parts := make([]string, len(c))
	for i, class := range c {
		parts[i] = class.String()
	}
	return strings.Join(parts, ",")
}

type Expr interface {
	MarshalJSON() ([]byte, error)
}

type And []Expr
type Or []Expr
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

func EQ(field, value string) ComparisonExpr    { return c(field, OpEQ.String(), value) }
func NE(field, value string) ComparisonExpr    { return c(field, OpNE.String(), value) }
func LT(field, value string) ComparisonExpr    { return c(field, OpLT.String(), value) }
func GT(field, value string) ComparisonExpr    { return c(field, OpGT.String(), value) }
func LE(field, value string) ComparisonExpr    { return c(field, OpLE.String(), value) }
func GE(field, value string) ComparisonExpr    { return c(field, OpGE.String(), value) }
func RG(field, min, max string) ComparisonExpr { return c(field, OpRG.String(), min, max) }
func RE(field, value string) ComparisonExpr    { return c(field, OpRE.String(), value) }
func DF(field string) ComparisonExpr           { return c(field, OpDF.String()) }
func ND(field string) ComparisonExpr           { return c(field, OpND.String()) }

type FieldSet map[string]struct{}

func NewFieldSet(fields ...Field) FieldSet {
	fs := FieldSet{}
	for _, f := range fields {
		fs.Add(f)
	}
	return fs
}

func (fs FieldSet) Add(field Field) {
	fs[field.String()] = struct{}{}
}
func (fs FieldSet) AddFields(fields []Field) {
	for _, f := range fields {
		fs.Add(f)
	}
}
func (fs FieldSet) Remove(field Field) {
	delete(fs, field.String())
}
func (fs FieldSet) List() []string {
	out := make([]string, 0, len(fs))
	for f := range fs {
		out = append(out, f)
	}
	sort.Strings(out) // Optional: stable ordering
	return out
}
func (fs FieldSet) String() string {
	return strings.Join(fs.List(), ",")
}

type Filter struct {
	Fields         FieldSet
	Limit          uint
	LimitFrom      uint
	NumberedStatus SBNS
	Kind           SBKind
	Group          SGGroup
	// Classes, limit results by up to 3 orbital classes
	// Refer to orbit class table at https://ssd-api.jpl.nasa.gov/doc/sbdb_filter.html
	Classes SBClasses
	// MustHaveSatellite, when true, will filter for bodies with at least one know satellite.
	MustHaveSatellite bool
	// ExcludeFragments, when true, will exclude all comet fragments from results.
	ExcludeFragments bool
	FieldConstraints Expr
}

func (f Filter) Values() (url.Values, error) {
	if f.Fields == nil || len(f.Fields) <= 0 {
		return nil, errors.New("must provide at least one field")
	}
	if len(f.Classes) > 3 {
		return nil, fmt.Errorf("len(SBClasses) = %d, max = 3", len(f.Classes))
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
