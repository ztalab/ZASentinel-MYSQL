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

package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"runtime"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/config"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/mysqlproxy/client"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/mysqlproxy/mysql"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/mysqlproxy/server"
	"github.com/ztalab/ZASentinel-MYSQL/utils"
)

const bufSize = 4096

type MysqlTCPProxy struct {
	ctx        context.Context
	server     *server.Server
	credential server.CredentialProvider
	closed     atomic.Value
	mysql      *config.Mysql
}

func Start(ctx context.Context) {
	p := &MysqlTCPProxy{
		ctx: ctx,
	}
	// fake-identity
	p.credential = server.NewInMemoryProvider()
	logrus.Printf("%#v", config.Get())
	p.credential.AddUser(config.Get().FakeIdentity.Username, config.Get().FakeIdentity.Password)
	// secret
	p.mysql = config.GetRemoteMysql()

	p.server = server.NewDefaultServer()

	var err error

	p.closed.Store(true)
	ln, err := net.Listen("tcp4", config.Get().Server.Addr)
	if err != nil {
		logrus.Errorf("listening port error：%v", err)
		return
	}
	utils.GoWithRecover(func() {
		if <-ctx.Done(); true {
			p.closed.Store(true)
			_ = ln.Close()
		}
	}, nil)
	p.closed.Store(false)
	logrus.Infof("start ZASentinel-MYSQL, listen addr: %s", config.Get().Server.Addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			if p.closed.Load().(bool) || errors.Is(err, io.EOF) {
				return
			}
			logrus.Errorf("mysql proxy listening error：%v", err)
			return
		}
		go p.handle(conn)
	}
}

func (p *MysqlTCPProxy) handle(conn net.Conn) {
	defer func() {
		err := recover()
		if err != nil {
			conn.Close()
			buf := make([]byte, bufSize)
			buf = buf[:runtime.Stack(buf, false)]
			logrus.Errorf("panic err:%s", string(buf))
		}
	}()

	var remoteConn *client.Conn
	clientConn, err := server.NewCustomizedConn(conn, p.server, p.credential, func(conn *server.Conn) error {
		var err error
		remoteConn, err = client.Connect(fmt.Sprintf("%s:%d", p.mysql.Host, p.mysql.Port), p.mysql.Username, p.mysql.Password, p.mysql.DBName, func(rconn *client.Conn) {
			if conn.Charset() > 0 {
				rconn.SetCollationID(conn.Charset())
			}
			capa := conn.Capability()
			if capa&mysql.CLIENT_MULTI_RESULTS > 0 {
				rconn.SetCapability(mysql.CLIENT_MULTI_RESULTS)
			}
			if capa&mysql.CLIENT_MULTI_STATEMENTS > 0 {
				rconn.SetCapability(mysql.CLIENT_MULTI_STATEMENTS)
			}
			if capa&mysql.CLIENT_PS_MULTI_RESULTS > 0 {
				rconn.SetCapability(mysql.CLIENT_PS_MULTI_RESULTS)
			}
		})
		if err != nil {
			return fmt.Errorf("failed to connect to remote mysql: %v", err)
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("mysql connection error：%v", err)
		return
	}
	defer func() {
		remoteConn.Close()
		clientConn.Close()
	}()

	stop := make(chan struct{}, 2)
	ioCopy := func(dst, src net.Conn) {
		buf := utils.ByteSliceGet(bufSize)
		defer utils.ByteSlicePut(buf)
		_, _ = io.CopyBuffer(dst, src, buf)
		stop <- struct{}{}
	}
	go ioCopy(remoteConn.Conn.Conn, clientConn.Conn.Conn)
	go ioCopy(clientConn.Conn.Conn, remoteConn.Conn.Conn)
	select {
	case <-stop:
	case <-p.ctx.Done():
	}
}
