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

package mysql

import (
	"fmt"

	"github.com/pingcap/errors"
)

var (
	ErrBadConn       = errors.New("connection was bad")
	ErrMalformPacket = errors.New("Malform packet error")

	ErrTxDone = errors.New("sql: Transaction has already been committed or rolled back")
)

type MyError struct {
	Code    uint16
	Message string
	State   string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("ERROR %d (%s): %s", e.Code, e.State, e.Message)
}

// NewDefaultError: default mysql error, must adapt errname message format
func NewDefaultError(errCode uint16, args ...interface{}) *MyError {
	e := new(MyError)
	e.Code = errCode

	if s, ok := MySQLState[errCode]; ok {
		e.State = s
	} else {
		e.State = DEFAULT_MYSQL_STATE
	}

	if format, ok := MySQLErrName[errCode]; ok {
		e.Message = fmt.Sprintf(format, args...)
	} else {
		e.Message = fmt.Sprint(args...)
	}

	return e
}

func NewError(errCode uint16, message string) *MyError {
	e := new(MyError)
	e.Code = errCode

	if s, ok := MySQLState[errCode]; ok {
		e.State = s
	} else {
		e.State = DEFAULT_MYSQL_STATE
	}

	e.Message = message

	return e
}

func ErrorCode(errMsg string) (code int) {
	var tmpStr string
	// golang scanf doesn't support %*,so I used a temporary variable
	fmt.Sscanf(errMsg, "%s%d", &tmpStr, &code)
	return
}
