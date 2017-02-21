// Copyright 2017-Present Pivotal Software, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package content

import (
	"errors"
	"fmt"
)

type Ranger struct {
	numHunks int
}

func NewRanger(hunks int) Ranger {
	return Ranger{
		numHunks: hunks,
	}
}

func (r Ranger) BuildRange(contentLength int64) ([]string, error) {
	if contentLength == 0 {
		return []string{}, errors.New("content length cannot be zero")
	}

	var ranges []string
	hunkSize := contentLength / int64(r.numHunks)
	if hunkSize == 0 {
		hunkSize = 2
	}

	iterations := (contentLength / hunkSize)
	remainder := contentLength % int64(hunkSize)

	for i := int64(0); i < int64(iterations); i++ {
		lowerByte := i * hunkSize
		upperByte := ((i + 1) * hunkSize) - 1
		if i == int64(iterations-1) {
			upperByte += remainder
		}
		bytes := fmt.Sprintf("%d-%d", lowerByte, upperByte)
		ranges = append(ranges, bytes)
	}

	return ranges, nil
}
