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

type Bitmap uint64

const MAX_BITMAP_NET = 6
const MAX_BITMAP_SIZE = 1 << MAX_BITMAP_NET
const MAX_BITMAP_HOST_MASK = (1 << MAX_BITMAP_NET) - 1

var hostmask = [MAX_BITMAP_NET + 1]Bitmap{}

func init() {
	m := Bitmap(1)
	for i := 0; i <= MAX_BITMAP_NET; i++ {
		// m:=Bitmap(1<<(1<<i))-1
		hostmask[i] = m
		m = (m+1)*(m+1) - 1
	}
}

func bitmapHostMask(reqsize int) Bitmap {
	return hostmask[bitmapHostBits(reqsize)]
}

func bitmapHostBits(reqsize int) int {
	return MAX_BITMAP_NET - reqsize
}

func bitmapHostSize(reqsize int) int {
	return 1 << bitmapHostBits(reqsize)
}

func (this Bitmap) canAllocate(reqsize int) int {
	s := bitmapHostSize(reqsize)
	m := bitmapHostMask(reqsize)

	for c := 0; c <= MAX_BITMAP_SIZE/s; c++ {
		masked := this & m
		if masked == 0 {
			return c * s
		}
		m <<= s
	}
	return -1
}

func (this *Bitmap) allocate(reqsize int) int {
	i := (*this).canAllocate(reqsize)
	if i >= 0 {
		(*this) |= bitmapHostMask(reqsize) << i
	}
	return i
}

func (this Bitmap) isAllocated(addr, reqsize int) bool {
	m := bitmapHostMask(reqsize) << addr
	return this&m == m
}

func (this Bitmap) isFree(addr, reqsize int) bool {
	m := bitmapHostMask(reqsize) << addr
	return this&m == 0
}

func (this *Bitmap) busy(addr, reqsize int) bool {
	m := bitmapHostMask(reqsize) << addr
	if (*this)&m != 0 {
		return false
	}
	(*this) |= m
	return true
}

func (this *Bitmap) free(addr, reqsize int) bool {
	m := bitmapHostMask(reqsize) << addr

	if (*this)&m != m {
		return false
	}
	(*this) &= ^m
	return true
}

func (this *Bitmap) set(addr, reqsize int, busy bool) bool {
	if busy {
		return this.busy(addr, reqsize)
	}
	return this.free(addr, reqsize)
}
