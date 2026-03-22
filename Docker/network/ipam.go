package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

// const ipamDefaultAllocatorPath = "/var/run/mydocker/network/ipam/subnet.json"
const ipamDefaultAllocatorPath = "/tmp/mydocker-test/network/ipam/subnet.json"

// IPAM 存放IP地址分配信息
type IPAM struct {
	// 分配文件存放路径
	SubnetAllocatorPath string
	// 网段和位图算法的数组map，key是网络，value是分配的位图数组
	Subnets *map[string]string
}

func NewIPAM() *IPAM {
	ipam := &IPAM{
		Subnets: &map[string]string{},
	}
	_ = ipam.load()
	return ipam
}

// 初始化一个IPAMd对象
var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

// 加载网段地址分配信息
func (ipam *IPAM) load() error {
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	if err != nil {
		return err
	}
	defer subnetConfigFile.Close()
	data, err := os.ReadFile(ipam.SubnetAllocatorPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, ipam.Subnets)
	if err != nil {
		// fmt.Println(string(data))
		return fmt.Errorf("dump allocation info err: %v", err)
	}
	log.Infof("load ipam file from: %s", subnetConfigFile.Name())
	return nil
}

// 保存网段地址分配信息
func (ipam *IPAM) dump() error {
	ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigFileDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(ipamConfigFileDir, 0755)
		} else {
			return err
		}
	}

	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath,
		os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Errorf("70 error!!!!!!!!!")
		return err
	}
	defer subnetConfigFile.Close()

	ipamConfigJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		return err
	}
	_, err = subnetConfigFile.Write(ipamConfigJson)
	if err != nil {
		return err
	}
	log.Infof("dump ipam file from: %s", ipamConfigFileDir)
	return nil
}

// 在网段中分配一个可用的IP地址
func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	if subnet == nil || subnet.IP == nil || subnet.Mask == nil {
		return nil, fmt.Errorf("invalid subnet")
	}

	// 存储网段中地址分配信息的数组
	if ipam.Subnets == nil {
		ipam.Subnets = &map[string]string{}
	}

	//从文件中加载已经分配了的网络信息
	if err = ipam.load(); err != nil {
		return nil, fmt.Errorf("load subnet file err: %v", err)
	}

	// net.ipnet.mask.size() 返回网段的子网掩码的总长度和网段前面的固定位的长度
	// 比如 127.0.0.0/8 网段的子网掩码是 255.0.0.0
	// 返回的是前面255所对应的位数和总位数，即8和24
	one, size := subnet.Mask.Size()

	// 如果之前没有分配过这个网段，则初始化网段的分配配置
	if _, exist := (*ipam.Subnets)[subnet.String()]; !exist {
		// 用0填满这个网段的配置， 1<<uint8(size-one)表示这个网段中有多少个可用的地址
		// size - one 是子网掩码后面的网络位数，2^（size-one）(即1<<uint8(size-one))表示可用的IP数
		(*ipam.Subnets)[subnet.String()] = strings.Repeat("0", 1<<uint8(size-one))
	}

	// 遍历网段的位图数组
	for c := range (*ipam.Subnets)[subnet.String()] {
		// 设置当前的序号值为1，即分配这个IP
		if (*ipam.Subnets)[subnet.String()][c] == '0' {
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])

			// Go中字符串不能修改，通过转成byte数组，再转成字符串赋能
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)

			// 这里的IP为初始IP，比如192.168.0.0/16，这里就是192.168.0.0
			ip = subnet.IP

			// 通过网段的IP与上面的偏移相加计算出分配的IP地址，由于IP地址是uint的一个数组
			// 需要通过数组中的每一项加所需要的值，比如网段172.16.0.0/12，数组序号是65555
			// 那么在[172,16,0,0〔方案選單〕]上依次加[uint8(65555 >> 24)、[uint8(65555 >> 16)、[uint8(65555 >> 8)、[uint8(65555 >> 8)
			// 即[0, 1, 0, 19],那么最后得到的172.17.0.19
			for t := uint(4); t > 0; t-- {
				[]byte(ip)[4-t] += uint8(c >> (8 * (t - 1)))
			}

			// 由于IP是从1开始的，所以加1
			ip[3] += 1
			break
		}
	}
	log.Infof("allocate subnet: %s, ip: %s", subnet.String(), ip.String())
	return ip, ipam.dump()
}

func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	log.Infof("release subnet: %s, ip: %s", subnet.String(), ipaddr.String())
	ipam.Subnets = &map[string]string{}

	// 从文件中加载网段分配信息
	err := ipam.load()
	if err != nil {
		return fmt.Errorf("load subnet file err: %v", err)
	}

	// 计算IP地址在网段位图数组中的索引位置
	index := 0
	// 将IP地址转换成4个字节的表现形式
	releaseIp := ipaddr.To4()
	// 由于IP是从1开始分配的，所以转换成索引应减一
	releaseIp[3] -= 1
	// 与分配IP相反，释放IP获得索引的方式是将IP地址的每一位相减后分别左移将对应的数值加到索引上
	for t := uint(4); t > 0; t -= 1 {
		index += int(releaseIp[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}
	log.Infof("release index: %d", index)

	// 将分配的位图索引中的位置的值为0
	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[index] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)

	// 保存释放IP后的配置信息
	return ipam.dump()
}
