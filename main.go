package main

import(
	"gleipnir/kernel"
)

func main() {
	var core kernel.Kernel
	core.Init()
	defer core.Shutdown()
	core.Run()
}