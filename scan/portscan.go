package scan

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type Addr struct {
	ip   string
	port int
}

func PortScan(hosts []string, ports []int, timeout int64) []string {
	var AliveAddress []string

	workers := Workers
	Addrs := make(chan Addr, 100)
	results := make(chan string, 100)
	var wg sync.WaitGroup

	//接收结果
	go func() {
		for found := range results {
			AliveAddress = append(AliveAddress, found)
			wg.Done()
		}
	}()

	//多线程扫描
	for i := 0; i < workers; i++ {
		go func() {
			for addr := range Addrs {
				PortConnect(addr, results, timeout, &wg)
				wg.Done()
			}
		}()
	}

	//添加扫描目标
	for _, port := range ports {
		if port < 1 || port > 65535 {
			log.Printf("端口超出范围(1~65535): %d, 忽略扫描!", port)
			continue
		}
		// log.Printf("******* 端口扫描: %d *******", port)
		for _, host := range hosts {
			wg.Add(1)
			Addrs <- Addr{host, port}
		}
	}
	wg.Wait()
	close(Addrs)
	close(results)
	return AliveAddress
}

func PortConnect(addr Addr, respondingHosts chan<- string, adjustedTimeout int64, wg *sync.WaitGroup) {
	host, port := addr.ip, addr.port
	conn, err := WrapperTcpWithTimeout("tcp4", fmt.Sprintf("%s:%v", host, port), time.Duration(adjustedTimeout)*time.Second)
	if err == nil {
		defer conn.Close()
		address := host + ":" + strconv.Itoa(port)
		// result := fmt.Sprintf("%s open", address)
		// fmt.Println(result)
		wg.Add(1)
		respondingHosts <- address
	}
}
