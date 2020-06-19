/*
 * Copyright 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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
 *
 *
 */

package simple

import (
	"encoding/json"

	"github.com/gardener/controller-manager-library/pkg/types/infodata"
)

const T_STRINGARRAY = infodata.TypeVersion("StringArray")

func init() {
	infodata.Register(T_STRINGARRAY, infodata.UnmarshalFunc((StringArray)(nil)))
}

type StringArray []string

func (this StringArray) TypeVersion() infodata.TypeVersion {
	return T_STRINGARRAY
}

func (this StringArray) Marshal() ([]byte, error) {
	return json.Marshal(&this)
}
