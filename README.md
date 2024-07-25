## 构建


#### 二进制打包

**amd64:**

```
CGO_ENABLE=0  GOOS=linux  GOARCH=amd64 go build -a -o map_lan_ip
```

#### Tar打包

```
./build.sh
```

## 配置文件
```
# 扫描/更新时间间隔，单位：s
interval: 60

# 需要扫描的网段，尽量缩小扫描范围，扫描越多耗时越长
# 192.168.2.0/24 256 个 IP + 2 个端口      约 3s
# 192.168.0.0/16 65536 个 IP + 2 个端口    约 10m40s
# 172.16.0.0/12  1048576 个 IP + 2 个端口  约 16*(10m40s) = 170m40s = 2h50m40s
# 10.0.0.0/8     16777216 个 IP + 2 个端口 约 256*(10m40s) = 2730m40s = 45h30m40s
cidrs:
    # - 10.0.0.0/8
    # - 172.16.0.0/12
    # - 192.168.0.0/16
    - 192.168.20.0/24

# 需要扫描的端口，支持 node_exporter 和 fping 两种插件类型的端口
ports:
    - port: 9800
      type: node_exporter
    - port: 9605
      type: fping

# fping 的 ping 目标地址
pingaddr: www.baidu.com

# Prometheus 动态加载的文件路径，最终的文件名称根据端口类型分为 targets_node_exporter.yml 和 targets_fping.yml
target: /root/promethues/prometheus-2.53.1.linux-amd64/targets.yml

# 日志文件路径
logfile: /var/log/map_lan_ip.log
```

## 安装部署

```
tar zxvf map_lan_ip_amd64.tar.gz
cd map_lan_ip_amd64
# 记得修改你的配置文件
vim.tiny map_lan_ip.yml
bash install.sh
```