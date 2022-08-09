/*
Copyright 2022-present The Ztalab Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import "sync"

var (
	byteSlicePool = sync.Pool{
		New: func() interface{} {
			return []byte{}
		},
	}
	byteSliceChan = make(chan []byte, 10)
)

func ByteSliceGet(length int) (data []byte) {
	select {
	case data = <-byteSliceChan:
	default:
		data = byteSlicePool.Get().([]byte)[:0]
	}

	if cap(data) < length {
		data = make([]byte, length)
	} else {
		data = data[:length]
	}

	return data
}

func ByteSlicePut(data []byte) {
	select {
	case byteSliceChan <- data:
	default:
		byteSlicePool.Put(data) //nolint:staticcheck
	}
}
