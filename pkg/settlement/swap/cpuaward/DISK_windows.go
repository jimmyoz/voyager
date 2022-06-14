package cpuaward
import (
//	"fmt"
//	"github.com/shirou/gopsutil/v3/mem"
	// "github.com/shirou/gopsutil/mem"  // to use v2
	//"fmt"
	"github.com/shirou/gopsutil/v3/disk"
)
type DiskStatus struct {

	All uint64 `json:"all"`

	Used uint64 `json:"used"`

	Free uint64 `json:"free"`

}

func getDiskUsage()(Disk DiskStatus) {
	states,err:=disk.Partitions(false)
	n:=0
	total:=0
	free:=0
	if err==nil {
		n=len(states)
		for i:=0;i<n;i++ {
			//	 println(states[i].Device)
			SSD,_:=disk.Usage(states[i].Device)
			total+=int(SSD.Total)
			free+=int(SSD.Free)
		//	fmt.Printf("%s Total: %fG, Free:%fG, UsedPercent:%f%%\n", states[i].Device,float32(SSD.Total)/1024.0/1024.0/1024.0, float32(SSD.Free)/1024/1024/1024, SSD.UsedPercent)
		//	println("total=",total," ","free=",free)
		}
	}
    Disk.All=uint64(total)
	Disk.Used=uint64(total-free)
	Disk.Free=uint64(free)
	return
}
/*func main() {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	println("mem info")
	fmt.Printf("Total: %vB, Free:%vB, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	fmt.Printf("Total: %vK, Free:%vK, UsedPercent:%f%%\n", v.Total/1024, v.Free/1024, v.UsedPercent)
	fmt.Printf("Total: %vM, Free:%vM, UsedPercent:%f%%\n", v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)
	fmt.Printf("Total: %vG, Free:%vG, UsedPercent:%f%%\n", float64(v.Total)/1024/1024/1024, float64(v.Free)/1024/1024/1024, v.UsedPercent)
	println("disk info")


	states,err:=disk.Partitions(false)
	n:=0
	total:=0
	free:=0
	if err==nil {
		n=len(states)
		for i:=0;i<n;i++ {
			//	 println(states[i].Device)
			SSD,_:=disk.Usage(states[i].Device)
			total+=int(SSD.Total)
			free+=int(SSD.Free)
			fmt.Printf("%s Total: %fG, Free:%fG, UsedPercent:%f%%\n", states[i].Device,float32(SSD.Total)/1024.0/1024.0/1024.0, float32(SSD.Free)/1024/1024/1024, SSD.UsedPercent)
			println("total=",total," ","free=",free)
		}
	}


	// convert to JSON. String() is also implemented
	//fmt.Println(v)
}*/
