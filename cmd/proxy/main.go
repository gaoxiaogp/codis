// Copyright 2014 Wandoujia Inc. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package main

import (
	"bufio"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/diditaxi/codis/pkg/proxy/router"
	"github.com/diditaxi/codis/pkg/utils"

	"github.com/docopt/docopt-go"
	log "github.com/ngaut/logging"
)

var (
	cpus          = 2
	addr          = ":9000"
	httpAddr      = ":9001"
	configFile    = "config.ini"
	whitelistFile = ""
)

var usage = `usage: proxy [-c <config_file>] [-w <whitelist_file>] [-L <log_file>] [--log-level=<loglevel>] [--cpu=<cpu_num>] [--addr=<proxy_listen_addr>] [--http-addr=<debug_http_server_addr>]

options:
   -c	set config file
   -w	set ip whitelist file
   -L	set output log file, default is stdout
   --log-level=<loglevel>	set log level: info, warn, error, debug [default: info]
   --cpu=<cpu_num>		num of cpu cores that proxy can use
   --addr=<proxy_listen_addr>		proxy listen address, example: 0.0.0.0:9000
   --http-addr=<debug_http_server_addr>		debug vars http server
`

var banner string = `
  _____  ____    ____/ /  (_)  _____
 / ___/ / __ \  / __  /  / /  / ___/
/ /__  / /_/ / / /_/ /  / /  (__  )
\___/  \____/  \__,_/  /_/  /____/

`

func handleSetLogLevel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	level := r.Form.Get("level")
	log.SetLevelByString(level)
	log.Info("set log level to", level)
}
func handleChangeLog(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	log.SetOutputByName(name)
	log.Info("set log name to", name)
}

func readWhiteList(fPath string) map[string]string {
	whitelist := make(map[string]string)

	f, err := os.Open(fPath)
	defer f.Close()

	if err == nil {
		buff := bufio.NewReader(f)
		for {
			line, err := buff.ReadString('\n')
			if err != nil {
				break
			}
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			line = line[0 : len(line)-1]
			whitelist[line] = ""
		}
		return whitelist
	}

	return nil
}

func main() {
	fmt.Print(banner)
	log.SetLevelByString("info")

	args, err := docopt.Parse(usage, nil, true, "codis proxy v0.1", true)
	if err != nil {
		log.Error(err)
	}

	// set config file
	if args["-c"] != nil {
		configFile = args["-c"].(string)
	}

	// set whitelist file
	if args["-w"] != nil {
		whitelistFile = args["-w"].(string)
	}

	// set output log file
	if args["-L"] != nil {
		log.SetOutputByName(args["-L"].(string))
	}

	// set log level
	if args["--log-level"] != nil {
		log.SetLevelByString(args["--log-level"].(string))
	}

	// set cpu
	if args["--cpu"] != nil {
		cpus, err = strconv.Atoi(args["--cpu"].(string))
		if err != nil {
			log.Fatal(err)
		}
	}

	// set addr
	if args["--addr"] != nil {
		addr = args["--addr"].(string)
	}

	// set http addr
	if args["--http-addr"] != nil {
		httpAddr = args["--http-addr"].(string)
	}

	dumppath := utils.GetExecutorPath()

	log.Info("dump file path:", dumppath)
	log.CrashLog(path.Join(dumppath, "codis-proxy.dump"))

	router.CheckUlimit(1024)
	runtime.GOMAXPROCS(cpus)

	http.HandleFunc("/setloglevel", handleSetLogLevel)
	http.HandleFunc("/changelogname", handleChangeLog)
	go http.ListenAndServe(httpAddr, nil)
	log.Info("running on ", addr)
	conf, err := router.LoadConf(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var ipwhielist map[string]string
	if whitelistFile != "" {
		ipwhielist = readWhiteList(whitelistFile)
	}

	s := router.NewServer(addr, httpAddr, conf, ipwhielist)
	s.Run()
	log.Warning("exit")
}
