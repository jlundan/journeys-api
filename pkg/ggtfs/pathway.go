package ggtfs

type Pathway struct {
	Id                  string
	FromStopId          string
	ToStopId            string
	PathwayMode         uint
	IsBidirectional     uint
	Length              float64
	TraversalTime       uint
	StairCount          uint
	MaxSlope            float64
	MinWidth            float64
	SignpostedAs        string
	ReverseSignpostedAs string
	LineNumber          int
}
