package network

import (
	"encoding/json"
	"fmt"
	"mini-docker/Docker/container"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

var (
	defaultNetworkPath = "/var/run/mydocker/network/network/"
	drivers            = map[string]NetworkDriver{}
	networks           = map[string]NetWork{}
)

type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"device"`
	IpAddress   net.IP           `json:"ip_address"`
	MacAddress  net.HardwareAddr `json:"mac_address"`
	Network     *NetWork         `json:"network"`
	PortMapping []string
}

type NetWork struct {
	Name      string
	IpRange   *net.IPNet
	Driver    string
	GatewayIp net.IP
	Subnet    string
}

type NetworkDriver interface {
	Name() string
	Create(subnet string, name string) (*NetWork, error)
	Delete(networkName string) error
	Connect(network *NetWork, endpoint *Endpoint) error
	Disconnect(network NetWork, endpoint *Endpoint) error
}

/*
将网络对象（NetWork 类型）的配置以 JSON 格式持久化（dump/保存）
*/
func (nw *NetWork) dump(dumpPath string) error {
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(dumpPath, 0755)
		} else {
			return fmt.Errorf("dump path err: %w", err)
		}
	}

	nwPath := path.Join(dumpPath, nw.Name)
	/*
		O_TRUNC：如果文件已存在 → 清空内容（截断为 0 字节）
		O_WRONLY：只写模式
		O_CREATE：不存在则创建
		→ 效果：总是创建一个新文件（或覆盖旧文件）
	*/
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open file err: %w", err)
	}
	defer nwFile.Close()

	// MarshalIndent 有缩进、换行，人类易读
	// Marshal 无多余空格、缩进、换行
	nwJson, err := json.MarshalIndent(nw, "", " ")
	if err != nil {
		return fmt.Errorf("save network config json err: %w", err)
	}

	_, err = nwFile.Write(nwJson)
	if err != nil {
		return fmt.Errorf("save network config json err : %w", err)
	}
	return nil
}

func remove(dumpPath string) error {

}

func (nw *NetWork) load(dumpPath string) error {
	return nil
}

func Init() error {
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver

	if _, err := os.Stat(defaultNetworkPath); err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(defaultNetworkPath, 0755)
		} else {
			return err
		}
	}

	_ = filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(nwPath, "/") {
			return nil
		}

		_, nwName := path.Split(nwPath)
		nw := NetWork{
			Name: nwName,
		}

		if err := nw.load(nwPath); err != nil {
			log.Errorf("error load network: %v", err)
		}

		networks[nwName] = nw
		return nil
	})
	return nil
}

// 配置容器网络端点的地址和路由
func configEndpointIpAddressAndRoute(ep *Endpoint, cinfo *container.ContainerInfo) error {
}

// 将容器的网络端点加入到容器的网络空间中
// 并锁定当前程序所执行的线程，使当前线程进入到容器的网络空间
// 返回值是一个函数指针，执行这个返回函数才会退出容器的网络空间，回归到宿主机的网络空间
// 这个函数中引用了之前介绍的github.com/vishvananda/netns类库来做namespace操作
func enterContainer(enLink *netlink.Link, ep *Endpoint, cinfo *container.ContainerInfo) error {
}

// 配置端口映射
func configPortMapping(ep *Endpoint, cinfo *container.ContainerInfo) error {
}

func Connect(networkName string, cinfo *container.ContainerInfo) error {

}

func CreateNetwork(driver, subnet, name string) error {
}
func ListNetwork() error {

}

func DeleteNetwork(networkName string) error {
}
