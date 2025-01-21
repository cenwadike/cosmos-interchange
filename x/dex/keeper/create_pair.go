package keeper

import (
	"errors"

	"interchange/x/dex/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
)

// TransmitCreatePairPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitCreatePairPacket(
	ctx sdk.Context,
	packetData types.CreatePairPacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	channelCap, ok := k.ScopedKeeper().GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %s", err)
	}

	return k.ibcKeeperFn().ChannelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnRecvCreatePairPacket processes packet reception
func (k Keeper) OnRecvCreatePairPacket(ctx sdk.Context, packet channeltypes.Packet, data types.CreatePairPacketData) (packetAck types.CreatePairPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// Get an order book index
	pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.SourceDenom, data.TargetDenom)

	// If an order book is found, return an error
	_, found := k.GetBuyOrderBook(ctx, pairIndex)
	if found {
		return packetAck, errors.New("the pair already exist")
	}

	// Create a new buy order book for source and target denoms
	book := types.NewBuyOrderBook(data.SourceDenom, data.TargetDenom)

	// Assign order book index
	book.Index = pairIndex

	// Save the order book to the store
	k.SetBuyOrderBook(ctx, book)
	return packetAck, nil
}

// OnAcknowledgementCreatePairPacket responds to the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementCreatePairPacket(ctx sdk.Context, packet channeltypes.Packet, data types.CreatePairPacketData, ack channeltypes.Acknowledgement) error {
	// switch dispatchedAck := ack.Response.(type) {
	// case *channeltypes.Acknowledgement_Error:

	// 	// TODO: failed acknowledgement logic
	// 	_ = dispatchedAck.Error

	// 	return nil
	// case *channeltypes.Acknowledgement_Result:
	// 	// Decode the packet acknowledgment
	// 	var packetAck types.CreatePairPacketAck

	// 	if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
	// 		// The counter-party module doesn't implement the correct acknowledgment format
	// 		return errors.New("cannot unmarshal acknowledgment")
	// 	}

	// 	// TODO: successful acknowledgement logic

	// 	return nil
	// default:
	// 	// The counter-party module doesn't implement the correct acknowledgment format
	// 	return errors.New("invalid acknowledgment format")
	// }

	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.CreatePairPacketAck
		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// Set the sell order book
		pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.SourceDenom, data.TargetDenom)
		book := types.NewSellOrderBook(data.SourceDenom, data.TargetDenom)
		book.Index = pairIndex
		k.SetSellOrderBook(ctx, book)

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutCreatePairPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutCreatePairPacket(ctx sdk.Context, packet channeltypes.Packet, data types.CreatePairPacketData) error {

	// TODO: packet timeout logic

	return nil
}
