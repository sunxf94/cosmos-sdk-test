package nameservice

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSetName struct {
	Name  string
	Value string
	Owner sdk.AccAddress
}

func NewMsgSetName(name, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}

func (msg MsgSetName) Route() string {
	return "nameservice"
}

func (msg MsgSetName) Type() string {
	return "set_name"
}

func (msg MsgSetName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if msg.Name == "" || msg.Value == "" {
		return sdk.ErrUnknownRequest("Name and/or Value be empty")
	}

	return nil
}

func (msg MsgSetName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgBuyName struct {
	Name  string
	Bid   sdk.Coins
	Buyer sdk.AccAddress
}

func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

func (msg MsgBuyName) Route() string {
	return "nameservice"
}

func (msg MsgBuyName) Type() string {
	return "bug_name"
}

func (msg MsgBuyName) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}

	if msg.Name == "" {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}

	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}

	return nil
}

func (msg MsgBuyName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
