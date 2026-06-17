package provider

import (
	"encoding/xml"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ip_tag tagSetDiff", func() {
	Context("when the desired set adds and removes tags relative to current", func() {
		It("returns the additions and removals", func() {
			toAdd, toRemove := tagSetDiff([]string{"web", "db"}, []string{"web", "prod"})

			Expect(toAdd).To(ConsistOf("prod"))
			Expect(toRemove).To(ConsistOf("db"))
		})
	})

	Context("when the desired set only adds tags", func() {
		It("returns additions and no removals", func() {
			toAdd, toRemove := tagSetDiff([]string{"web"}, []string{"web", "prod"})

			Expect(toAdd).To(ConsistOf("prod"))
			Expect(toRemove).To(BeEmpty())
		})
	})

	Context("when the desired set only removes tags", func() {
		It("returns removals and no additions", func() {
			toAdd, toRemove := tagSetDiff([]string{"web", "db"}, []string{"web"})

			Expect(toAdd).To(BeEmpty())
			Expect(toRemove).To(ConsistOf("db"))
		})
	})

	Context("when current and desired are identical", func() {
		It("returns no changes", func() {
			toAdd, toRemove := tagSetDiff([]string{"web", "db"}, []string{"db", "web"})

			Expect(toAdd).To(BeEmpty())
			Expect(toRemove).To(BeEmpty())
		})
	})

	Context("when inputs contain duplicates", func() {
		It("treats them as sets and does not emit duplicates", func() {
			toAdd, toRemove := tagSetDiff([]string{"web", "web"}, []string{"prod", "prod"})

			Expect(toAdd).To(ConsistOf("prod"))
			Expect(toRemove).To(ConsistOf("web"))
		})
	})
})

var _ = Describe("ip_tag reconcileManagedTags", func() {
	Context("when some managed tags drifted away and unmanaged tags are present", func() {
		It("keeps only managed tags still present on the firewall", func() {
			// Managed {web,prod}; firewall has {web,db}: prod drifted away, db is
			// unmanaged. State should reconcile to {web}.
			result := reconcileManagedTags([]string{"web", "prod"}, []string{"web", "db"})

			Expect(result).To(ConsistOf("web"))
		})
	})

	Context("when none of the managed tags are present", func() {
		It("returns an empty set", func() {
			result := reconcileManagedTags([]string{"web", "prod"}, []string{"db"})

			Expect(result).To(BeEmpty())
		})
	})
})

var _ = Describe("ip_tag user-id command marshaling", func() {
	Context("when registering tags on an IP", func() {
		It("marshals a uid-message with a register payload", func() {
			out, err := xml.Marshal(registerIpTagCommand("1.1.1.1", []string{"web", "prod"}))

			Expect(err).ToNot(HaveOccurred())
			Expect(string(out)).To(Equal(
				`<uid-message><version>1.0</version><type>update</type>` +
					`<payload><register><entry ip="1.1.1.1">` +
					`<tag><member>web</member><member>prod</member></tag>` +
					`</entry></register></payload></uid-message>`))
		})
	})

	Context("when unregistering tags from an IP", func() {
		It("marshals a uid-message with an unregister payload", func() {
			out, err := xml.Marshal(unregisterIpTagCommand("1.1.1.1", []string{"db"}))

			Expect(err).ToNot(HaveOccurred())
			Expect(string(out)).To(Equal(
				`<uid-message><version>1.0</version><type>update</type>` +
					`<payload><unregister><entry ip="1.1.1.1">` +
					`<tag><member>db</member></tag>` +
					`</entry></unregister></payload></uid-message>`))
		})
	})
})

