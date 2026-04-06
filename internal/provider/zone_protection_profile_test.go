package provider_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/PaloAltoNetworks/terraform-provider-panos/internal/provider"
	zp "github.com/PaloAltoNetworks/pango/network/profiles/zone_protection"
)

// helpers
func boolPtr(b bool) *bool   { return &b }
func strPtr(s string) *string { return &s }
func int64Ptr(i int64) *int64 { return &i }

func nullFloodObject() types.Object {
	var f *provider.ZoneProtectionProfileFloodObject
	obj := types.ObjectNull(f.AttributeTypes())
	return obj
}

var _ = Describe("ZoneProtectionProfile provider model", func() {
	ctx := context.Background()

	Describe("CopyToPango", func() {
		Context("top-level bool fields", func() {
			It("copies discard_strict_source_routing and discard_loose_source_routing", func() {
				model := &provider.ZoneProtectionProfileResourceModel{
					Name:                       types.StringValue("test"),
					Flood:                      nullFloodObject(),
					DiscardIpSpoof:             types.BoolValue(true),
					DiscardStrictSourceRouting: types.BoolValue(true),
					DiscardLooseSourceRouting:  types.BoolValue(false),
					DiscardMalformedOption:     types.BoolNull(),
					RemoveTcpTimestamp:         types.BoolNull(),
					DiscardIpFrag:              types.BoolNull(),
					TcpSynWithData:             types.BoolNull(),
					StripTcpFastOpenAndData:    types.BoolNull(),
					StripMptcpOption:           types.StringNull(),
					Description:                types.StringNull(),
				}

				var obj *zp.Entry
				diags := model.CopyToPango(ctx, nil, nil, &obj, nil)
				Expect(diags.HasError()).To(BeFalse())
				Expect(obj).ToNot(BeNil())
				Expect(obj.DiscardStrictSourceRouting).ToNot(BeNil())
				Expect(*obj.DiscardStrictSourceRouting).To(BeTrue())
				Expect(obj.DiscardLooseSourceRouting).ToNot(BeNil())
				Expect(*obj.DiscardLooseSourceRouting).To(BeFalse())
			})
		})
	})

	Describe("CopyFromPango", func() {
		Context("top-level bool fields", func() {
			It("copies discard_strict_source_routing and discard_loose_source_routing from Entry", func() {
				entry := &zp.Entry{
					Name:                       "test",
					DiscardStrictSourceRouting: boolPtr(true),
					DiscardLooseSourceRouting:  boolPtr(false),
				}

				model := &provider.ZoneProtectionProfileResourceModel{
					Flood: nullFloodObject(),
				}
				diags := model.CopyFromPango(ctx, nil, nil, entry, nil)
				Expect(diags.HasError()).To(BeFalse())
				Expect(model.DiscardStrictSourceRouting.ValueBool()).To(BeTrue())
				Expect(model.DiscardLooseSourceRouting.ValueBool()).To(BeFalse())
			})
		})

		Context("flood fields", func() {
			It("copies all five flood protocols including other", func() {
				entry := &zp.Entry{
					Name: "test",
					Flood: &zp.Flood{
						Icmp: &zp.FloodProtocol{
							Enable: boolPtr(true),
							Red: &zp.FloodRates{
								AlarmRate:    int64Ptr(10000),
								ActivateRate: int64Ptr(10000),
								MaximalRate:  int64Ptr(40000),
							},
						},
						Other: &zp.FloodProtocol{
							Enable: boolPtr(true),
						},
					},
				}

				model := &provider.ZoneProtectionProfileResourceModel{
					Flood: nullFloodObject(),
				}
				diags := model.CopyFromPango(ctx, nil, nil, entry, nil)
				Expect(diags.HasError()).To(BeFalse())
				Expect(model.Flood.IsNull()).To(BeFalse())

				var floodObj provider.ZoneProtectionProfileFloodObject
				diags = model.Flood.As(ctx, &floodObj, basetypes.ObjectAsOptions{})
				Expect(diags.HasError()).To(BeFalse())
				Expect(floodObj.Other.IsNull()).To(BeFalse())
			})
		})
	})

	Describe("AttributeTypes", func() {
		It("resource model AttributeTypes includes all expected keys", func() {
			var m provider.ZoneProtectionProfileResourceModel
			attrTypes := m.AttributeTypes()

			expectedKeys := []string{
				"location", "name", "description", "flood",
				"discard_ip_spoof",
				"discard_strict_source_routing",
				"discard_loose_source_routing",
				"discard_malformed_option",
				"remove_tcp_timestamp",
				"discard_ip_frag",
				"tcp_syn_with_data",
				"strip_tcp_fast_open_and_data",
				"strip_mptcp_option",
			}

			for _, key := range expectedKeys {
				_, ok := attrTypes[key]
				Expect(ok).To(BeTrue(), "missing key: %s", key)
			}
		})

		It("datasource model AttributeTypes includes all expected keys", func() {
			var m provider.ZoneProtectionProfileDataSourceModel
			attrTypes := m.AttributeTypes()

			Expect(attrTypes).To(HaveKey("discard_strict_source_routing"))
			Expect(attrTypes).To(HaveKey("discard_loose_source_routing"))
			Expect(attrTypes["discard_strict_source_routing"]).To(Equal(attr.Type(types.BoolType)))
		})
	})
})
