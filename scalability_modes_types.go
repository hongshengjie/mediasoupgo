package mediasoupgo

// ScalabilityMode represents the scalability mode for SVC (Scalable Video Coding).
type ScalabilityMode struct {
	// Number of spatial layers.
	SpatialLayers int
	// Number of temporal layers.
	TemporalLayers int
	// Whether key-frame SVC (KSVC) is enabled.
	Ksvc bool
}
