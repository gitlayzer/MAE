package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "mae",
	Short: "MAE is a webhook service for image acceleration provided to K8S",
	Long:  "This acceleration service can achieve domestic acceleration of Docker Hub through a WebHook strategy",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				klog.Error(err)
			}
		}
		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}
