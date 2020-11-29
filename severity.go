package manta

const (
	OK   = "ok"
	Info = "info"
	Warn = "warn"
	High = "high"
	Crit = "crit"
)

var SeverityValue = map[string]int{
	OK:   0,
	Info: 1,
	Warn: 2,
	High: 3,
	Crit: 4,
}

var Severities = []string{
	OK,
	Info,
	Warn,
	High,
	Crit,
}
