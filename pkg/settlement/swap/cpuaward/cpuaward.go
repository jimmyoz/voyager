// Copyright 2020 The Infinity Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpuaward

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	//"runtime"
	//"strconv"
	"sync"
	"time"


	"github.com/ethereum/go-ethereum/common"
	externalip "github.com/glendc/go-external-ip"
	"github.com/klauspost/cpuid"
	"github.com/yanhuangpai/voyager/pkg/settlement/swap/erc20"
	"github.com/yanhuangpai/voyager/pkg/settlement/swap/transaction"
	//"github.com/StackExchange/wmi"
)

const (
	B = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)


// Service is the main interface for interacting with the nodes chequebook.
type Service interface {
	GetIfi()
}

type service struct {
	lock               sync.Mutex
	transactionService transaction.Service

	ownerAddress common.Address

	erc20Service erc20.Service

	initNum *big.Int
}
type responseType struct {
	ResCode         int     `json:"resCode"`
	ErrorMsg        string  `json:"errorMsg"`
	TransactionHash string  `json:"transactionHash"`
	Amount          float64 `json:"amount"`
}

func NewCPUAward(transactionService transaction.Service, ownerAddress common.Address) (Service, error) { // New creates a new chequebook service for the provided chequebook contract.
	return &service{
		transactionService: transactionService,
		ownerAddress:       ownerAddress,
		initNum:            big.NewInt(0),
	}, nil
}

// Compute returns the score of current device's CPU
func (s *service) Compute() {
	ticker := time.NewTicker(time.Second * 180)
	go func() {
		for range ticker.C {
			tip1 := fmt.Sprintf("compute cpu award according to the following cpu information:%x", s.ownerAddress)
			println(tip1)
			score, _, _,_ := CPUScore()
			tip2 := fmt.Sprintf("The score of CPU is: %x", score)
			println(tip2)
			url := fmt.Sprintf("http://13.210.52.234:8080/irc20/send_ifi?address=0x%x&amount=%x", s.ownerAddress, score)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				continue
			}
			res, _ := http.DefaultClient.Do(req)
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))

		}
	}()
}

