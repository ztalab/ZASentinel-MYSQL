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

package server

import (
	"github.com/ztalab/ZASentinel-MYSQL/pkg/secret"
	//asyncruntime "github.com/ztalab/ZASentinel-MYSQL/pkg/runtime"
	//"gitlab.oneitfarm.com/bifrost/cilog"
)

const (
	VaultMysqlUsername = "USERNAME"
	VaultMysqlPassword = "PASSWORD"
	VaultMysqlHost     = "HOST"
	VaultMysqlPort     = "PORT"
	VaultMysqlDbname   = "NAME"
)

// username == password == dbname  == mysqltag
type VaultProvider struct {
	store secret.SecretStore // username -> password
}

func NewVaultProvider(scr secret.SecretStore) *VaultProvider {
	return &VaultProvider{store: scr}
}

func (m *VaultProvider) CheckUsername(username string) (found bool, err error) {
	val := m.store.Get(username)
	return val != "", nil
}

func (m *VaultProvider) GetCredential(username string) (password string, found bool, err error) {
	val := m.store.Get(username)
	return val, val != "", nil
}

func (m *VaultProvider) AddUser(username, password string) {
	m.store.Set(username, password)
}
