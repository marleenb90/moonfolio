package runner

import (
	"log"

	"github.com/tebben/moonfolio/ui"

	"strconv"

	"github.com/jroimartin/gocui"
	"github.com/tebben/moonfolio/coindata"
	"github.com/tebben/moonfolio/configuration"
	"github.com/tebben/moonfolio/cryptocompare"
	"github.com/tebben/moonfolio/transactions"
)

var (
	store     transactions.IStore
	allCoins  *cryptocompare.CoinList
	userCoins map[string]*coindata.CoinData
	config    *configuration.Config
	gui       *gocui.Gui
)

// Start starts the runner
func Start(c *configuration.Config, g *gocui.Gui) {
	config = c
	gui = g

	createStore()
	getCoinList()
	createCoinData()
	go startUpdater(int64(c.IntervalSeconds * 1000))
	ui.ReDraw()
}

func createStore() {
	var err error
	store, err = transactions.NewBoltStore(config.TransactionsFile)
	if err != nil {
		log.Fatalf("Unable to open or create transaction store: %v", err)
	}
}

func getCoinList() {
	c, err := cryptocompare.GetCoinList()
	if err != nil {
		log.Fatal("Error getting coin list")
		// ToDo set error in interface
	}

	allCoins = c
}

func createCoinData() {
	ta, err := store.GetTransactions()
	if err != nil {
		log.Fatalf("unable to get transactions: %v", err)
	}

	userCoins = make(map[string]*coindata.CoinData, 0)

	// add transaction to a user coin, create user coin if there is no coin yet
	for i := range ta {
		t := ta[i]
		if uc, ok := userCoins[t.CoinID]; !ok {
			uc = &coindata.CoinData{}

			// not created, find coin in allCoins
			if ccCoin, ok2 := allCoins.Data[t.CoinID]; ok2 {
				uc.SetName(ccCoin.FullName)
				uc.SetSymbol(ccCoin.Symbol)
				rank, _ := strconv.Atoi(ccCoin.SortOrder)
				uc.SetRank(rank)
			} else { // coins from cryptocompare did not find the coin referred in the transaction
				uc.SetName(t.CoinID)
				uc.SetSymbol(t.CoinID)
			}

			uc.AddTransaction(t)
			userCoins[uc.Symbol] = uc
		} else {
			userCoins[t.CoinID].AddTransaction(t)
		}
	}
}

// set error in GUI
func setError() {

}
