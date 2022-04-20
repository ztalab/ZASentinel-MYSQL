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
	"github.com/ztalab/ZASentinel-MYSQL/utils"
)

func (c *Conn) writeCommand(command byte) error {
	c.ResetSequence()

	return c.WritePacket([]byte{
		0x01, // 1 bytes long
		0x00,
		0x00,
		0x00, // sequence
		command,
	})
}

func (c *Conn) writeCommandBuf(command byte, arg []byte) error {
	c.ResetSequence()

	length := len(arg) + 1
	data := utils.ByteSliceGet(length + 4)
	data[4] = command

	copy(data[5:], arg)

	err := c.WritePacket(data)

	utils.ByteSlicePut(data)

	return err
}

func (c *Conn) writeCommandStr(command byte, arg string) error {
	return c.writeCommandBuf(command, utils.StringToByteSlice(arg))
}

func (c *Conn) writeCommandUint32(command byte, arg uint32) error {
	c.ResetSequence()

	return c.WritePacket([]byte{
		0x05, // 5 bytes long
		0x00,
		0x00,
		0x00, // sequence

		command,

		byte(arg),
		byte(arg >> 8),
		byte(arg >> 16),
		byte(arg >> 24),
	})
}

func (c *Conn) writeCommandStrStr(command byte, arg1 string, arg2 string) error {
	c.ResetSequence()

	data := make([]byte, 4, 6+len(arg1)+len(arg2))

	data = append(data, command)
	data = append(data, arg1...)
	data = append(data, 0)
	data = append(data, arg2...)

	return c.WritePacket(data)
}
