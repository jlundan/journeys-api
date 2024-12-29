package ggtfs

type Frequency struct {
	TripId      string
	StartTime   string
	EndTime     string
	HeadwaySecs uint
	ExactTimes  *int
	LineNumber  int
}
