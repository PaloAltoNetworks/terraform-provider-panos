package provider

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ = Describe("Position", func() {
	var resp resource.ValidateConfigResponse
	JustBeforeEach(func() {
		resp.Diagnostics = diag.Diagnostics{}
	})

	Context("when validating the 'where' attribute", func() {
		DescribeTable("allowed values",
			func(whereValue string, withPivot bool, expectError bool) {
				pivot := types.StringNull()
				directly := types.BoolNull()
				if withPivot {
					pivot = types.StringValue("mocked")
					directly = types.BoolValue(false)
				}

				position := TerraformPositionObject{
					Where:    types.StringValue(whereValue),
					Pivot:    pivot,
					Directly: directly,
				}
				position.ValidateConfig(&resp)

				if expectError {
					Expect(resp.Diagnostics.Errors()).To(HaveExactElements([]diag.Diagnostic{
						diag.NewAttributeErrorDiagnostic(
							path.Root("position").AtName("where"),
							"Missing attribute configuration",
							fmt.Sprintf("where attribute must be one of the valid values: first, last, before, after, found: '%s'", whereValue)),
					}))
					Expect(resp.Diagnostics.Warnings()).To(BeEmpty())
				} else {
					Expect(resp.Diagnostics.Errors()).To(BeEmpty())
					Expect(resp.Diagnostics.Warnings()).To(BeEmpty())
				}
			},
			Entry("should accept 'first'", "first", false, false),
			Entry("should accept 'last'", "last", false, false),
			Entry("should accept 'before'", "before", true, false),
			Entry("should accept 'after'", "after", true, false),
			Entry("should reject 'invalid'", "invalid", false, true),
			Entry("should reject empty string", "", false, true),
		)
	})
	Context("and where attribute is 'before' or 'after'", func() {
		Context("and pivot is null", func() {
			It("should return a diagnostic error about pivot not being set", func() {
				position := TerraformPositionObject{
					Where: types.StringValue("before"),
					Pivot: types.StringNull(),
				}
				position.ValidateConfig(&resp)
				Expect(resp.Diagnostics.Errors()).To(HaveExactElements([]diag.Diagnostic{
					diag.NewAttributeErrorDiagnostic(
						path.Root("position").AtName("pivot"),
						"Missing attribute configuration",
						"position pivot attribute must be set to a valid object when where attribute is set to either 'after' or 'before'"),
				}))
				Expect(resp.Diagnostics.Warnings()).To(HaveLen(0))

			})
		})
		Context("and pivot is an empty string", func() {
			It("should return a diagnostic error about pivot not being set", func() {
				position := TerraformPositionObject{
					Where: types.StringValue("after"),
					Pivot: types.StringValue(""),
				}
				position.ValidateConfig(&resp)
				Expect(resp.Diagnostics.Errors()).To(HaveExactElements([]diag.Diagnostic{
					diag.NewAttributeErrorDiagnostic(
						path.Root("position").AtName("pivot"),
						"Missing attribute configuration",
						"position pivot attribute must be set to a valid object when where attribute is set to either 'after' or 'before'"),
				}))
				Expect(resp.Diagnostics.Warnings()).To(HaveLen(0))
			})
		})
		Context("and pivot is a valid string", func() {
			Context("and directly is null", func() {
				It("should return an error about directly not being configured properly", func() {
					position := TerraformPositionObject{
						Where:    types.StringValue("after"),
						Pivot:    types.StringValue("mocked"),
						Directly: types.BoolNull(),
					}
					position.ValidateConfig(&resp)
					Expect(resp.Diagnostics.Contains(diag.NewAttributeErrorDiagnostic(
						path.Root("position").AtName("directly"),
						"Missing attribute configuration",
						"Expected directly to be configured with pivot"),
					))
					Expect(resp.Diagnostics.Warnings()).To(HaveLen(0))
				})
			})
		})
		Context("and directly is set", func() {
			Context("and pivot is null", func() {
				It("should return an error about directly not being configured properly", func() {
					position := TerraformPositionObject{
						Where:    types.StringValue("after"),
						Pivot:    types.StringNull(),
						Directly: types.BoolValue(true),
					}
					position.ValidateConfig(&resp)
					Expect(resp.Diagnostics.Contains(diag.NewAttributeErrorDiagnostic(
						path.Root("position").AtName("pivot"),
						"Missing attribute configuration",
						"Expected pivot to be configured with directly"),
					))
					Expect(resp.Diagnostics.Warnings()).To(HaveLen(0))
				})
			})
		})
	})
	Context("when 'where' is 'first' or 'last' and 'pivot' is set", func() {
		DescribeTable("should return a warning about pivot being ignored",
			func(whereValue string) {
				position := TerraformPositionObject{
					Where:    types.StringValue(whereValue),
					Pivot:    types.StringValue("some-pivot"),
					Directly: types.BoolValue(false),
				}
				position.ValidateConfig(&resp)

				Expect(resp.Diagnostics.Errors()).To(BeEmpty())
				Expect(resp.Diagnostics.Warnings()).To(HaveExactElements([]diag.Diagnostic{
					diag.NewAttributeWarningDiagnostic(
						path.Root("position").AtName("pivot"),
						"Unexpected attribute configuration",
						"pivot attribute is ignored when where is set to 'first' or 'last'"),
					diag.NewAttributeWarningDiagnostic(
						path.Root("position").AtName("directly"),
						"Unexpected attribute configuration",
						"directly attribute is ignored when where is set to 'first' or 'last'"),
				}))
			},
			Entry("when 'where' is 'first'", "first"),
			Entry("when 'where' is 'last'", "last"),
		)
	})
	Context("when 'where' is 'before' or 'after' and both 'pivot' and 'directly' are set", func() {
		DescribeTable("should not return any errors or warnings",
			func(whereValue string) {
				position := TerraformPositionObject{
					Where:    types.StringValue(whereValue),
					Pivot:    types.StringValue("some-pivot"),
					Directly: types.BoolValue(true),
				}
				position.ValidateConfig(&resp)

				Expect(resp.Diagnostics.Errors()).To(BeEmpty())
				Expect(resp.Diagnostics.Warnings()).To(BeEmpty())
			},
			Entry("when 'where' is 'before'", "before"),
			Entry("when 'where' is 'after'", "after"),
		)
	})

	Context("when 'where' is 'before' or 'after' and 'pivot' is set but 'directly' is not", func() {
		DescribeTable("should return an error about directly not being configured",
			func(whereValue string) {
				position := TerraformPositionObject{
					Where:    types.StringValue(whereValue),
					Pivot:    types.StringValue("some-pivot"),
					Directly: types.BoolNull(),
				}
				position.ValidateConfig(&resp)

				Expect(resp.Diagnostics.Errors()).To(HaveExactElements([]diag.Diagnostic{
					diag.NewAttributeErrorDiagnostic(
						path.Root("position").AtName("directly"),
						"Missing attribute configuration",
						"Expected directly to be configured with pivot"),
				}))
				Expect(resp.Diagnostics.Warnings()).To(BeEmpty())
			},
			Entry("when 'where' is 'before'", "before"),
			Entry("when 'where' is 'after'", "after"),
		)
	})

	Context("when 'where' is 'before' or 'after' and 'directly' is set but 'pivot' is not", func() {
		DescribeTable("should return an error about pivot not being configured",
			func(whereValue string) {
				position := TerraformPositionObject{
					Where:    types.StringValue(whereValue),
					Pivot:    types.StringNull(),
					Directly: types.BoolValue(true),
				}
				position.ValidateConfig(&resp)

				Expect(resp.Diagnostics.Errors()).To(HaveExactElements([]diag.Diagnostic{
					diag.NewAttributeErrorDiagnostic(
						path.Root("position").AtName("pivot"),
						"Missing attribute configuration",
						"position pivot attribute must be set to a valid object when where attribute is set to either 'after' or 'before'"),
				}))
				Expect(resp.Diagnostics.Warnings()).To(BeEmpty())
			},
			Entry("when 'where' is 'before'", "before"),
			Entry("when 'where' is 'after'", "after"),
		)
	})

	Context("when attributes are set to unknown", func() {
		Context("when 'where' is unknown", func() {
			It("should not return any errors or warnings", func() {
				position := TerraformPositionObject{
					Where:    types.StringUnknown(),
					Pivot:    types.StringNull(),
					Directly: types.BoolNull(),
				}
				position.ValidateConfig(&resp)

				Expect(resp.Diagnostics.Errors()).To(BeEmpty())
				Expect(resp.Diagnostics.Warnings()).To(BeEmpty())
			})
		})

		Context("when 'pivot' is unknown", func() {
			It("should not validate pivot-related conditions", func() {
				position := TerraformPositionObject{
					Where:    types.StringValue("before"),
					Pivot:    types.StringUnknown(),
					Directly: types.BoolValue(true),
				}
				position.ValidateConfig(&resp)

				Expect(resp.Diagnostics.Errors()).To(BeEmpty())
				Expect(resp.Diagnostics.Warnings()).To(BeEmpty())
			})
		})

		Context("when 'directly' is unknown", func() {
			It("should not validate directly-related conditions", func() {
				position := TerraformPositionObject{
					Where:    types.StringValue("after"),
					Pivot:    types.StringValue("some-pivot"),
					Directly: types.BoolUnknown(),
				}
				position.ValidateConfig(&resp)

				Expect(resp.Diagnostics.Errors()).To(BeEmpty())
				Expect(resp.Diagnostics.Warnings()).To(BeEmpty())
			})
		})

		Context("when both 'pivot' and 'directly' are unknown", func() {
			It("should not validate pivot and directly related conditions", func() {
				position := TerraformPositionObject{
					Where:    types.StringValue("before"),
					Pivot:    types.StringUnknown(),
					Directly: types.BoolUnknown(),
				}
				position.ValidateConfig(&resp)

				Expect(resp.Diagnostics.Errors()).To(BeEmpty())
				Expect(resp.Diagnostics.Warnings()).To(BeEmpty())
			})
		})
	})
})
