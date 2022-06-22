package cpuaward

import (
	"syscall"
	"unsafe"
)

var (
	DedicateSystemMemory uint32
	DedicateVideoMemory  uint32
	SharedSystemMemory   uint32
	//deviceDesc           string
	//deviceID uint32
	//Revision uint32
	//SubsysId uint32
	//VendorId uint32
)

func getGPUSize()uint32{
	//println(C.add(1, 2))
	//str := make([]string, syscall.MAX_PATH)
	b := make([]uint32, syscall.MAX_PATH)
	lib_getGPU_info(&b[0], &b[1], &b[2], &b[3], &b[4],
		&b[5], &b[6])
	DedicateSystemMemory=b[0]
	DedicateVideoMemory=b[1]
	SharedSystemMemory=b[2]
	//deviceDesc           string
//	deviceID=b[3]
//	Revision=b[4]
//	SubsysId=b[5]
//	VendorId=b[6]
	//println(b[0] + b[1] + b[2])
	return DedicateSystemMemory+DedicateVideoMemory+SharedSystemMemory

}

func IntPtr(n *uint32) uintptr {
	return uintptr(unsafe.Pointer(n))
}
func StrPtr(s string) uintptr {
	pointer, _ := syscall.BytePtrFromString(s)
	return uintptr(unsafe.Pointer(pointer))
}
func lib_getGPU_info(DedicateSystemMemory *uint32, DedicateVideoMemory *uint32, SharedSystemMemory *uint32, deviceID *uint32, Revision *uint32,
	SubsysId *uint32, VendorId *uint32) bool {
	lib := syscall.NewLazyDLL("GPU_info_windows.dll")
	//fmt.Println("dll:", lib.Name)
	getGPU_info := lib.NewProc("getGPU_info")
	//fmt.Println("+++++++NewProc:", getGPU_info, "+++++++")
	_, _, err := getGPU_info.Call(IntPtr(DedicateSystemMemory), IntPtr(DedicateVideoMemory), IntPtr(SharedSystemMemory), IntPtr(deviceID), IntPtr(Revision), IntPtr(SubsysId), IntPtr(VendorId))
	if err != nil {
		return false
	}
	return true
}
