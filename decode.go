package sbdb

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Decode parses an SBDB JSON payload from r.
// A json.Decoder is used with UseNumber so numeric fields are decoded as
// json.Number instead of default float64 values.
func Decode(r io.Reader) (*Payload, error) {
	if r == nil {
		return nil, errors.New("nil reader")
	}
	if _, ok := r.(*bufio.Reader); !ok {
		r = bufio.NewReader(r)
	}
	var p = Payload{}
	dec := json.NewDecoder(r)
	dec.UseNumber()
	if err := dec.Decode(&p); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	return &p, nil
}

// Payload is a raw SBDB response containing records and metadata.
type Payload struct {
	Signature struct {
		Version string `json:"version"`
		Source  string `json:"source"`
	} `json:"signature"`
	Fields []string `json:"fields"`
	Data   [][]any  `json:"data"`
	Count  int      `json:"count"`
}

// Records returns the payload data as a slice of generic Records.
// Because Decode uses json.Decoder.UseNumber, any numeric values may be
// json.Number, which callers should handle appropriately
func (p *Payload) Records() ([]Record, error) {
	records := make([]Record, len(p.Data))
	for i, b := range p.Data {
		if len(b) != len(p.Fields) {
			return nil, fmt.Errorf("data element %d has %d fields, expected %d", i, len(b), len(p.Fields))
		}
		records[i] = make(Record)
		for j, v := range b {
			records[i][Field(p.Fields[j])] = v
		}
	}

	return records, nil
}

// Bodies converts the payload data into strongly typed Body structs.
func (p *Payload) Bodies() ([]Body, error) {
	records, err := p.Records()
	if err != nil {
		return nil, err
	}
	bodies := make([]Body, len(records))
	for i, r := range records {
		bodies[i] = Body{
			Identity:    r.identity(),
			Orbit:       r.orbit(),
			Uncertainty: r.uncertainty(),
			Solution:    r.solution(),
			Quality:     r.quality(),
			NonGrav:     r.nonGrav(),
			Physical:    r.physical(),
		}
	}
	return bodies, nil
}

// Record represents a single result row as a map of field names to values.
type Record map[Field]any

func (r Record) identity() Identity {
	return Identity{
		SpkID:       r.getInt(SpkID),
		FullName:    r.getString(FullName),
		Kind:        r.getString(Kind),
		PDES:        r.getString(PDES),
		Name:        r.getString(Name),
		Prefix:      r.getString(Prefix),
		Class:       r.getString(Class),
		NEO:         r.getBool(NEO),
		PHA:         r.getBool(PHA),
		Sats:        r.getInt(Sats),
		TJupiter:    r.getFloat(TJupiter),
		MOID:        r.getFloat(MOID),
		MOIDLD:      r.getFloat(MOIDLD),
		MOIDJupiter: r.getFloat(MOIDJupiter),
	}
}
func (r Record) orbit() Orbit {
	return Orbit{
		OrbitID:          r.getString(OrbitID),
		Epoch:            r.getFloat(Epoch),
		EpochMJD:         r.getFloat(EpochMJD),
		EpochCal:         r.getString(EpochCal),
		Equinox:          r.getString(Equinox),
		Eccentricity:     r.getFloat(Eccentricity),
		SemimajorAxis:    r.getFloat(SemimajorAxis),
		PerihelionDist:   r.getFloat(PerihelionDist),
		Inclination:      r.getFloat(Inclination),
		AscNode:          r.getFloat(AscNode),
		PeriapsisArg:     r.getFloat(PeriapsisArg),
		MeanAnomaly:      r.getFloat(MeanAnomaly),
		PeriapsisTime:    r.getFloat(PeriapsisTime),
		PeriapsisTimeCal: r.getString(PeriapsisTimeCal),
		OrbitalPeriod:    r.getFloat(OrbitalPeriod),
		OrbitalPeriodYr:  r.getFloat(OrbitalPeriodYr),
		MeanMotion:       r.getFloat(MeanMotion),
		AphelionDist:     r.getFloat(AphelionDist),
	}
}
func (r Record) uncertainty() Uncertainty {
	return Uncertainty{
		SigmaEcc:     r.getFloat(SigmaEcc),
		SigmaA:       r.getFloat(SigmaA),
		SigmaQ:       r.getFloat(SigmaQ),
		SigmaI:       r.getFloat(SigmaI),
		SigmaAscNode: r.getFloat(SigmaAscNode),
		SigmaPeriArg: r.getFloat(SigmaPeriArg),
		SigmaTP:      r.getFloat(SigmaTP),
		SigmaMA:      r.getFloat(SigmaMA),
		SigmaPeriod:  r.getFloat(SigmaPeriod),
		SigmaN:       r.getFloat(SigmaN),
		SigmaAD:      r.getFloat(SigmaAD),
	}
}
func (r Record) solution() Solution {
	return Solution{
		Source:         r.getString(Source),
		SolutionDate:   r.getString(SolutionDate),
		Producer:       r.getString(Producer),
		DataArc:        r.getInt(DataArc),
		FirstObs:       r.getString(FirstObs),
		LastObs:        r.getString(LastObs),
		ObsUsed:        r.getInt(ObsUsed),
		DelayObsUsed:   r.getInt(DelayObsUsed),
		DopplerObsUsed: r.getInt(DopplerObsUsed),
	}
}
func (r Record) quality() Quality {
	return Quality{
		TwoBody:       r.getBool(TwoBody),
		PEUsed:        r.getString(PEUsed),
		SBUsed:        r.getString(SBUsed),
		ConditionCode: r.getInt(ConditionCode),
		RMS:           r.getFloat(RMS),
	}
}
func (r Record) nonGrav() NonGrav {
	return NonGrav{
		A1:      r.getFloat(A1),
		A2:      r.getFloat(A2),
		A3:      r.getFloat(A3),
		DT:      r.getFloat(DT),
		S0:      r.getFloat(S0),
		A1Sigma: r.getFloat(A1Sigma),
		A2Sigma: r.getFloat(A2Sigma),
		A3Sigma: r.getFloat(A3Sigma),
		DTSigma: r.getFloat(DTSigma),
		S0Sigma: r.getFloat(S0Sigma),
	}
}
func (r Record) physical() Physical {
	return Physical{
		H:             r.getFloat(H),
		G:             r.getFloat(G),
		M1:            r.getFloat(M1),
		K1:            r.getFloat(K1),
		M2:            r.getFloat(M2),
		K2:            r.getFloat(K2),
		PC:            r.getFloat(PC),
		HSigma:        r.getFloat(HSigma),
		Diameter:      r.getFloat(Diameter),
		Extent:        r.getString(Extent),
		GM:            r.getFloat(GM),
		Density:       r.getFloat(Density),
		RotPer:        r.getFloat(RotPer),
		Pole:          r.getString(Pole),
		Albedo:        r.getFloat(Albedo),
		BV:            r.getFloat(BV),
		UB:            r.getFloat(UB),
		IR:            r.getFloat(IR),
		SpecT:         r.getString(SpecT),
		SpecB:         r.getString(SpecB),
		DiameterSigma: r.getFloat(DiameterSigma),
	}
}

