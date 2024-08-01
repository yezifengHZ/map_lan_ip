package config

import (
	"log"
	"map_lan_ip/scan"
	"map_lan_ip/utils"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type PromethuesTargets struct {
	Targets []string         `yaml:"targets"` // Prometheus监控目标地址[IP:Port]
	Labels  PrometheusLabels `yaml:"labels"`  // 标签
}

type PrometheusLabels struct {
	PingAddr string `yaml:"pingaddr"` // fping目标地址
}

func ReadTargets(file string) ([]PromethuesTargets, error) {
	data, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var targets []PromethuesTargets
	err = yaml.Unmarshal(data, &targets)
	if err != nil {
		return nil, err
	}

	return targets, nil
}

func WriteTargets(file string, targetsList []PromethuesTargets) error {
	data, err := yaml.Marshal(&targetsList)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTargets(file string, newTargets []PromethuesTargets, usePorts map[int]bool) error {
	// 读取现有配置
	data, err := ReadTargets(file)
	if err != nil {
		return err
	}

	changed := false

	// 扫描未发现的Targets
	newDataTargets := []PromethuesTargets{}
	var MapNewTargets map[string]struct{}
	MapNewTargetsPingaddr := map[string]bool{}
	for _, t := range newTargets {
		MapNewTargets = utils.ConvertStrSliceToMap(t.Targets)
		MapNewTargetsPingaddr[t.Labels.PingAddr] = true
	}
	for _, d := range data {
		for _, addr := range d.Targets {
			// pingaddr 不一致，忽略
			if !MapNewTargetsPingaddr[d.Labels.PingAddr] {
				log.Println("Ping目标地址已改变,移除监控项:", addr, "Ping目标地址:", d.Labels.PingAddr)
				changed = true
				continue
			}
			if utils.ContainsInMap(MapNewTargets, addr) {
				newDataTargets = append(newDataTargets, PromethuesTargets{Targets: []string{addr}, Labels: PrometheusLabels{PingAddr: d.Labels.PingAddr}})
				continue
			}
			// 二次确认Targets是否存在
			hostport := strings.Split(addr, ":")
			if len(hostport) == 2 {
				host := hostport[0]
				port, _ := strconv.Atoi(hostport[1])
				if !usePorts[port] {
					log.Println("监控地址端口已改变,移除监控项:", addr)
					changed = true
					continue
				}
				aliveAddress := scan.PortScan([]string{host}, []int{port}, scan.Timeout)
				if len(aliveAddress) == 0 {
					log.Println("监控地址端口未发现,移除监控项:", addr)
					changed = true
					continue
				}
			}
			newDataTargets = append(newDataTargets, PromethuesTargets{Targets: []string{addr}, Labels: PrometheusLabels{PingAddr: d.Labels.PingAddr}})
		}
	}

	// 新增Targets
	var MapData map[string]struct{}
	for _, t := range newDataTargets {
		MapData = utils.ConvertStrSliceToMap(t.Targets)
	}
	for _, n := range newTargets {
		for _, addr := range n.Targets {
			if utils.ContainsInMap(MapData, addr) {
				continue
			}
			newDataTargets = append(newDataTargets, PromethuesTargets{Targets: []string{addr}, Labels: PrometheusLabels{PingAddr: n.Labels.PingAddr}})
			changed = true
		}
	}

	// 更新配置
	if changed {
		err = WriteTargets(file, newDataTargets)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateNodeExporterTargets(file string, newTargets PromethuesTargets, port []Ports) error {
	targetType := "node_exporter"
	nodeExpPorts := map[int]bool{}
	for _, p := range port {
		if p.Type == targetType {
			nodeExpPorts[p.Prot] = true
		}
	}
	targets := newTargets.Targets
	nodeTargets := []string{}
	for _, t := range targets {
		hostport := strings.Split(t, ":")
		if len(hostport) == 2 {
			p, _ := strconv.Atoi(hostport[1])
			if nodeExpPorts[p] {
				nodeTargets = append(nodeTargets, t)
			}
		}
	}

	newNodeExpTargets := []PromethuesTargets{{Targets: nodeTargets}}

	nodeExpFile := strings.Replace(file, ".yml", "_"+targetType+".yml", 1)

	err := UpdateTargets(nodeExpFile, newNodeExpTargets, nodeExpPorts)
	if err != nil {
		return err
	}

	return nil
}

func UpdateFpingTargets(file string, newTargets PromethuesTargets, port []Ports, pingTargets []string) error {
	targetType := "fping"
	fpingPorts := map[int]bool{}
	for _, p := range port {
		if p.Type == targetType {
			fpingPorts[p.Prot] = true
		}
	}
	targets := newTargets.Targets
	fpingTargets := []PromethuesTargets{}
	for _, t := range targets {
		hostport := strings.Split(t, ":")
		if len(hostport) == 2 {
			p, _ := strconv.Atoi(hostport[1])
			if fpingPorts[p] {
				for _, pingTarget := range pingTargets {
					fpingTargets = append(fpingTargets, PromethuesTargets{Targets: []string{t}, Labels: PrometheusLabels{PingAddr: pingTarget}})
				}
			}
		}
	}

	fpingFile := strings.Replace(file, ".yml", "_"+targetType+".yml", 1)

	err := UpdateTargets(fpingFile, fpingTargets, fpingPorts)
	if err != nil {
		return err
	}

	return nil
}
