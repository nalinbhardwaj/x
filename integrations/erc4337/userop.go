package erc4337

import (
	"context"

	"github.com/indexsupply/x/contrib/erc4337"
	"github.com/indexsupply/x/e2pg"

	"github.com/jackc/pgx/v5"
)

type integration struct {
	name string
}

var Integration = integration{
	name: "ERC4337 UserOperationEvent",
}

func (i integration) Events(ctx context.Context) [][]byte {
	return [][]byte{erc4337.UserOperationEventSignatureHash}
}

func (i integration) Delete(ctx context.Context, pg e2pg.PG, h []byte) error {
	const q = `
		delete from erc4337_userops
		where task_id = $1
		and chain_id = $2
		and block_hash = $3
	`
	_, err := pg.Exec(ctx, q, e2pg.TaskID(ctx), e2pg.ChainID(ctx), h)
	return err
}

func (i integration) Insert(ctx context.Context, pg e2pg.PG, blocks []e2pg.Block) (int64, error) {
	var rows = make([][]any, 0, 1<<12) // 2^12 max batch size?
	for bidx := 0; bidx < len(blocks); bidx++ {
		for ridx := 0; ridx < blocks[bidx].Receipts.Len(); ridx++ {
			r := blocks[bidx].Receipts.At(ridx)
			for lidx := 0; lidx < r.Logs.Len(); lidx++ {
				l := r.Logs.At(lidx)
				xfr, err := erc4337.MatchUserOperationEvent(l)
				if err != nil {
					continue
				}
				rows = append(rows, []any{
					e2pg.TaskID(ctx),
					e2pg.ChainID(ctx),
					blocks[bidx].Num(),
					blocks[bidx].Hash(),
					blocks[bidx].Transactions.At(ridx).Hash(),
					ridx,
					lidx,
					l.Address,

					xfr.UserOpHash[:],
					xfr.Sender[:],
					xfr.Paymaster[:],
					xfr.Nonce.String(),
					xfr.Success,
					xfr.ActualGasCost.String(),
					xfr.ActualGasUsed.String(),
				})
				xfr.Done()
			}
		}
	}
	return pg.CopyFrom(
		context.Background(),
		pgx.Identifier{"erc4337_userops"},
		[]string{
			"task_id",
			"chain_id",
			"block_number",
			"block_hash",
			"transaction_hash",
			"transaction_index",
			"log_index",
			"contract",
			
			"userOpHash",
			"userOpSender",
			"userOpPaymaster",
			"userOpNonce",
			"userOpSuccess",
			"userOpActualGasCost",
			"userOpActualGasUsed",
		},
		pgx.CopyFromRows(rows),
	)
}
