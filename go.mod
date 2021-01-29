module github.com/terraform-providers/terraform-provider-panos

require (
	github.com/PaloAltoNetworks/pango v0.5.1
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.4.0
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586
)

//replace github.com/PaloAltoNetworks/pango => ../pango

go 1.13
