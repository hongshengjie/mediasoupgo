package mediasoupgo

type ScalabilityMode struct {
	SpatialLayers  int  `json:"spatialLayers"`
	TemporalLayers int  `json:"temporalLayers"`
	Ksvc           bool `json:"ksvc"`
}
