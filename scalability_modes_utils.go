package mediasoupgo

import (
	"regexp"
	"strconv"
)

var ScalabilityModeRegex = regexp.MustCompile(`^[LS]([1-9]\d{0,1})T([1-9]\d{0,1})(_KEY)?`)

func parseScalabilityMode(scalabilityMode string) ScalabilityMode {
	match := ScalabilityModeRegex.FindStringSubmatch(scalabilityMode)

	if match != nil {
		spatialLayers, _ := strconv.Atoi(match[1])
		temporalLayers, _ := strconv.Atoi(match[2])
		ksvc := match[3] != ""
		return ScalabilityMode{
			SpatialLayers:  spatialLayers,
			TemporalLayers: temporalLayers,
			Ksvc:           ksvc,
		}
	} else {
		return ScalabilityMode{
			SpatialLayers:  1,
			TemporalLayers: 1,
			Ksvc:           false,
		}
	}
}
