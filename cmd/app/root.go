package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ppot7/haybaler/eodhdapi"
	"github.com/ppot7/haybaler/eodpostgres"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "haybaler",
		Short: "A Data Migration Tool",
		Long: `Haybaler allows the migration of data for
	importing, updating and extracting daily asset data`,
	}

	configFile string
)

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "f", "haybaler.config", "file containing db and api configurations")

	rootCmd.AddCommand(addCmd)
}

// config, err := eodpostgres.CreateDefaultConfiguration(envMap["eod.ps.host"], envMap["eod.ps.port"], envMap["eod.ps.user"],
// 	envMap["eod.ps.pwd"], envMap["eod.ps.db"])
// if err != nil {
// 	fmt.Println("error connecting to postgres db: ", err)
// 	return
// }
// conn, err := eodpostgres.ConnectToPsDatabase(context.TODO(), envMap["eod.ps.schema"], envMap["eod.ps.pv_table"],
// 	envMap["eod.ps.dividend_table"], envMap["eod.ps.split_table"], config)
// if err != nil {
// 	fmt.Println("error connecting to postgres db: ", err)
// 	return
// }
// defer conn.Close(context.TODO())

func createEodClient(configMap map[string]string) *eodhdapi.EodHdApiClient {
	return eodhdapi.CreateEodHdClient(configMap["eodhd.host"], configMap["endhd.token"], nil)
}

func createEodPsConnection(configMap map[string]string) (*eodpostgres.EodPsConnection, error) {
	config, err := eodpostgres.CreateDefaultConfiguration(configMap["eod.ps.host"], configMap["eod.ps.port"],
		configMap["eod.ps.user"], configMap["eod.ps.pwd"], configMap["eod.ps.db"])
	if err != nil {
		slog.Error("could not establish connection to ps database", "err", err)
		return nil, fmt.Errorf("could not establish connection to ps database")
	}

	conn, err := eodpostgres.ConnectToPsDatabase(context.TODO(), configMap["eod.ps.schema"], configMap["eod.ps.pv_table"],
		configMap["eod.ps.dividend_table"], configMap["eod.ps.split_table"], config)
	if err != nil {
		slog.Error("connection error ", "err", err)
		return nil, fmt.Errorf("connection error %s", err)
	}

	return conn, nil
}
