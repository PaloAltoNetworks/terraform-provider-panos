package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	decryption "github.com/PaloAltoNetworks/pango/objects/profiles/decryption"
	pangoutil "github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdkmanager "github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

// -----------------------------------------------------------------------
// DataSource
// -----------------------------------------------------------------------

var (
	_ datasource.DataSource              = &DecryptionProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &DecryptionProfileDataSource{}
)

func NewDecryptionProfileDataSource() datasource.DataSource {
	return &DecryptionProfileDataSource{}
}

type DecryptionProfileDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*decryption.Entry, decryption.Location, *decryption.Service]
}

type DecryptionProfileDataSourceModel struct {
	Location             types.Object `tfsdk:"location"`
	Name                 types.String `tfsdk:"name"`
	SslForwardProxy      types.Object `tfsdk:"ssl_forward_proxy"`
	SslInboundInspection types.Object `tfsdk:"ssl_inbound_inspection"`
	SslNoProxy           types.Object `tfsdk:"ssl_no_proxy"`
	SslProtocolSettings  types.Object `tfsdk:"ssl_protocol_settings"`
}

func (o DecryptionProfileDataSourceModel) AncestorName() string { return "" }
func (o DecryptionProfileDataSourceModel) EntryName() *string    { return nil }

func (o *DecryptionProfileDataSourceModel) CopyToPango(ctx context.Context, client pangoutil.PangoClient, ancestors []Ancestor, obj **decryption.Entry, ev *EncryptedValuesManager) diag.Diagnostics {
	var diags diag.Diagnostics
	if *obj == nil {
		*obj = new(decryption.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).SslForwardProxy = copySslForwardProxyToPango(ctx, o.SslForwardProxy, &diags)
	(*obj).SslInboundInspection = copySslInboundInspectionToPango(ctx, o.SslInboundInspection, &diags)
	(*obj).SslNoProxy = copySslNoProxyToPango(ctx, o.SslNoProxy, &diags)
	(*obj).SslProtocolSettings = copySslProtocolSettingsToPango(ctx, o.SslProtocolSettings, &diags)
	return diags
}

func (o *DecryptionProfileDataSourceModel) CopyFromPango(ctx context.Context, client pangoutil.PangoClient, ancestors []Ancestor, obj *decryption.Entry, ev *EncryptedValuesManager) diag.Diagnostics {
	var diags diag.Diagnostics
	o.Name = types.StringValue(obj.Name)

	var d diag.Diagnostics
	o.SslForwardProxy, d = copySslForwardProxyFromPango(ctx, obj.SslForwardProxy)
	diags.Append(d...)
	o.SslInboundInspection, d = copySslInboundInspectionFromPango(ctx, obj.SslInboundInspection)
	diags.Append(d...)
	o.SslNoProxy, d = copySslNoProxyFromPango(ctx, obj.SslNoProxy)
	diags.Append(d...)
	o.SslProtocolSettings, d = copySslProtocolSettingsFromPango(ctx, obj.SslProtocolSettings)
	diags.Append(d...)
	return diags
}

func (o *DecryptionProfileDataSourceModel) resourceXpathParentComponents() ([]string, error) {
	return []string{}, nil
}

func (o *DecryptionProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_decryption_profile"
}

func (o *DecryptionProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DecryptionProfileDataSourceSchema()
}

func (o *DecryptionProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	providerData := req.ProviderData.(*ProviderData)
	o.client = providerData.Client
	specifier, _, err := decryption.Versioning(o.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	batchSize := providerData.MultiConfigBatchSize
	o.manager = sdkmanager.NewEntryObjectManager[*decryption.Entry, decryption.Location, *decryption.Service](
		o.client, decryption.NewService(o.client), batchSize, specifier, decryption.SpecMatches,
	)
}

func (o *DecryptionProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DecryptionProfileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location decryption.Location
	resp.Diagnostics.Append(decryptionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "performing datasource read", map[string]any{
		"resource_name": "panos_decryption_profile_datasource",
		"function":      "Read",
		"name":          state.Name.ValueString(),
	})

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}
	object, err := o.manager.Read(ctx, location, components, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading entry", err.Error())
		return
	}

	resp.Diagnostics.Append(state.CopyFromPango(ctx, o.client, nil, object, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// -----------------------------------------------------------------------
// Resource
// -----------------------------------------------------------------------

var (
	_ resource.Resource                = &DecryptionProfileResource{}
	_ resource.ResourceWithConfigure   = &DecryptionProfileResource{}
	_ resource.ResourceWithImportState = &DecryptionProfileResource{}
)

func NewDecryptionProfileResource() resource.Resource {
	if _, found := resourceFuncMap["panos_decryption_profile"]; !found {
		resourceFuncMap["panos_decryption_profile"] = resourceFuncs{
			CreateImportId: DecryptionProfileImportStateCreator,
		}
	}
	return &DecryptionProfileResource{}
}

type DecryptionProfileResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*decryption.Entry, decryption.Location, *decryption.Service]
}

type DecryptionProfileResourceModel struct {
	Location             types.Object `tfsdk:"location"`
	Name                 types.String `tfsdk:"name"`
	SslForwardProxy      types.Object `tfsdk:"ssl_forward_proxy"`
	SslInboundInspection types.Object `tfsdk:"ssl_inbound_inspection"`
	SslNoProxy           types.Object `tfsdk:"ssl_no_proxy"`
	SslProtocolSettings  types.Object `tfsdk:"ssl_protocol_settings"`
}

func (o DecryptionProfileResourceModel) AncestorName() string { return "" }
func (o DecryptionProfileResourceModel) EntryName() *string    { return nil }

func (o *DecryptionProfileResourceModel) ValidateConfig(ctx context.Context, resp *resource.ValidateConfigResponse, p path.Path) {
}

func (o *DecryptionProfileResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var r DecryptionProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &r)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.ValidateConfig(ctx, resp, path.Empty())
}

func (o *DecryptionProfileResourceModel) CopyToPango(ctx context.Context, client pangoutil.PangoClient, ancestors []Ancestor, obj **decryption.Entry, ev *EncryptedValuesManager) diag.Diagnostics {
	var diags diag.Diagnostics
	if *obj == nil {
		*obj = new(decryption.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).SslForwardProxy = copySslForwardProxyToPango(ctx, o.SslForwardProxy, &diags)
	(*obj).SslInboundInspection = copySslInboundInspectionToPango(ctx, o.SslInboundInspection, &diags)
	(*obj).SslNoProxy = copySslNoProxyToPango(ctx, o.SslNoProxy, &diags)
	(*obj).SslProtocolSettings = copySslProtocolSettingsToPango(ctx, o.SslProtocolSettings, &diags)
	return diags
}

func (o *DecryptionProfileResourceModel) CopyFromPango(ctx context.Context, client pangoutil.PangoClient, ancestors []Ancestor, obj *decryption.Entry, ev *EncryptedValuesManager) diag.Diagnostics {
	var diags diag.Diagnostics
	o.Name = types.StringValue(obj.Name)

	var d diag.Diagnostics
	o.SslForwardProxy, d = copySslForwardProxyFromPango(ctx, obj.SslForwardProxy)
	diags.Append(d...)
	o.SslInboundInspection, d = copySslInboundInspectionFromPango(ctx, obj.SslInboundInspection)
	diags.Append(d...)
	o.SslNoProxy, d = copySslNoProxyFromPango(ctx, obj.SslNoProxy)
	diags.Append(d...)
	o.SslProtocolSettings, d = copySslProtocolSettingsFromPango(ctx, obj.SslProtocolSettings)
	diags.Append(d...)
	return diags
}

func (o *DecryptionProfileResourceModel) resourceXpathParentComponents() ([]string, error) {
	return []string{}, nil
}

func (o *DecryptionProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_decryption_profile"
}

func (o *DecryptionProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DecryptionProfileResourceSchema()
}

func (o *DecryptionProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	providerData := req.ProviderData.(*ProviderData)
	o.client = providerData.Client
	specifier, _, err := decryption.Versioning(o.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	batchSize := providerData.MultiConfigBatchSize
	o.manager = sdkmanager.NewEntryObjectManager[*decryption.Entry, decryption.Location, *decryption.Service](
		o.client, decryption.NewService(o.client), batchSize, specifier, decryption.SpecMatches,
	)
}

func (o *DecryptionProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state DecryptionProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ev, err := NewEncryptedValuesManager(nil, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed to init encrypted values manager", err.Error())
		return
	}

	var location decryption.Location
	resp.Diagnostics.Append(decryptionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	if o.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	var obj *decryption.Entry
	resp.Diagnostics.Append(state.CopyToPango(ctx, o.client, nil, &obj, ev)...)
	if resp.Diagnostics.HasError() {
		return
	}

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}
	created, err := o.manager.Create(ctx, location, components, obj)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	resp.Diagnostics.Append(state.CopyFromPango(ctx, o.client, nil, created, ev)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := json.Marshal(ev)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal encrypted values state", err.Error())
		return
	}
	resp.Private.SetKey(ctx, "encrypted_values", payload)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (o *DecryptionProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DecryptionProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	encryptedValues, diags := req.Private.GetKey(ctx, "encrypted_values")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ev, err := NewEncryptedValuesManager(encryptedValues, true)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read encrypted values from private state", err.Error())
		return
	}

	var location decryption.Location
	resp.Diagnostics.Append(decryptionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_decryption_profile_resource",
		"function":      "Read",
		"name":          state.Name.ValueString(),
	})

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}
	object, err := o.manager.Read(ctx, location, components, state.Name.ValueString())
	if err != nil {
		if errors.Is(err, sdkmanager.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error reading entry", err.Error())
		}
		return
	}

	resp.Diagnostics.Append(state.CopyFromPango(ctx, o.client, nil, object, ev)...)

	state.Location = state.Location

	payload, err := json.Marshal(ev)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal encrypted values state", err.Error())
		return
	}
	resp.Private.SetKey(ctx, "encrypted_values", payload)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (o *DecryptionProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DecryptionProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	encryptedValues, diags := req.Private.GetKey(ctx, "encrypted_values")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ev, err := NewEncryptedValuesManager(encryptedValues, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read encrypted values from private state", err.Error())
		return
	}

	var location decryption.Location
	resp.Diagnostics.Append(decryptionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_decryption_profile_resource",
		"function":      "Update",
	})

	if o.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}

	var obj *decryption.Entry
	if state.Name.ValueString() != plan.Name.ValueString() {
		obj, err = o.manager.Read(ctx, location, components, state.Name.ValueString())
	} else {
		obj, err = o.manager.Read(ctx, location, components, plan.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.CopyToPango(ctx, o.client, nil, &obj, ev)...)
	if resp.Diagnostics.HasError() {
		return
	}

	components, err = plan.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}

	var newName string
	if state.Name.ValueString() != plan.Name.ValueString() {
		newName = plan.Name.ValueString()
		obj.Name = state.Name.ValueString()
	}

	updated, err := o.manager.Update(ctx, location, components, obj, newName)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.CopyFromPango(ctx, o.client, nil, updated, ev)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Location = state.Location

	payload, err := json.Marshal(ev)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal encrypted values state", err.Error())
		return
	}
	resp.Private.SetKey(ctx, "encrypted_values", payload)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (o *DecryptionProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DecryptionProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location decryption.Location
	resp.Diagnostics.Append(decryptionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}
	err = o.manager.Delete(ctx, location, components, []string{state.Name.ValueString()})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

// -----------------------------------------------------------------------
// ImportState
// -----------------------------------------------------------------------

type DecryptionProfileImportState struct {
	Location types.Object `json:"location"`
	Name     types.String `json:"name"`
}

func (o DecryptionProfileImportState) MarshalJSON() ([]byte, error) {
	type shadow struct {
		Location interface{} `json:"location"`
		Name     *string     `json:"name"`
	}
	var location_object interface{}
	{
		var err error
		location_object, err = TypesObjectToMap(o.Location, DecryptionProfileResourceLocationSchema())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal location into JSON document: %w", err)
		}
	}
	obj := shadow{
		Location: location_object,
		Name:     o.Name.ValueStringPointer(),
	}
	return json.Marshal(obj)
}