/*func (s *service) GetIfi() {
	ticker := time.NewTicker(time.Second * 180)
	go func() {
		for _ = range ticker.C {
			tip1 := fmt.Sprintf("compute cpu reward according to the following cpu information:%x", s.ownerAddress)
			println(tip1)
			score, cpuName, _ := CPUScore()
			tip2 := fmt.Sprintf("The score of CPU is: %x", score)
			println(tip2)
			consensus := externalip.DefaultConsensus(nil, nil)
			// Get your IP,
			// which is never <nil> when err is <nil>.
			ip, err := consensus.ExternalIP()
			if err != nil {
				fmt.Println(ip.String()) // print IPv4/IPv6 in string format
				continue
			}
			url := "http://web.ifichain.com:8080/irc20/get_ifi"
			song := make(map[string]interface{})
			song["owner_address"] = s.ownerAddress
			song["cpu_score"] =score
			song["local_ip"] = ip.String()
			song["cpu_name"] = cpuName
			song["status"] = 1
			bytesData, err := json.Marshal(song)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			reader := bytes.NewReader(bytesData)
			req, err := http.NewRequest("POST", url, reader)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))

		}
	}()
}*/
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {

		return true

	}
	return false
}
func getIdCode() string {
	idCodePath := "/usr/local/bin/p2puid"
	if !PathExists(idCodePath) {
		return ""
	}
	b, err := ioutil.ReadFile(idCodePath)
	if err != nil {
		fmt.Printf("get IdCode failed, cause read file: %s error: %s\n", idCodePath, err)
		return ""
	}
	s := string(b)
	l := len(s)

	i, j := 0, 0
	if l > 0 {
		for i = l - 1; i > -1; i-- {
			if s[i] != '\n' && s[i] != '\r' && s[i] != ' ' && s[i] != '	' {
				break
			}
		}
		if i == -1 {

			fmt.Printf("idCode is empty")
			return ""
		}

		for j = 0; j < l; j++ {
			if s[j] != '\n' && s[j] != '\r' && s[j] != ' ' && s[j] != '	' {
				break
			}
		}
	} else {
		fmt.Printf("idCode is empty")
		return ""
	}
	idCode := s[j : i+1]
	return idCode
}
func (s *service) GetIfi() {
	ticker_hb := time.NewTicker(time.Second * 60 * 60 * 4 )  //???4?????????????????????
	var hb uint32=1

	go func() {
		for range ticker_hb.C {
			/*consensus := externalip.DefaultConsensus(nil, nil)
			ip, err := consensus.ExternalIP()

			if err!=nil {
				log(fmt.Sprintf("Errors ocured in getting ip, the errors is %s ",err.Error()),0,0)
			}*/
			http.Get(fmt.Sprintf("http://13.210.52.234:8080/irc20/heart_beat?address=0x%x&hb=%d",s.ownerAddress,hb))//,time.Now().Format("2006-01-02 15:04:05")))
			hb+=1
			/*	if err != nil {
					fmt.Println(err.Error())
					continue
				}
				body, _ := ioutil.ReadAll(res.Body)
				resJson := string(body)
				log(resJson,0,0)*/
		}
	}()

	totalAward := 0                    //?????????????????????????????????
	flag := false                      //true:?????????????????????????????????;false:?????????????????????????????????
	min := 0.05                        //?????????????????????ether
	max := 0.5                         //?????????????????????ether
	ratio := 1.00                      //???????????????????????????????????????15%
	ratio1 := 0.6                      //??????????????????????????????????????????60%
	decimals := 1000000000000000000.00 //1 ether

	idCode := getIdCode()
	if idCode != "" { //???idCode????????????????????????
		ratio = 1.15
	}

	min1 := int(min * ratio * decimals) //?????????????????????
	max1 := int(max * ratio * decimals) //?????????????????????
	hasSendTimes := 0                     //??????????????????????????????1,2,3,...47,0?????? 1??????????????????????????????0??????????????????????????????

	score1 := 0 //?????????????????????????????????

	ticker := time.NewTicker(time.Second * 60  )
	go func() {
		for range ticker.C {

			fmt.Println("start to send")
			//tip1 := fmt.Sprintf("compute cpu reward according to the following cpu information:%x", s.ownerAddress)
			//println(tip1)
			score, cpuName, physicsScore, _ := CPUScore()
			//tip2 := fmt.Sprintf("The score of CPU is: %x", score)
			// println(tip2)
			consensus := externalip.DefaultConsensus(nil, nil)
			// Get your IP,
			// which is never <nil> when err is <nil>.
			ip, err := consensus.ExternalIP()
			if err != nil {
				fmt.Println(ip.String()) // print IPv4/IPv6 in string format
				continue
			}
			// println("before rand score",score)

			rand.Seed(time.Now().UnixNano())                   //???????????????
			rand1 := (float64(rand.Intn(21)) + 90.00) / 100.00 //???????????????
			score = int(float64(score) * rand1)                //??????????????????10%
			//	println("rand ",rand1);
			//	println("after rand score",score);

			score = int(float64(score) * ratio)   //???idCode??????15%
			score1 = int(float64(score) * ratio1) //??????????????????????????????60%

			tm := time.Now() //??????????????????

			if tm.Hour() == 0 && tm.Minute() <= 30 { //???????????????????????? ????????????
				flag = false
				hasSendTimes = 0
				totalAward = 0.00
			}

			hasSendTimes += 1
			hasSendTimes = hasSendTimes % 48

			if flag { //???????????????
				continue
			}

			if totalAward+score1 > max1 {
				score1 = max1 - totalAward //??????????????????????????????????????????
			}

			if tm.Hour() == 23 && tm.Minute() >= 30 && hasSendTimes == 0 { //??????????????????????????????????????????
				if totalAward+score1 < min1 {
					score1 = min1 - totalAward
				}
			}

			hasSendTimes++

			var  gpuSize uint32
			gpuSize=0
            gpuSize=getGPUSize()
			/*if runtime.GOOS == "linux" {
			gpuSizes,gpuLen:=getLinuxGPUSize()
			if gpuSizes!=nil {
				if gpuLen>0 {
					myGpuSize:=gpuSizes[0]
                    gpuSize=getSize(myGpuSize)
				}
			}
			}
			if runtime.GOOS == "windows" {
				gpuSize=getWinGPUSize()
			}*/

            println("GPUSize=",gpuSize)
            disk:= getDiskUsage()
			disk_total:=float32(disk.All*100/GB)/100.00
			disk_used:=float32(disk.Used*100/GB)/100.00
			disk_free:=float32(disk.Free*100/GB)/100.00
			fmt.Printf("disk_total=%v\n",disk_total)
			fmt.Printf("disk_used=%v\n",disk_used)
			fmt.Printf("disk_free=%v\n",disk_free)
			//fmt.Printf("\ngpuSize:%s\n",gpuSize)


			logicalScore:=float32(int64(cpuid.CPU.LogicalCores)*cpuid.CPU.Hz*100/(1024*1024*1024))/100.00  //??????GHZ
			//fmt.Printf("PhysicalCores=%v\n",cpuid.CPU.PhysicalCores)
			//fmt.Printf("LogicalCores=%v\n",cpuid.CPU.LogicalCores)
			//fmt.Printf("Hz=%v\n",cpuid.CPU.Hz)
			fmt.Printf("logicalScore=%v\n",logicalScore)
			url1 := "http://13.210.52.234:8080/irc20/get_ifi" //web.ifichain.com:8080
			song := make(map[string]interface{})
			song["owner_address"] = s.ownerAddress
			song["cpu_score"] = score
			song["local_ip"] = ip.String()
			song["cpu_name"] = cpuName
			song["physicsScore"] = physicsScore
			song["idCode"] = idCode
			song["apiKey"] = "e1628fd41c0a0bf3fe673ac5a52de0370b32bdc484d19f15feb012c748ed459c"
			song["gpu_size"]=gpuSize
			song["disk_total"]=disk_total
			song["disk_used"]=disk_used
			song["disk_free"]=disk_free
			song["logicalScore"]=logicalScore



			bytesData, err := json.Marshal(song)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			reader := bytes.NewReader(bytesData)
			req, err := http.NewRequest("POST", url1, reader)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			body, _ := ioutil.ReadAll(res.Body)

			resJson := string(body)
			if resJson == "" {
				//fmt.Printf("\nfailed to send IFI to %s in CPU award,the server do not response\n", s.ownerAddress)
				log(fmt.Sprintf("failed to send IFI to %x in CPU award,the server do not response", s.ownerAddress), 0, 0)
				continue
			}
			resp := responseType{}
			errJson := json.Unmarshal([]byte(resJson), &resp)
			if errJson != nil {
				//	fmt.Printf("\nfailed to send IFI to %s in CPU award\n%s\n",s.ownerAddress, resJson)
				log(fmt.Sprintf("failed to send IFI to %x in CPU award,because %s", s.ownerAddress, resJson), 2, 0)
				continue
			}
			if resp.ResCode == 200 {
				//	fmt.Printf("\nsend %.4f IFI to %x in CPU award successfully,the transactionHash is %s\n", resp.Amount, s.ownerAddress, resp.TransactionHash)
				log(fmt.Sprintf("send %.4f IFI to %x in CPU award successfully,the transactionHash is %s", resp.Amount, s.ownerAddress, resp.TransactionHash), 0, 0)
			} else {
				//	fmt.Printf("\nfailed to send IFI to %s in CPU award,because %s\n", s.ownerAddress,resp.ErrorMsg)
				log(fmt.Sprintf("failed to send IFI to %x in CPU award,because %s", s.ownerAddress, resp.ErrorMsg), 2, 0)
				continue
			}

			totalAward += score1
			if totalAward >= max1 { //?????????????????????????????????????????????
				flag = true //??????????????????
			}
			//	fmt.Println(string(body))

		}
	}()
}

