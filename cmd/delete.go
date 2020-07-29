package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/jqlu/ackctl/client"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	deleteCmd.AddCommand(newDeleteClusterCmd())
	deleteCmd.AddCommand(newDeleteNodeCmd())
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func newDeleteClusterCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "cluster",
		Short: "Delete a cluster",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cluster, err := client.FindClusterByPrefix(args[0])
			if err != nil {
				log.Fatalf("Unable to find cluster: %v", err)
			}

			confirm := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Are you sure to delete cluster %s(%s)? Cluster cannot be restored after deletion",
					cluster.Name, cluster.ClusterID),
				Default: false,
			}
			_ = survey.AskOne(prompt, &confirm)

			if !confirm {
				return
			}

			csClient := client.GetCsClient()
			if err := csClient.DeleteKubernetesCluster(cluster.ClusterID); err != nil {
				log.Fatalf("Failed to delete cluster %s(%s): %v", cluster.Name, cluster.ClusterID, err)
			}

			fmt.Printf("Starting to delete cluster %s(%s)\n", cluster.Name, cluster.ClusterID)
		},
	}

	return command
}

func newDeleteNodeCmd() *cobra.Command {
	var (
		clusterId string
		drain     bool
		release   bool
	)

	command := &cobra.Command{
		Use:   "node",
		Short: "Delete a node",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatalf("Missing node name")
			}

			nodeName := args[0]
			if err := client.RemoveNode(clusterId, nodeName, release, drain); err != nil {
				log.Fatalf("Failed to delete node %s of cluster %s:%v", nodeName, clusterId, err)
			}

			fmt.Printf("Staring to delete node %s of cluster %s\n", nodeName, clusterId)
		},
	}

	command.Flags().StringVarP(&clusterId, "cluster-id", "c", "", "specify cluster id")
	command.Flags().BoolVar(&drain, "drain", false, "whether to drain the node before deletion")
	command.Flags().BoolVar(&release, "release", false, "whether to release the ECS instance after deletion")

	return command
}
