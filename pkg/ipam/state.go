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

func NewIPAM(cidr *net.IPNet) *IPAM {
	copy := *cidr
	if len(cidr.Mask) == net.IPv4len {
		copy.IP = cidr.IP.To4()
	} else {
		copy.IP = cidr.IP.To16()
	}
	block := &Block{
		cidr: &copy,
	}
	return &IPAM{
		block: block,
	}
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
	size := 0
	b := this.block

	for b != nil {
		s := b.Size()
		if b.canAlloc(reqsize) {
			if found == nil || s > size {
				found = b
				size = s
				if s == reqsize {
					break
				}
			}
		}
		b = b.next
	}
	if found == nil {
		return nil
	}
	for size < reqsize && found.canSplit() {
		found.split()
		size++
	}
	cidr := found.cidr
	found.busy = true
	this.cleanup(found)
	return cidr
}

func (this *IPAM) cleanup(b *Block) {
	for b != nil {
		b = b.join()
	}
}

func (this *IPAM) Busy(cidr *net.IPNet) bool {
	return this.set(cidr, true)
}

func (this *IPAM) Free(cidr *net.IPNet) bool {
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
	if b.busy == busy {
		return false
	}

	size := b.Size()
	for size < reqsize && b.canSplit() {
		upper := b.split()
		if upper.cidr.Contains(cidr.IP) {
			b = upper
		}
		size++
	}
	if size > reqsize {
		return false
	}

	b.set(cidr, busy)
	this.cleanup(b)
	return true
}

////////////////////////////////////////////////////////////////////////////////

type Block struct {
	busy bool
	cidr *net.IPNet
	prev *Block
	next *Block
}

func (this *Block) canAlloc(reqsize int) bool {
   return !this.busy && this.Size()<=reqsize
}

func (this *Block) canSplit() bool {
	return len(this.cidr.IP)*8-this.Size()>1
}

func (this *Block) matchState(b *Block) bool {
	return this.busy == b.busy
}

func (this *Block) set(cidr *net.IPNet, busy bool) bool {
	s, _ := cidr.Mask.Size()
	if s!=this.Size() || !this.cidr.Contains(cidr.IP) {
		return false
	}
	if this.busy == busy {
		return false
	}
	this.busy = busy
	return true
}

func (this *Block) alloc(reqsize int) *net.IPNet {
	if reqsize!=this.Size() {
		return nil
	}
	cidr := this.cidr
	this.busy = true
	return cidr
}

func (this *Block) split() *Block {
	ones, bits := this.cidr.Mask.Size()
	if bits == ones {
		return nil
	}

	mask := net.CIDRMask(ones+1, bits)
	delta := sub(mask, this.cidr.Mask)
	upper := &Block{
		cidr: &net.IPNet{
			IP:   net.IP(or(this.cidr.IP, delta)),
			Mask: mask,
		},
		busy: this.busy,
		prev: this,
		next: this.next,
	}
	if this.next != nil {
		this.next.prev = upper
	}
	this.next = upper
	this.cidr = &net.IPNet{
		IP:   this.cidr.IP,
		Mask: mask,
	}
	return upper
}

func (this *Block) pair() (*Block, *Block) {
	if this.IsUpper() {
		return this.prev, this
	}
	return this, this.next
}

func (this *Block) join() *Block {
	lower, upper := this.pair()
	if lower == nil || upper == nil {
		return nil
	}
	ones, bits := lower.cidr.Mask.Size()
	if ones != upper.Size() {
		return nil
	}
	if !lower.matchState(upper) {
		return nil
	}

	if upper.next != nil {
		upper.next.prev = lower
	}
	lower.next = upper.next

	mask := net.CIDRMask(ones-1, bits)
	lower.cidr = &net.IPNet{
		IP:   lower.cidr.IP,
		Mask: mask,
	}
	return lower
}

func (this *Block) String() string {
	msg := "free"
	if this.busy {
		msg = "busy"
	}
	return fmt.Sprintf("%s[%s]", this.cidr.String(), msg)
}

func (this *Block) Next() *Block {
	return this.next
}

func (this *Block) Prev() *Block {
	return this.prev
}

func (this *Block) Size() int {
	ones, _ := this.cidr.Mask.Size()
	return ones
}

func (this *Block) IsUpper() bool {
	ones, bits := this.cidr.Mask.Size()
	delta := sub(this.cidr.Mask, net.CIDRMask(ones-1, bits))
	return !isZero(and(delta, this.cidr.IP))
}

////////////////////////////////////////////////////////////////////////////////

func sub(a, b []byte) []byte {
	if len(a) != len(b) {
		return nil
	}
	new := make(net.IPMask, len(a), len(a))
	for i, v := range a {
		new[i] = v - b[i]
	}
	return new
}

func or(a, b []byte) []byte {
	if len(a) != len(b) {
		return nil
	}
	new := make(net.IPMask, len(a), len(a))
	for i, v := range a {
		new[i] = v | b[i]
	}
	return new
}

func and(a, b []byte) []byte {
	if len(a) != len(b) {
		return nil
	}
	new := make(net.IPMask, len(a), len(a))
	for i, v := range a {
		new[i] = v & b[i]
	}
	return new
}

func isZero(a []byte) bool {
	for _, v := range a {
		if v != 0 {
			return false
		}
	}
	return true
}