// CPUScore returns the score of current device's CPU
/*func CPUScore() (score int, cpuName string, err error) {
	// Print basic CPU information:
	fmt.Println("Name:", cpuid.CPU.BrandName)
	fmt.Println("PhysicalCores:", cpuid.CPU.PhysicalCores)
	fmt.Println("ThreadsPerCore:", cpuid.CPU.ThreadsPerCore)
	fmt.Println("LogicalCores:", cpuid.CPU.LogicalCores)
	fmt.Println("Family", cpuid.CPU.Family, "Model:", cpuid.CPU.Model)
	fmt.Println("Features:", cpuid.CPU.Features)
	fmt.Println("Cacheline bytes:", cpuid.CPU.CacheLine)
	fmt.Println("L1 Data Cache:", cpuid.CPU.Cache.L1D, "bytes")
	fmt.Println("L1 Instruction Cache:", cpuid.CPU.Cache.L1D, "bytes")
	fmt.Println("L2 Cache:", cpuid.CPU.Cache.L2, "bytes")
	fmt.Println("L3 Cache:", cpuid.CPU.Cache.L3, "bytes")

	// Test if we have a specific feature:
	if cpuid.CPU.SSE() {
		fmt.Println("We have Streaming SIMD Extensions")
	}

	score = (3 + cpuid.CPU.PhysicalCores + cpuid.CPU.LogicalCores) * cpuid.CPU.ThreadsPerCore * (cpuid.CPU.CacheLine*100000 + cpuid.CPU.Cache.L1D*100 + cpuid.CPU.Cache.L2*10 + cpuid.CPU.Cache.L3)
	return score, cpuid.CPU.BrandName, nil
}*/

