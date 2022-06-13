package cpuaward
import (
	"github.com/StackExchange/wmi"
)


type gpuInfo struct {
	AdapterRam uint32
	Name string
}


func getGPUSize()(uint32) {
	var gpuinfo []gpuInfo
	err := wmi.Query("Select * from Win32_VideoController", &gpuinfo)
	if err != nil {
		return 0
	}
	size:=gpuinfo[0].AdapterRam
	size=size/1024
	//fmt.Println("GPU:=", size)
	return size
}