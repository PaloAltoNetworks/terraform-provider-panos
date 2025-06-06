package provider_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PaloAltoNetworks/terraform-provider-panos/internal/provider"
)

type AncestorMock struct {
	name      string
	entryName string
}

func (o AncestorMock) AncestorName() string {
	return o.name
}

func (o AncestorMock) EntryName() *string {
	if o.entryName == "" {
		return nil
	}
	return &o.entryName
}

var _ = Describe("CreateXpathForParameterWithAncestors", func() {
	Context("When no ancestors are provided", func() {
		It("should generate a single element xpath", func() {
			var ancestors []provider.Ancestor

			xpath, err := provider.CreateXpathForAttributeWithAncestors(ancestors, "attr-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(xpath).To(Equal("/attr-1"))
		})
	})
	Context("When single ancestor of nested type is provided", func() {
		It("should generate a single element xpath", func() {
			ancestors := []provider.Ancestor{
				&AncestorMock{
					name: "attr-1",
				},
			}

			xpath, err := provider.CreateXpathForAttributeWithAncestors(ancestors, "attr-2")
			Expect(err).ToNot(HaveOccurred())
			Expect(xpath).To(Equal("/attr-1/attr-2"))
		})
	})
	Context("When multiple ancestors are provided", func() {
		Context("and one of ancestors is a list", func() {
			It("should generate a single element xpath", func() {
				ancestors := []provider.Ancestor{
					&AncestorMock{
						name: "attr-1",
					},
					&AncestorMock{
						name:      "attr-2",
						entryName: "element-1",
					},
				}

				xpath, err := provider.CreateXpathForAttributeWithAncestors(ancestors, "attr-3")
				Expect(err).ToNot(HaveOccurred())
				Expect(xpath).To(Equal(`/attr-1/attr-2/entry[@name="element-1"]/attr-3`))
			})
		})

	})
})
