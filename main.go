// This file is part of BravoSystem.
// Copyright(c) 2016 by chnvideo.com
// All rights reserved.
//
// Author
//      - fanbingqiang <fanbingqiang@bravovcloud.com>

package main

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"bufio"
    "flag"
	ocore "github.com/ossrs/go-oryx-lib/logger"
    "time"
    "runtime"
)

type mp4Handler struct {}

// the interval to check the parent pid.
const checkParentInterval = time.Second * 1

func (*mp4Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()
	dir, _ := os.Getwd()                // get the current dir
	path = dir + "/objs/nginx/html" + path                   // get the full path of the file
	if _, err := os.Stat(path); err != nil {
        io.WriteString(w, path)
		io.WriteString(w, "\n404 NOT Found!\n")
	}
	// io.WriteString(w, path)
	
	buf := make([]byte, 1024)
	f, _ := os.Open(path)
	defer f.Close()
	
	req := bufio.NewReader(f)
	for {
		n, _ := req.Read(buf)
		if n == 0 {
			break
		}
		w.Write(buf[0:n])
	}
}

func watchParent() {
	ppid := os.Getppid()

	go func() {
		for {
			if v := os.Getppid(); v == 1 || v != ppid {
				ocore.Error.Println(nil, "quit for parent problem, ppid is", v)
				os.Exit(0)
			}
			//ocore.Trace.Println(nil, "parent pid", ppid, "ok")

			time.Sleep(checkParentInterval)
		}
	}()
	ocore.Trace.Println(nil, "parent process watching, ppid is", ppid)
}

func main() {
//	fmt.Println(fmt.Sprintf("BRAVO VOD/%v is Bravo Vod System", version))
//	fmt.Println(fmt.Sprintf("%v, Copyright(c) 2016 by chnvideo.com", server))
	
    var listen string
	var nbCpus int
    flag.StringVar(&listen, "listen", ":8090", "the bravo http mp4 server listen at")
	flag.IntVar(&nbCpus, "cpu", 1, "the cpus to use")
    
	flag.Usage = func() {
		fmt.Println(fmt.Sprintf("Usage: %v [--listen=string] [--cpu=int] [-h|--help]", os.Args[0]))
		fmt.Println(fmt.Sprintf("	listen, the listen [host]:port. default :8090"))
		fmt.Println(fmt.Sprintf("	cpus, the cpu to use. default 1"))
		fmt.Println(fmt.Sprintf("	help, show this help and exit"))
		fmt.Println(fmt.Sprintf("For example:"))
		fmt.Println(fmt.Sprintf("	%v --listen=:8090", os.Args[0]))
	}
	flag.Parse()
    
    watchParent()
	
	runtime.GOMAXPROCS(nbCpus)
	ocore.Trace.Println(nil, fmt.Sprintf("Bravo vod pid:%v, ppid:%v, cpus:%v", os.Getpid(), os.Getppid(), nbCpus))
	
	err := http.ListenAndServe(listen, &mp4Handler{})
	if err != nil {
		ocore.Error.Println(nil, "VOD listen at", listen, "failed, err is", err)
		return
	}
	
	return
}
