package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	path "path/filepath"
	"regexp"
)

var cmdVersion = &Command{
	UsageLine: "version",
	Short:     "show the Uni, Unigo and Go version",
	Long: `
show the Uni, Unigo and Go version

bee version
    Uni   :0.0.1
    Unigo :0.0.1
    Go    :go version go1.3.3 linux/amd64

`,
}

func init() {
	cmdVersion.Run = versionCmd
}

func versionCmd(cmd *Command, args []string) int {
	fmt.Println("uni   :" + version)
	fmt.Println("unigo :" + getunigoVersion())
	//fmt.Println("Go    :" + runtime.Version())
	goversion, err := exec.Command("go", "version").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Go    :" + string(goversion))
	return 0
}

func getunigoVersion() string {
	gopath := os.Getenv("GOPATH")
	re, err := regexp.Compile(`const VERSION = "([0-9.]+)"`)
	if err != nil {
		return ""
	}
	if gopath == "" {
		err = fmt.Errorf("you should set GOPATH in the env")
		return ""
	}
	wgopath := path.SplitList(gopath)
	for _, wg := range wgopath {
		wg, _ = path.EvalSymlinks(path.Join(wg, "src", "unigo"))
		filename := path.Join(wg, "unigo.go")
		_, err := os.Stat(filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			ColorLog("[ERRO] get unigo.go has error\n")
		}
		fd, err := os.Open(filename)
		if err != nil {
			ColorLog("[ERRO] open unigo.go has error\n")
			continue
		}
		reader := bufio.NewReader(fd)
		for {
			byteLine, _, er := reader.ReadLine()
			if er != nil && er != io.EOF {
				return ""
			}
			if er == io.EOF {
				break
			}
			line := string(byteLine)
			s := re.FindStringSubmatch(line)
			if len(s) >= 2 {
				return s[1]
			}
		}

	}
	return "you don't install unigo"
}
