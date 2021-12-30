/*
Package pnrm is the client.Panorama namespace.
*/
package pnrm

import (
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/PaloAltoNetworks/pango/pnrm/dg"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/gcp/account"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/gcp/gke/cluster"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/gcp/gke/cluster/group"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/swfwlicense/bootstrapdef"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/swfwlicense/manager"
	"github.com/PaloAltoNetworks/pango/pnrm/template"
	"github.com/PaloAltoNetworks/pango/pnrm/template/stack"
	"github.com/PaloAltoNetworks/pango/pnrm/template/variable"
)

// Panorama is the panorama.Panorama namespace.
type Panorama struct {
	DeviceGroup                *dg.Panorama
	GcpAccount                 *account.Panorama
	GkeCluster                 *cluster.Panorama
	GkeClusterGroup            *group.Panorama
	LicenseBootstrapDefinition *bootstrapdef.Panorama
	LicenseManager             *manager.Panorama
	Template                   *template.Panorama
	TemplateStack              *stack.Panorama
	TemplateVariable           *variable.Panorama
}

func PanoramaNamespace(x util.XapiClient) *Panorama {
	return &Panorama{
		DeviceGroup:                dg.PanoramaNamespace(x),
		GcpAccount:                 account.PanoramaNamespace(x),
		GkeCluster:                 cluster.PanoramaNamespace(x),
		GkeClusterGroup:            group.PanoramaNamespace(x),
		LicenseBootstrapDefinition: bootstrapdef.PanoramaNamespace(x),
		LicenseManager:             manager.PanoramaNamespace(x),
		Template:                   template.PanoramaNamespace(x),
		TemplateStack:              stack.PanoramaNamespace(x),
		TemplateVariable:           variable.PanoramaNamespace(x),
	}
}
