package cli

import (
	"fmt"
	"strconv"

	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
	"gopkg.in/urfave/cli.v2"

	"github.com/filecoin-project/go-lotus/chain/address"
	"github.com/filecoin-project/go-lotus/chain/types"
)

var clientCmd = &cli.Command{
	Name:  "client",
	Usage: "Make deals, store data, retrieve data",
	Subcommands: []*cli.Command{
		clientImportCmd,
		clientLocalCmd,
		clientDealCmd,
	},
}

var clientImportCmd = &cli.Command{
	Name:  "import",
	Usage: "Import data",
	Action: func(cctx *cli.Context) error {
		api, err := GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		ctx := ReqContext(cctx)

		c, err := api.ClientImport(ctx, cctx.Args().First())
		if err != nil {
			return err
		}
		fmt.Println(c.String())
		return nil
	},
}

var clientLocalCmd = &cli.Command{
	Name:  "local",
	Usage: "List locally imported data",
	Action: func(cctx *cli.Context) error {
		api, err := GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		ctx := ReqContext(cctx)

		list, err := api.ClientListImports(ctx)
		if err != nil {
			return err
		}
		for _, v := range list {
			fmt.Printf("%s %s %d %s\n", v.Key, v.FilePath, v.Size, v.Status)
		}
		return nil
	},
}

var clientDealCmd = &cli.Command{
	Name:  "deal",
	Usage: "Initialize storage deal with a miner",
	Action: func(cctx *cli.Context) error {
		api, err := GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		ctx := ReqContext(cctx)

		if cctx.NArg() != 4 {
			return xerrors.New("expected 4 args: dataCid, miner, price, duration")
		}

		// [data, miner, dur]

		data, err := cid.Parse(cctx.Args().Get(0))
		if err != nil {
			return err
		}

		miner, err := address.NewFromString(cctx.Args().Get(1))
		if err != nil {
			return err
		}

		// TODO: parse bigint
		price, err := strconv.ParseInt(cctx.Args().Get(2), 10, 32)
		if err != nil {
			return err
		}

		dur, err := strconv.ParseInt(cctx.Args().Get(3), 10, 32)
		if err != nil {
			return err
		}

		proposal, err := api.ClientStartDeal(ctx, data, miner, types.NewInt(uint64(price)), uint64(dur))
		if err != nil {
			return err
		}

		fmt.Println(proposal)
		return nil
	},
}