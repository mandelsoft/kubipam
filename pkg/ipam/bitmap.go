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

const MAX_NET = 6
const MAX_BITS = 1 << MAX_NET

var covermask = [MAX_NET + 1]Bitmap{}

func init() {
	m := Bitmap(1)
	for i := 0; i <= MAX_NET; i++ {
		// m:=Bitmap(1<<(1<<i))-1
		covermask[i] = m
		m = (m+1)*(m+1) - 1
	}
}

func coverMask(reqsize int) Bitmap {
	return covermask[coverBits(reqsize)]
}

func coverBits(reqsize int) int {
	return MAX_NET - reqsize
}

func size(reqsize int) int {
	return 1 << coverBits(reqsize)
}

func (this Bitmap) canAllocate(reqsize int) int {
	s := size(reqsize)
	m := coverMask(reqsize)

	for c := 0; c <= MAX_BITS/s; c++ {
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
		(*this) |= coverMask(reqsize) << i
	}
	return i
}

func (this Bitmap) isAllocated(addr, reqsize int) bool {
	m := coverMask(reqsize) << addr
	return this&m == m
}

func (this Bitmap) isFree(addr, reqsize int) bool {
	m := coverMask(reqsize) << addr
	return this&m == 0
}

func (this *Bitmap) busy(addr, reqsize int) bool {
	m := coverMask(reqsize) << addr
	if (*this)&m != 0 {
		return false
	}
	(*this) |= m
	return true
}

func (this *Bitmap) free(addr, reqsize int) bool {
	m := coverMask(reqsize) << addr

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
