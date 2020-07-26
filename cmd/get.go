package cmd

import (
	"fmt"
	"github.com/jqlu/ackctl/client"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"os"
	"text/tabwriter"
)

func init() {
	getCmd.AddCommand(newGetClusterCmd())
	getCmd.AddCommand(newGetNodePoolCmd())
	getCmd.AddCommand(newGetNodeCmd())
	rootCmd.AddCommand(getCmd)
}

const clusterDetailFormat = `Cluster:
Name:       {{.Name}}
Id:         {{.ClusterID}}
State:      {{.State}}
Region:     {{.RegionID}}
Created At: {{.Created}}
`

var clusterDetailTpl = template.Must(template.New("cluster_detail").Parse(clusterDetailFormat))

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resource(s)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func newGetClusterCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "cluster",
		Short: "Get cluster(s)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			csClient := client.GetCsClient()

			if len(args) == 0 {
				clusters, err := csClient.DescribeClusters("")
				if err != nil {
					panic(err)
				}

				w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
				fmt.Fprintln(w, "ID\tState\tRegion\tType\tName")
				for _, c := range clusters {
					fmt.Fprintf(w, "%s\t%v\t%v\t%v\t%v\n", c.ClusterID, c.State, c.RegionID, c.ClusterType, c.Name)
				}
				w.Flush()
				return
			}

			cluster, err := client.FindClusterByPrefix(args[0])
			if err != nil {
				log.Fatalf("Unable to find cluster: %v", err)
			}

			_ = clusterDetailTpl.Execute(os.Stdout, cluster)
		},
	}

	return command
}

func newGetNodePoolCmd() *cobra.Command {
	var (
		clusterId string
	)
	command := &cobra.Command{
		Use:   "nodepool",
		Short: "Get node pool(s) of a cluster",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			nodepools, err := client.ListNodePools(clusterId)
			if err != nil {
				log.Fatalf("failed to list node pools for cluster %s: %v", clusterId, err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			fmt.Fprintln(w, "Name\tId\tState\tTotal\tServing\tOffline")
			for _, n := range nodepools {
				fmt.Fprintf(w, "%s\t%v\t%v\t%v\t%v\t%v\n", n.Name, n.NodePoolId, n.State, n.TotalNodes, n.ServingNodes, n.OfflineNodes)
			}
			w.Flush()
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "specify cluster id")

	return command
}

func newGetNodeCmd() *cobra.Command {
	var (
		clusterId string
		nodePoolId string
	)
	command := &cobra.Command{
		Use:   "node",
		Short: "List nodes of a node pool",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			nodes, err := client.ListNodes(clusterId, nodePoolId)
			if err != nil {
				log.Fatalf("failed to list nodes in node pool %s for cluster %s: %v", nodePoolId, clusterId, err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			fmt.Fprintln(w, "Instance Id\tNode Name\tInstance Status\tRole\tInstance Type\tNode Status")
			for _, n := range nodes {
				fmt.Fprintf(w, "%s\t%v\t%v\t%v\t%v\t%v\n", n.InstanceId, n.NodeName, n.NodeStatus, n.InstanceRole, n.InstanceType, n.State)
			}
			w.Flush()
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "specify cluster id")
	command.Flags().StringVarP(&nodePoolId, "node-pool-id", "p", "", "specify node pool id")

	return command
}
