package task

import (
	"log"
	"map_lan_ip/config"
	"map_lan_ip/scan"
	"time"
)

func MapLanIp() {
	c, err := config.ReadConfig()
	if err != nil {
		log.Fatal("打开配置文件失败:", err.Error())
		return
	}

	interval := c.Interval
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	do := func(){
		start := time.Now()

		// 扫描端口
		scanPorts := []int{}
		for _, port := range c.Ports {
			scanPorts = append(scanPorts, port.Prot)
		}
		aliveAddress := scan.Scan(c.CIDRs, scanPorts)

		// 更新Targets.yml
		newTargets := config.PromethuesTargets{Targets: aliveAddress, Labels: config.PrometheusLabels{PingAddr: c.PingAddr}}
		err := config.UpdateNodeExporterTargets(c.Target, newTargets, c.Ports)
		if err != nil {
			log.Fatal("更新 NodeExporterTargets 失败:", err.Error())
			return
		}

		err = config.UpdateFpingTargets(c.Target, newTargets, c.Ports)
		if err != nil {
			log.Fatal("更新 FpingTargets 失败:", err.Error())
			return
		}

		log.Printf("[*] 扫描结束,耗时: %s\n", time.Since(start))
	}

	// 初次扫描
	do()

	for {
		select {
		case <-ticker.C:
			do()
		}
	}
}
