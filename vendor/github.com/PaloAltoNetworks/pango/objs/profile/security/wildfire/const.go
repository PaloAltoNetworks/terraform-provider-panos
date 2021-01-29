package wildfire

// Valid values for Rule.Direction.
const (
	DirectionUpload   = "upload"
	DirectionDownload = "download"
	DirectionBoth     = "both"
)

// Valid values for Rule.Analysis.
const (
	AnalysisPublicCloud  = "public-cloud"
	AnalysisPrivateCloud = "private-cloud"
)

const (
	singular = "wildfire analysis security profile"
	plural   = "wildfire analysis security profiles"
)
