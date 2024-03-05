package provider

import (
	"context"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type genericTfid struct {
	Name     *string        `json:"name,omitempty"`
	Names    []string       `json:"names,omitempty"`
	Location map[string]any `json:"location"`
}

func (g genericTfid) IsValid() error { return nil }

// Data source.
var (
	_ datasource.DataSource              = &tfidDataSource{}
	_ datasource.DataSourceWithConfigure = &tfidDataSource{}
)

func NewTfidDataSource() datasource.DataSource {
	return &tfidDataSource{}
}

type tfidDataSource struct {
	client *pango.XmlApiClient
}

type tfidDsModel struct {
	Location  types.String `tfsdk:"location"`
	Variables types.Map    `tfsdk:"variables"`
	Name      types.String `tfsdk:"name"`
	Names     types.List   `tfsdk:"names"`

	Tfid types.String `tfsdk:"tfid"`
}

func (d *tfidDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tfid"
}

func (d *tfidDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Helper data source: create a tfid from the given information. Note that the tfid ouptut from this data source may not exactly match what a resource uses, but it will still be a valid ID to use for resource imports.",

		Attributes: map[string]dsschema.Attribute{
			"location": dsschema.StringAttribute{
				Description: "The location path name.",
				Required:    true,
			},
			"variables": dsschema.MapAttribute{
				Description: "The variables and values for the specified location.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"name": dsschema.StringAttribute{
				Description: "(Singleton resource) The config's name.",
				Optional:    true,
			},
			"names": dsschema.ListAttribute{
				Description: "(Grouping resources) The names of the configs.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"tfid": dsschema.StringAttribute{
				Description: "The tfid created from the given parts.",
				Computed:    true,
			},
		},
	}
}

func (d *tfidDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.XmlApiClient)
}

func (d *tfidDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tfidDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source read", map[string]any{
		"data_source_name": "panos_tfid",
	})

	vars := make(map[string]types.String, len(state.Variables.Elements()))
	resp.Diagnostics.Append(state.Variables.ElementsAs(ctx, &vars, false).Errors()...)

	var names []string
	resp.Diagnostics.Append(state.Names.ElementsAs(ctx, &names, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	loc := genericTfid{
		Name:  state.Name.ValueStringPointer(),
		Names: append([]string(nil), names...),
	}

	if len(vars) == 0 {
		loc.Location = map[string]any{
			state.Location.ValueString(): true,
		}
	} else {
		content := make(map[string]string)
		for key, value := range vars {
			content[key] = value.ValueString()
		}
		loc.Location = map[string]any{
			state.Location.ValueString(): content,
		}
	}

	// Encode the tfid from the info given.
	idstr, err := EncodeLocation(loc)
	if err != nil {
		resp.Diagnostics.AddError("error encoding tfid", err.Error())
		return
	}

	// Set the tfid param.
	state.Tfid = types.StringValue(idstr)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
