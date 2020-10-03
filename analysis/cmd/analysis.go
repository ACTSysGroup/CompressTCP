package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"tcpanalysis/pkg/pcap"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
     "tcpanalysis/common"
)


var dataset_map map[string]common.DatasetMeta

var (
	schemafile string
    dataset_key string
	rootCmd    = &cobra.Command{
		Use:   "tcpana",
		Short: "tcpana is a small tool to analysis tcp payload",
		Run: func(cmd *cobra.Command, args []string) {
            var schema_list  []common.DatasetMeta
            dataset_map = make(map[string]common.DatasetMeta)
            if data,err:=ioutil.ReadFile(schemafile);err!=nil{
                fmt.Printf ("no schemafile found, err:%s",err.Error())
            } else{
                yaml.Unmarshal(data, &schema_list)
                for _,schema := range schema_list{
                    dataset_map[schema.Key] = schema
                }
                if schema,ok:=dataset_map[dataset_key];ok{
                    mpcap.LoadPayload(schema)
                } else{
                    fmt.Printf("no dataset found. available:\n")
                    for k,_ := range dataset_map{
                        fmt.Printf("\t%s\n",k)
                    }
                }
            }
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&schemafile, "schemafile", "s", "", "schema file")
	rootCmd.Flags().StringVarP(&dataset_key, "dataset", "d", "", "key of dataset")
	rootCmd.MarkFlagRequired("schemafile")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
