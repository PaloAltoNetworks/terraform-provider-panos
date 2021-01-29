package action

// Valid values for ActionType.
const (
	ActionTypeTagging     = "tagging"
	ActionTypeIntegration = "integration"
)

// Valid values for Action.
const (
	ActionAddTag    = "add-tag"
	ActionRemoveTag = "remove-tag"
	ActionAzure     = "Azure-Security-Center-Integration"
)

// Valid values for Target.
const (
	TargetSource      = "source-address"
	TargetDestination = "destination-address"
)

// Valid values for Registration.
const (
	RegistrationLocal    = "localhost"
	RegistrationRemote   = "remote"
	RegistrationPanorama = "panorama"
)

const (
	singular = "log forwarding profile match list action"
	plural   = "log forwarding profile match list actions"
)
