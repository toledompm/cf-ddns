package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/toledompm/cloudflare-ddns/config"
	"github.com/toledompm/cloudflare-ddns/handlers"
)

var (
	cfg     *config.Config
	cfgFile string
	cf      handlers.ICloudflare
	nw      handlers.Network
	rootCmd = &cobra.Command{
		Use:   "cf-ddns",
		Short: "Cloudflare DDNS Updater",
		Long:  `Cloudflare DDNS Updater is a simple tool to update Cloudflare DNS records with your current IP address.`,
		Run: func(cmd *cobra.Command, args []string) {
			for _, record := range cfg.Records {
				if cfg.IPV6.Enabled {
					ipv6, err := nw.GetIPV6()
					if err != nil {
						fmt.Printf("Error getting IP: %s\n", err)
						os.Exit(1)
					}

					err = cf.UpdateRecord(ipv6, record.Name, cfg.Cloudflare.ZoneID, record.Proxy, "AAAA")
					if err != nil {
						fmt.Printf("Error updating record: %s\n", err)
						os.Exit(1)
					}
				}

				if cfg.IPV4.Enabled {
					ipv4, err := nw.GetIPV4()
					if err != nil {
						fmt.Printf("Error getting IP: %s\n", err)
						os.Exit(1)
					}

					err = cf.UpdateRecord(ipv4, record.Name, cfg.Cloudflare.ZoneID, record.Proxy, "A")
					if err != nil {
						fmt.Printf("Error updating record: %s\n", err)
						os.Exit(1)
					}
				}
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "(required) path to config file")
}

func initConfig() {
	if cfgFile == "" {
		fmt.Printf("No config file specified, use --help to see required flags\n")
		os.Exit(1)
	}

	cfgFileBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	cfg = config.New()
	cfg = config.MustParseConfig(cfgFileBytes, cfg)

	cf = handlers.NewCloudflare(cfg.Cloudflare.Token)
	nw = handlers.NewNetwork(cfg.IPV6.FetchAddress, cfg.IPV4.FetchAddress)
}
