package network

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type BridgeNetworkDriver struct {
}

func (b *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (b *BridgeNetworkDriver) Create(subnet string, name string) (*NetWork, error) {
	return nil, nil
}

func (b *BridgeNetworkDriver) Delete(network NetWork) error {
	return nil
}

func (b *BridgeNetworkDriver) Connect(network *NetWork, endpoint *Endpoint) error {
	return nil
}

func (b *BridgeNetworkDriver) Disconnect(network NetWork, endpoint *Endpoint) error {
	return nil
}

func (b *BridgeNetworkDriver) initBridge(n *NetWork) error {

}

func (b *BridgeNetworkDriver) deleteBridge(n *NetWork) error {

}

func (b *BridgeNetworkDriver) createBridgeInterface(n string) error {
}

func setInterfaceUp(interfaceName string) error {

}

func setInterfaceIp(interfaceName string) error {
}

/*
 */
func setupIpTables(bridgeName string, subnet *net.IPNet) error {
	iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("iptables err %v, %w", output, err)
	}
	return nil
}