func (o *DecryptionProfileImportState) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Location interface{} `json:"location"`
		Name     *string     `json:"name"`
	}
	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}
	var location_object types.Object
	{
		location_map, ok := shadow.Location.(map[string]interface{})
		if !ok {
			return NewDiagnosticsError("Failed to unmarshal JSON document into location: expected map[string]interface{}", nil)
		}
		var err error
		location_object, err = MapToTypesObject(location_map, DecryptionProfileResourceLocationSchema())
		if err != nil {
			return fmt.Errorf("failed to unmarshal location from JSON: %w", err)
		}
	}
	o.Location = location_object
	o.Name = types.StringPointerValue(shadow.Name)
	return nil
}

func DecryptionProfileImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}
	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}
	var location types.Object
	switch value := locationAttr.(type) {
	case types.Object:
		location = value
	default:
		return nil, fmt.Errorf("location attribute expected to be an object")
	}
	nameAttr, ok := attrs["name"]
	if !ok {
		return nil, fmt.Errorf("name attribute missing")
	}
	var name types.String
	switch value := nameAttr.(type) {
	case types.String:
		name = value
	default:
		return nil, fmt.Errorf("name attribute expected to be a string")
	}
	importStruct := DecryptionProfileImportState{
		Location: location,
		Name:     name,
	}
	return json.Marshal(importStruct)
}

