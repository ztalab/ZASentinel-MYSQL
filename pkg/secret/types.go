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

package secret

import "sync"

type SecretStore interface {
	Get(key string) string
	Set(key, value string)
	Delete(key string)
}

type DefaultSecretStore map[string]string

var defaultSecretLock sync.RWMutex

func (d DefaultSecretStore) Get(key string) string {
	defaultSecretLock.RLock()
	v := d[key]
	defaultSecretLock.RUnlock()
	return v
}

func (d DefaultSecretStore) Set(key, val string) {
	defaultSecretLock.Lock()
	d[key] = val
	defaultSecretLock.Unlock()
}

func (d DefaultSecretStore) Delete(key string) {
	defaultSecretLock.Lock()
	delete(d, key)
	defaultSecretLock.Unlock()
}
