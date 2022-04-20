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

// see: https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_handshake_v10.html
func (c *Conn) writeInitialHandshake() error {
	data := make([]byte, 4)

	// min version 10
	data = append(data, 10)

	// server version[00]
	data = append(data, c.serverConf.serverVersion...)
	data = append(data, 0x00)

	// connection id
	data = append(data, byte(c.connectionID), byte(c.connectionID>>8), byte(c.connectionID>>16), byte(c.connectionID>>24))

	// auth-plugin-data-part-1
	data = append(data, c.salt[0:8]...)

	// filter 0x00 byte, terminating the first part of a scramble
	data = append(data, 0x00)

	defaultFlag := c.serverConf.capability
	// capability flag lower 2 bytes, using default capability here
	data = append(data, byte(defaultFlag), byte(defaultFlag>>8))

	// charset
	data = append(data, c.serverConf.collationId)

	// status
	data = append(data, byte(c.status), byte(c.status>>8))

	// capability flag upper 2 bytes, using default capability here
	data = append(data, byte(defaultFlag>>16), byte(defaultFlag>>24))

	// server supports CLIENT_PLUGIN_AUTH and CLIENT_SECURE_CONNECTION
	data = append(data, byte(8+12+1))

	// reserved 10 [00]
	data = append(data, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

	// auth-plugin-data-part-2
	data = append(data, c.salt[8:]...)
	// second part of the password cipher [mininum 13 bytes],
	// where len=MAX(13, length of auth-plugin-data - 8)
	// add \NUL to terminate the string
	data = append(data, 0x00)

	// auth plugin name
	data = append(data, c.serverConf.defaultAuthMethod...)

	// EOF if MySQL version (>= 5.5.7 and < 5.5.10) or (>= 5.6.0 and < 5.6.2)
	// \NUL otherwise, so we use \NUL
	data = append(data, 0)

	return c.WritePacket(data)
}
