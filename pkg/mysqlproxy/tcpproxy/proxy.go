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

package tcpproxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/config"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/mysqlproxy/client"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/mysqlproxy/mysql"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/mysqlproxy/server"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/secret"
	"github.com/ztalab/ZASentinel-MYSQL/utils"
	"io"
	"net"
	"runtime"
	"sync/atomic"
)

const bufSize = 4096

var (
	// number of current active connections
	activeConnCount int64
	// total number of historical connections
	totalConnCount int64
)

type MysqlTCPProxy struct {
	ctx        context.Context
	server     *server.Server
	credential server.CredentialProvider
	closed     atomic.Value
	secret     secret.SecretStore
}

func Start(ctx context.Context, conf *config.Config, store secret.SecretStore) {
	p := &MysqlTCPProxy{
		ctx:    ctx,
		secret: store,
	}
	p.server = server.NewDefaultServer()
	p.credential = server.NewVaultProvider(store)

	var err error

	p.closed.Store(true)
	ln, err := net.Listen("tcp4", conf.Server.Addr)
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
	logrus.Info("start ZASentinel-MYSQL...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			if p.closed.Load().(bool) || errors.Is(err, io.EOF) {
				return
			}
			logrus.Errorf("mysql代理监听错误：%v", err)
			return
		}
		go p.handle(conn)
	}
}

func (p *MysqlTCPProxy) handle(conn net.Conn) {
	atomic.AddInt64(&totalConnCount, 1)
	atomic.AddInt64(&activeConnCount, 1)
	defer func() {
		atomic.AddInt64(&activeConnCount, -1)
		err := recover()
		if err != nil {
			conn.Close()
			buf := make([]byte, bufSize)
			buf = buf[:runtime.Stack(buf, false)] // 获得当前goroutine的stacktrace
			logrus.Errorf("panic错误:%s", string(buf))
		}
	}()

	var remoteConn *client.Conn
	clientConn, err := server.NewCustomizedConn(conn, p.server, p.credential, func(conn *server.Conn) error {
		var addr, user, pass, dbname string
		var err error
		addr, user, pass, dbname, err = p.getAddrInfo(conn.GetUser())
		if err != nil {
			return fmt.Errorf("获取机密存储mysql连接信息失败: %v", err)
		}
		remoteConn, err = client.Connect(addr, user, pass, dbname, func(rconn *client.Conn) {
			// 转发客户端能力标识
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
			return fmt.Errorf("连接远程mysql失败: %v", err)
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("mysql连接错误：%v", err)
		return
	}
	defer func() {
		remoteConn.Close()
		clientConn.Close()
	}()

	errc := make(chan error, 2)
	ioCopy := func(dst, src net.Conn) {
		buf := utils.ByteSliceGet(bufSize)
		defer utils.ByteSlicePut(buf)
		_, err := io.CopyBuffer(dst, src, buf)
		errc <- err
	}
	go ioCopy(remoteConn.Conn.Conn, clientConn.Conn.Conn)
	go ioCopy(clientConn.Conn.Conn, remoteConn.Conn.Conn)
	select {
	case <-errc:
	case <-p.ctx.Done():
	}
}

// 获取
func (p *MysqlTCPProxy) getAddrInfo(alias string) (addr, user, pass, dbname string, err error) {
	// 连接池不存在则创建
	//var kv map[string]string
	//kv, err = p.secret.Get(fmt.Sprintf("mysql_%s", alias), false)
	//if err != nil {
	//	return
	//}
	//addr = fmt.Sprintf("%s:%s", kv[server.VaultMysqlHost], kv[server.VaultMysqlPort])
	//user = kv[server.VaultMysqlUsername]
	//pass = kv[server.VaultMysqlPassword]
	//dbname = kv[server.VaultMysqlDbname]
	return
}

func Metrics() map[string]interface{} {
	return map[string]interface{}{
		"active_conn_count": atomic.LoadInt64(&activeConnCount),
		"total_conn_count":  atomic.LoadInt64(&totalConnCount),
	}
}