var _ = Describe("ip_tag registered-ip request marshaling", func() {
	Context("when filtering by IP only", func() {
		It("marshals a paginated show request with a 500 page limit", func() {
			out, err := xml.Marshal(registeredIpRequest("1.1.1.1", "", 1))

			Expect(err).ToNot(HaveOccurred())
			Expect(string(out)).To(Equal(
				`<show><object><registered-ip>` +
					`<ip>1.1.1.1</ip><limit>500</limit><start-point>1</start-point>` +
					`</registered-ip></object></show>`))
		})
	})

	Context("when filtering by both IP and tag", func() {
		It("includes a tag entry filter and the requested start point", func() {
			out, err := xml.Marshal(registeredIpRequest("1.1.1.1", "web", 501))

			Expect(err).ToNot(HaveOccurred())
			Expect(string(out)).To(Equal(
				`<show><object><registered-ip>` +
					`<tag><entry name="web"></entry></tag>` +
					`<ip>1.1.1.1</ip><limit>500</limit><start-point>501</start-point>` +
					`</registered-ip></object></show>`))
		})
	})
})

var _ = Describe("ip_tag parseRegisteredIpResponse", func() {
	Context("when PAN-OS returns registered IP entries", func() {
		It("parses each IP with its members into a map", func() {
			body := []byte(`<response status="success"><result>` +
				`<entry ip="1.1.1.1"><tag><member>web</member><member>db</member></tag></entry>` +
				`<entry ip="2.2.2.2"><tag><member>prod</member></tag></entry>` +
				`</result></response>`)

			result, err := parseRegisteredIpResponse(body)

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(HaveLen(2))
			Expect(result["1.1.1.1"]).To(ConsistOf("web", "db"))
			Expect(result["2.2.2.2"]).To(ConsistOf("prod"))
		})
	})

	Context("when PAN-OS returns an outfile instead of entries", func() {
		It("returns an error pointing at the unsupported response shape", func() {
			body := []byte(`<response status="success"><result>` +
				`<msg><line><outfile>regip.txt</outfile></line></msg>` +
				`</result></response>`)

			_, err := parseRegisteredIpResponse(body)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("outfile"))
		})
	})
})

var _ = Describe("ip_tag collectRegisteredIps pagination", func() {
	Context("when the first page is already shorter than the page limit", func() {
		It("returns that page and stops after one fetch", func() {
			var starts []int
			fetch := func(startPoint int) ([]ipTagRespEntry, error) {
				starts = append(starts, startPoint)
				return []ipTagRespEntry{{Ip: "1.1.1.1", Tags: []string{"web"}}}, nil
			}

			result, err := collectRegisteredIps(fetch, 2)

			Expect(err).ToNot(HaveOccurred())
			Expect(starts).To(Equal([]int{1}))
			Expect(result["1.1.1.1"]).To(ConsistOf("web"))
		})
	})

	Context("when a full page is followed by a short page", func() {
		It("assembles entries across pages and advances the start point", func() {
			pages := map[int][]ipTagRespEntry{
				1: {{Ip: "1.1.1.1", Tags: []string{"a"}}, {Ip: "2.2.2.2", Tags: []string{"b"}}},
				3: {{Ip: "3.3.3.3", Tags: []string{"c"}}},
			}
			var starts []int
			fetch := func(startPoint int) ([]ipTagRespEntry, error) {
				starts = append(starts, startPoint)
				return pages[startPoint], nil
			}

			result, err := collectRegisteredIps(fetch, 2)

			Expect(err).ToNot(HaveOccurred())
			Expect(starts).To(Equal([]int{1, 3}))
			Expect(result).To(HaveLen(3))
			Expect(result["3.3.3.3"]).To(ConsistOf("c"))
		})
	})

	Context("when a full page is followed by an empty page", func() {
		It("stops after the empty page without losing earlier entries", func() {
			pages := map[int][]ipTagRespEntry{
				1: {{Ip: "1.1.1.1", Tags: []string{"a"}}, {Ip: "2.2.2.2", Tags: []string{"b"}}},
				3: {},
			}
			var starts []int
			fetch := func(startPoint int) ([]ipTagRespEntry, error) {
				starts = append(starts, startPoint)
				return pages[startPoint], nil
			}

			result, err := collectRegisteredIps(fetch, 2)

			Expect(err).ToNot(HaveOccurred())
			Expect(starts).To(Equal([]int{1, 3}))
			Expect(result).To(HaveLen(2))
		})
	})
})
