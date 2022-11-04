package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sudhakar-mamillapalli/checkers/x/checkers/rules"
	"github.com/sudhakar-mamillapalli/checkers/x/checkers/types"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message

    // 1. Find the game.
    storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
    if !found {
        return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
    }

    // 2. check player legitimate
    // Note This uses the certainty that the MsgPlayMove.Creator has been
    // verified by its signature.  Create is creator if this message, it
    // must match either the black or red player
    isBlack := storedGame.Black == msg.Creator
    isRed := storedGame.Red == msg.Creator
    var player rules.Player
    if !isBlack && !isRed {
        return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
    } else if isBlack && isRed {
        // I think this is the case when both black and red are same player
        // i.e. playing themselves.
        player = rules.StringPieces[storedGame.Turn].Player
    } else if isBlack {
        player = rules.BLACK_PLAYER
    } else {
        player = rules.RED_PLAYER
    }

    // 3. Instantiate the board in order to implement the rules
    // ParseGame in full_game
    game, err := storedGame.ParseGame()
    if err != nil {
        panic(err.Error())
    }

    // 4. s it the player's turn? Check using the rules file's own TurnIs
    if !game.TurnIs(player) {
        return nil, sdkerrors.Wrapf(types.ErrNotPlayerTurn, "%s", player)
    }

    //5. Do the move
    captured, moveErr := game.Move(
        rules.Pos{
            X: int(msg.FromX),
            Y: int(msg.FromY),
        },
        rules.Pos{
            X: int(msg.ToX),
            Y: int(msg.ToY),
        },
    )
    if moveErr != nil {
        return nil, sdkerrors.Wrapf(types.ErrWrongMove, moveErr.Error())
    }

    //6. Stored updated board
    storedGame.Board = game.String()
    storedGame.Turn = rules.PieceStrings[game.Turn]
    k.Keeper.SetStoredGame(ctx, storedGame)

    //7. Return move result. The Captured and Winner information would be lost
    //if you did not get it out of the function one way or another. More
    //accurately, one would have to replay the transaction to discover the
    //values. It is best to make this information easily accessible.

    return &types.MsgPlayMoveResponse{
        CapturedX: int32(captured.X),
        CapturedY: int32(captured.Y),
        Winner:    rules.PieceStrings[game.Winner()],
    }, nil








	return &types.MsgPlayMoveResponse{}, nil
}
