package tre

import (
	"errors"
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/pkg/csv"
	"os"
)

const MunicipalityFileName = "municipalities.txt"
const GtfsEnvKey = "JOURNEYS_GTFS_PATH"

type municipalityData struct {
	municipalityHeaders map[string]uint8
	municipalityRows    [][]string
}

func NewContext() (model.Context, []error, []error) {
	var (
		errs     []error
		warnings []error
	)

	if os.Getenv(GtfsEnvKey) == "" {
		errs = append(errs, errors.New(fmt.Sprintf("%v not set in environment", GtfsEnvKey)))
	}

	if len(errs) > 0 {
		return Context{}, errs, warnings
	}

	g, gtfsErrors := NewGTFSContextForDirectory(os.Getenv(GtfsEnvKey))
	if len(gtfsErrors) > 0 {
		errs = append(errs, gtfsErrors...)
	}
	gtfsWarnings := Validate(g)
	if len(gtfsWarnings) > 0 {
		warnings = append(warnings, gtfsWarnings...)
	}

	m, err := readMunicipalities()
	if err != nil {
		errs = append(errs, err)
	}

	if m == nil && g == nil {
		return Context{}, errs, warnings
	}

	if m == nil {
		return Context{
			lines:  buildLines(*g),
			routes: buildRoutes(*g),
		}, errs, warnings
	}

	municipalities := buildMunicipalities(*m)

	if g == nil {
		return Context{
			municipalities: municipalities,
		}, errs, warnings
	}

	stopPoints := buildStopPoints(*g, municipalities)
	lines := buildLines(*g)
	routes := buildRoutes(*g)
	journeys, journeyPatterns := buildJourneys(*g, lines, routes, stopPoints)

	return Context{
		lines:           lines,
		journeyPatterns: journeyPatterns,
		stopPoints:      stopPoints,
		municipalities:  municipalities,
		journeys:        journeys,
		routes:          routes,
	}, errs, warnings
}

type Context struct {
	lines           Lines
	journeyPatterns JourneyPatterns
	stopPoints      StopPoints
	municipalities  Municipalities
	journeys        Journeys
	routes          Routes
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

func readMunicipalities() (*municipalityData, error) {
	var err error
	m := &municipalityData{}
	m.municipalityHeaders, m.municipalityRows, err = csv.ParseFile(fmt.Sprintf("%v/%v", os.Getenv(GtfsEnvKey), MunicipalityFileName), true)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return m, nil
}
