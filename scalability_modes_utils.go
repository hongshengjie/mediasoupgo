package mediasoupgo

import (
	"regexp"
	"strconv"
)

// ScalabilityModeRegex is a compiled regular expression for parsing scalability modes.
// It matches patterns like "L1T1", "S3T2", "L10T5_KEY", etc.
var ScalabilityModeRegex = regexp.MustCompile(`^[LS]([1-9]\d{0,1})T([1-9]\d{0,1})(_KEY)?$`)

// ParseScalabilityMode parses a scalability mode string into a ScalabilityMode struct.
// The input string should match the format defined by ScalabilityModeRegex (e.g., "L3T2", "S1T1_KEY").
// If the input is empty or invalid, it returns a default ScalabilityMode with 1 spatial layer,
// 1 temporal layer, and KSVC disabled.
func ParseScalabilityMode(scalabilityMode string) ScalabilityMode {
	// If scalabilityMode is empty, use an empty string as fallback (equivalent to TypeScript's ?? '')
	if scalabilityMode == "" {
		scalabilityMode = ""
	}

	// Execute the regex match
	match := ScalabilityModeRegex.FindStringSubmatch(scalabilityMode)

	if match != nil {
		// Convert matched groups to integers
		spatialLayers, _ := strconv.Atoi(match[1])  // First group: spatial layers
		temporalLayers, _ := strconv.Atoi(match[2]) // Second group: temporal layers
		ksvc := match[3] != ""                      // Third group: "_KEY" (optional)

		return ScalabilityMode{
			SpatialLayers:  spatialLayers,
			TemporalLayers: temporalLayers,
			Ksvc:           ksvc,
		}
	}

	// Default values if no match
	return ScalabilityMode{
		SpatialLayers:  1,
		TemporalLayers: 1,
		Ksvc:           false,
	}
}
