package query

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type NumberedStatus uint

const (
	NumStatusAny NumberedStatus = iota
	NumStatusNumbered
	NumStatusUnnumbered
)

func (n NumberedStatus) String() string {
	switch n {
	case NumStatusAny:
		return ""
	case NumStatusNumbered:
		return "n"
	case NumStatusUnnumbered:
		return "u"
	default:
		return fmt.Sprintf("Invalid NumberedStatus(%d)", n)
	}
}

type Kind uint

const (
	KindAny Kind = iota
	KindAsteroid
	KindComet
)

func (k Kind) String() string {
	switch k {
	case KindAny:
		return ""
	case KindAsteroid:
		return "a"
	case KindComet:
		return "c"
	default:
		return fmt.Sprintf("Invalid Kind(%d)", k)
	}
}

type Group uint

const (
	GroupAny Group = iota
	GroupNEO
	GroupPHA
)

func (g Group) String() string {
	switch g {
	case GroupAny:
		return ""
	case GroupNEO:
		return "neo"
	case GroupPHA:
		return "pha"
	default:
		return fmt.Sprintf("Invalid Group(%d)", g)
	}
}

type Class uint
type Classes []Class

const (
	IEO Class = iota + 1 // Atira
	ATE                  // Aten
	APO                  // Apollo
	AMO                  // Amor
	MCA                  // Mars-crossing Asteroid
	IMB                  // Inner Main-belt Asteroid
	MBA                  // Main-belt Asteroid
	OMB                  // Outer Main-belt Asteroid
	TJN                  // Jupiter Trojan
	AST                  // Asteroid
	CEN                  // Centaur
	TNO                  // TransNeptunian Object
	PAA                  // Parabolic “Asteroid”
	HYA                  // Hyperbolic “Asteroid”
	ETc                  // Encke-type Comet
	JFc                  // Jupiter-family Comet
	JFC                  // Jupiter-family Comet*
	CTc                  // Chiron-type Comet
	HTC                  // Halley-type Comet*
	PAR                  // Parabolic Comet
	HYP                  // Hyperbolic Comet
	COM                  // Comet
)

var classCodes = map[Class]string{
	IEO: "IEO", ATE: "ATE", APO: "APO", AMO: "AMO", MCA: "MCA", IMB: "IMB",
	MBA: "MBA", OMB: "OMB", TJN: "TJN", AST: "AST", CEN: "CEN", TNO: "TNO",
	PAA: "PAA", HYA: "HYA", ETc: "ETc", JFc: "JFc", JFC: "JFC", CTc: "CTc",
	HTC: "HTC", PAR: "PAR", HYP: "HYP", COM: "COM",
}

func (c Class) String() string {
	if s, ok := classCodes[c]; ok {
		return s
	}
	return fmt.Sprintf("Invalid Class(%d)", c)
}

func (c Classes) String() string {
	n := len(c)
	if n > 3 {
		n = 3
	}
	parts := make([]string, n)
	for i, class := range c[:n] {
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

type Operator uint

const (
	OpEQ Operator = iota + 1 //	equal
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

var opCode = map[Operator]string{
	OpEQ: "EQ", OpNE: "NE", OpLT: "LT", OpGT: "GT", OpLE: "LE",
	OpGE: "GE", OpRG: "RG", OpRE: "RE", OpDF: "DF", OpND: "ND",
}

func (op Operator) String() string {
	if s, ok := opCode[op]; ok {
		return s
	}
	return fmt.Sprintf("InvalidOp(%d)", op)
}

func C(field string, args ...string) ComparisonExpr {
	return ComparisonExpr(strings.Join(append([]string{field}, args...), "|"))
}
func EQ(field, value string) ComparisonExpr    { return C(field, OpEQ.String(), value) }
func NE(field, value string) ComparisonExpr    { return C(field, OpNE.String(), value) }
func LT(field, value string) ComparisonExpr    { return C(field, OpLT.String(), value) }
func GT(field, value string) ComparisonExpr    { return C(field, OpGT.String(), value) }
func LE(field, value string) ComparisonExpr    { return C(field, OpLE.String(), value) }
func GE(field, value string) ComparisonExpr    { return C(field, OpGE.String(), value) }
func RG(field, min, max string) ComparisonExpr { return C(field, OpRG.String(), min, max) }
func RE(field, value string) ComparisonExpr    { return C(field, OpRE.String(), value) }
func DF(field string) ComparisonExpr           { return C(field, OpDF.String()) }
func ND(field string) ComparisonExpr           { return C(field, OpND.String()) }

type Filter struct {
	NumberedStatus NumberedStatus
	Kind           Kind
	Group          Group
	// Class, limit results by up to 3 orbital classes
	// Refer to orbit class table at https://ssd-api.jpl.nasa.gov/doc/sbdb_filter.html
	Classes Classes
	// MustHaveSatellite, when true, will filter for bodies with at least one know satellite.
	MustHaveSatellite bool
	// ExcludeFragments, when true, will exclude all comet fragments from results.
	ExcludeFragments bool
	FieldConstraints *Expr
}

func (f Filter) Values() (url.Values, error) {
	v := url.Values{}
	if f.NumberedStatus > NumStatusAny {
		v.Set("sb-ns", f.NumberedStatus.String())
	}
	if f.Kind > KindAny {
		v.Set("sb-kind", f.Kind.String())
	}
	if f.Group > GroupAny {
		v.Set("sb-group", f.Group.String())
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
