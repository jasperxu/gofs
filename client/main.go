package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func init() {
	// 修正Win下的目录问题
	file, _ := exec.LookPath(os.Args[0])
	dir, _ := path.Split(strings.Replace(file, "\\", "/", -1))
	os.Chdir(dir)
}

func main() {
	g := GoFS{URL: "http://localhost:8080/"}

	url, err := g.Upload("./微信截图_20170912143114.png", "/Jasper/1.png")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url)
	url, err = g.Upload("./微信截图_20170912143114.png", "/Jasper/2.png")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url)
	err = g.Download("/Jasper/2.png", "./2.png")
	if err != nil {
		fmt.Println(err)
	}
	err = g.Delete("/Jasper/2.png")
	if err != nil {
		fmt.Println(err)
	}
}
