package cmd

import (
	"fmt"
	"github.com/jqlu/ackctl/client"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func init() {
	createCmd.AddCommand(newCreateClusterCmd())
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resource(s)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func newCreateClusterCmd() *cobra.Command {
	var (
		file string
	)
	command := &cobra.Command{
		Use:   "cluster",
		Short: "Create a cluster",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			path, err := filepath.Abs(file)
			if err != nil {
				log.Fatalf("failed to parse file path: %v", err)
			}

			params, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("failed to read file %s: %v", path, err)
			}

			response, err := client.CreateCluster(params)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create cluster: %v", err)
				os.Exit(1)
			}

			fmt.Printf("Starting to create cluster %s\n", response.ClusterId)
		},
	}

	command.Flags().StringVarP(&file, "file", "f", "", "specify cluster creation parameters file")

	return command
}
