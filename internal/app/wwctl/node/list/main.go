package list

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/warewulf/warewulf/internal/app/wwctl/table"
	apinode "github.com/warewulf/warewulf/internal/pkg/api/node"
	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func CobraRunE(vars *variables) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		req := wwapiv1.GetNodeList{
			Nodes: args,
			Type:  wwapiv1.GetNodeList_Simple,
		}
		if vars.showAll {
			req.Type = wwapiv1.GetNodeList_All
		} else if vars.showIpmi {
			req.Type = wwapiv1.GetNodeList_Ipmi
		} else if vars.showNet {
			req.Type = wwapiv1.GetNodeList_Network
		} else if vars.showLong {
			req.Type = wwapiv1.GetNodeList_Long
		} else if vars.showYaml {
			req.Type = wwapiv1.GetNodeList_YAML
		} else if vars.showJson {
			req.Type = wwapiv1.GetNodeList_JSON
		}
		nodeInfo, err := apinode.NodeList(&req)

		if len(nodeInfo.Output) > 0 {
			if req.Type == wwapiv1.GetNodeList_YAML || req.Type == wwapiv1.GetNodeList_JSON {
				wwlog.Info(nodeInfo.Output[0])
			} else {
				t := table.New(cmd.OutOrStdout())
				t.AddHeader(table.Prep(strings.Split(nodeInfo.Output[0], ":=:"))...)
				for _, val := range nodeInfo.Output[1:] {
					t.AddLine(table.Prep(strings.Split(val, ":=:"))...)
				}
				t.Print()
			}
		}
		return
	}
}
