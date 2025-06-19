package sbdb

// Field represents a SBDB Field, used to build queries and to process responses.
type Field string

// String returns the string representation of the Field
func (f Field) String() string {
	return string(f)
}

// SBDB field names representing identity response fields.
const (
	SpkID       Field = "spkid"     // SPICE identifier for the body
	FullName    Field = "full_name" // Complete object designation
	Kind        Field = "kind"      // Body kind, e.g. asteroid or comet
	PDes        Field = "pdes"      // Primary designation
	Name        Field = "name"      // IAU name
	Prefix      Field = "prefix"    // Numbered prefix
	Class       Field = "class"     // Dynamical class
	NEO         Field = "neo"       // Near Earth Object flag
	PHA         Field = "pha"       // Potentially Hazardous Asteroid flag
	Sats        Field = "sats"      // Number of known satellites
	TJupiter    Field = "t_jup"     // Tisserand parameter w.r.t. Jupiter
	MOID        Field = "moid"      // Earth minimum orbit intersection distance (au)
	MOIDLD      Field = "moid_ld"   // Earth MOID in lunar distances
	MOIDJupiter Field = "moid_jup"  // Jupiter MOID (au)
)

// IdentityFields returns all fields from Identity
func IdentityFields() []Field {
	return []Field{
		SpkID, FullName, Kind, PDes, Name, Prefix, Class, NEO, PHA, Sats, TJupiter, MOID, MOIDLD, MOIDJupiter,
	}
}

// SBDB field names representing orbit response fields.
const (
	OrbitID          Field = "orbit_id"  // Orbit solution identifier
	Epoch            Field = "epoch"     // Reference epoch (JD)
	EpochMJD         Field = "epoch_mjd" // Reference epoch (MJD)
	EpochCal         Field = "epoch_cal" // Reference epoch (calendar)
	Equinox          Field = "equinox"   // Reference frame
	Eccentricity     Field = "e"         // Orbital eccentricity
	SemimajorAxis    Field = "a"         // Semi-major axis (au)
	PerihelionDist   Field = "q"         // Perihelion distance (au)
	Inclination      Field = "i"         // Inclination to the ecliptic (deg)
	AscNode          Field = "om"        // Longitude of ascending node (deg)
	PeriapsisArg     Field = "w"         // Argument of periapsis (deg)
	MeanAnomaly      Field = "ma"        // Mean anomaly at epoch (deg)
	PeriapsisTime    Field = "tp"        // Time of periapsis passage (JD)
	PeriapsisTimeCal Field = "tp_cal"    // Time of periapsis passage (calendar)
	OrbitalPeriod    Field = "per"       // Orbital period (days)
	OrbitalPeriodYr  Field = "per_y"     // Orbital period (years)
	MeanMotion       Field = "n"         // Mean motion (deg/day)
	AphelionDist     Field = "ad"        // Aphelion distance (au)
)

// OrbitFields returns all fields from Orbit
func OrbitFields() []Field {
	return []Field{
		OrbitID, Epoch, EpochMJD, EpochCal, Equinox, Eccentricity, SemimajorAxis, PerihelionDist, Inclination, AscNode,
		PeriapsisArg, MeanAnomaly, PeriapsisTime, PeriapsisTimeCal, OrbitalPeriod, OrbitalPeriodYr, MeanMotion,
		AphelionDist,
	}
}

// SBDB field names representing uncertainty response fields.
const (
	SigmaEcc     Field = "sigma_e"   // 1-sigma uncertainty of eccentricity
	SigmaA       Field = "sigma_a"   // 1-sigma uncertainty of semi-major axis (au)
	SigmaQ       Field = "sigma_q"   // 1-sigma uncertainty of perihelion distance (au)
	SigmaI       Field = "sigma_i"   // 1-sigma uncertainty of inclination (deg)
	SigmaAscNode Field = "sigma_om"  // 1-sigma uncertainty of ascending node (deg)
	SigmaPeriArg Field = "sigma_w"   // 1-sigma uncertainty of periapsis argument (deg)
	SigmaTP      Field = "sigma_tp"  // 1-sigma uncertainty of time of periapsis (JD)
	SigmaMA      Field = "sigma_ma"  // 1-sigma uncertainty of mean anomaly (deg)
	SigmaPeriod  Field = "sigma_per" // 1-sigma uncertainty of orbital period (days)
	SigmaN       Field = "sigma_n"   // 1-sigma uncertainty of mean motion (deg/day)
	SigmaAD      Field = "sigma_ad"  // 1-sigma uncertainty of aphelion distance (au)
)

