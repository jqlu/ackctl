package cmd

import (
	"fmt"
	"github.com/jqlu/ackctl/client"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	scaleCmd.AddCommand(newScaleNodePoolCmd())
	rootCmd.AddCommand(scaleCmd)
}

var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale resource(s)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func newScaleNodePoolCmd() *cobra.Command {
	var (
		clusterId string
		increment int
	)

	command := &cobra.Command{
		Use:   "nodepool",
		Short: "Scale a node pool",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatalf("Missing node pool id")
			}

			nodePoolId := args[0]
			if err := client.ScaleNodePool(clusterId, nodePoolId, increment); err != nil {
				log.Fatalf("Failed to scale node pool %s of cluster %s:%v", nodePoolId, clusterId, err)
			}

			fmt.Printf("Staring to scale node pool %s of cluster %s\n", nodePoolId, clusterId)
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "specify cluster id")
	command.Flags().IntVar(&increment, "increment", 0, "increment of node(s)")

	return command
}
