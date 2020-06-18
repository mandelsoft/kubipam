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

package main

import (
	"fmt"
	"net"

	ipam "github.com/mandelsoft/kubipam/pkg/ipam"
)

func main() {
	_, cidr, _ := net.ParseCIDR("100.64.0.0/16")
	ipam, _ := ipam.NewIPAM(cidr)

	fmt.Printf("initial: %s\n", ipam)

	a1 := ipam.Alloc(23)
	a2 := ipam.Alloc(24)
	a3 := ipam.Alloc(23)
	a4 := ipam.Alloc(24)
	fmt.Printf("alloc  : %s\n", ipam)

	ipam.Free(a1)
	fmt.Printf("free a1: %s\n", ipam)
	ipam.Free(a2)
	fmt.Printf("free a2: %s\n", ipam)
	ipam.Free(a4)
	fmt.Printf("free a4: %s\n", ipam)
	ipam.Free(a3)
	fmt.Printf("free a3: %s\n", ipam)
}