// UncertaintyFields returns all fields from Uncertainty
func UncertaintyFields() []Field {
	return []Field{
		SigmaEcc, SigmaA, SigmaQ, SigmaI, SigmaAscNode, SigmaPeriArg, SigmaTP, SigmaMA, SigmaPeriod, SigmaN, SigmaAD,
	}
}

// SBDB field names representing solution response fields.
const (
	Source         Field = "source"         // Source of orbit solution
	SolutionDate   Field = "soln_date"      // Solution date
	Producer       Field = "producer"       // Producer of orbit solution
	DataArc        Field = "data_arc"       // Data-arc span (days)
	FirstObs       Field = "first_obs"      // First observation date
	LastObs        Field = "last_obs"       // Last observation date
	ObsUsed        Field = "n_obs_used"     // Number of observations used
	DelayObsUsed   Field = "n_del_obs_used" // Number of delay observations used
	DopplerObsUsed Field = "n_dop_obs_used" // Number of Doppler observations used
	TwoBody        Field = "two_body"       // Two-body approximation flag
	PEUsed         Field = "pe_used"        // Planetary ephemeris used
	SBUsed         Field = "sb_used"        // Small-body perturbers used
	ConditionCode  Field = "condition_code" // Orbit uncertainty condition code
	RMS            Field = "rms"            // RMS residual (arcsec)
)

// SourceFields returns all fields from Source
func SourceFields() []Field {
	return []Field{
		Source, SolutionDate, Producer, DataArc, FirstObs, LastObs, ObsUsed, DelayObsUsed, DopplerObsUsed, TwoBody,
		PEUsed, SBUsed, ConditionCode, RMS,
	}
}

// SBDB field names representing non-gravitational response fields.
const (
	A1      Field = "A1"       // Non-gravitational acceleration parameter A1 (au/d^2)
	A2      Field = "A2"       // Non-gravitational acceleration parameter A2 (au/d^2)
	A3      Field = "A3"       // Non-gravitational acceleration parameter A3 (au/d^2)
	DT      Field = "DT"       // Non-gravitational time parameter (days)
	S0      Field = "S0"       // Non-gravitational scale factor
	A1Sigma Field = "A1_sigma" // 1-sigma uncertainty of A1 (au/d^2)
	A2Sigma Field = "A2_sigma" // 1-sigma uncertainty of A2 (au/d^2)
	A3Sigma Field = "A3_sigma" // 1-sigma uncertainty of A3 (au/d^2)
	DTSigma Field = "DT_sigma" // 1-sigma uncertainty of DT (days)
	S0Sigma Field = "S0_sigma" // 1-sigma uncertainty of S0
)

// NonGravFields returns all fields from NonGrav
func NonGravFields() []Field {
	return []Field{
		A1, A2, A3, DT, S0, A1Sigma, A2Sigma, A3Sigma, DTSigma, S0Sigma,
	}
}

// Field names representing SBDB physical properties response fields.
const (
	H             Field = "H"              // Absolute magnitude H
	G             Field = "G"              // Photometric slope parameter G
	M1            Field = "M1"             // Photometric parameter M1
	K1            Field = "K1"             // Photometric parameter K1
	M2            Field = "M2"             // Photometric parameter M2
	K2            Field = "K2"             // Photometric parameter K2
	PC            Field = "PC"             // Photometric color index PC
	HSigma        Field = "H_sigma"        // 1-sigma uncertainty of H
	Diameter      Field = "diameter"       // Diameter (km)
	Extent        Field = "extent"         // Physical extent (km)
	GM            Field = "GM"             // Gravitational parameter (km^3/s^2)
	Density       Field = "density"        // Bulk density (g/cm^3)
	RotPer        Field = "rot_per"        // Rotation period (hours)
	Pole          Field = "pole"           // Pole orientation (deg)
	Albedo        Field = "albedo"         // Geometric albedo
	BV            Field = "BV"             // B-V color index
	UB            Field = "UB"             // U-B color index
	IR            Field = "IR"             // Infrared color index
	SpecT         Field = "spec_T"         // Spectral taxonomy
	SpecB         Field = "spec_B"         // Spectral bin
	DiameterSigma Field = "diameter_sigma" // 1-sigma uncertainty of diameter (km)
)

// PhysicalFields returns all fields from Physical
func PhysicalFields() []Field {
	return []Field{
		H, G, M1, K1, M2, K2, PC, HSigma, Diameter, Extent, GM, Density, RotPer, Pole, Albedo, BV, UB, IR, SpecT, SpecB,
		DiameterSigma,
	}
}

