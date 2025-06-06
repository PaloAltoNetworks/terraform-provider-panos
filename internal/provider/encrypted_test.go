package provider

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EncryptedValuesManager", func() {
	payload := []byte(`{"values":{"/attr-1":{"hashing_type":"solo","encrypted":"$enc$value","plaintext":"value"}}}`)
	Context("when creating encrypted values manager from existing payload", func() {
		It("should return correct values", func() {
			ev, err := NewEncryptedValuesManager(payload, false)
			Expect(err).ToNot(HaveOccurred())

			value, found := ev.GetPlaintextValue("/attr-1")
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("value"))

			value, found = ev.GetEncryptedValue("/attr-1")
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("$enc$value"))

			marshalled, err := json.Marshal(ev)
			Expect(err).ToNot(HaveOccurred())
			Expect(marshalled).To(Equal(payload))
		})
	})
	Context("when creating encrypted values manager with no existing payload", func() {
		var payload []byte
		Context("and inserting plaintext and encrypted values for a given xpath", func() {
			It("should return expected values back when requested", func() {
				ev, err := NewEncryptedValuesManager(payload, false)
				Expect(err).ToNot(HaveOccurred())

				err = ev.StorePlaintextValue("/attr-1", HashingSoloType, "value")
				Expect(err).ToNot(HaveOccurred())

				err = ev.StoreEncryptedValue("/attr-1", HashingSoloType, "$enc$value")
				Expect(err).ToNot(HaveOccurred())

				value, found := ev.GetPlaintextValue("/attr-1")
				Expect(found).To(BeTrue())
				Expect(value).To(Equal("value"))

				value, found = ev.GetEncryptedValue("/attr-1")
				Expect(found).To(BeTrue())
				Expect(value).To(Equal("$enc$value"))

				expected := []byte(`{"values":{"/attr-1":{"hashing_type":"solo","encrypted":"$enc$value","plaintext":"value"}}}`)
				marshalled, err := json.Marshal(ev)
				Expect(err).ToNot(HaveOccurred())
				Expect(marshalled).To(Equal(expected))
			})
		})
	})
})
