package tre

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs"
)

const GtfsEnvKey = "JOURNEYS_GTFS_PATH"

func NewContext(gtfsPath string) model.Context {
	g := NewGTFSContextForDirectory(gtfsPath)

	municipalities := buildMunicipalities(*g.Municipalities)
	stopPoints := buildStopPoints(*g, municipalities)
	lines := buildLines(*g)
	routes := buildRoutes(*g)
	journeys, journeyPatterns := buildJourneys(*g, lines, routes, stopPoints)

	return Context{
		lines:             lines,
		journeyPatterns:   journeyPatterns,
		stopPoints:        stopPoints,
		municipalities:    municipalities,
		journeys:          journeys,
		routes:            routes,
		parseErrors:       g.Errors,
		validationNotices: g.ValidationNotices,
	}
}

type Context struct {
	lines             Lines
	journeyPatterns   JourneyPatterns
	stopPoints        StopPoints
	municipalities    Municipalities
	journeys          Journeys
	routes            Routes
	parseErrors       []error
	validationNotices []ggtfs.ValidationNotice
}

func (context Context) Lines() model.Lines {
	return context.lines
}
func (context Context) JourneyPatterns() model.JourneyPatterns {
	return context.journeyPatterns
}
func (context Context) StopPoints() model.StopPoints {
	return context.stopPoints
}
func (context Context) Municipalities() model.Municipalities {
	return context.municipalities
}
func (context Context) Journeys() model.Journeys {
	return context.journeys
}
func (context Context) Routes() model.Routes {
	return context.routes
}
func (context Context) GetParseErrors() []string {
	if len(context.parseErrors) == 0 {
		return []string{}
	}

	errs := make([]string, len(context.parseErrors))
	for i, e := range context.parseErrors {

		errs[i] = e.Error()
	}

	return errs
}
func (context Context) GetViolations() []string {
	if len(context.validationNotices) == 0 {
		return []string{}
	}

	var violations []string
	for i, v := range context.validationNotices {
		if v.Severity() == ggtfs.SeverityViolation {
			violations[i] = fmt.Sprintf("[%v] %s", v.Severity(), v.Code())
		}
	}

	return violations
}
func (context Context) GetRecommendations() []string {
	if len(context.validationNotices) == 0 {
		return []string{}
	}

	var recommendations []string
	for i, v := range context.validationNotices {
		if v.Severity() == ggtfs.SeverityRecommendation {
			recommendations[i] = fmt.Sprintf("[%v] %s", v.Severity(), v.Code())
		}
	}

	return recommendations
}

func (context Context) GetInfos() []string {
	if len(context.validationNotices) == 0 {
		return []string{}
	}

	var infos []string
	for i, v := range context.validationNotices {
		if v.Severity() == ggtfs.SeverityInfo {
			infos[i] = fmt.Sprintf("[%v] %s", v.Severity(), v.Code())
		}
	}

	return infos
}