// Body represents a small-body record from the SBDB Query API.
type Body struct {
	Identity    Identity    // Basic identifying information
	Orbit       Orbit       // Orbital elements at Epoch
	Uncertainty Uncertainty // Uncertainties for orbital elements
	Solution    Solution    // Orbit solution metadata
	Quality     Quality     // Fit and model information
	NonGrav     NonGrav     // Non-gravitational parameters
	Physical    Physical    // Physical characteristics
}

// Identity groups name and classification data.
type Identity struct {
	SpkID       *int     `json:"spkid,omitempty"`     // SPICE identifier
	FullName    *string  `json:"full_name,omitempty"` // Complete object designation
	Kind        *string  `json:"kind,omitempty"`      // Body kind
	PDES        *string  `json:"pdes,omitempty"`      // Primary designation
	Name        *string  `json:"name,omitempty"`      // IAU name
	Prefix      *string  `json:"prefix,omitempty"`    // Numbered prefix
	Class       *string  `json:"class,omitempty"`     // Dynamical class
	NEO         *bool    `json:"neo,omitempty"`       // Near Earth Object flag
	PHA         *bool    `json:"pha,omitempty"`       // Potentially Hazardous flag
	Sats        *int     `json:"sats,omitempty"`      // Number of satellites
	TJupiter    *float64 `json:"t_jup,omitempty"`     // Tisserand parameter w.r.t. Jupiter
	MOID        *float64 `json:"moid,omitempty"`      // Earth MOID (au)
	MOIDLD      *float64 `json:"moid_ld,omitempty"`   // Earth MOID (LD)
	MOIDJupiter *float64 `json:"moid_jup,omitempty"`  // Jupiter MOID (au)
}

// Orbit holds the osculating orbital elements.
type Orbit struct {
	OrbitID          *string  `json:"orbit_id,omitempty"`  // Orbit solution identifier
	Epoch            *float64 `json:"epoch,omitempty"`     // Epoch of osculation (JD)
	EpochMJD         *float64 `json:"epoch_mjd,omitempty"` // Epoch of osculation (MJD)
	EpochCal         *string  `json:"epoch_cal,omitempty"` // Epoch of osculation (calendar)
	Equinox          *string  `json:"equinox,omitempty"`   // Reference equinox
	Eccentricity     *float64 `json:"e,omitempty"`         // Orbital eccentricity
	SemimajorAxis    *float64 `json:"a,omitempty"`         // Semi-major axis (au)
	PerihelionDist   *float64 `json:"q,omitempty"`         // Perihelion distance (au)
	Inclination      *float64 `json:"i,omitempty"`         // Inclination (deg)
	AscNode          *float64 `json:"om,omitempty"`        // Longitude of ascending node (deg)
	PeriapsisArg     *float64 `json:"w,omitempty"`         // Argument of periapsis (deg)
	MeanAnomaly      *float64 `json:"ma,omitempty"`        // Mean anomaly at epoch (deg)
	PeriapsisTime    *float64 `json:"tp,omitempty"`        // Time of periapsis (JD)
	PeriapsisTimeCal *string  `json:"tp_cal,omitempty"`    // Time of periapsis (calendar)
	OrbitalPeriod    *float64 `json:"per,omitempty"`       // Orbital period (days)
	OrbitalPeriodYr  *float64 `json:"per_y,omitempty"`     // Orbital period (years)
	MeanMotion       *float64 `json:"n,omitempty"`         // Mean motion (deg/day)
	AphelionDist     *float64 `json:"ad,omitempty"`        // Aphelion distance (au)
}

// Uncertainty lists one-sigma uncertainties for the orbital elements.
type Uncertainty struct {
	SigmaEcc     *float64 `json:"sigma_e,omitempty"`   // Uncertainty of eccentricity
	SigmaA       *float64 `json:"sigma_a,omitempty"`   // Uncertainty of semi-major axis (au)
	SigmaQ       *float64 `json:"sigma_q,omitempty"`   // Uncertainty of perihelion distance (au)
	SigmaI       *float64 `json:"sigma_i,omitempty"`   // Uncertainty of inclination (deg)
	SigmaAscNode *float64 `json:"sigma_om,omitempty"`  // Uncertainty of ascending node (deg)
	SigmaPeriArg *float64 `json:"sigma_w,omitempty"`   // Uncertainty of periapsis argument (deg)
	SigmaTP      *float64 `json:"sigma_tp,omitempty"`  // Uncertainty of time of periapsis (JD)
	SigmaMA      *float64 `json:"sigma_ma,omitempty"`  // Uncertainty of mean anomaly (deg)
	SigmaPeriod  *float64 `json:"sigma_per,omitempty"` // Uncertainty of orbital period (days)
	SigmaN       *float64 `json:"sigma_n,omitempty"`   // Uncertainty of mean motion (deg/day)
	SigmaAD      *float64 `json:"sigma_ad,omitempty"`  // Uncertainty of aphelion distance (au)
}

