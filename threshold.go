package manta

type ThresholdType string

const (
	NoDate    = "nodata"
	GreatThan = "gt"
	Equal     = "eq"
	NotEqual  = "ne"
	LessThan  = "lt"
	Inside    = "inside"
	Outside   = "outside"
)

var (
	thresholdTypes = []string{
		NoDate,
		GreatThan,
		Equal,
		NotEqual,
		LessThan,
		Inside,
		Outside,
	}
)
