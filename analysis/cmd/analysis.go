package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
    "tcpanalysis/pkg/pcap"
)

var (
	pcapfile string
	rootCmd  = &cobra.Command{
		Use:   "tcpana",
		Short: "tcpana is a small tool to analysis tcp payload",
		Run: func(cmd *cobra.Command, args []string) {
            mpcap.LoadPayload(pcapfile)
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&pcapfile, "pcapfile", "f", "","pcap file")
    rootCmd.MarkFlagRequired("pcapfile")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
