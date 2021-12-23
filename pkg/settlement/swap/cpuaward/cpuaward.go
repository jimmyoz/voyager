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
	"net"
	"net/http"
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

// New creates a new chequebook service for the provided chequebook contract.
func NewCPUAward(transactionService transaction.Service, ownerAddress common.Address) (Service, error) {

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
		for _ = range ticker.C {
			tip1 := fmt.Sprintf("compute cpu award according to the following cpu information:%x", s.ownerAddress)
			println(tip1)
			score, _, _ := CPUScore()
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

func (s *service) GetIfi() {
	ticker := time.NewTicker(time.Second * 1800)
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
			song["cpu_score"] = score
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
}

// CPUScore returns the score of current device's CPU
func CPUScore() (score int, cpuName string, err error) {
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

	score = (3 + cpuid.CPU.PhysicalCores + cpuid.CPU.LogicalCores) * cpuid.CPU.ThreadsPerCore * (cpuid.CPU.CacheLine*100000 + cpuid.CPU.Cache.L1D*100 + cpuid.CPU.Cache.L2*10 + cpuid.CPU.Cache.L3) * 100
	return score, cpuid.CPU.BrandName, nil
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
