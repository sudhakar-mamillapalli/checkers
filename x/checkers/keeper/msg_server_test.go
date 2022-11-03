package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/sudhakar-mamillapalli/checkers/testutil/keeper"
	"github.com/sudhakar-mamillapalli/checkers/x/checkers/keeper"
	"github.com/sudhakar-mamillapalli/checkers/x/checkers/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