func (o *DecryptionProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var obj DecryptionProfileImportState
	data, err := base64.StdEncoding.DecodeString(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to decode Import ID", err.Error())
		return
	}
	if err = json.Unmarshal(data, &obj); err != nil {
		var diagsErr *DiagnosticsError
		if errors.As(err, &diagsErr) {
			resp.Diagnostics.Append(diagsErr.Diagnostics()...)
		} else {
			resp.Diagnostics.AddError("Failed to unmarshal Import ID", err.Error())
		}
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("location"), obj.Location)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), obj.Name)...)
}

// -----------------------------------------------------------------------
// Schemas
// -----------------------------------------------------------------------

func DecryptionProfileDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{
			"location": DecryptionProfileDataSourceLocationSchema(),
			"name": dsschema.StringAttribute{
				Description: "Decryption profile name.",
				Required:    true,
			},
			"ssl_forward_proxy":      dsschema.SingleNestedAttribute{Description: "SSL forward proxy settings.", Optional: true, Computed: true, Attributes: decryptionProfileSslForwardProxyDsSchema()},
			"ssl_inbound_inspection": dsschema.SingleNestedAttribute{Description: "SSL inbound proxy (inbound inspection) settings.", Optional: true, Computed: true, Attributes: decryptionProfileSslInboundInspectionDsSchema()},
			"ssl_no_proxy":           dsschema.SingleNestedAttribute{Description: "SSL no-proxy settings.", Optional: true, Computed: true, Attributes: decryptionProfileSslNoProxyDsSchema()},
			"ssl_protocol_settings":  dsschema.SingleNestedAttribute{Description: "SSL protocol settings.", Optional: true, Computed: true, Attributes: decryptionProfileSslProtocolSettingsDsSchema()},
		},
	}
}

func DecryptionProfileResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{
			"location": DecryptionProfileResourceLocationSchema(),
			"name": rsschema.StringAttribute{
				Description: "Decryption profile name.",
				Required:    true,
			},
			"ssl_forward_proxy":      rsschema.SingleNestedAttribute{Description: "SSL forward proxy settings.", Optional: true, Attributes: decryptionProfileSslForwardProxyRsSchema()},
			"ssl_inbound_inspection": rsschema.SingleNestedAttribute{Description: "SSL inbound proxy (inbound inspection) settings.", Optional: true, Attributes: decryptionProfileSslInboundInspectionRsSchema()},
			"ssl_no_proxy":           rsschema.SingleNestedAttribute{Description: "SSL no-proxy settings.", Optional: true, Attributes: decryptionProfileSslNoProxyRsSchema()},
			"ssl_protocol_settings":  rsschema.SingleNestedAttribute{Description: "SSL protocol settings.", Optional: true, Attributes: decryptionProfileSslProtocolSettingsRsSchema()},
		},
	}
}

func decryptionProfileSslForwardProxyDsSchema() map[string]dsschema.Attribute {
	return map[string]dsschema.Attribute{
		"auto_include_altname":              dsschema.BoolAttribute{Description: "Automatically include the subject alternative name.", Optional: true, Computed: true},
		"block_client_cert":                 dsschema.BoolAttribute{Description: "Block sessions if client certificate is required but not received.", Optional: true, Computed: true},
		"block_expired_certificate":         dsschema.BoolAttribute{Description: "Block sessions with expired server certificates.", Optional: true, Computed: true},
		"block_timeout_cert":                dsschema.BoolAttribute{Description: "Block sessions if OCSP/CRL signer certificate is unavailable.", Optional: true, Computed: true},
		"block_tls13_downgrade_no_resource": dsschema.BoolAttribute{Description: "Block TLS 1.3 sessions that cannot be decrypted due to resource constraints.", Optional: true, Computed: true},
		"block_unknown_cert":                dsschema.BoolAttribute{Description: "Block sessions if certificate status cannot be determined.", Optional: true, Computed: true},
		"block_unsupported_cipher":          dsschema.BoolAttribute{Description: "Block sessions if cipher suite is not supported.", Optional: true, Computed: true},
		"block_unsupported_version":         dsschema.BoolAttribute{Description: "Block sessions if protocol version is not supported.", Optional: true, Computed: true},
		"block_untrusted_issuer":            dsschema.BoolAttribute{Description: "Block sessions with untrusted server certificate issuers.", Optional: true, Computed: true},
		"restrict_cert_exts":                dsschema.BoolAttribute{Description: "Restrict certificate extensions to those specified.", Optional: true, Computed: true},
		"strip_alpn":                        dsschema.BoolAttribute{Description: "Strip the ALPN extension from ClientHello.", Optional: true, Computed: true},
	}
}

func decryptionProfileSslForwardProxyRsSchema() map[string]rsschema.Attribute {
	return map[string]rsschema.Attribute{
		"auto_include_altname":              rsschema.BoolAttribute{Description: "Automatically include the subject alternative name.", Optional: true},
		"block_client_cert":                 rsschema.BoolAttribute{Description: "Block sessions if client certificate is required but not received.", Optional: true},
		"block_expired_certificate":         rsschema.BoolAttribute{Description: "Block sessions with expired server certificates.", Optional: true},
		"block_timeout_cert":                rsschema.BoolAttribute{Description: "Block sessions if OCSP/CRL signer certificate is unavailable.", Optional: true},
		"block_tls13_downgrade_no_resource": rsschema.BoolAttribute{Description: "Block TLS 1.3 sessions that cannot be decrypted due to resource constraints.", Optional: true},
		"block_unknown_cert":                rsschema.BoolAttribute{Description: "Block sessions if certificate status cannot be determined.", Optional: true},
		"block_unsupported_cipher":          rsschema.BoolAttribute{Description: "Block sessions if cipher suite is not supported.", Optional: true},
		"block_unsupported_version":         rsschema.BoolAttribute{Description: "Block sessions if protocol version is not supported.", Optional: true},
		"block_untrusted_issuer":            rsschema.BoolAttribute{Description: "Block sessions with untrusted server certificate issuers.", Optional: true},
		"restrict_cert_exts":                rsschema.BoolAttribute{Description: "Restrict certificate extensions to those specified.", Optional: true},
		"strip_alpn":                        rsschema.BoolAttribute{Description: "Strip the ALPN extension from ClientHello.", Optional: true},
	}
}

