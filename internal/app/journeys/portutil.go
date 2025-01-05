package journeys

import (
	"strconv"
)

func ParsePort(envPort string, defaultPort int) (int, error) {
	var serverPort int

	if envPort == "" {
		serverPort = defaultPort
	} else {
		p, err := strconv.Atoi(envPort)
		if err != nil {
			return 0, err
		}
		serverPort = p
	}

	return serverPort, nil
}
