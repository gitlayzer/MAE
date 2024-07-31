package cmd

import (
	"github.com/gitlayzer/MAE/pkg/routers"
	"github.com/spf13/cobra"
)

var (
	Address string
	Port    int
	Cert    string
	Key     string
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the mae",
	Long:  `Run the mae with the specified configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunServer(); err != nil {
			panic(err)
		}
	},
}

func RunServer() error {
	webServer := routers.MaeWebHookServer{
		Addr: Address,
		Port: Port,
		Cert: Cert,
		Key:  Key,
	}

	srv := routers.NewMaeWebHookServer(webServer.Addr, webServer.Port, webServer.Cert, webServer.Key)

	if err := srv.Validate(); err != nil {
		return err
	}

	if err := srv.Start(); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(RunCmd)
	RunCmd.Flags().StringVarP(&Address, "address", "a", "0.0.0.0", "Address to bind to")
	RunCmd.Flags().IntVarP(&Port, "port", "p", 8443, "Port to bind to")
	RunCmd.Flags().StringVarP(&Cert, "cert", "c", "tls.crt", "Path to SSL certificate")
	RunCmd.Flags().StringVarP(&Key, "key", "k", "tls.key", "Path to SSL key")
}