// Solution tracks the provenance of the orbital solution.
type Solution struct {
	Source         *string `json:"source,omitempty"`         // Orbit solution source
	SolutionDate   *string `json:"soln_date,omitempty"`      // Date solution was computed
	Producer       *string `json:"producer,omitempty"`       // Orbit producer
	DataArc        *int    `json:"data_arc,omitempty"`       // Data-arc length (days)
	FirstObs       *string `json:"first_obs,omitempty"`      // First observation date
	LastObs        *string `json:"last_obs,omitempty"`       // Last observation date
	ObsUsed        *int    `json:"n_obs_used,omitempty"`     // Number of observations used
	DelayObsUsed   *int    `json:"n_del_obs_used,omitempty"` // Number of delay observations used
	DopplerObsUsed *int    `json:"n_dop_obs_used,omitempty"` // Number of Doppler observations used
}

// Quality describes the orbit fit and modeling options.
type Quality struct {
	TwoBody       *bool    `json:"two_body,omitempty"`       // Two-body approximation flag
	PEUsed        *string  `json:"pe_used,omitempty"`        // Planetary ephemeris used
	SBUsed        *string  `json:"sb_used,omitempty"`        // Small-body perturbers used
	ConditionCode *int     `json:"condition_code,omitempty"` // Orbit condition code
	RMS           *float64 `json:"rms,omitempty"`            // Fit RMS residual (arcsec)
}

// NonGrav holds the non-gravitational parameters.
type NonGrav struct {
	A1      *float64 `json:"A1,omitempty"`       // Nongravitational parameter A1 (au/d^2)
	A2      *float64 `json:"A2,omitempty"`       // Nongravitational parameter A2 (au/d^2)
	A3      *float64 `json:"A3,omitempty"`       // Nongravitational parameter A3 (au/d^2)
	DT      *float64 `json:"DT,omitempty"`       // Nongravitational time parameter (days)
	S0      *float64 `json:"S0,omitempty"`       // Nongravitational scaling factor
	A1Sigma *float64 `json:"A1_sigma,omitempty"` // Uncertainty of A1 (au/d^2)
	A2Sigma *float64 `json:"A2_sigma,omitempty"` // Uncertainty of A2 (au/d^2)
	A3Sigma *float64 `json:"A3_sigma,omitempty"` // Uncertainty of A3 (au/d^2)
	DTSigma *float64 `json:"DT_sigma,omitempty"` // Uncertainty of DT (days)
	S0Sigma *float64 `json:"S0_sigma,omitempty"` // Uncertainty of S0
}

// Physical contains physical and photometric parameters.
type Physical struct {
	H             *float64 `json:"H,omitempty"`              // Absolute magnitude
	G             *float64 `json:"G,omitempty"`              // Photometric slope parameter
	M1            *float64 `json:"M1,omitempty"`             // Photometric parameter M1
	K1            *float64 `json:"K1,omitempty"`             // Photometric parameter K1
	M2            *float64 `json:"M2,omitempty"`             // Photometric parameter M2
	K2            *float64 `json:"K2,omitempty"`             // Photometric parameter K2
	PC            *float64 `json:"PC,omitempty"`             // Photometric color index
	HSigma        *float64 `json:"H_sigma,omitempty"`        // Uncertainty in H
	Diameter      *float64 `json:"diameter,omitempty"`       // Effective diameter (km)
	Extent        *string  `json:"extent,omitempty"`         // Physical extent (km)
	GM            *float64 `json:"GM,omitempty"`             // Gravitational parameter (km^3/s^2)
	Density       *float64 `json:"density,omitempty"`        // Bulk density (g/cm^3)
	RotPer        *float64 `json:"rot_per,omitempty"`        // Rotation period (hours)
	Pole          *string  `json:"pole,omitempty"`           // Spin axis orientation
	Albedo        *float64 `json:"albedo,omitempty"`         // Geometric albedo
	BV            *float64 `json:"BV,omitempty"`             // B-V color index
	UB            *float64 `json:"UB,omitempty"`             // U-B color index
	IR            *float64 `json:"IR,omitempty"`             // IR color index
	SpecT         *string  `json:"spec_T,omitempty"`         // Spectral taxonomy
	SpecB         *string  `json:"spec_B,omitempty"`         // Spectral bin
	DiameterSigma *float64 `json:"diameter_sigma,omitempty"` // Uncertainty of diameter (km)
}
