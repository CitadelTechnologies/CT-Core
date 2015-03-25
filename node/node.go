package node

import(
    "fmt"
    "runtime"
)

func Initialize(c chan string){

    cores := runtime.NumCPU()

    var memoryStats runtime.MemStats

    runtime.ReadMemStats(&memoryStats)

    c <- "This machine has " + fmt.Sprintf("%d", cores) + " cores and the program uses " + fmt.Sprintf("%d", memoryStats.TotalAlloc) + " bytes"

}
