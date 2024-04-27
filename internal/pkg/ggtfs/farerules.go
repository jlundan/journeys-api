package ggtfs

type FareRule struct {
	Id            string
	RouteId       *string
	OriginId      *string
	DestinationId *string
	ContainsId    *string
	LineNumber    int
}