func decryptionProfileSslInboundInspectionDsSchema() map[string]dsschema.Attribute {
	return map[string]dsschema.Attribute{
		"block_if_hsm_unavailable":          dsschema.BoolAttribute{Description: "Block sessions if HSM is unavailable.", Optional: true, Computed: true},
		"block_if_no_resource":              dsschema.BoolAttribute{Description: "Block sessions if decryption resources are unavailable.", Optional: true, Computed: true},
		"block_tls13_downgrade_no_resource": dsschema.BoolAttribute{Description: "Block TLS 1.3 sessions that cannot be decrypted due to resource constraints.", Optional: true, Computed: true},
		"block_unsupported_cipher":          dsschema.BoolAttribute{Description: "Block sessions if cipher suite is not supported.", Optional: true, Computed: true},
		"block_unsupported_version":         dsschema.BoolAttribute{Description: "Block sessions if protocol version is not supported.", Optional: true, Computed: true},
	}
}

func decryptionProfileSslInboundInspectionRsSchema() map[string]rsschema.Attribute {
	return map[string]rsschema.Attribute{
		"block_if_hsm_unavailable":          rsschema.BoolAttribute{Description: "Block sessions if HSM is unavailable.", Optional: true},
		"block_if_no_resource":              rsschema.BoolAttribute{Description: "Block sessions if decryption resources are unavailable.", Optional: true},
		"block_tls13_downgrade_no_resource": rsschema.BoolAttribute{Description: "Block TLS 1.3 sessions that cannot be decrypted due to resource constraints.", Optional: true},
		"block_unsupported_cipher":          rsschema.BoolAttribute{Description: "Block sessions if cipher suite is not supported.", Optional: true},
		"block_unsupported_version":         rsschema.BoolAttribute{Description: "Block sessions if protocol version is not supported.", Optional: true},
	}
}

func decryptionProfileSslNoProxyDsSchema() map[string]dsschema.Attribute {
	return map[string]dsschema.Attribute{
		"block_client_cert":         dsschema.BoolAttribute{Description: "Block sessions if client certificate is required but not received.", Optional: true, Computed: true},
		"block_expired_certificate": dsschema.BoolAttribute{Description: "Block sessions with expired server certificates.", Optional: true, Computed: true},
		"block_timeout_cert":        dsschema.BoolAttribute{Description: "Block sessions if OCSP/CRL signer certificate is unavailable.", Optional: true, Computed: true},
		"block_unknown_cert":        dsschema.BoolAttribute{Description: "Block sessions if certificate status cannot be determined.", Optional: true, Computed: true},
		"block_unsupported_cipher":  dsschema.BoolAttribute{Description: "Block sessions if cipher suite is not supported.", Optional: true, Computed: true},
		"block_unsupported_version": dsschema.BoolAttribute{Description: "Block sessions if protocol version is not supported.", Optional: true, Computed: true},
		"block_untrusted_issuer":    dsschema.BoolAttribute{Description: "Block sessions with untrusted server certificate issuers.", Optional: true, Computed: true},
	}
}

func decryptionProfileSslNoProxyRsSchema() map[string]rsschema.Attribute {
	return map[string]rsschema.Attribute{
		"block_client_cert":         rsschema.BoolAttribute{Description: "Block sessions if client certificate is required but not received.", Optional: true},
		"block_expired_certificate": rsschema.BoolAttribute{Description: "Block sessions with expired server certificates.", Optional: true},
		"block_timeout_cert":        rsschema.BoolAttribute{Description: "Block sessions if OCSP/CRL signer certificate is unavailable.", Optional: true},
		"block_unknown_cert":        rsschema.BoolAttribute{Description: "Block sessions if certificate status cannot be determined.", Optional: true},
		"block_unsupported_cipher":  rsschema.BoolAttribute{Description: "Block sessions if cipher suite is not supported.", Optional: true},
		"block_unsupported_version": rsschema.BoolAttribute{Description: "Block sessions if protocol version is not supported.", Optional: true},
		"block_untrusted_issuer":    rsschema.BoolAttribute{Description: "Block sessions with untrusted server certificate issuers.", Optional: true},
	}
}

func decryptionProfileSslProtocolSettingsDsSchema() map[string]dsschema.Attribute {
	return map[string]dsschema.Attribute{
		"auth_algo_md5":       dsschema.BoolAttribute{Description: "Allow MD5 authentication algorithm.", Optional: true, Computed: true},
		"auth_algo_sha1":      dsschema.BoolAttribute{Description: "Allow SHA1 authentication algorithm.", Optional: true, Computed: true},
		"auth_algo_sha256":    dsschema.BoolAttribute{Description: "Allow SHA256 authentication algorithm.", Optional: true, Computed: true},
		"auth_algo_sha384":    dsschema.BoolAttribute{Description: "Allow SHA384 authentication algorithm.", Optional: true, Computed: true},
		"enc_algo_3des":       dsschema.BoolAttribute{Description: "Allow 3DES encryption algorithm.", Optional: true, Computed: true},
		"enc_algo_aes_128_cbc": dsschema.BoolAttribute{Description: "Allow AES-128-CBC encryption algorithm.", Optional: true, Computed: true},
		"enc_algo_aes_128_gcm": dsschema.BoolAttribute{Description: "Allow AES-128-GCM encryption algorithm.", Optional: true, Computed: true},
		"enc_algo_aes_256_cbc": dsschema.BoolAttribute{Description: "Allow AES-256-CBC encryption algorithm.", Optional: true, Computed: true},
		"enc_algo_aes_256_gcm": dsschema.BoolAttribute{Description: "Allow AES-256-GCM encryption algorithm.", Optional: true, Computed: true},
		"enc_algo_rc4":        dsschema.BoolAttribute{Description: "Allow RC4 encryption algorithm.", Optional: true, Computed: true},
		"keyxchg_algo_dhe":    dsschema.BoolAttribute{Description: "Allow DHE key exchange algorithm.", Optional: true, Computed: true},
		"keyxchg_algo_ecdhe":  dsschema.BoolAttribute{Description: "Allow ECDHE key exchange algorithm.", Optional: true, Computed: true},
		"keyxchg_algo_rsa":    dsschema.BoolAttribute{Description: "Allow RSA key exchange algorithm.", Optional: true, Computed: true},
		"max_version": dsschema.StringAttribute{Description: "Maximum TLS protocol version: tls1-0, tls1-1, tls1-2, tls1-3, or max.", Optional: true, Computed: true},
		"min_version": dsschema.StringAttribute{Description: "Minimum TLS protocol version: tls1-0, tls1-1, tls1-2, or tls1-3.", Optional: true, Computed: true},
	}
}

