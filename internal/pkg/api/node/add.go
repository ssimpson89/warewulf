package apinode

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/hostlist"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/util"
	"github.com/warewulf/warewulf/internal/pkg/warewulfd"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
	"gopkg.in/yaml.v3"
)

// NodeAdd adds nodes for management by Warewulf.
func NodeAdd(nap *wwapiv1.NodeAddParameter) (err error) {

	if nap == nil {
		return fmt.Errorf("NodeAddParameter is nil")
	}

	nodeDB, err := node.New()
	if err != nil {
		return fmt.Errorf("failed to open node database: %w", err)
	}
	dbHash := nodeDB.Hash()
	if hex.EncodeToString(dbHash[:]) != nap.Hash && !nap.Force {
		return fmt.Errorf("got wrong hash, not modifying node database")
	}
	node_args := hostlist.Expand(nap.NodeNames)
	var ipv4, ipmiaddr net.IP
	for _, a := range node_args {
		n, err := nodeDB.AddNode(a)
		if err != nil {
			return fmt.Errorf("failed to add node: %w", err)
		}
		err = yaml.Unmarshal([]byte(nap.NodeConfYaml), &n)
		if err != nil {
			return fmt.Errorf("failed to decode nodeConf: %w", err)
		}
		wwlog.Info("Added node: %s", a)
		for _, dev := range n.NetDevs {
			if !ipv4.IsUnspecified() && ipv4 != nil {
				// if more nodes are added increment IPv4 address
				ipv4 = util.IncrementIPv4(ipv4, 1)
				wwlog.Verbose("Incremented IP addr to %s", ipv4)
				dev.Ipaddr = ipv4

			} else if !dev.Ipaddr.IsUnspecified() {
				ipv4 = dev.Ipaddr
			}
		}
		if n.Ipmi != nil {
			if !ipmiaddr.IsUnspecified() && ipmiaddr != nil {
				ipmiaddr = util.IncrementIPv4(ipmiaddr, 1)
				wwlog.Verbose("Incremented ipmi IP addr to %s", ipmiaddr)
				n.Ipmi.Ipaddr = ipmiaddr
			} else if !n.Ipmi.Ipaddr.IsUnspecified() {
				ipmiaddr = n.Ipmi.Ipaddr
			}
		}
	}

	err = nodeDB.Persist()
	if err != nil {
		return fmt.Errorf("failed to persist new node: %w", err)
	}

	err = warewulfd.DaemonReload()
	if err != nil {
		return fmt.Errorf("failed to reload warewulf daemon: %w", err)
	}
	return
}
