// Copyright 2020 The Infinity Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpuaward

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	externalip "github.com/glendc/go-external-ip"
	"github.com/klauspost/cpuid"
	"github.com/yanhuangpai/voyager/pkg/settlement/swap/erc20"
	"github.com/yanhuangpai/voyager/pkg/settlement/swap/transaction"
)

const ()

var ()

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
			url := fmt.Sprintf("http://web.ifichain.com:8080/irc20/send_ifi?address=0x%x&amount=%x", s.ownerAddress, score)
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

	totalAward := 0                    //每天给挖矿者总的激励数
	flag := false                      //true:当天发送量已达到最大值;false:当天发送量未达到最大值
	min := 0.05                        //最小激励，单位ether
	max := 0.5                         //最大激励，单位ether
	ratio := 1.00                      //用于控制客户自家机器激励多15%
	ratio1 := 0.6                      //用于控制挖矿者激励为总激励的60%
	decimals := 1000000000000000000.00 //1 ether

	idCode := getIdCode()
	if idCode != "" { //有idCode，为客户自家机器
		ratio = 1.15
	}

	min1 := (int)(min * ratio * decimals) //扩大最小值范围
	max1 := (int)(max * ratio * decimals) //扩大最大值范围
	hasSendTimes := 0                     //发送次数，取值区间【1,2,3,...47,0】， 1表示当天第一次发送，0表示当天最后一次发送

	score1 := 0 //当天应发给挖矿者的激励

	ticker := time.NewTicker(time.Second * 60 * 30)
	go func() {
		for range ticker.C {

			fmt.Println("start to send");
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

			rand.Seed(time.Now().UnixNano())                   //随机数种子
			rand1 := (float64(rand.Intn(21)) + 90.00) / 100.00 //生成随机数
			score = int(float64(score) * rand1)                //随机上下浮动10%
			//	println("rand ",rand1);
			//	println("after rand score",score);

			score = int(float64(score) * ratio)   //有idCode增加15%
			score1 = int(float64(score) * ratio1) //挖矿者激励为总激励的60%

			tm := time.Now() //获取当前时间

			if tm.Hour() == 0 && tm.Minute() <= 30 { //第二天第一次发送 ，初始化
				flag = false
				hasSendTimes = 0
				totalAward = 0.00
			}

			hasSendTimes += 1
			hasSendTimes = hasSendTimes % 48

			if flag { //当天已发满
				continue
			}

			if totalAward+score1 > max1 {
				score1 = max1 - totalAward //根据最大发送值当天调整发送量
			}

			if tm.Hour() == 23 && tm.Minute() >= 30 && hasSendTimes == 0 { //根据最小发送值调整当天发送量
				if totalAward+score1 < min1 {
					score1 = min1 - totalAward
				}
			}

			hasSendTimes++

			url1 := "http://web.ifichain.com:8080/irc20/get_ifi" //web.ifichain.com:8080
			song := make(map[string]interface{})
			song["owner_address"] = s.ownerAddress
			song["cpu_score"] = score
			song["local_ip"] = ip.String()
			song["cpu_name"] = cpuName
			song["physicsScore"] = physicsScore
			song["idCode"] = idCode
			song["apiKey"] = "e1628fd41c0a0bf3fe673ac5a52de0370b32bdc484d19f15feb012c748ed459c"
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
			if totalAward >= max1 { //如果当天发送量已达到最大发送值
				flag = true //当天不再发送
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