func decryptionProfileSslProtocolSettingsRsSchema() map[string]rsschema.Attribute {
	return map[string]rsschema.Attribute{
		"auth_algo_md5":       rsschema.BoolAttribute{Description: "Allow MD5 authentication algorithm.", Optional: true},
		"auth_algo_sha1":      rsschema.BoolAttribute{Description: "Allow SHA1 authentication algorithm.", Optional: true},
		"auth_algo_sha256":    rsschema.BoolAttribute{Description: "Allow SHA256 authentication algorithm.", Optional: true},
		"auth_algo_sha384":    rsschema.BoolAttribute{Description: "Allow SHA384 authentication algorithm.", Optional: true},
		"enc_algo_3des":       rsschema.BoolAttribute{Description: "Allow 3DES encryption algorithm.", Optional: true},
		"enc_algo_aes_128_cbc": rsschema.BoolAttribute{Description: "Allow AES-128-CBC encryption algorithm.", Optional: true},
		"enc_algo_aes_128_gcm": rsschema.BoolAttribute{Description: "Allow AES-128-GCM encryption algorithm.", Optional: true},
		"enc_algo_aes_256_cbc": rsschema.BoolAttribute{Description: "Allow AES-256-CBC encryption algorithm.", Optional: true},
		"enc_algo_aes_256_gcm": rsschema.BoolAttribute{Description: "Allow AES-256-GCM encryption algorithm.", Optional: true},
		"enc_algo_rc4":        rsschema.BoolAttribute{Description: "Allow RC4 encryption algorithm.", Optional: true},
		"keyxchg_algo_dhe":    rsschema.BoolAttribute{Description: "Allow DHE key exchange algorithm.", Optional: true},
		"keyxchg_algo_ecdhe":  rsschema.BoolAttribute{Description: "Allow ECDHE key exchange algorithm.", Optional: true},
		"keyxchg_algo_rsa":    rsschema.BoolAttribute{Description: "Allow RSA key exchange algorithm.", Optional: true},
		"max_version": rsschema.StringAttribute{Description: "Maximum TLS protocol version: tls1-0, tls1-1, tls1-2, tls1-3, or max.", Optional: true},
		"min_version": rsschema.StringAttribute{Description: "Minimum TLS protocol version: tls1-0, tls1-1, tls1-2, or tls1-3.", Optional: true},
	}
}

// -----------------------------------------------------------------------
// Location types
// -----------------------------------------------------------------------

type DecryptionProfileSharedLocation struct{}

type DecryptionProfileDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}

type DecryptionProfileLocation struct {
	Shared      types.Object `tfsdk:"shared"`
	DeviceGroup types.Object `tfsdk:"device_group"`
}

func (o DecryptionProfileLocation) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"shared": types.ObjectType{AttrTypes: map[string]attr.Type{}},
		"device_group": types.ObjectType{AttrTypes: map[string]attr.Type{
			"panorama_device": types.StringType,
			"name":            types.StringType,
		}},
	}
}

