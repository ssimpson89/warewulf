package wicked

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/app/wwctl/overlay/show"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_wickedOverlay(t *testing.T) {
	env := testenv.New(t)
	defer env.RemoveAll()
	env.ImportFile("var/lib/warewulf/overlays/wicked/rootfs/etc/wicked/ifconfig/ifcfg.xml.ww", "../rootfs/etc/wicked/ifconfig/ifcfg.xml.ww")

	tests := []struct {
		name       string
		nodes_conf string
		args       []string
		log        string
	}{
		{
			name:       "wicked",
			nodes_conf: "nodes.conf",
			args:       []string{"--render", "node1", "wicked", "etc/wicked/ifconfig/ifcfg.xml.ww"},
			log:        wicked,
		},
		{
			name:       "wicked-vlans",
			nodes_conf: "nodes.conf-vlan",
			args:       []string{"--render", "node1", "wicked", "etc/wicked/ifconfig/ifcfg.xml.ww"},
			log:        wicked_vlans,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env.ImportFile("etc/warewulf/nodes.conf", tt.nodes_conf)
			cmd := show.GetCommand()
			cmd.SetArgs(tt.args)
			stdout := bytes.NewBufferString("")
			stderr := bytes.NewBufferString("")
			logbuf := bytes.NewBufferString("")
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			wwlog.SetLogWriter(logbuf)
			err := cmd.Execute()
			assert.NoError(t, err)
			assert.Empty(t, stdout.String())
			assert.Empty(t, stderr.String())
			assert.Equal(t, tt.log, logbuf.String())
		})
	}
}

const wicked string = `backupFile: true
writeFile: true
Filename: ifcfg-default.xml

<!--
This file is autogenerated by warewulf
-->
<interface origin="static generated warewulf config">
  <name>wwnet0</name>
  <link-type>ethernet</link-type>
  <control>
    <mode>boot</mode>
  </control>
  <firewall/>
  <link/>
  <ipv4>
    <enabled>true</enabled>
    <arp-verify>true</arp-verify>
  </ipv4>
  <ipv4:static>
    <address>
      <local>192.168.3.21/24</local>
    </address>
    <route>
      <nexthop>
        <gateway>192.168.3.1</gateway>
      </nexthop>
    </route>
  </ipv4:static>
  <ipv6>
    <enabled>true</enabled>
    <privacy>prefer-public</privacy>
    <accept-redirects>false</accept-redirects>
  </ipv6>
</interface>

backupFile: true
writeFile: true
Filename: ifcfg-secondary.xml
<!--
This file is autogenerated by warewulf
-->
<interface origin="static generated warewulf config">
  <name>wwnet1</name>
  <link-type>ethernet</link-type>
  <control>
    <mode>boot</mode>
  </control>
  <firewall/>
  <link/>
  <ipv4>
    <enabled>true</enabled>
    <arp-verify>true</arp-verify>
  </ipv4>
  <ipv4:static>
    <address>
      <local>192.168.3.22/24</local>
    </address>
    <route>
      <nexthop>
        <gateway>192.168.3.1</gateway>
      </nexthop>
    </route>
  </ipv4:static>
  <ipv6>
    <enabled>true</enabled>
    <privacy>prefer-public</privacy>
    <accept-redirects>false</accept-redirects>
  </ipv6>
</interface>
`

const wicked_vlans string = `backupFile: true
writeFile: true
Filename: ifcfg-tagged.xml

<!--
This file is autogenerated by warewulf
-->
<interface origin="static generated warewulf config">
  <name>eth0.902</name>
  <link-type>vlan</link-type>
  <vlan>
    <device>eth0</device>
    <tag>902</tag>
    <protocol>ieee802-1Q</protocol>
  </vlan>
  <control>
    <mode>boot</mode>
  </control>
  <firewall/>
  <link/>
  <ipv4>
    <enabled>true</enabled>
    <arp-verify>true</arp-verify>
  </ipv4>
  <ipv4:static>
    <address>
      <local></local>
    </address>
    <route>
      <destination>192.168.1.0/24</destination>
      <nexthop>
        <gateway>192.168.2.254</gateway>
      </nexthop>
    </route>
  </ipv4:static>
  <ipv6>
    <enabled>true</enabled>
    <privacy>prefer-public</privacy>
    <accept-redirects>false</accept-redirects>
  </ipv6>
</interface>

backupFile: true
writeFile: true
Filename: ifcfg-untagged.xml
<!--
This file is autogenerated by warewulf
-->
<interface origin="static generated warewulf config">
  <name>eth0</name>
  <link-type>ethernet</link-type>
  <control>
    <mode>boot</mode>
  </control>
  <firewall/>
  <link/>
  <ipv4>
    <enabled>true</enabled>
    <arp-verify>true</arp-verify>
  </ipv4>
  <ipv4:static>
    <address>
      <local></local>
    </address>
  </ipv4:static>
  <ipv6>
    <enabled>true</enabled>
    <privacy>prefer-public</privacy>
    <accept-redirects>false</accept-redirects>
  </ipv6>
</interface>
`
