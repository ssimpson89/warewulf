package add

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/warewulfd"
)

func Test_Add(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		stdout  string
		inDb    string
		outDb   string
	}{
		{name: "single node add",
			args:    []string{"n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01: {}
`},
		{name: "single node add with default profile",
			args:    []string{"n01"},
			wantErr: false,
			stdout:  "",
			inDb: `
nodeprofiles:
  default: {}`,
			outDb: `
nodeprofiles:
  default: {}
nodes:
  n01:
    profiles:
    - default`},
		{name: "single node add, profile foo",
			args:    []string{"--profile=foo", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    profiles:
    - foo
`},
		{name: "single node add, discoverable true, explicit",
			args:    []string{"--discoverable=true", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    discoverable: "true"
`},
		{name: "single node add, discoverable true, implicit",
			args:    []string{"--discoverable", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    discoverable: "true"
`},
		{name: "single node add, discoverable wrong argument",
			args:    []string{"--discoverable=maybe", "n01"},
			wantErr: true,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes: {}
`},
		{name: "single node add, discoverable false",
			args:    []string{"--discoverable=false", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    discoverable: "false"
`},
		{name: "single node add with Kernel args",
			args:    []string{"--kernelargs=foo", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    kernel:
      args:
      - foo
`},
		{name: "double node add explicit",
			args:    []string{"n01", "n02"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01: {}
  n02: {}
`},
		{name: "single node with ipaddr6",
			args:    []string{"--ipaddr6=fdaa::1", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    network devices:
      default:
        ip6addr: fdaa::1
`},
		{name: "single node with ipaddr",
			args:    []string{"--ipaddr=10.0.0.1", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    network devices:
      default:
        ipaddr: 10.0.0.1
`},
		{name: "single node with malformed ipaddr",
			args:    []string{"--ipaddr=10.0.1", "n01"},
			wantErr: true,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes: {}
`},
		{name: "three nodes with ipaddr",
			args:    []string{"--ipaddr=10.10.0.1", "n[01-02,03]"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    network devices:
      default:
        ipaddr: 10.10.0.1
  n02:
    network devices:
      default:
        ipaddr: 10.10.0.2
  n03:
    network devices:
      default:
        ipaddr: 10.10.0.3
`},
		{name: "three nodes with ipaddr different network",
			args:    []string{"--ipaddr=10.10.0.1", "--netname=foo", "n[01-03]"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    network devices:
      foo:
        ipaddr: 10.10.0.1
  n02:
    network devices:
      foo:
        ipaddr: 10.10.0.2
  n03:
    network devices:
      foo:
        ipaddr: 10.10.0.3
`},
		{name: "three nodes with ipaddr different network, with ipmiaddr",
			args:    []string{"--ipaddr=10.10.0.1", "--netname=foo", "--ipmiaddr=10.20.0.1", "n[01-03]"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    ipmi:
      ipaddr: 10.20.0.1
    network devices:
      foo:
        ipaddr: 10.10.0.1
  n02:
    ipmi:
      ipaddr: 10.20.0.2
    network devices:
      foo:
        ipaddr: 10.10.0.2
  n03:
    ipmi:
      ipaddr: 10.20.0.3
    network devices:
      foo:
        ipaddr: 10.10.0.3
`},
		{name: "one node with filesystem",
			args:    []string{"--fsname=/dev/vda1", "--fspath=/var", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    filesystems:
      /dev/vda1:
        path: /var
`},
		{name: "one node with filesystem",
			args:    []string{"--fsname=dev/vda1", "--fspath=/var", "n01"},
			wantErr: true,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes: {}
`},
		{name: "one node with filesystem and partition ",
			args:    []string{"--fsname=var", "--fspath=/var", "--partname=var", "--diskname=/dev/vda", "--partnumber=1", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    disks:
      /dev/vda:
        partitions:
          var:
            number: "1"
    filesystems:
      /dev/disk/by-partlabel/var:
        path: /var
`},
		{name: "one node with filesystem with btrfs and partition ",
			args:    []string{"--fsname=var", "--fspath=/var", "--fsformat=btrfs", "--partname=var", "--diskname=/dev/vda", "--partnumber=1", "n01"},
			wantErr: false,
			stdout:  "",
			outDb: `nodeprofiles: {}
nodes:
  n01:
    disks:
      /dev/vda:
        partitions:
          var:
            number: "1"
    filesystems:
      /dev/disk/by-partlabel/var:
        format: btrfs
        path: /var
`},
	}
	warewulfd.SetNoDaemon()
	for _, tt := range tests {
		env := testenv.New(t)
		env.WriteFile("etc/warewulf/nodes.conf", tt.inDb)
		var err error
		t.Run(tt.name, func(t *testing.T) {
			baseCmd := GetCommand()
			baseCmd.SetArgs(tt.args)
			buf := new(bytes.Buffer)
			baseCmd.SetOut(buf)
			baseCmd.SetErr(buf)
			err = baseCmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			config, configErr := node.New()
			assert.NoError(t, configErr)
			dumpBytes, _ := config.Dump()
			assert.YAMLEq(t, tt.outDb, string(dumpBytes))
		})
	}
}