func (r Record) getFloat(field Field) *float64 {
	if r[field] == nil {
		return nil
	}
	switch v := r[field].(type) {
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			logFailedTypeAssert("getFloat(json.Number)", field, r[field])
			return nil
		}
		return &f
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			logFailedTypeAssert("getFloat(string)", field, r[field])
			return nil
		}
		return &f
	default:
		logFailedTypeAssert("getFloat", field, r[field])
		return nil
	}
}
func (r Record) getInt(field Field) *int {
	if r[field] == nil {
		return nil
	}
	switch v := r[field].(type) {
	case json.Number:
		i64, err := v.Int64()
		if errors.Is(err, strconv.ErrSyntax) {
			f, err := v.Float64()
			if err != nil {
				logFailedTypeAssert("getFloat(json.Number)", field, r[field])
				return nil
			}
			i := int(f)
			return &i
		}
		if err != nil {
			logFailedTypeAssert("getInt(json.Number)", field, r[field])
			return nil
		}
		i := int(i64)
		return &i
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			logFailedTypeAssert("getInt(string)", field, r[field])
			return nil
		}
		return &i
	default:
		logFailedTypeAssert("getInt", field, r[field])
		return nil
	}
}
func (r Record) getString(field Field) *string {
	if r[field] == nil {
		return nil
	}
	switch v := r[field].(type) {
	case string:
		if field == FullName {
			v = strings.TrimSpace(v)
		}
		return &v
	case json.Number:
		s := v.String()
		return &s
	default:
		log.Debug("Type assertion failed, defaulting to fmt.Sprint", "fn", "getString", "field", field, "value", r[field])
		s := fmt.Sprint(r[field])
		return &s
	}
}
func (r Record) getBool(field Field) *bool {
	if r[field] == nil {
		return nil
	}
	switch v := r[field].(type) {
	case bool:
		return &v
	case string:
		switch strings.ToUpper(v) {
		case "Y", "T":
			v := true
			return &v
		case "N", "F":
			v := false
			return &v
		default:
			logFailedTypeAssert("getBool", field, r[field])
			return nil
		}
	default:
		logFailedTypeAssert("getBool", field, r[field])
		return nil
	}
}
