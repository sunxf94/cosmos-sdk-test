package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	"nameservice/x/nameservice/client/cli"
)

type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	nameSvcQueryCmd := &cobra.Command{
		Use:   "nameservice",
		Short: "Querying commands for the nameservice module",
	}

	nameSvcQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdResolveName(mc.storeKey, mc.cdc),
		cli.GetCmdWhois(mc.storeKey, mc.cdc),
	)...)

	return nameSvcQueryCmd
}

func (mc ModuleClient) GetTxCmd() *cobra.Command {
	nameSvcTxCmd := &cobra.Command{
		Use:   "nameservice",
		Short: "Nameservice transactions subcommands",
	}

	nameSvcTxCmd.AddCommand(client.GetCommands(
		cli.GetCmdBuyName(mc.cdc),
		cli.GetCmdSetName(mc.cdc),
	)...)

	return nameSvcTxCmd
}