func DecryptionProfileLocationSchema() rsschema.Attribute {
	return rsschema.SingleNestedAttribute{
		Description: "The location of this object.",
		Required:    true,
		Attributes: map[string]rsschema.Attribute{
			"shared": rsschema.SingleNestedAttribute{
				Description: "Panorama shared object.",
				Optional:    true,
				Attributes:  map[string]rsschema.Attribute{},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRelative().AtParent().AtName("shared"),
						path.MatchRelative().AtParent().AtName("device_group"),
					}...),
				},
			},
			"device_group": rsschema.SingleNestedAttribute{
				Description: "Located in a specific Device Group.",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"panorama_device": rsschema.StringAttribute{
						Description: "Panorama device name.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("localhost.localdomain"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": rsschema.StringAttribute{
						Description: "Device Group name.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func DecryptionProfileResourceLocationSchema() rsschema.Attribute {
	return DecryptionProfileLocationSchema()
}

func DecryptionProfileDataSourceLocationSchema() dsschema.Attribute {
	return dsschema.SingleNestedAttribute{
		Description: "The location of this object.",
		Required:    true,
		Attributes: map[string]dsschema.Attribute{
			"shared": dsschema.SingleNestedAttribute{
				Description: "Panorama shared object.",
				Optional:    true,
				Attributes:  map[string]dsschema.Attribute{},
			},
			"device_group": dsschema.SingleNestedAttribute{
				Description: "Located in a specific Device Group.",
				Optional:    true,
				Attributes: map[string]dsschema.Attribute{
					"panorama_device": dsschema.StringAttribute{
						Description: "Panorama device name.",
						Optional:    true,
						Computed:    true,
					},
					"name": dsschema.StringAttribute{
						Description: "Device Group name.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (o DecryptionProfileSharedLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{}{})
}

func (o *DecryptionProfileSharedLocation) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &struct{}{})
}

func (o DecryptionProfileDeviceGroupLocation) MarshalJSON() ([]byte, error) {
	type shadow struct {
		PanoramaDevice *string `json:"panorama_device,omitempty"`
		Name           *string `json:"name,omitempty"`
	}
	return json.Marshal(shadow{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
		Name:           o.Name.ValueStringPointer(),
	})
}

func (o *DecryptionProfileDeviceGroupLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		PanoramaDevice *string `json:"panorama_device,omitempty"`
		Name           *string `json:"name,omitempty"`
	}
	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}
	o.PanoramaDevice = types.StringPointerValue(shadow.PanoramaDevice)
	o.Name = types.StringPointerValue(shadow.Name)
	return nil
}

func (o DecryptionProfileLocation) MarshalJSON() ([]byte, error) {
	type shadow struct {
		Shared      *DecryptionProfileSharedLocation      `json:"shared,omitempty"`
		DeviceGroup *DecryptionProfileDeviceGroupLocation `json:"device_group,omitempty"`
	}
	var shared_object *DecryptionProfileSharedLocation
	{
		diags := o.Shared.As(context.TODO(), &shared_object, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, NewDiagnosticsError("Failed to marshal shared into JSON document", diags.Errors())
		}
	}
	var device_group_object *DecryptionProfileDeviceGroupLocation
	{
		diags := o.DeviceGroup.As(context.TODO(), &device_group_object, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, NewDiagnosticsError("Failed to marshal device_group into JSON document", diags.Errors())
		}
	}
	return json.Marshal(shadow{
		Shared:      shared_object,
		DeviceGroup: device_group_object,
	})
}

func (o *DecryptionProfileLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Shared      *DecryptionProfileSharedLocation      `json:"shared,omitempty"`
		DeviceGroup *DecryptionProfileDeviceGroupLocation `json:"device_group,omitempty"`
	}
	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}
	{
		if shadow.Shared != nil {
			var diags diag.Diagnostics
			o.Shared, diags = types.ObjectValueFrom(context.TODO(), map[string]attr.Type{}, shadow.Shared)
			if diags.HasError() {
				return NewDiagnosticsError("Failed to unmarshal shared from JSON document", diags.Errors())
			}
		} else {
			o.Shared = types.ObjectNull(map[string]attr.Type{})
		}
	}
	{
		deviceGroupAttrTypes := map[string]attr.Type{
			"panorama_device": types.StringType,
			"name":            types.StringType,
		}
		if shadow.DeviceGroup != nil {
			var diags diag.Diagnostics
			o.DeviceGroup, diags = types.ObjectValueFrom(context.TODO(), deviceGroupAttrTypes, shadow.DeviceGroup)
			if diags.HasError() {
				return NewDiagnosticsError("Failed to unmarshal device_group from JSON document", diags.Errors())
			}
		} else {
			o.DeviceGroup = types.ObjectNull(deviceGroupAttrTypes)
		}
	}
	return nil
}

func decryptionProfileLocationFromTF(ctx context.Context, locationObj types.Object, location *decryption.Location) diag.Diagnostics {
	var diags diag.Diagnostics

	var terraformLocation DecryptionProfileLocation
	diags.Append(locationObj.As(ctx, &terraformLocation, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	if !terraformLocation.Shared.IsNull() {
		location.Shared = &decryption.SharedLocation{}
	}

	if !terraformLocation.DeviceGroup.IsNull() {
		location.DeviceGroup = &decryption.DeviceGroupLocation{}
		var inner DecryptionProfileDeviceGroupLocation
		diags.Append(terraformLocation.DeviceGroup.As(ctx, &inner, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return diags
		}
		location.DeviceGroup.PanoramaDevice = inner.PanoramaDevice.ValueString()
		location.DeviceGroup.DeviceGroup = inner.Name.ValueString()
	}

	return diags
}

// -----------------------------------------------------------------------
// Nested object TF types
// -----------------------------------------------------------------------

type DecryptionProfileSslForwardProxyObject struct {
	AutoIncludeAltname            types.Bool `tfsdk:"auto_include_altname"`
	BlockClientCert               types.Bool `tfsdk:"block_client_cert"`
	BlockExpiredCertificate       types.Bool `tfsdk:"block_expired_certificate"`
	BlockTimeoutCert              types.Bool `tfsdk:"block_timeout_cert"`
	BlockTls13DowngradeNoResource types.Bool `tfsdk:"block_tls13_downgrade_no_resource"`
	BlockUnknownCert              types.Bool `tfsdk:"block_unknown_cert"`
	BlockUnsupportedCipher        types.Bool `tfsdk:"block_unsupported_cipher"`
	BlockUnsupportedVersion       types.Bool `tfsdk:"block_unsupported_version"`
	BlockUntrustedIssuer          types.Bool `tfsdk:"block_untrusted_issuer"`
	RestrictCertExts              types.Bool `tfsdk:"restrict_cert_exts"`
	StripAlpn                     types.Bool `tfsdk:"strip_alpn"`
}

func (o *DecryptionProfileSslForwardProxyObject) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"auto_include_altname":              types.BoolType,
		"block_client_cert":                 types.BoolType,
		"block_expired_certificate":         types.BoolType,
		"block_timeout_cert":                types.BoolType,
		"block_tls13_downgrade_no_resource": types.BoolType,
		"block_unknown_cert":                types.BoolType,
		"block_unsupported_cipher":          types.BoolType,
		"block_unsupported_version":         types.BoolType,
		"block_untrusted_issuer":            types.BoolType,
		"restrict_cert_exts":                types.BoolType,
		"strip_alpn":                        types.BoolType,
	}
}

type DecryptionProfileSslInboundInspectionObject struct {
	BlockIfHsmUnavailable         types.Bool `tfsdk:"block_if_hsm_unavailable"`
	BlockIfNoResource             types.Bool `tfsdk:"block_if_no_resource"`
	BlockTls13DowngradeNoResource types.Bool `tfsdk:"block_tls13_downgrade_no_resource"`
	BlockUnsupportedCipher        types.Bool `tfsdk:"block_unsupported_cipher"`
	BlockUnsupportedVersion       types.Bool `tfsdk:"block_unsupported_version"`
}

func (o *DecryptionProfileSslInboundInspectionObject) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"block_if_hsm_unavailable":          types.BoolType,
		"block_if_no_resource":              types.BoolType,
		"block_tls13_downgrade_no_resource": types.BoolType,
		"block_unsupported_cipher":          types.BoolType,
		"block_unsupported_version":         types.BoolType,
	}
}

type DecryptionProfileSslNoProxyObject struct {
	BlockClientCert         types.Bool `tfsdk:"block_client_cert"`
	BlockExpiredCertificate types.Bool `tfsdk:"block_expired_certificate"`
	BlockTimeoutCert        types.Bool `tfsdk:"block_timeout_cert"`
	BlockUnknownCert        types.Bool `tfsdk:"block_unknown_cert"`
	BlockUnsupportedCipher  types.Bool `tfsdk:"block_unsupported_cipher"`
	BlockUnsupportedVersion types.Bool `tfsdk:"block_unsupported_version"`
	BlockUntrustedIssuer    types.Bool `tfsdk:"block_untrusted_issuer"`
}

func (o *DecryptionProfileSslNoProxyObject) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"block_client_cert":         types.BoolType,
		"block_expired_certificate": types.BoolType,
		"block_timeout_cert":        types.BoolType,
		"block_unknown_cert":        types.BoolType,
		"block_unsupported_cipher":  types.BoolType,
		"block_unsupported_version": types.BoolType,
		"block_untrusted_issuer":    types.BoolType,
	}
}

