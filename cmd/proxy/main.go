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

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/ztalab/ZASentinel-MYSQL/pkg/config"
	"github.com/ztalab/ZASentinel-MYSQL/proxy"
)

func main() {
	app := cli.NewApp()

	app.Name = "ZASentinel-MYSQL"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "config, c",
			Usage:    "config file path",
			Value:    "config/config.yaml",
			Required: false,
		},
		cli.StringFlag{
			Name:  "log,l",
			Usage: "log level: debug,info,warning,error",
			Value: "debug",
		},
	}

	app.Before = func(c *cli.Context) error {
		// init log
		lv, err := logrus.ParseLevel(c.String("log"))
		if err != nil {
			return err
		}
		logrus.SetLevel(lv)
		// init config
		confPath := c.String("config")
		config.InitConfig(confPath)

		return nil
	}

	app.Action = func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.TODO())
		stop := make(chan struct{})
		go func() {
			proxy.Start(ctx)
			stop <- struct{}{}
		}()

		return exitSignal(cancel, stop)
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func exitSignal(cancel context.CancelFunc, stop chan struct{}) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-sigs:
			cancel()
			select {
			case <-stop:
				fmt.Println("shutdown！！！！")
			case <-time.After(time.Second * 5):
				fmt.Println("timeout forced exit！！！！")
			}
			os.Exit(0)
		case <-stop:
			cancel()
			fmt.Println("shutdown！！！！")
			os.Exit(0)
		}
	}
	return nil
}
