package rpc

import (
	"encoding/xml"
	"fmt"

	"github.com/czerwonk/junos_exporter/alarm"
	"github.com/czerwonk/junos_exporter/bgp"
	"github.com/czerwonk/junos_exporter/connector"
	"github.com/czerwonk/junos_exporter/interfaces"
)

type RpcClient struct {
	conn *connector.SshConnection
}

func NewClient(ssh *connector.SshConnection) *RpcClient {
	return &RpcClient{conn: ssh}
}

func (c *RpcClient) AlarmCounter() (*alarm.AlarmCounter, error) {
	red := 0
	yellow := 0

	cmds := []string{"show system alarms", "show chassis alarms"}

	for _, cmd := range cmds {
		var a = AlarmRpc{}
		err := c.runCommandAndParse(cmd, &a)
		if err != nil {
			return nil, err
		}

		for _, d := range a.Information.Details {
			if d.Class == "Major" {
				red++
			} else if d.Class == "Minor" {
				yellow++
			}
		}
	}

	return &alarm.AlarmCounter{RedCount: float64(red), YellowCount: float64(yellow)}, nil
}

func (*RpcClient) InterfaceStats() ([]*interfaces.InterfaceStats, error) {
	return make([]*interfaces.InterfaceStats, 0), nil
}

func (*RpcClient) BgpSessions() ([]*bgp.BgpSession, error) {
	return make([]*bgp.BgpSession, 0), nil
}

func (c *RpcClient) runCommandAndParse(cmd string, obj interface{}) error {
	b, err := c.conn.RunCommand(fmt.Sprintf("%s | display xml", cmd))
	if err != nil {
		return err
	}

	err = xml.Unmarshal(b, obj)
	return err
}
