package provider

import (
	"context"

	"github.com/PaloAltoNetworks/pango/panorama/template"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TemplateCustom stores state shared between PreCreate and PostCreate hooks.
type TemplateCustom struct {
	savedDefaultVsys *string
}

func NewTemplateCustom(data *ProviderData) (*TemplateCustom, error) {
	return &TemplateCustom{}, nil
}

// PreCreate saves and strips default_vsys from the SDK object before Create,
// because PAN-OS cannot set this field during initial creation.
func (o *TemplateResource) PreCreate(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
	state *TemplateResourceModel,
	location template.Location,
	obj *template.Entry,
	ev *EncryptedValuesManager,
) {
	o.custom.savedDefaultVsys = obj.DefaultVsys
	obj.DefaultVsys = nil
}

// PostCreate creates the vsys referenced by default_vsys inside the template,
// then sets default_vsys via an Update call. The vsys is added directly to the
// SDK entry's Config struct so the Update's edit action includes it.
func (o *TemplateResource) PostCreate(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
	state *TemplateResourceModel,
	location template.Location,
	obj *template.Entry,
	ev *EncryptedValuesManager,
) {
	if o.custom.savedDefaultVsys == nil {
		return
	}

	defaultVsys := *o.custom.savedDefaultVsys
	templateName := state.Name.ValueString()

	tflog.Info(ctx, "performing post-create update to set default_vsys", map[string]any{
		"resource_name": "panos_template",
		"name":          templateName,
		"default_vsys":  defaultVsys,
	})

	// Populate the Config struct with the vsys entry so the SDK's edit
	// action includes it. PAN-OS requires the vsys to exist before it
	// accepts it as a valid default_vsys reference.
	obj.Config = &template.Config{
		Devices: []template.ConfigDevices{
			{
				Name: "localhost.localdomain",
				Vsys: []template.ConfigDevicesVsys{
					{Name: defaultVsys},
				},
			},
		},
	}
	obj.DefaultVsys = o.custom.savedDefaultVsys

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath for post-create update", err.Error())
		return
	}

	updated, err := o.manager.Update(ctx, location, components, obj, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting default_vsys after create",
			"Template created successfully but setting default_vsys failed. "+
				"Run terraform apply again to retry. Error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(state.CopyFromPango(ctx, o.client, nil, updated, ev)...)
}