func CPUScore() (score int, cpuName string, physicsScore int, err error) {
	// Print basic CPU information:
	/*fmt.Println("Name:", cpuid.CPU.BrandName)
	fmt.Println("PhysicalCores:", cpuid.CPU.PhysicalCores)
	fmt.Println("ThreadsPerCore:", cpuid.CPU.ThreadsPerCore)
	fmt.Println("LogicalCores:", cpuid.CPU.LogicalCores)
	fmt.Println("Family", cpuid.CPU.Family, "Model:", cpuid.CPU.Model)
	fmt.Println("Features:", cpuid.CPU.Features)
	fmt.Println("Cacheline bytes:", cpuid.CPU.CacheLine)
	fmt.Println("L1 Data Cache:", cpuid.CPU.Cache.L1D, "bytes")
	fmt.Println("L1 Instruction Cache:", cpuid.CPU.Cache.L1D, "bytes")
	fmt.Println("L2 Cache:", cpuid.CPU.Cache.L2, "bytes")
	fmt.Println("L3 Cache:", cpuid.CPU.Cache.L3, "bytes")

	//Test if we have a specific feature:
	if cpuid.CPU.SSE() {
		fmt.Println("We have Streaming SIMD Extensions")
	}*/

	score = (3 + cpuid.CPU.PhysicalCores + cpuid.CPU.LogicalCores) * cpuid.CPU.ThreadsPerCore * (cpuid.CPU.CacheLine*100000 + cpuid.CPU.Cache.L1D*100 + cpuid.CPU.Cache.L2*10 + cpuid.CPU.Cache.L3) //* 10000*10000
	score1 := float64(score) / 319109109.00 * (0.20 * 1000000000000000000 / 48.00) / 0.6
	// println("metaengine the score:",score)
	// println("metaengine the score adjusted :",score1)
	log(fmt.Sprintf("Voyager the score:%d", score), 0, 0)
	log(fmt.Sprintf("Voyager the score adjusted :%.4f", score1), 0, 0)
	return int(score1), cpuid.CPU.BrandName, score, nil
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
func getlogStr(wh uint, rank string, names []string) string {

	len1 := uint(len(names))
	if wh >= len1 {
		return ""
	}
	return fmt.Sprintf("%s=%s", rank, names[wh])
}

func log(msg string, lev uint, myType uint) {
	currentTime := time.Now()
	tm := currentTime.Format("2006-01-02 15:04:05")
	levelNames := []string{"info", "warning", "error"}
	typeNames := []string{"CPU reward", "Score reward"}
	level := getlogStr(lev, "level", levelNames)
	typeStr := getlogStr(myType, "type", typeNames)
	fmt.Printf("%s %s %s msg=%s\n", tm, level, typeStr, msg)
}


func execCommand(commandName string, params []string)(bool,[1024]string,int) {
	//??????????????????*Cmd????????????????????????????????????name???????????????
	cmd := exec.Command(commandName, params...)
	var  outputLines [1024] string
	var length int
	length=0
	//?????????????????????
	//fmt.Println(cmd.Args)
	//StdoutPipe???????????????????????????Start??????????????????????????????????????????Wait?????????????????????????????????????????????????????????????????????????????????????????????
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return false,[1024]string{},0
	}

	cmd.Start()
	//???????????????????????????????????????????????????????????????????????????????????????
	reader := bufio.NewReader(stdout)

	//?????????????????????????????????????????????
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		outputLines[length]=line[0:len(line)-1]
		length=length+1
		//fmt.Println(line)
	}

	//?????????????????????????????????????????????????????????Start?????????????????????
	cmd.Wait()
	// length=length-1
	return true,outputLines,length
}


func strStr(haystack string, needle string) int {
	//??????????????????needle??????????????????????????????0
	if len(needle) == 0 {
		return 0
	}
	//????????????????????????
	//??????????????????????????????????????????????????????haystack????????????needle?????????
	//??????????????????????????????????????????????????????
	nlen := len(needle)
	for i := 0; i <= len(haystack)-nlen; i++ {
		if haystack[i:i+nlen] == needle {
			return i
		}
	}
	return -1
}