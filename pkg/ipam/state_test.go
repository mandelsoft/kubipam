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
	_, cidr, _:= net.ParseCIDR("10.0.0.0/8")

	It("initializes ipam correctly", func() {
		ipam:=NewIPAM(cidr)

		Expect(ipam.block.next).To(BeNil())
		Expect(ipam.block.prev).To(BeNil())
	})

	It("initializes splits", func() {
		ipam:=NewIPAM(cidr)

		r:=ipam.Alloc(9)
		Expect(r.String()).To(Equal("10.0.0.0/9"))

		r=ipam.Alloc(10)
		Expect(r.String()).To(Equal("10.128.0.0/10"))
	})

	It("free", func() {
		ipam:=NewIPAM(cidr)

		r1:=ipam.Alloc(9)
		Expect(r1.String()).To(Equal("10.0.0.0/9"))

		Expect(ipam.Free(r1)).To(BeTrue())

		Expect(ipam.block.next).To(BeNil())
		Expect(ipam.block.prev).To(BeNil())
	})

	/*
	It("scenario", func() {
		ipam:=NewIPAM(cidr)

		r1:=ipam.Alloc(9)
		Expect(r1.String()).To(Equal("10.0.0.0/9"))

		r2:=ipam.Alloc(10)
		Expect(r2.String()).To(Equal("10.128.0.0/10"))

		r3:=ipam.Alloc(12)
		Expect(r3.String()).To(Equal("10.192.0.0/12"))

		r4:=ipam.Alloc(11)
		Expect(r4.String()).To(Equal("10.224.0.0/11"))


		Expect(ipam.Free(r1)).To(BeTrue())
		Expect(ipam.Free(r3)).To(BeTrue())
		Expect(ipam.Free(r2)).To(BeTrue())
		Expect(ipam.Free(r4)).To(BeTrue())

		Expect(ipam.block.next).To(BeNil())
		Expect(ipam.block.prev).To(BeNil())
	})
	 */

})