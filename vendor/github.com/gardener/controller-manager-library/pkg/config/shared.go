/*
 * Copyright 2019 SAP SE or an SAP affiliate company. All rights reserved.
 * This file is licensed under the Apache Software License, v. 2 except as noted
 * otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 *
 */

package config

import (
	"fmt"
	"reflect"
)

type SharedOptionSet struct {
	*DefaultOptionSet
	unshared map[string]bool
	shared   OptionSet

	descriptionMapper StringMapper
}

var _ OptionGroup = (*SharedOptionSet)(nil)

func NewSharedOptionSet(name, prefix string, descMapper StringMapper) *SharedOptionSet {
	if descMapper == nil {
		descMapper = IdenityStringMapper
	}
	s := &SharedOptionSet{
		DefaultOptionSet:  NewDefaultOptionSet(name, prefix),
		unshared:          map[string]bool{},
		descriptionMapper: descMapper,
	}
	return s
}

func (this *SharedOptionSet) Unshare(name string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.unshared[name] = true
}

func (this *SharedOptionSet) AddOptionsToSet(set OptionSet) {
	this.Complete()
	this.lock.Lock()
	defer this.lock.Unlock()

	this.shared = set
	for name, o := range this.arbitraryOptions {
		unshared := this.unshared[name]
		if this.prefix != "" || unshared {
			this.addOptionToSet(o, set, this.descriptionMapper(o.Description))
		}
		if !unshared {
			if old := set.GetOption(name); old != nil {
				if o.Type != old.Type {
					panic(fmt.Sprintf("type mismatch for shared option %s (%s)", name, this.prefix))
				}
			} else {
				set.AddOption(o.Type, nil, o.Name, o.Flag().Shorthand, nil, o.Description)
			}
		}
	}
}

func (this *SharedOptionSet) evalShared() {
	this.lock.Lock()
	defer this.lock.Unlock()

	// fmt.Printf("eval shared %s\n", this.prefix)
	for name, o := range this.arbitraryOptions {
		if !this.unshared[name] && !o.Changed() {
			// fmt.Printf("eval shared %s\n", name)
			shared := this.shared.GetOption(name)
			if shared.Changed() {
				value := reflect.ValueOf(shared.Target).Elem()
				// fmt.Printf("   %s changed shared\n", name)
				o.Flag().Changed = true
				reflect.ValueOf(o.Target).Elem().Set(value)
			}
		}
	}
}

func (this *SharedOptionSet) Evaluate() error {
	this.evalShared()
	return this.DefaultOptionSet.Evaluate()
}
