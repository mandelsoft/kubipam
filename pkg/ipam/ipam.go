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
	"fmt"
	"net"
)

type IPAM struct {
	block *Block
}

func NewIPAM(cidr *net.IPNet, ranges ...*IPRange) (*IPAM, error) {
	copy := *cidr
	if len(cidr.Mask) == net.IPv4len {
		copy.IP = cidr.IP.To4()
	} else {
		copy.IP = cidr.IP.To16()
	}
	block := &Block{
		cidr: &copy,
	}
	ipam := &IPAM{
		block: block,
	}
	if len(ranges) > 0 {
		cidrs, err := Excludes(cidr, ranges...)
		if err != nil {
			return nil, err
		}
		for _, c := range cidrs {
			ipam.Busy(c)
		}

		for b := ipam.block; b != nil; b = b.next {
			if b.isCompletelyBusy() {
				if b.prev != nil {
					b.prev.next = b.next
				}
				if b.next != nil {
					b.next.prev = b.prev
				}
				if b == ipam.block {
					ipam.block = b.next
				}
			}
		}
		if ipam.block == nil {
			return nil, fmt.Errorf("no available IP addresses")
		}
	}
	return ipam, nil
}

func NewIPAMForRanges(ranges IPRanges) (*IPAM, error) {
	if len(ranges) == 0 {
		return nil, fmt.Errorf("no ranges specified for IPAM")
	}
	ipam := &IPAM{}
	cidrs, err := Includes(ranges...)
	if err != nil {
		return nil, err
	}

	ipv4 := true
	for _, cidr := range cidrs {
		if cidr.IP.To4() == nil {
			ipv4 = false
		}
		break
	}

	last := ipam.block
	for _, cidr := range cidrs {
		var b *Block
		if ipv4 {
			cidr = CIDRto4(cidr)
		} else {
			cidr = CIDRto16(cidr)
		}
		if CIDRHostMaskSize(cidr) < MAX_BITMAP_NET {
			n := cidr
			for CIDRHostMaskSize(n) < MAX_BITMAP_NET {
				n = CIDRExtend(n)
			}
			if last != nil && CIDREqual(n, last.cidr) {
				last.set(cidr, false)
				continue
			}
			b = &Block{cidr: n, busy: BITMAP_BUSY}
			b.set(cidr, false)
		} else {
			b = &Block{cidr: cidr}
		}
		b.prev = last
		if last != nil {
			last.next = b
		} else {
			ipam.block = b
		}
		last = b
	}
	return ipam, nil
}

func (this *IPAM) Bits() int {
	return CIDRBits(this.block.cidr)
}

func (this *IPAM) String() string {
	s := ""
	sep := ""
	b := this.block
	for b != nil {
		s = fmt.Sprintf("%s%s%s", s, sep, b)
		sep = ", "
		b = b.next
	}
	return s
}

func (this *IPAM) Alloc(reqsize int) *net.IPNet {
	var found *Block
	b := this.block

	for b != nil {
		s := b.Size()
		if b.canAlloc(reqsize) {
			if found == nil || s > found.Size() {
				found = b
				if found.matchSize(reqsize) {
					break
				}
			}
		}
		b = b.next
	}
	if found == nil {
		return nil
	}
	found = this.split(found, reqsize)

	cidr := found.alloc(reqsize)
	if cidr != nil {
		this.join(found)
	}
	return cidr
}

func (this *IPAM) split(b *Block, reqsize int) *Block {
	for b.Size() < reqsize && b.canSplit() {
		b.split()
	}
	return b
}

func (this *IPAM) join(b *Block) {
	for b != nil {
		b = b.join()
	}
}

func (this *IPAM) Busy(cidr *net.IPNet) bool {
	cidr=CIDRAlign(cidr, this.Bits())
	if cidr == nil {
		return false
	}
	return this.set(cidr, true)
}

func (this *IPAM) Free(cidr *net.IPNet) bool {
	cidr=CIDRAlign(cidr, this.Bits())
	if cidr == nil {
		return false
	}
	return this.set(cidr, false)
}

func (this *IPAM) set(cidr *net.IPNet, busy bool) bool {
	reqsize, _ := cidr.Mask.Size()
	b := this.block
	for b != nil && !b.cidr.Contains(cidr.IP) {
		b = b.next
	}
	if b == nil {
		return false
	}

	size := b.Size()
	if b.canSplit() {
		if b.isBusy() == busy {
			return false
		}
		for size < reqsize && b.canSplit() {
			upper := b.split()
			if upper.cidr.Contains(cidr.IP) {
				b = upper
			}
			size++
		}
	}

	if size > reqsize {
		return false
	}

	if !b.set(cidr, busy) {
		return false
	}
	this.join(b)
	return true
}
