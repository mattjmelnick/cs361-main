package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// MAIN PAGE
	mainTitle := tview.NewTextView().
		SetText("Go Finance TUI")

	mainDescription := tview.NewTextView().
		// SetText("Track your investments and savings\nin one convenient location\n\nSearch and add stocks, cryptocurrencies\nand other sources of income")
		SetText(`Track your investments and savings
in one convenient location
		
Search and add stocks, cryptocurrencies,
and other sources of income`)

	mainCommands := tview.NewTextView().
		SetText(`COMMANDS
	portfolio       View portfolio screen
	search-stocks   Search for stocks
	add-stocks      Add stocks to portfolio
	search-crypto   Search for cryptocurrencies
	add-crypto      Add cryptocurrencies to portfolio
	add-other       Add other income to portfolio
	quit            Quit the application`).
		SetTextAlign(tview.AlignLeft)

	mainInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainTitle, 3, 1, false).
		AddItem(mainDescription, 5, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(mainCommands, 9, 1, false).
		AddItem(mainInput, 1, 1, true)

	// PORTFOLIO PAGE
	portfolioTitle := tview.NewTextView().
		SetText("Portfolio")

	portfolioDescription := tview.NewTextView().
		SetText(`View your stocks, cryptocurrencies, and
other sources of income

Stocks and cryptocurrency totals are
rounded down to 2 decimal places`)

	portfolioCommands := tview.NewTextView().
		SetText(`COMMANDS
	main			Go to main screen
	search-stocks   Search for stocks
	add-stocks      Add stocks to portfolio
	search-crypto   Search for cryptocurrencies
	add-crypto      Add cryptocurrencies to portfolio
	add-other       Add other income to portfolio
	edit <name>		Edit stock, crypto, or income
	delete <name>	Delete stock, crypto, or income
	quit            Quit the application`)

	portfolioInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	portfolioLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(portfolioTitle, 3, 1, false).
		AddItem(portfolioDescription, 5, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(portfolioCommands, 11, 1, false).
		AddItem(portfolioInput, 1, 1, true)

	// SEARCH STOCKS PAGE
	searchStocksTitle := tview.NewTextView().
		SetText("Search Stocks")

	searchStocksDescription := tview.NewTextView().
		SetText(`Search for stocks using the ticker symbol or company name

Example: search $AAPL
		 search Apple`)

	searchStocksCommands := tview.NewTextView().
		SetText(`COMMANDS
	search $TICKER		Search for company
	search <company>	Search for company by name
	show-more			Show additional price details
	add-stocks      	Add stocks to portfolio
	portfolio			Go to portfolio screen
	main				Go to main screen
	quit            	Quit the application`)

	searchStocksInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	searchStocksLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchStocksTitle, 3, 1, false).
		AddItem(searchStocksDescription, 5, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(searchStocksCommands, 9, 1, false).
		AddItem(searchStocksInput, 1, 1, true)

	// ADD STOCKS PAGE
	addStocksTitle := tview.NewTextView().
		SetText("Add Stocks")

	addStocksDescription := tview.NewTextView().
		SetText(`Add stocks to your portfolio
		
Example: add $AAPL 50 -> Add 50 shares of $AAPL`)

	addStocksCommands := tview.NewTextView().
		SetText(`COMMANDS
	add $TICKER <number>	Add <number> shares to your portfolio
	edit $TICKER			Edit stock/quantity
	delete $TICKER			Delete stock
	search-stocks			Go to search stocks screen
	portfolio				Go to portfolio screen
	main					Go to main screen
	quit            		Quit the application`)

	addStocksInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	addStocksLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(addStocksTitle, 3, 1, false).
		AddItem(addStocksDescription, 4, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(addStocksCommands, 9, 1, false).
		AddItem(addStocksInput, 1, 1, true)

	// SEARCH CRYPTO PAGE
	searchCryptoTitle := tview.NewTextView().
		SetText("Search Cryptocurrencies")

	searchCryptoDescription := tview.NewTextView().
		SetText(`Search for cryptocurrencies using their abbreviations
		
Example: search $BTC`)

	searchCryptoCommands := tview.NewTextView().
		SetText(`COMMANDS
	search $COIN		Search for cryptocurrency
	add-crypto		Go to add cryptocurrency screen
	portfolio			Go to portfolio screen
	main				Go to main screen
	quit            	Quit the application`)

	searchCryptoInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	searchCryptoLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchCryptoTitle, 3, 1, false).
		AddItem(searchCryptoDescription, 5, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(searchCryptoCommands, 7, 1, false).
		AddItem(searchCryptoInput, 1, 1, true)

	// ADD CRYPTO PAGE
	addCryptoTitle := tview.NewTextView().
		SetText("Add Cryptocurrencies")

	addCryptoDescription := tview.NewTextView().
		SetText(`Add cryptocurrencies to your portfolio
		
Example: add $BTC 0.0001 -> Add 0.0001 worth of $BTC to your portfolio`)

	addCryptoCommands := tview.NewTextView().
		SetText(`COMMANDS
	add $COIN <number>	Add <number> coins to your portfolio
	edit $COIN			Edit cryptocurrency/quantity
	delete $COIN		Delete cryptocurrency 
	search-crypto		Go to search cryptocurrencies screen
	portfolio			Go to portfolio screen
	main				Go to main screen
	quit            	Quit the application`)

	addCryptoInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	addCryptoLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(addCryptoTitle, 3, 1, false).
		AddItem(addCryptoDescription, 4, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(addCryptoCommands, 9, 1, false).
		AddItem(addCryptoInput, 1, 1, true)

	// ADD OTHER INCOME PAGE
	addIncomeTitle := tview.NewTextView().
		SetText("Add Other Income")

	addIncomeDescription := tview.NewTextView().
		SetText(`Add other sources of income to your portfolio
		
Example: add checkings 1000`)

	addIncomeCommands := tview.NewTextView().
		SetText(`COMMANDS
	add <name> <number>			Add income description and amount
	edit <description>			Edit income source
	delete <description>		Delete income source
	portfolio					Go to portfolio screen
	main						Go to main screen
	quit						Quit the appliation`)

	addIncomeInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	addIncomeLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(addIncomeTitle, 3, 1, false).
		AddItem(addIncomeDescription, 4, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(addIncomeCommands, 8, 1, false).
		AddItem(addIncomeInput, 1, 1, true)

	// PAGE ROUTES
	pages := tview.NewPages().
		AddPage("main", mainLayout, true, true).
		AddPage("portfolio", portfolioLayout, true, false).
		AddPage("searchStocks", searchStocksLayout, true, false).
		AddPage("addStocks", addStocksLayout, true, false).
		AddPage("searchCrypto", searchCryptoLayout, true, false).
		AddPage("addCrypto", addCryptoLayout, true, false).
		AddPage("addIncome", addIncomeLayout, true, false)

	// INPUTS
	mainInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := mainInput.GetText()
			switch cmd {
			case "portfolio":
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "search-stocks":
				pages.SwitchToPage("searchStocks")
				app.SetFocus(searchStocksInput)
			case "add-stocks":
				pages.SwitchToPage("addStocks")
				app.SetFocus(addStocksInput)
			case "search-crypto":
				pages.SwitchToPage("searchCrypto")
				app.SetFocus(searchCryptoInput)
			case "add-crypto":
				pages.SwitchToPage("addCrypto")
				app.SetFocus(addCryptoInput)
			case "add-other":
				pages.SwitchToPage("addIncome")
				app.SetFocus(addIncomeInput)
			case "quit":
				app.Stop()
			default:
			}
			mainInput.SetText("")
		}
	})

	portfolioInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := portfolioInput.GetText()
			switch cmd {
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "search-stocks":
				pages.SwitchToPage("searchStocks")
				app.SetFocus(searchStocksInput)
			case "add-stocks":
				pages.SwitchToPage("addStocks")
				app.SetFocus(addStocksInput)
			case "search-crypto":
				pages.SwitchToPage("searchCrypto")
				app.SetFocus(searchCryptoInput)
			case "add-crypto":
				pages.SwitchToPage("addCrypto")
				app.SetFocus(addCryptoInput)
			case "add-other":
				pages.SwitchToPage("addIncome")
				app.SetFocus(addIncomeInput)
			case "quit":
				app.Stop()
			default:
			}
			portfolioInput.SetText("")
		}
	})

	searchStocksInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := searchStocksInput.GetText()
			switch cmd {
			case "add-stocks":
				pages.SwitchToPage("addStocks")
				app.SetFocus(addStocksInput)
			case "portfolio":
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				app.Stop()
			default:
			}
			searchStocksInput.SetText("")
		}
	})

	addStocksInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := addStocksInput.GetText()
			switch cmd {
			case "search-stocks":
				pages.SwitchToPage("searchStocks")
				app.SetFocus(searchStocksInput)
			case "portfolio":
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				app.Stop()
			default:
			}
			addStocksInput.SetText("")
		}
	})

	searchCryptoInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := searchCryptoInput.GetText()
			switch cmd {
			case "add-crypto":
				pages.SwitchToPage("addCrypto")
				app.SetFocus(addCryptoInput)
			case "portfolio":
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				app.Stop()
			default:
			}
			searchCryptoInput.SetText("")
		}
	})

	addCryptoInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := addCryptoInput.GetText()
			switch cmd {
			case "search-crypto":
				pages.SwitchToPage("searchCrypto")
				app.SetFocus(searchCryptoInput)
			case "portfolio":
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				app.Stop()
			default:
			}
			addCryptoInput.SetText("")
		}
	})

	addIncomeInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := addIncomeInput.GetText()
			switch cmd {
			case "portfolio":
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				app.Stop()
			default:
			}
			addIncomeInput.SetText("")
		}
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
