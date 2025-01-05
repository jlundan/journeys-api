package tre

import (
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs"
)

const GtfsEnvKey = "JOURNEYS_GTFS_PATH"

func NewContext(gtfsPath string) model.Context {
	ctx := NewGTFSContextForDirectory(gtfsPath)

	municipalities := buildMunicipalities(*ctx.Municipalities)
	stopPoints := buildStopPoints(*ctx, municipalities)
	lines := buildLines(*ctx)
	routes := buildRoutes(*ctx)
	journeys, journeyPatterns := buildJourneys(*ctx, lines, routes, stopPoints)

	var errs []string
	for i, e := range ctx.Errors {
		errs[i] = e.Error()
	}

	var violations []string
	var recommendations []string
	var infos []string

	for i, v := range ctx.ValidationNotices {
		switch v.Severity() {
		case ggtfs.SeverityViolation:
			violations[i] = fmt.Sprintf("[%v] %s", v.Severity(), v.Code())
		case ggtfs.SeverityRecommendation:
			recommendations[i] = fmt.Sprintf("[%v] %s", v.Severity(), v.Code())
		case ggtfs.SeverityInfo:
			infos[i] = fmt.Sprintf("[%v] %s", v.Severity(), v.Code())
		}
	}

	return Context{
		lines:           lines,
		journeyPatterns: journeyPatterns,
		stopPoints:      stopPoints,
		municipalities:  municipalities,
		journeys:        journeys,
		routes:          routes,
		parseErrors:     errs,
		violations:      violations,
		recommendations: recommendations,
		infos:           infos,
	}
}

type Context struct {
	lines           Lines
	journeyPatterns JourneyPatterns
	stopPoints      StopPoints
	municipalities  Municipalities
	journeys        Journeys
	routes          Routes
	parseErrors     []string
	violations      []string
	recommendations []string
	infos           []string
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
	return context.parseErrors
}
func (context Context) GetViolations() []string {
	if len(context.violations) == 0 {
		return []string{}
	}
	return context.violations
}
func (context Context) GetRecommendations() []string {
	if len(context.recommendations) == 0 {
		return []string{}
	}
	return context.recommendations
}
func (context Context) GetInfos() []string {
	if len(context.infos) == 0 {
		return []string{}
	}
	return context.infos
}
