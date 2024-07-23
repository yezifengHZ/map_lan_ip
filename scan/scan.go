package scan

// 线程数
var Workers int = 600

// 端口扫描超时时间
var Timeout int64 = 3

// Socks代理
var Socks5Proxy string

func Scan(cidrs []string, ports []int) []string {
	// log.Println("******* 开始扫描 *******")

	hosts, _ := ParseIP(cidrs)

	aliveAddress := PortScan(hosts, ports, Timeout)

	return aliveAddress
}
