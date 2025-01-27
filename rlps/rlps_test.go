package rlps

import (
	"encoding/hex"
	"net/http/httptest"
	"testing"

	"github.com/indexsupply/x/e2pg"
	"github.com/indexsupply/x/geth"
	"github.com/indexsupply/x/geth/gethtest"
	"kr.dev/diff"
)

func check(tb testing.TB, err error) {
	if err != nil {
		tb.Fatal(err)
	}
}

func h2b(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func TestHash(t *testing.T) {
	gtest := gethtest.New(t, "http://zeus:8545")
	defer gtest.Done()

	var (
		srv = NewServer(gtest.FileCache, gtest.Client)
		ts  = httptest.NewServer(srv)
		cli = NewClient(ts.URL)
	)
	defer ts.Close()

	h, err := cli.Hash(16000000)
	diff.Test(t, t.Fatalf, err, nil)
	diff.Test(t, t.Fatalf, hex.EncodeToString(h), "3dc4ef568ae2635db1419c5fec55c4a9322c05302ae527cd40bff380c1d465dd")
}

func TestLatest(t *testing.T) {
	gtest := gethtest.New(t, "http://zeus:8545")
	gtest.SetLatest(16000000, h2b("3dc4ef568ae2635db1419c5fec55c4a9322c05302ae527cd40bff380c1d465dd"))
	defer gtest.Done()

	var (
		srv = NewServer(gtest.FileCache, gtest.Client)
		ts  = httptest.NewServer(srv)
		cli = NewClient(ts.URL)
	)
	defer ts.Close()

	n, h, err := cli.Latest()
	diff.Test(t, t.Fatalf, err, nil)
	diff.Test(t, t.Fatalf, n, uint64(16000000))
	diff.Test(t, t.Fatalf, hex.EncodeToString(h), "3dc4ef568ae2635db1419c5fec55c4a9322c05302ae527cd40bff380c1d465dd")
}

func TestLoadBlocks(t *testing.T) {
	gtest := gethtest.New(t, "http://zeus:8545")
	gtest.SetLatest(16000000, h2b("3dc4ef568ae2635db1419c5fec55c4a9322c05302ae527cd40bff380c1d465dd"))
	defer gtest.Done()

	var (
		srv = NewServer(gtest.FileCache, gtest.Client)
		ts  = httptest.NewServer(srv)
		cli = NewClient(ts.URL)
	)
	defer ts.Close()

	var (
		buffers = []geth.Buffer{geth.Buffer{Number: 16000000}}
		blocks  = make([]e2pg.Block, 1)
	)
	err := cli.LoadBlocks(nil, buffers, blocks)
	diff.Test(t, t.Fatalf, err, nil)
	diff.Test(t, t.Fatalf, blocks[0].Num(), uint64(16000000))
	diff.Test(t, t.Fatalf, blocks[0].Hash(), h2b("3dc4ef568ae2635db1419c5fec55c4a9322c05302ae527cd40bff380c1d465dd"))
	diff.Test(t, t.Fatalf, blocks[0].Transactions.Len(), 211)
}

func repeat(b byte, n int) []byte {
	var res = make([]byte, n)
	for i := range res {
		res[i] = b
	}
	return res
}

func TestLoadBlocks_Filter(t *testing.T) {
	gtest := gethtest.New(t, "http://zeus:8545")
	gtest.SetLatest(16000000, h2b("3dc4ef568ae2635db1419c5fec55c4a9322c05302ae527cd40bff380c1d465dd"))
	defer gtest.Done()

	var (
		srv = NewServer(gtest.FileCache, gtest.Client)
		ts  = httptest.NewServer(srv)
		cli = NewClient(ts.URL)
	)
	defer ts.Close()

	var (
		buffers = []geth.Buffer{geth.Buffer{Number: 16000000}}
		blocks  = make([]e2pg.Block, 1)
		filter  = [][]byte{repeat('2', 32)}
	)
	err := cli.LoadBlocks(filter, buffers, blocks)
	diff.Test(t, t.Fatalf, err, nil)
	diff.Test(t, t.Fatalf, blocks[0].Num(), uint64(16000000))
	diff.Test(t, t.Fatalf, blocks[0].Hash(), h2b("3dc4ef568ae2635db1419c5fec55c4a9322c05302ae527cd40bff380c1d465dd"))
	diff.Test(t, t.Fatalf, blocks[0].Transactions.Len(), 0)
}

func TestLoadBlocks_Filter_Error(t *testing.T) {
	gtest := gethtest.New(t, "http://zeus:8545")
	gtest.SetLatest(16000000, h2b("3dc4ef568ae2635db1419c5fec55c4a9322c05302ae527cd40bff380c1d465dd"))
	defer gtest.Done()

	var (
		srv = NewServer(gtest.FileCache, gtest.Client)
		ts  = httptest.NewServer(srv)
		cli = NewClient(ts.URL)
	)
	defer ts.Close()

	var (
		buffers = []geth.Buffer{geth.Buffer{Number: 16000000}}
		blocks  = make([]e2pg.Block, 1)
		filter  = [][]byte{repeat('2', 33)}
	)
	err := cli.LoadBlocks(filter, buffers, blocks)
	const want = "rlps error: filter item must be 32 bytes. got: 33\n"
	diff.Test(t, t.Fatalf, err.Error(), want)
}
