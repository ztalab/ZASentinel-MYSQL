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

package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

const (
	YamlStorage  = "yaml"
	VaultStorage = "vault"
)

var DefaultConfig *Config

type Config struct {
	Server       ServerS       `json:"server" yaml:"server"`
	Confidential ConfidentialS `json:"confidential" yaml:"confidential"`
	FakeIdentity FakeIdentityS `json:"fakeIdentity" yaml:"fakeIdentity"`
}

type ServerS struct {
	Addr string `json:"addr" yaml:"addr"`
}

type ConfidentialS struct {
	Storage string `json:"storage" yaml:"storage"`
	Vault   VaultS `json:"vault" yaml:"vault"`
	Yaml    Mysql  `json:"yaml" yaml:"yaml"`
}

type VaultS struct {
	Addr     string `json:"addr" yaml:"addr"`
	Token    string `json:"token" yaml:"token"`
	DataPath string `json:"dataPath" yaml:"dataPath"`
}

type Mysql struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	DBName   string `json:"dbname" yaml:"dbname"`
}

type FakeIdentityS struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func InitConfig(path string) {
	if DefaultConfig != nil {
		return
	}
	DefaultConfig = &Config{}
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(DefaultConfig); err != nil {
		panic(err)
	}

	switch DefaultConfig.Confidential.Storage {
	case YamlStorage:
	case VaultStorage:
	default:
		panic("Unsupported storage type!")
	}

	if len(DefaultConfig.FakeIdentity.Username) == 0 || len(DefaultConfig.FakeIdentity.Password) == 0 {
		panic("fake-identity empty")
	}
}

func Get() *Config {
	return DefaultConfig
}

func GetRemoteMysql() *Mysql {
	if DefaultConfig.Confidential.Storage == YamlStorage {
		return &DefaultConfig.Confidential.Yaml
	}

	return getVaultData()
}

func getVaultData() *Mysql {
	url := fmt.Sprintf("%s/v1/secret/data/%s",
		strings.TrimRight(DefaultConfig.Confidential.Vault.Addr, "/"),
		DefaultConfig.Confidential.Vault.DataPath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Vault-Token", DefaultConfig.Confidential.Vault.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("request vault err: %v", err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("read vault response body err: %v", err))
	}

	res := gjson.GetBytes(body, "errors")
	if len(res.String()) > 0 {
		panic(fmt.Errorf("request vault err: %s", res.String()))
	}
	res = gjson.GetBytes(body, "data.data")
	if len(res.Map()) == 0 {
		panic(fmt.Errorf("request vault err: %s", res.String()))
	}
	result := res.Map()

	mysql := &Mysql{
		Host:     result["host"].String(),
		Port:     int(result["port"].Int()),
		Username: result["username"].String(),
		Password: result["password"].String(),
		DBName:   result["dbname"].String(),
	}

	if len(mysql.Host) == 0 || mysql.Port <= 0 || len(mysql.Username) == 0 || len(mysql.Password) == 0 {
		panic(fmt.Errorf("mysql config err: %s", res.String()))
	}
	return mysql
}