type DecryptionProfileSslProtocolSettingsObject struct {
	AuthAlgoMd5        types.Bool   `tfsdk:"auth_algo_md5"`
	AuthAlgoSha1       types.Bool   `tfsdk:"auth_algo_sha1"`
	AuthAlgoSha256     types.Bool   `tfsdk:"auth_algo_sha256"`
	AuthAlgoSha384     types.Bool   `tfsdk:"auth_algo_sha384"`
	EncAlgo3des        types.Bool   `tfsdk:"enc_algo_3des"`
	EncAlgoAes128Cbc   types.Bool   `tfsdk:"enc_algo_aes_128_cbc"`
	EncAlgoAes128Gcm   types.Bool   `tfsdk:"enc_algo_aes_128_gcm"`
	EncAlgoAes256Cbc   types.Bool   `tfsdk:"enc_algo_aes_256_cbc"`
	EncAlgoAes256Gcm   types.Bool   `tfsdk:"enc_algo_aes_256_gcm"`
	EncAlgoRc4         types.Bool   `tfsdk:"enc_algo_rc4"`
	KeyxchgAlgoDhe     types.Bool   `tfsdk:"keyxchg_algo_dhe"`
	KeyxchgAlgoEcdhe   types.Bool   `tfsdk:"keyxchg_algo_ecdhe"`
	KeyxchgAlgoRsa     types.Bool   `tfsdk:"keyxchg_algo_rsa"`
	MaxVersion         types.String `tfsdk:"max_version"`
	MinVersion         types.String `tfsdk:"min_version"`
}

func (o *DecryptionProfileSslProtocolSettingsObject) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"auth_algo_md5":        types.BoolType,
		"auth_algo_sha1":       types.BoolType,
		"auth_algo_sha256":     types.BoolType,
		"auth_algo_sha384":     types.BoolType,
		"enc_algo_3des":        types.BoolType,
		"enc_algo_aes_128_cbc": types.BoolType,
		"enc_algo_aes_128_gcm": types.BoolType,
		"enc_algo_aes_256_cbc": types.BoolType,
		"enc_algo_aes_256_gcm": types.BoolType,
		"enc_algo_rc4":         types.BoolType,
		"keyxchg_algo_dhe":     types.BoolType,
		"keyxchg_algo_ecdhe":   types.BoolType,
		"keyxchg_algo_rsa":     types.BoolType,
		"max_version":          types.StringType,
		"min_version":          types.StringType,
	}
}

// -----------------------------------------------------------------------
// Conversion helpers
// -----------------------------------------------------------------------

