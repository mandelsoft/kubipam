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
)

func MustParseCIDR(s string) *net.IPNet {
	ip, cidr, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return &net.IPNet{
		IP:   ip,
		Mask: cidr.Mask,
	}
}

func CIDRNetMaskSize(cidr *net.IPNet) int {
	s, _ := cidr.Mask.Size()
	return s
}

func CIDRHostMaskSize(cidr *net.IPNet) int {
	s, l := cidr.Mask.Size()
	return l - s
}

////////////////////////////////////////////////////////////////////////////////

func sub(a, b net.IPMask) []byte {
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

func CloneIP(ip net.IP) net.IP {
	return append(ip[:0:0], ip...)
}

func CloneMask(mask net.IPMask) net.IPMask {
	return append(mask[:0:0], mask...)
}

func CloneCIDR(cidr *net.IPNet) *net.IPNet {
	return &net.IPNet{
		IP:   CloneIP(cidr.IP),
		Mask: CloneMask(cidr.Mask),
	}
}

func CIDRSubIP(cidr *net.IPNet, n int64) net.IP {
	if n < 0 || n >= 1<<CIDRHostMaskSize(cidr) {
		return nil
	}
	return IPAdd(CIDRFirstIP(cidr), n)
}

func CIDRNet(cidr *net.IPNet) *net.IPNet {
	net := *cidr
	net.IP = CIDRFirstIP(cidr)
	return &net
}

func CIDRFirstIP(cidr *net.IPNet) net.IP {
	return cidr.IP.Mask(cidr.Mask)
}

func CIDRLastIP(cidr *net.IPNet) net.IP {
	return CIDRSubIP(cidr, (1<<CIDRHostMaskSize(cidr))-1)
}

func IPAdd(ip net.IP, n int64) net.IP {
	ip = CloneIP(ip)
	for i := len(ip) - 1; n > 0; i-- {
		n += int64(ip[i])
		ip[i] = uint8(n & 0xff)
		n >>= 8
	}
	return ip
}

func IPDiff(a, b net.IP) int64 {
	var d int64
	a = a.To16()
	b = b.To16()
	for i, _ := range a {
		db := int64(a[i]) - int64(b[i])
		d = d*256 + db
	}
	return d
}

func CIDRSplit(cidr *net.IPNet) (lower, upper *net.IPNet) {
	ones, bits := cidr.Mask.Size()
	if ones == bits {
		return nil, nil
	}
	mask := net.CIDRMask(ones+1, bits)
	delta := sub(mask, cidr.Mask)
	upper = &net.IPNet{
		IP:   net.IP(or(cidr.IP, delta)),
		Mask: mask,
	}
	lower = &net.IPNet{
		IP:   cidr.IP,
		Mask: mask,
	}
	return
}

func CIDRExtend(cidr *net.IPNet) *net.IPNet {
	ones, bits := cidr.Mask.Size()
	if ones == 0 {
		return nil
	}
	mask := net.CIDRMask(ones-1, bits)
	return &net.IPNet{
		IP:   net.IP(and(cidr.IP, mask)),
		Mask: mask,
	}
}

func IPtoCIDR(ip net.IP) *net.IPNet {
	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(len(ip)*8, len(ip)*8),
	}
}
