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

func ReadTargets(file string) (PromethuesTargets, error) {
	data, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		return PromethuesTargets{}, nil
	}
	if err != nil {
		return PromethuesTargets{}, err
	}

	var targets []PromethuesTargets
	err = yaml.Unmarshal(data, &targets)
	if err != nil {
		return PromethuesTargets{}, err
	}

	return targets[0], nil
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

func UpdateTargets(file string, newTargets PromethuesTargets) error {
	// 读取现有配置
	data, err := ReadTargets(file)
	if err != nil {
		return err
	}

	changed := false

	// 扫描未发现的Targets
	newDataTargets := []string{}
	MapNewTargets := utils.ConvertStrSliceToMap(newTargets.Targets)
	for _, addr := range data.Targets {
		if utils.ContainsInMap(MapNewTargets, addr) {
			newDataTargets = append(newDataTargets, addr)
			continue
		}
		// 二次确认Targets是否存在
		hostport := strings.Split(addr, ":")
		if len(hostport) == 2 {
			host := hostport[0]
			port, _ := strconv.Atoi(hostport[1])
			aliveAddress := scan.PortScan([]string{host}, []int{port}, scan.Timeout)
			if len(aliveAddress) == 0 {
				log.Println("监控地址端口未发现,移除监控项:", addr)
				changed = true
				continue
			}
		}
		newDataTargets = append(newDataTargets, addr)
	}
	data.Targets = newDataTargets

	// 新增Targets
	MapData := utils.ConvertStrSliceToMap(data.Targets)
	for _, addr := range newTargets.Targets {
		if utils.ContainsInMap(MapData, addr) {
			continue
		}
		data.Targets = append(data.Targets, addr)
		changed = true
	}
	data.Labels.PingAddr = newTargets.Labels.PingAddr

	// 更新配置
	if changed {
		err = WriteTargets(file, []PromethuesTargets{data})
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

	newNodeExpTargets := PromethuesTargets{Targets: nodeTargets}

	nodeExpFile := strings.Replace(file, ".yml", "_"+targetType+".yml", 1)

	err := UpdateTargets(nodeExpFile, newNodeExpTargets)
	if err != nil {
		return err
	}

	return nil
}

func UpdateFpingTargets(file string, newTargets PromethuesTargets, port []Ports) error {
	targetType := "fping"
	fpingPorts := map[int]bool{}
	for _, p := range port {
		if p.Type == targetType {
			fpingPorts[p.Prot] = true
		}
	}
	targets := newTargets.Targets
	fpingTargets := []string{}
	for _, t := range targets {
		hostport := strings.Split(t, ":")
		if len(hostport) == 2 {
			p, _ := strconv.Atoi(hostport[1])
			if fpingPorts[p] {
				fpingTargets = append(fpingTargets, t)
			}
		}
	}

	newFpingTargets := PromethuesTargets{Targets: fpingTargets, Labels: PrometheusLabels{PingAddr: newTargets.Labels.PingAddr}}

	fpingFile := strings.Replace(file, ".yml", "_"+targetType+".yml", 1)

	err := UpdateTargets(fpingFile, newFpingTargets)
	if err != nil {
		return err
	}

	return nil
}
