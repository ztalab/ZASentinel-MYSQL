/*
 *   Licensed to the Apache Software Foundation (ASF) under one or more
 *  contributor license agreements.  See the NOTICE file distributed with
 *  this work for additional information regarding copyright ownership.
 *  The ASF licenses this file to You under the Apache License, Version 2.0
 *  (the "License"); you may not use this file except in compliance with
 *  the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package utils

import (
	"bytes"
	"sync"
)

var (
	bytesBufferPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	bytesBufferChan = make(chan *bytes.Buffer, 10)
)

func BytesBufferGet() (data *bytes.Buffer) {
	select {
	case data = <-bytesBufferChan:
	default:
		data = bytesBufferPool.Get().(*bytes.Buffer)
	}

	data.Reset()

	return data
}

func BytesBufferPut(data *bytes.Buffer) {
	select {
	case bytesBufferChan <- data:
	default:
		bytesBufferPool.Put(data)
	}
}
