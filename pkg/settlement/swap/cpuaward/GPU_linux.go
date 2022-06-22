package cpuaward

import (
	"fmt"
	"strconv"
	//	"os/exec"
)

func getGPUSize() uint32 {
	var gpuSize uint32
	gpuSize=0
		gpuSizes,gpuLen:=getGPULinuxSize()
		if gpuSizes!=nil {
			if gpuLen>0 {
				myGpuSize:=gpuSizes[0]
				gpuSize=getSize(myGpuSize)/1024
			}
		}
    return gpuSize
}
func getGPULinuxSize()([]string,int) {
	// getGPUInfo()
	gpuNum:=""
	numIndex:=0
	tempLen:=0
	size:=""
	var gpuSizes [] string=make([]string,1024)
	var gpuLen int=0
	status,outputLines,length:=execCommand("/bin/bash",[] string{"-c","lspci | grep -i vga"})
	//println("length=",length)
	if status {
		/*for i:=0;i<length;i++ {
		     println(outputLines[i])
		}*/
		if length>0 {
			tempLen=len(outputLines[0])
			for j:=0;j<tempLen;j++ {
				if outputLines[0][j]==' ' {
					numIndex=j;
					break
				}
			}
			gpuNum=outputLines[0][0:numIndex]
			println("gpuNum=",gpuNum)
			status1,outputLines1,length1:=execCommand("/bin/bash",[] string{"-c",fmt.Sprintf("lspci -v -s %s",gpuNum)})
			if status1  {
				for i:=0;i<length1;i++ {
					// println(outputLines1[i])

					startIndex:=strStr(outputLines1[i],"[size=")
					if startIndex<0 {
						continue
					}
					tempStr:=outputLines1[i][startIndex:]
					startIndex1:=strStr(tempStr,"=")
					endIndex1:=strStr(tempStr,"]")
					size=tempStr[startIndex1+1:endIndex1]
					//println(size)
					gpuSizes[gpuLen]=size
					gpuLen=gpuLen+1
				}
			}
		}
	}

	if len(gpuSizes)>0 {
		return gpuSizes,gpuLen
	}
	return nil,gpuLen
}
func getSize(strSize string) uint32 {
	if strSize != "" {
		len := len(strSize)
		unit := strSize[len-1]
		size, err := strconv.Atoi(strSize[0 : len-1])
		if err != nil {
			return 0
		}
		var times uint32 = 1

		switch unit {
		case 'G':
			times = 1024 * 1024
			break
		case 'M':
			times = 1024
			break
		case 'K':
			times = 1
			break
		default:
			times = 1
			break
		}
		return uint32(size) * times

	}
	return 0
}