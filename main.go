package main

import(
	"os"
	"runtime"
	"gleipnir/kernel"
)

var PathSeparator string
var Gopath string

func main() {

	definePaths()

	var core kernel.Kernel
	core.Init()
	core.Run()
}

func definePaths() {

	if runtime.GOOS == "windows" {
		PathSeparator = "\\"
	} else {
		PathSeparator = "/"
	}

	Gopath = os.Getenv("GOPATH")

}