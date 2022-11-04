package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sudhakar-mamillapalli/checkers/x/checkers/types"
)

func (k msgServer) RejectGame(goCtx context.Context, msg *types.MsgRejectGame) (*types.MsgRejectGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message

    // Note player can't reject  a game once they begin to play When loading a
    // StoredGame from storage you have no way of knowing whether a player
    // already played or not. To access this information add a new field to the
    // StoredGame called MoveCount. In proto/checkers/stored_game.proto:

    storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
    if !found {
        return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
    }

    // Note black plays first
    if storedGame.Black == msg.Creator {
        if 0 < storedGame.MoveCount { // Notice the use of the new field
            return nil, types.ErrBlackAlreadyPlayed
        }
    } else if storedGame.Red == msg.Creator {
        if 1 < storedGame.MoveCount { // Notice the use of the new field
            return nil, types.ErrRedAlreadyPlayed
        }
    } else {
        return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
    }

    // Remove the game using the Keeper.RemoveStoredGame
    k.Keeper.RemoveStoredGame(ctx, msg.GameIndex)

    // emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(types.GameRejectedEventType,
        sdk.NewAttribute(types.GameRejectedEventCreator, msg.Creator),
        sdk.NewAttribute(types.GameRejectedEventGameIndex, msg.GameIndex),
        ),
    )


	return &types.MsgRejectGameResponse{}, nil
}
