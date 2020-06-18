/*
 * Copyright 2020 Mandelsoft. All rights reserved.
 *  This file is licensed under the Apache Software License, v. 2 except as noted
 *  otherwise in the LICENSE file
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package ipam

import (
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IPAM", func() {
	Context("using complete blocks", func() {
		_, cidr, _ := net.ParseCIDR("10.0.0.0/8")

		It("initializes ipam correctly", func() {
			ipam, _ := NewIPAM(cidr)

			Expect(ipam.block.next).To(BeNil())
			Expect(ipam.block.prev).To(BeNil())
		})

		It("initializes splits", func() {
			ipam, _ := NewIPAM(cidr)

			r := ipam.Alloc(9)
			Expect(r.String()).To(Equal("10.0.0.0/9"))

			r = ipam.Alloc(10)
			Expect(r.String()).To(Equal("10.128.0.0/10"))

			Expect(ipam.String()).To(Equal("10.0.0.0/9[busy], 10.128.0.0/10[busy], 10.192.0.0/10[free]"))
		})

		It("free", func() {
			ipam, _ := NewIPAM(cidr)

			r1 := ipam.Alloc(9)
			Expect(r1.String()).To(Equal("10.0.0.0/9"))

			Expect(ipam.Free(r1)).To(BeTrue())

			Expect(ipam.block.next).To(BeNil())
			Expect(ipam.block.prev).To(BeNil())
		})

		It("scenario", func() {
			ipam, _ := NewIPAM(cidr)

			r1 := ipam.Alloc(9)
			Expect(r1.String()).To(Equal("10.0.0.0/9"))

			r2 := ipam.Alloc(10)
			Expect(r2.String()).To(Equal("10.128.0.0/10"))

			r3 := ipam.Alloc(12)
			Expect(r3.String()).To(Equal("10.192.0.0/12"))

			r4 := ipam.Alloc(11)
			Expect(r4.String()).To(Equal("10.224.0.0/11"))

			Expect(ipam.String()).To(Equal("10.0.0.0/9[busy], 10.128.0.0/10[busy], 10.192.0.0/12[busy], 10.208.0.0/12[free], 10.224.0.0/11[busy]"))

			Expect(ipam.Free(r1)).To(BeTrue())
			Expect(ipam.Free(r3)).To(BeTrue())
			Expect(ipam.Free(r2)).To(BeTrue())
			Expect(ipam.Free(r4)).To(BeTrue())

			Expect(ipam.block.next).To(BeNil())
			Expect(ipam.block.prev).To(BeNil())
		})
	})

	Context("using bitmaps", func() {
		_, cidr, _ := net.ParseCIDR("10.0.0.0/26")

		It("initializes ipam correctly", func() {
			ipam, _ := NewIPAM(cidr)

			Expect(ipam.block.next).To(BeNil())
			Expect(ipam.block.prev).To(BeNil())
		})

		It("check 28", func() {
			ipam, _ := NewIPAM(cidr)

			r := ipam.Alloc(28)
			Expect(r.String()).To(Equal("10.0.0.0/28"))

			Expect(ipam.block.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 11111111 11111111]"))
		})

		It("check 28/30/28", func() {
			ipam, _ := NewIPAM(cidr)

			r1 := ipam.Alloc(28)
			Expect(r1.String()).To(Equal("10.0.0.0/28"))
			r2 := ipam.Alloc(30)
			Expect(r2.String()).To(Equal("10.0.0.16/30"))
			r3 := ipam.Alloc(28)
			Expect(r3.String()).To(Equal("10.0.0.32/28"))

			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 11111111 11111111 00000000 00001111 11111111 11111111]"))
		})

		It("free 28", func() {
			ipam, _ := NewIPAM(cidr)

			r1 := ipam.Alloc(28)
			Expect(r1.String()).To(Equal("10.0.0.0/28"))

			Expect(ipam.Free(r1)).To(BeTrue())
			Expect(ipam.String()).To(Equal("10.0.0.0/26[free]"))
		})

		It("scenario", func() {
			ipam, _ := NewIPAM(cidr)

			r1 := ipam.Alloc(28)
			Expect(r1.String()).To(Equal("10.0.0.0/28"))
			r2 := ipam.Alloc(30)
			Expect(r2.String()).To(Equal("10.0.0.16/30"))
			r3 := ipam.Alloc(28)
			Expect(r3.String()).To(Equal("10.0.0.32/28"))

			Expect(ipam.Free(r1)).To(BeTrue())
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 11111111 11111111 00000000 00001111 00000000 00000000]"))
			Expect(ipam.Free(r3)).To(BeTrue())
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00001111 00000000 00000000]"))
			Expect(ipam.Free(r2)).To(BeTrue())
			Expect(ipam.String()).To(Equal("10.0.0.0/26[free]"))

		})
	})

	Context("mixed", func() {
		_, cidr, _ := net.ParseCIDR("10.0.0.0/24")

		It("initializes ipam correctly", func() {
			ipam, _ := NewIPAM(cidr)

			Expect(ipam.block.next).To(BeNil())
			Expect(ipam.block.prev).To(BeNil())

			Expect(ipam.String()).To(Equal("10.0.0.0/24[free]"))
		})

		It("check 32/25", func() {
			ipam, _ := NewIPAM(cidr)

			r1 := ipam.Alloc(32)
			Expect(r1.String()).To(Equal("10.0.0.0/32"))
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001], 10.0.0.64/26[free], 10.0.0.128/25[free]"))

			r2 := ipam.Alloc(25)
			Expect(r2.String()).To(Equal("10.0.0.128/25"))
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001], 10.0.0.64/26[free], 10.0.0.128/25[busy]"))
		})

		It("scenario", func() {
			ipam, _ := NewIPAM(cidr)

			r1 := ipam.Alloc(32)
			Expect(r1.String()).To(Equal("10.0.0.0/32"))
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001], 10.0.0.64/26[free], 10.0.0.128/25[free]"))

			r2 := ipam.Alloc(25)
			Expect(r2.String()).To(Equal("10.0.0.128/25"))
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001], 10.0.0.64/26[free], 10.0.0.128/25[busy]"))

			Expect(ipam.Free(r1)).To(BeTrue())
			Expect(ipam.String()).To(Equal("10.0.0.0/25[free], 10.0.0.128/25[busy]"))

			Expect(ipam.Free(r2)).To(BeTrue())
			Expect(ipam.String()).To(Equal("10.0.0.0/24[free]"))
		})

		It("scenario 1", func() {
			ipam, _ := NewIPAM(cidr)
			ipam.Busy(MustParseCIDR("10.0.0.127/32"))
			Expect(ipam.String()).To(Equal("10.0.0.0/26[free], 10.0.0.64/26[10000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000], 10.0.0.128/25[free]"))
		})

		It("scenario 2", func() {
			ipam, _ := NewIPAM(cidr)
			ipam.Busy(MustParseCIDR("10.0.0.0/29"))
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000000 11111111], 10.0.0.64/26[free], 10.0.0.128/25[free]"))
			ipam.Busy(MustParseCIDR("10.0.0.8/32"))
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000001 11111111], 10.0.0.64/26[free], 10.0.0.128/25[free]"))
			ipam.Busy(MustParseCIDR("10.0.0.128/25"))
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000001 11111111], 10.0.0.64/26[free], 10.0.0.128/25[busy]"))
			ipam.Busy(MustParseCIDR("10.0.0.127/32"))
			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000001 11111111], 10.0.0.64/26[10000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000], 10.0.0.128/25[busy]"))
		})
	})

	Context("range", func() {
		_, cidr, _ := net.ParseCIDR("10.0.0.0/24")

		It("initializes ipam correctly", func() {
			ipam, _ := NewIPAM(cidr, MustParseIPRange("10.0.0.10-10.0.0.250"))

			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000011 11111111], 10.0.0.64/26[free], 10.0.0.128/26[free], 10.0.0.192/26[11111000 00000000 00000000 00000000 00000000 00000000 00000000 00000000]"))
		})

		It("initializes ipam correctly with sparse range", func() {
			ipam, _ := NewIPAM(cidr, MustParseIPRange("10.0.0.10-10.0.0.126"))

			Expect(ipam.String()).To(Equal("10.0.0.0/26[00000000 00000000 00000000 00000000 00000000 00000000 00000011 11111111], 10.0.0.64/26[10000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000]"))
		})
	})
})
