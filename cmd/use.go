package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/denverdino/aliyungo/cs"
	"github.com/jqlu/ackctl/client"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func init() {
	useCmd.AddCommand(newUseClusterCmd())
	rootCmd.AddCommand(useCmd)
}

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Configure kubeconfig for cluster(s)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func newUseClusterCmd() *cobra.Command {
	var (
		useAll bool
	)

	var command = &cobra.Command{
		Use:   "cluster",
		Short: "Configure kubeconfig for cluster(s)",
		Long:  ``,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clusters := make([]*cs.ClusterType, 0)
			if useAll {
				cl, err := client.GetAccessibleClusters()
				if err != nil {
					log.Fatalf("Failed to list accessible clusters: %v", err)
				}
				clusters = cl
			} else {
				cluster, err := client.FindClusterByPrefix(args[0])
				if err != nil {
					log.Fatalf("Unable to find cluster: %v", err)
				}
				clusters = append(clusters, cluster)
			}

			cb, err := client.GetMergedKubeConfig(clusters)
			if err != nil {
				log.Fatalf("Failed to get kubeConfig: %v", err)
			}

			userHome, err := os.UserHomeDir()
			if err != nil {
				return errors.New("failed to get user home")
			}
			configPath := filepath.Join(userHome, ".kube", "config")

			_, err = os.Stat(configPath)
			confirm := true
			if !os.IsNotExist(err) {
				prompt := &survey.Confirm{
					Message: fmt.Sprintf("Config %v exists, overwrite?", configPath),
					Default: true,
				}
				_ = survey.AskOne(prompt, &confirm)
			}

			if !confirm {
				return nil
			}

			err = ioutil.WriteFile(configPath, cb, 440)
			if err != nil {
				log.Fatalf("failed to write kubeConfig to file")
			}

			if useAll {
				fmt.Printf("Merged kubeConfigs of %v clusters into: %v.\n", len(clusters), configPath)
				fmt.Printf("Use 'kubectl config get-contexts' to list contexts.\n")
				fmt.Printf("Use 'kubectl config use-context' to select context.\n")
			} else {
				fmt.Printf("%v updated to use cluster\n", configPath)
			}

			return nil
		},
	}

	command.Flags().BoolVar(&useAll, "all", false, "get kubeconfig of all clusters and merge into one file")
	return command
}
