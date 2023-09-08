package main

import (
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	app "nameservice"
	nsclient "nameservice/x/nameservice/client"
	nsrest "nameservice/x/nameservice/client/rest"
)

const (
	storeAcc = "acc"
	storeNs  = "nameservice"
)

var defaultCLIHome = os.ExpandEnv("$HOME/.nscli")

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	mc := []sdk.ModuleClients{
		nsclient.NewModuleClient(storeNs, cdc),
	}

	rootCmd := &cobra.Command{
		Use:   "nscli",
		Short: "Nameservice Client",
	}
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentFlags().String(client.FlagFrom, "", "From mark owner of coins")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(defaultCLIHome),
		queryCmd(cdc, mc),
		txCmd(cdc, mc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
	)

	executor := cli.PrepareMainCmd(rootCmd, "NS", defaultCLIHome)
	if err := executor.Execute(); err != nil {
		panic(err)
	}
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err = os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)
		if err = viper.ReadInConfig(); err != nil {
			return err
		}
	}

	if err = viper.BindPFlag(client.FlagChainID,
		cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err = viper.BindPFlag(client.FlagFrom,
		cmd.PersistentFlags().Lookup(client.FlagFrom)); err != nil {
		return err
	}
	if err = viper.BindPFlag(cli.EncodingFlag,
		cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	if err = viper.BindPFlag(cli.OutputFlag,
		cmd.PersistentFlags().Lookup(cli.OutputFlag)); err != nil {
		return err
	}

	return nil
}

func queryCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		tx.SearchTxCmd(cdc),
		tx.QueryTxCmd(cdc),
		client.LineBreak,
		authcmd.GetAccountCmd(storeAcc, cdc),
	)

	for _, m := range mc {
		cmd.AddCommand(m.GetQueryCmd())
	}

	return cmd
}

func txCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	cmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		tx.GetBroadcastCommand(cdc),
		client.LineBreak,
	)

	for _, m := range mc {
		cmd.AddCommand(m.GetTxCmd())
	}

	return cmd
}

func registerRoutes(rs *lcd.RestServer) {
	rs.CliCtx = rs.CliCtx.WithAccountDecoder(rs.Cdc)
	rpc.RegisterRoutes(rs.CliCtx, rs.Mux)
	tx.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	auth.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeAcc)
	bank.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	nsrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeNs)
}
