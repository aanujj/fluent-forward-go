/*
MIT License

Copyright contributors to the fluent-forward-go project

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package protocol_test

import (
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aanujj/fluent-forward-go/fluent/protocol"
)

var _ = Describe("Transport", func() {
	Describe("EventTime", func() {
		var (
			ent protocol.EntryExt
		)

		BeforeEach(func() {
			ent = protocol.EntryExt{
				Timestamp: protocol.EventTime{
					Time: time.Unix(int64(1257894000), int64(12340000)),
				},
			}
		})

		// This covers both MarshalBinaryTo() and UnmarshalBinary()
		It("Marshals and unmarshals correctly", func() {
			b, err := ent.MarshalMsg(nil)

			Expect(err).NotTo(HaveOccurred())

			// This is the msgpack fixext8 encoding for the timestamp
			// per the fluent-forward spec:
			// D7 == fixext8
			// 00 == type 0
			// 4AF9F070 == 1257894000
			// 00BC4B20 == 12340000
			Expect(
				strings.Contains(fmt.Sprintf("%X", b), "D7004AF9F07000BC4B20"),
			).To(BeTrue())

			var unment protocol.EntryExt
			_, err = unment.UnmarshalMsg(b)
			Expect(err).NotTo(HaveOccurred())

			Expect(unment.Timestamp.Time.Equal(ent.Timestamp.Time)).To(BeTrue())
		})
	})

	Describe("EntryList", func() {
		var (
			e1 protocol.EntryList
			et time.Time
		)

		BeforeEach(func() {
			et = time.Now()
			e1 = protocol.EntryList{
				{
					Timestamp: protocol.EventTime{et},
					Record: map[string]interface{}{
						"foo":    "bar",
						"george": "jungle",
					},
				},
				{
					Timestamp: protocol.EventTime{et},
					Record: map[string]interface{}{
						"foo":    "kablooie",
						"george": "frank",
					},
				},
			}
		})

		Describe("Un/MarshalPacked", func() {
			var (
				e2 protocol.EntryList
			)

			BeforeEach(func() {
				e2 = protocol.EntryList{
					{
						Timestamp: protocol.EventTime{et},
						Record: map[string]interface{}{
							"foo":    "bar",
							"george": "jungle",
						},
					},
					{
						Timestamp: protocol.EventTime{et},
						Record: map[string]interface{}{
							"foo":    "kablooie",
							"george": "frank",
						},
					},
				}
			})

			It("Can marshal and unmarshal packed entries", func() {
				b, err := e2.MarshalPacked()
				Expect(err).ToNot(HaveOccurred())

				el := protocol.EntryList{}
				_, err = el.UnmarshalPacked(b)
				Expect(err).ToNot(HaveOccurred())
				Expect(el.Equal(e2)).To(BeTrue())
			})
		})

		Describe("Equal", func() {
			var (
				e2 protocol.EntryList
			)

			BeforeEach(func() {
				e2 = protocol.EntryList{
					{
						Timestamp: protocol.EventTime{et},
						Record: map[string]interface{}{
							"foo":    "bar",
							"george": "jungle",
						},
					},
					{
						Timestamp: protocol.EventTime{et},
						Record: map[string]interface{}{
							"foo":    "kablooie",
							"george": "frank",
						},
					},
				}
			})

			It("Returns true", func() {
				Expect(e1.Equal(e2)).To(BeTrue())
			})

			Context("When the lists have different element counts", func() {
				BeforeEach(func() {
					e2 = e2[:1]
				})

				It("Returns false", func() {
					Expect(e1.Equal(e2)).To(BeFalse())
				})
			})

			Context("When the lists have differing elements", func() {
				BeforeEach(func() {
					e2[0].Timestamp = protocol.EventTime{et.Add(5 * time.Second)}
				})

				It("Returns false", func() {
					Expect(e1.Equal(e2)).To(BeFalse())
				})
			})
		})
	})

	Describe("NewPackedForwardMessage", func() {
		var (
			tag     string
			entries protocol.EntryList
		)

		BeforeEach(func() {
			tag = "foo.bar"
			entries = protocol.EntryList{
				{
					Timestamp: protocol.EventTime{time.Now()},
					Record: map[string]interface{}{
						"foo":    "bar",
						"george": "jungle",
					},
				},
				{
					Timestamp: protocol.EventTime{time.Now()},
					Record: map[string]interface{}{
						"foo":    "kablooie",
						"george": "frank",
					},
				},
			}
		})

		It("Returns a PackedForwardMessage", func() {
			msg, err := protocol.NewPackedForwardMessage(tag, entries)
			Expect(err).NotTo(HaveOccurred())
			Expect(msg).NotTo(BeNil())
			Expect(*msg.Options.Size).To(Equal(len(entries)))
			Expect(msg.Options.Compressed).To(BeEmpty())
		})

		It("Correctly encodes the entries into a bytestream", func() {
			msg, err := protocol.NewPackedForwardMessage(tag, entries)
			Expect(err).NotTo(HaveOccurred())
			elist := make(protocol.EntryList, 2)
			_, err = elist.UnmarshalPacked(msg.EventStream)
			Expect(err).NotTo(HaveOccurred())
			Expect(elist.Equal(entries)).To(BeTrue())
		})
	})

	Describe("NewCompressedPackedForwardMessage", func() {
		var (
			tag     string
			entries []protocol.EntryExt
		)

		BeforeEach(func() {
			tag = "foo.bar"
			entries = []protocol.EntryExt{
				{
					Timestamp: protocol.EventTime{time.Now()},
					Record: map[string]interface{}{
						"foo":    "bar",
						"george": "jungle",
					},
				},
				{
					Timestamp: protocol.EventTime{time.Now()},
					Record: map[string]interface{}{
						"foo":    "kablooie",
						"george": "frank",
					},
				},
			}
		})

		It("Returns a message with a gzip-compressed event stream", func() {
			msg, err := protocol.NewCompressedPackedForwardMessage(tag, entries)
			Expect(err).ToNot(HaveOccurred())
			Expect(msg).NotTo(BeNil())
			Expect(*msg.Options.Size).To(Equal(len(entries)))
			Expect(msg.Options.Compressed).To(Equal("gzip"))
		})
	})
})
