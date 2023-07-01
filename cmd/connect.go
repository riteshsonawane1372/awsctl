package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/surajincloud/awsctl/pkg/ec2"
)

var sshConnect = &cobra.Command{
	Use:   "ssh",
	Short: "Connect to EC2 instance",
	Long: `For example:
		awsctl ssh --ip 192.168.1.1
	`,
	RunE: connectEC2,
}

func connectEC2(cmd *cobra.Command, args []string) error {
	err := ec2.SshConnect(cmd, args)
	if err != nil {
		log.Fatal("Unable to connect via SSH", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(sshConnect)
	sshConnect.Flags().String("ip", "", `awsctl ssh --ip 192.168.1.1`)
	sshConnect.MarkFlagRequired("ip")
}
