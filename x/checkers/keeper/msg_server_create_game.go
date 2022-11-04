package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sudhakar-mamillapalli/checkers/x/checkers/rules"
	"github.com/sudhakar-mamillapalli/checkers/x/checkers/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message

	// 1. Get the current Id
	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}
	newIndex := strconv.FormatUint(systemInfo.NextId, 10)

	// Create the object to be stored
	newGame := rules.New()
	storedGame := types.StoredGame{
		Index: newIndex,
		Board: newGame.String(),
		Turn:  rules.PieceStrings[newGame.Turn], // convert Player to string
		Black: msg.Black,
		Red:   msg.Red,
        MoveCount: 0,
	}

	// This is helper function we wrote in full_game.go
	err := storedGame.Validate()
	if err != nil {
		// panic instead of error since players cannot stall the blockchain.
		// They can spam but they will have to pay gas fee to to this point
		return nil, err
	}

	/*
	 * .Red, and .Black need to be checked because they were copied as strings.
	 * You do not need to check .Creator because at this stage the message's
	 * signatures have been verified, and the creator is the signer.
	 */

	// save the StoredGame object. keeper/msg_server.go has MsgServer def
	k.Keeper.SetStoredGame(ctx, storedGame)

	// update systemInfo
	systemInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.GameCreatedEventType,
			sdk.NewAttribute(types.GameCreatedEventCreator, msg.Creator),
			sdk.NewAttribute(types.GameCreatedEventGameIndex, newIndex),
			sdk.NewAttribute(types.GameCreatedEventBlack, msg.Black),
			sdk.NewAttribute(types.GameCreatedEventRed, msg.Red),
		),
	)

	return &types.MsgCreateGameResponse{
		GameIndex: newIndex,
	}, nil
}
