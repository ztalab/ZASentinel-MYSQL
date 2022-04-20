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

package client

import (
	"crypto/tls"
	"crypto/x509"
)

// NewClientTLSConfig: generate TLS config for client side
// if insecureSkipVerify is set to true, serverName will not be validated
func NewClientTLSConfig(caPem, certPem, keyPem []byte, insecureSkipVerify bool, serverName string) *tls.Config {
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPem) {
		panic("failed to add ca PEM")
	}

	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		panic(err)
	}

	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            pool,
		InsecureSkipVerify: insecureSkipVerify,
		ServerName:         serverName,
	}
	return config
}