func copySslForwardProxyToPango(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *decryption.SslForwardProxy {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	var tfObj DecryptionProfileSslForwardProxyObject
	diags.Append(obj.As(ctx, &tfObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &decryption.SslForwardProxy{
		AutoIncludeAltname:            tfObj.AutoIncludeAltname.ValueBoolPointer(),
		BlockClientCert:               tfObj.BlockClientCert.ValueBoolPointer(),
		BlockExpiredCertificate:       tfObj.BlockExpiredCertificate.ValueBoolPointer(),
		BlockTimeoutCert:              tfObj.BlockTimeoutCert.ValueBoolPointer(),
		BlockTls13DowngradeNoResource: tfObj.BlockTls13DowngradeNoResource.ValueBoolPointer(),
		BlockUnknownCert:              tfObj.BlockUnknownCert.ValueBoolPointer(),
		BlockUnsupportedCipher:        tfObj.BlockUnsupportedCipher.ValueBoolPointer(),
		BlockUnsupportedVersion:       tfObj.BlockUnsupportedVersion.ValueBoolPointer(),
		BlockUntrustedIssuer:          tfObj.BlockUntrustedIssuer.ValueBoolPointer(),
		RestrictCertExts:              tfObj.RestrictCertExts.ValueBoolPointer(),
		StripAlpn:                     tfObj.StripAlpn.ValueBoolPointer(),
	}
}

func copySslForwardProxyFromPango(ctx context.Context, s *decryption.SslForwardProxy) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var obj *DecryptionProfileSslForwardProxyObject
	if s == nil {
		return types.ObjectNull(obj.AttributeTypes()), diags
	}
	tfObj := DecryptionProfileSslForwardProxyObject{
		AutoIncludeAltname:            boolPtrToType(s.AutoIncludeAltname),
		BlockClientCert:               boolPtrToType(s.BlockClientCert),
		BlockExpiredCertificate:       boolPtrToType(s.BlockExpiredCertificate),
		BlockTimeoutCert:              boolPtrToType(s.BlockTimeoutCert),
		BlockTls13DowngradeNoResource: boolPtrToType(s.BlockTls13DowngradeNoResource),
		BlockUnknownCert:              boolPtrToType(s.BlockUnknownCert),
		BlockUnsupportedCipher:        boolPtrToType(s.BlockUnsupportedCipher),
		BlockUnsupportedVersion:       boolPtrToType(s.BlockUnsupportedVersion),
		BlockUntrustedIssuer:          boolPtrToType(s.BlockUntrustedIssuer),
		RestrictCertExts:              boolPtrToType(s.RestrictCertExts),
		StripAlpn:                     boolPtrToType(s.StripAlpn),
	}
	result, d := types.ObjectValueFrom(ctx, tfObj.AttributeTypes(), tfObj)
	diags.Append(d...)
	return result, diags
}

func copySslInboundInspectionToPango(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *decryption.SslInboundInspection {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	var tfObj DecryptionProfileSslInboundInspectionObject
	diags.Append(obj.As(ctx, &tfObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &decryption.SslInboundInspection{
		BlockIfHsmUnavailable:         tfObj.BlockIfHsmUnavailable.ValueBoolPointer(),
		BlockIfNoResource:             tfObj.BlockIfNoResource.ValueBoolPointer(),
		BlockTls13DowngradeNoResource: tfObj.BlockTls13DowngradeNoResource.ValueBoolPointer(),
		BlockUnsupportedCipher:        tfObj.BlockUnsupportedCipher.ValueBoolPointer(),
		BlockUnsupportedVersion:       tfObj.BlockUnsupportedVersion.ValueBoolPointer(),
	}
}

func copySslInboundInspectionFromPango(ctx context.Context, s *decryption.SslInboundInspection) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var obj *DecryptionProfileSslInboundInspectionObject
	if s == nil {
		return types.ObjectNull(obj.AttributeTypes()), diags
	}
	tfObj := DecryptionProfileSslInboundInspectionObject{
		BlockIfHsmUnavailable:         boolPtrToType(s.BlockIfHsmUnavailable),
		BlockIfNoResource:             boolPtrToType(s.BlockIfNoResource),
		BlockTls13DowngradeNoResource: boolPtrToType(s.BlockTls13DowngradeNoResource),
		BlockUnsupportedCipher:        boolPtrToType(s.BlockUnsupportedCipher),
		BlockUnsupportedVersion:       boolPtrToType(s.BlockUnsupportedVersion),
	}
	result, d := types.ObjectValueFrom(ctx, tfObj.AttributeTypes(), tfObj)
	diags.Append(d...)
	return result, diags
}

func copySslNoProxyToPango(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *decryption.SslNoProxy {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	var tfObj DecryptionProfileSslNoProxyObject
	diags.Append(obj.As(ctx, &tfObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &decryption.SslNoProxy{
		BlockClientCert:         tfObj.BlockClientCert.ValueBoolPointer(),
		BlockExpiredCertificate: tfObj.BlockExpiredCertificate.ValueBoolPointer(),
		BlockTimeoutCert:        tfObj.BlockTimeoutCert.ValueBoolPointer(),
		BlockUnknownCert:        tfObj.BlockUnknownCert.ValueBoolPointer(),
		BlockUnsupportedCipher:  tfObj.BlockUnsupportedCipher.ValueBoolPointer(),
		BlockUnsupportedVersion: tfObj.BlockUnsupportedVersion.ValueBoolPointer(),
		BlockUntrustedIssuer:    tfObj.BlockUntrustedIssuer.ValueBoolPointer(),
	}
}

func copySslNoProxyFromPango(ctx context.Context, s *decryption.SslNoProxy) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var obj *DecryptionProfileSslNoProxyObject
	if s == nil {
		return types.ObjectNull(obj.AttributeTypes()), diags
	}
	tfObj := DecryptionProfileSslNoProxyObject{
		BlockClientCert:         boolPtrToType(s.BlockClientCert),
		BlockExpiredCertificate: boolPtrToType(s.BlockExpiredCertificate),
		BlockTimeoutCert:        boolPtrToType(s.BlockTimeoutCert),
		BlockUnknownCert:        boolPtrToType(s.BlockUnknownCert),
		BlockUnsupportedCipher:  boolPtrToType(s.BlockUnsupportedCipher),
		BlockUnsupportedVersion: boolPtrToType(s.BlockUnsupportedVersion),
		BlockUntrustedIssuer:    boolPtrToType(s.BlockUntrustedIssuer),
	}
	result, d := types.ObjectValueFrom(ctx, tfObj.AttributeTypes(), tfObj)
	diags.Append(d...)
	return result, diags
}

func copySslProtocolSettingsToPango(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *decryption.SslProtocolSettings {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	var tfObj DecryptionProfileSslProtocolSettingsObject
	diags.Append(obj.As(ctx, &tfObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &decryption.SslProtocolSettings{
		AuthAlgoMd5:        tfObj.AuthAlgoMd5.ValueBoolPointer(),
		AuthAlgoSha1:       tfObj.AuthAlgoSha1.ValueBoolPointer(),
		AuthAlgoSha256:     tfObj.AuthAlgoSha256.ValueBoolPointer(),
		AuthAlgoSha384:     tfObj.AuthAlgoSha384.ValueBoolPointer(),
		EncAlgo3des:        tfObj.EncAlgo3des.ValueBoolPointer(),
		EncAlgoAes128Cbc:   tfObj.EncAlgoAes128Cbc.ValueBoolPointer(),
		EncAlgoAes128Gcm:   tfObj.EncAlgoAes128Gcm.ValueBoolPointer(),
		EncAlgoAes256Cbc:   tfObj.EncAlgoAes256Cbc.ValueBoolPointer(),
		EncAlgoAes256Gcm:   tfObj.EncAlgoAes256Gcm.ValueBoolPointer(),
		EncAlgoRc4:         tfObj.EncAlgoRc4.ValueBoolPointer(),
		KeyxchgAlgoDhe:     tfObj.KeyxchgAlgoDhe.ValueBoolPointer(),
		KeyxchgAlgoEcdhe:   tfObj.KeyxchgAlgoEcdhe.ValueBoolPointer(),
		KeyxchgAlgoRsa:     tfObj.KeyxchgAlgoRsa.ValueBoolPointer(),
		MaxVersion:         tfObj.MaxVersion.ValueStringPointer(),
		MinVersion:         tfObj.MinVersion.ValueStringPointer(),
	}
}

func copySslProtocolSettingsFromPango(ctx context.Context, s *decryption.SslProtocolSettings) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var obj *DecryptionProfileSslProtocolSettingsObject
	if s == nil {
		return types.ObjectNull(obj.AttributeTypes()), diags
	}
	tfObj := DecryptionProfileSslProtocolSettingsObject{
		AuthAlgoMd5:        boolPtrToType(s.AuthAlgoMd5),
		AuthAlgoSha1:       boolPtrToType(s.AuthAlgoSha1),
		AuthAlgoSha256:     boolPtrToType(s.AuthAlgoSha256),
		AuthAlgoSha384:     boolPtrToType(s.AuthAlgoSha384),
		EncAlgo3des:        boolPtrToType(s.EncAlgo3des),
		EncAlgoAes128Cbc:   boolPtrToType(s.EncAlgoAes128Cbc),
		EncAlgoAes128Gcm:   boolPtrToType(s.EncAlgoAes128Gcm),
		EncAlgoAes256Cbc:   boolPtrToType(s.EncAlgoAes256Cbc),
		EncAlgoAes256Gcm:   boolPtrToType(s.EncAlgoAes256Gcm),
		EncAlgoRc4:         boolPtrToType(s.EncAlgoRc4),
		KeyxchgAlgoDhe:     boolPtrToType(s.KeyxchgAlgoDhe),
		KeyxchgAlgoEcdhe:   boolPtrToType(s.KeyxchgAlgoEcdhe),
		KeyxchgAlgoRsa:     boolPtrToType(s.KeyxchgAlgoRsa),
		MaxVersion:         types.StringPointerValue(s.MaxVersion),
		MinVersion:         types.StringPointerValue(s.MinVersion),
	}
	result, d := types.ObjectValueFrom(ctx, tfObj.AttributeTypes(), tfObj)
	diags.Append(d...)
	return result, diags
}
