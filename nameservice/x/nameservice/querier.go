package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"
)

const (
	QueryResolve = "resolve"
	QueryWhois   = "whois"
	QueryNames   = "names"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryResolve:
			return queryResolve(ctx, path[1:], req, keeper)
		case QueryWhois:
			return queryWhois(ctx, path[1:], req, keeper)
		case QueryNames:
			return queryNames(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown nameservice query endpoint")
		}
	}
}

func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	name := path[0]
	value := keeper.ResolveName(ctx, name)
	if value == "" {
		return []byte{}, sdk.ErrUnknownRequest("Could not resolve name")
	}

	b, err := codec.MarshalJSONIndent(keeper.cdc, QueryResResolve{value})
	if err != nil {
		panic("Could not marshal result to jSON, err: " + err.Error())
	}

	return b, nil
}

type QueryResResolve struct {
	Value string `json:"value"`
}

func (r QueryResResolve) String() string {
	return r.Value
}

func queryWhois(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	name := path[0]
	whois := keeper.GetWhois(ctx, name)

	b, err := codec.MarshalJSONIndent(keeper.cdc, whois)
	if err != nil {
		panic("Could not marshal result to jSON, err: " + err.Error())
	}

	return b, nil
}

func queryNames(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	list := QueryResNames{}

	iterator := keeper.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		list = append(list, string(iterator.Key()))
	}

	b, err := codec.MarshalJSONIndent(keeper.cdc, list)
	if err != nil {
		panic("Could not marshal result to jSON, err: " + err.Error())
	}

	return b, nil
}

type QueryResNames []string

func (n QueryResNames) String() string {
	return strings.Join(n, "\n")
}
