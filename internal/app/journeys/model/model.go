package model

type Context interface {
	Lines() Lines
	JourneyPatterns() JourneyPatterns
	StopPoints() StopPoints
	Municipalities() Municipalities
	Journeys() Journeys
	Routes() Routes
	GetParseErrors() []string
	GetViolations() []string
	GetRecommendations() []string
	GetInfos() []string
}

type Lines interface {
	GetOne(string) (*Line, error)
	GetAll() []*Line
}

type JourneyPatterns interface {
	GetOne(string) (*JourneyPattern, error)
	GetAll() []*JourneyPattern
}

type StopPoints interface {
	GetOne(string) (*StopPoint, error)
	GetAll() []*StopPoint
}

type Municipalities interface {
	GetOne(string) (*Municipality, error)
	GetAll() []*Municipality
}

type Journeys interface {
	GetOne(string) (*Journey, error)
	GetOneByActivityId(string) (*Journey, error)
	GetAll() []*Journey
}

type Routes interface {
	GetOne(string) (*Route, error)
	GetAll() []*Route
}

type Line struct {
	Name        string
	Description string
}

type JourneyPattern struct {
	Id         string
	Name       string
	StopPoints []*StopPoint
	Route      *Route
	Journeys   []*Journey
}

type StopPoint struct {
	Name         string
	ShortName    string
	Latitude     float64
	Longitude    float64
	TariffZone   string
	Municipality *Municipality
}

type Municipality struct {
	PublicCode string
	Name       string
}

type Journey struct {
	Id                   string
	HeadSign             string
	Direction            string
	WheelchairAccessible bool
	GtfsInfo             *JourneyGtfsInfo
	DayTypes             []string
	DayTypeExceptions    []*DayTypeException
	Calls                []*JourneyCall
	Line                 *Line
	JourneyPattern       *JourneyPattern
	ValidFrom            string
	ValidTo              string
	Route                *Route
	ArrivalTime          string
	DepartureTime        string
	ActivityId           string
}

type JourneyGtfsInfo struct {
	TripId string
}

type DayTypeException struct {
	From string
	To   string
	Runs bool
}

type JourneyCall struct {
	DepartureTime string
	ArrivalTime   string
	StopPoint     *StopPoint
}

type Route struct {
	Id              string
	Line            *Line
	Name            string
	JourneyPatterns []*JourneyPattern
	Journeys        []*Journey
	GeoProjection   string
}
