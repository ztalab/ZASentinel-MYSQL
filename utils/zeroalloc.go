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

import "unsafe"

func StringToByteSlice(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func ByteSliceToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Uint64ToInt64(val uint64) int64 {
	return *(*int64)(unsafe.Pointer(&val))
}

func Uint64ToFloat64(val uint64) float64 {
	return *(*float64)(unsafe.Pointer(&val))
}

func Int64ToUint64(val int64) uint64 {
	return *(*uint64)(unsafe.Pointer(&val))
}

func Float64ToUint64(val float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&val))
}
