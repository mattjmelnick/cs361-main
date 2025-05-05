package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// COMMAND TEXTS
	mainCommandsText := (`COMMANDS
	portfolio       View portfolio screen
	search-stocks   Search for stocks
	add-stocks      Add stocks to portfolio
	search-crypto   Search for cryptocurrencies
	add-crypto      Add cryptocurrencies to portfolio
	add-other       Add other income to portfolio
	quit            Quit the application`)

	portfolioCommandsText := (`COMMANDS
	main					Go to main screen
	search-stocks   		Search for stocks
	add-stocks      		Add stocks to portfolio
	search-crypto   		Search for cryptocurrencies
	add-crypto      		Add cryptocurrencies to portfolio
	add-other       		Add other income to portfolio
	edit <name> <number>	Edit stock, crypto, or income
	delete <name>			Delete stock, crypto, or income
	quit            		Quit the application`)

	searchStocksCommandsText := (`COMMANDS
	search $TICKER		Search for company
	search <company>	Search for company by name
	show-more			Show additional price details
	add-stocks      	Add stocks to portfolio
	portfolio			Go to portfolio screen
	main				Go to main screen
	quit            	Quit the application`)

	addStocksCommandsText := (`COMMANDS
	add $TICKER <number>	Add <number> shares to your portfolio
	search-stocks			Go to search stocks screen
	portfolio				Go to portfolio screen
	main					Go to main screen
	quit            		Quit the application`)

	searchCryptoCommandsText := (`COMMANDS
	search $COIN		Search for cryptocurrency
	add-crypto		Go to add cryptocurrency screen
	portfolio			Go to portfolio screen
	main				Go to main screen
	quit            	Quit the application`)

	addCryptoCommandsText := (`COMMANDS
	add $COIN <number>	Add <number> coins to your portfolio
	search-crypto		Go to search cryptocurrencies screen
	portfolio			Go to portfolio screen
	main				Go to main screen
	quit            	Quit the application`)

	addIncomeCommandsText := (`COMMANDS
	add <name> <number>			Add income description and amount
	portfolio					Go to portfolio screen
	main						Go to main screen
	quit						Quit the appliation`)

	// MAIN PAGE
	mainTitle := tview.NewTextView().
		SetText("Go Finance TUI")

	mainDescription := tview.NewTextView().
		SetText(`Track your investments and savings
in one convenient location
		
Search and add stocks, cryptocurrencies,
and other sources of income`)

	mainCommands := tview.NewTextView().SetText(mainCommandsText)

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

	portfolioCommands := tview.NewTextView().SetText(portfolioCommandsText)

	portfolioInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	portfolioEditMessage := tview.NewTextView().SetText("")

	var portfolio = make(map[string]int)
	grandTotal := 0

	portfolioStockTableTitle := tview.NewTextView().SetText("STOCKS")
	portfolioStockTableTotal := tview.NewTextView().SetText(fmt.Sprintf("Total: %d", grandTotal))
	portfolioStockTable := tview.NewTable().SetBorders(true)

	updatePortfolioStockTable := func() {
		portfolioStockTable.Clear()

		// Set header
		portfolioStockTable.SetCell(0, 0, tview.NewTableCell("Ticker").
			SetAlign(tview.AlignCenter).SetSelectable(false))
		portfolioStockTable.SetCell(0, 1, tview.NewTableCell("Quantity").
			SetAlign(tview.AlignCenter).SetSelectable(false))
		portfolioStockTable.SetCell(0, 2, tview.NewTableCell("Price").
			SetAlign(tview.AlignCenter).SetSelectable(false))
		portfolioStockTable.SetCell(0, 3, tview.NewTableCell("Total").
			SetAlign(tview.AlignCenter).SetSelectable(false))

		row := 1
		grandTotal = 0

		// TODO: CHANGE TO DYNAMIC PRICES FROM MICROSERVICE
		for ticker, qty := range portfolio {
			total := 50 * qty
			portfolioStockTable.SetCell(row, 0, tview.NewTableCell(ticker))
			portfolioStockTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d", qty)).
				SetAlign(tview.AlignRight))
			portfolioStockTable.SetCell(row, 2, tview.NewTableCell("50").
				SetAlign(tview.AlignRight))
			portfolioStockTable.SetCell(row, 3, tview.NewTableCell(strconv.Itoa(total)).
				SetAlign(tview.AlignRight))
			grandTotal += total
			row++
		}
		portfolioStockTableTotal.SetText(fmt.Sprintf("Total: %d", grandTotal))
	}

	portfolioLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(portfolioTitle, 3, 1, false).
		AddItem(portfolioDescription, 5, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(portfolioCommands, 11, 1, false).
		AddItem(portfolioInput, 2, 1, true).
		AddItem(portfolioEditMessage, 2, 1, false).
		AddItem(portfolioStockTableTitle, 1, 0, false).
		AddItem(portfolioStockTableTotal, 1, 0, false).
		AddItem(portfolioStockTable, 0, 4, false)

	// SEARCH STOCKS PAGE
	searchStocksTitle := tview.NewTextView().
		SetText("Search Stocks")

	searchStocksDescription := tview.NewTextView().
		SetText(`Search for stocks using the ticker symbol or company name

Example: search $AAPL
		 search Apple`)

	searchStocksCommands := tview.NewTextView().SetText(searchStocksCommandsText)

	searchStocksInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	searchStocksTable := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 0)

	showMore := false

	searchStocksLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchStocksTitle, 3, 1, false).
		AddItem(searchStocksDescription, 5, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(searchStocksCommands, 9, 1, false).
		AddItem(searchStocksInput, 1, 1, true).
		AddItem(searchStocksTable, 0, 1, false)

	renderStocksTable := func() {
		searchStocksTable.Clear()

		// Full data
		// TODO: USE API VALUES FROM MICROSERVICE
		fullData := []string{"Apple Inc.", "AAPL", "$170.50", "$168.00", "$172.00"}
		headers := []string{"Company", "Ticker", "Price"}
		data := []string{fullData[0], fullData[1], fullData[2]}

		if showMore {
			headers = []string{"Company", "Ticker", "Price", "Open", "Close"}
			data = fullData
		}

		// Set headers
		for col, val := range headers {
			cell := tview.NewTableCell(val).
				SetTextColor(tcell.ColorYellow).
				SetSelectable(false).
				SetAlign(tview.AlignCenter).
				SetAttributes(tcell.AttrBold)
			searchStocksTable.SetCell(0, col, cell)
		}

		// Set data row
		for col, val := range data {
			cell := tview.NewTableCell(val).
				SetAlign(tview.AlignCenter)
			searchStocksTable.SetCell(1, col, cell)
		}

		searchStocksLayout.ResizeItem(searchStocksTable, 4, 1)
	}

	// ADD STOCKS PAGE
	addStocksTitle := tview.NewTextView().
		SetText("Add Stocks")

	addStocksDescription := tview.NewTextView().
		SetText(`Add stocks to your portfolio
		
Example: add $AAPL 50 -> Add 50 shares of $AAPL`)

	addStocksCommands := tview.NewTextView().SetText(addStocksCommandsText)

	addStocksInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	addStocksMessage := tview.NewTextView().SetText("")

	addStocksLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(addStocksTitle, 3, 1, false).
		AddItem(addStocksDescription, 4, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(addStocksCommands, 9, 1, false).
		AddItem(addStocksInput, 2, 1, true).
		AddItem(addStocksMessage, 2, 1, false)

	// SEARCH CRYPTO PAGE
	searchCryptoTitle := tview.NewTextView().
		SetText("Search Cryptocurrencies")

	searchCryptoDescription := tview.NewTextView().
		SetText(`Search for cryptocurrencies using their abbreviations
		
Example: search $BTC`)

	searchCryptoCommands := tview.NewTextView().SetText(searchCryptoCommandsText)

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

	addCryptoCommands := tview.NewTextView().SetText(addCryptoCommandsText)

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

	addIncomeCommands := tview.NewTextView().SetText(addIncomeCommandsText)

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
				portfolioEditMessage.SetText("")
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
				PromptQuit(app, mainLayout, mainCommands, mainInput, mainCommandsText)
			default:
			}
			mainInput.SetText("")
		}
	})

	var portfolioInputDoneFunc func(tcell.Key)

	portfolioInputDoneFunc = func(key tcell.Key) {
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
				PromptQuit(app, portfolioLayout, portfolioCommands, portfolioInput, portfolioCommandsText)
			default:
				// edit $TICKER <number>
				parts := strings.Split(cmd, " ")
				if len(parts) == 3 && strings.ToLower(parts[0]) == "edit" &&
					strings.HasPrefix(parts[1], "$") {
					ticker := strings.TrimPrefix(parts[1], "$")
					quantity, err := strconv.Atoi(parts[2])
					if err == nil && quantity > 0 {
						if _, exists := portfolio[ticker]; exists {
							portfolio[ticker] = quantity
							portfolioEditMessage.SetText(fmt.Sprintf("Updated $%s, Quantity: %d", ticker, quantity))
							updatePortfolioStockTable()
						} else {
							portfolioEditMessage.SetText(fmt.Sprintf("Ticker $%s not found in portfolio", ticker))
						}
					} else {
						portfolioEditMessage.SetText("Invalid entry")
					}
				} else if len(parts) == 2 && strings.ToLower(parts[0]) == "delete" &&
					strings.HasPrefix(parts[1], "$") {
					ticker := strings.TrimPrefix(parts[1], "$")
					portfolioEditMessage.SetText(fmt.Sprintf("Are you sure you want to delete %s? Type 'y' to confirm", ticker))
					portfolioInput.SetText("")

					portfolioInput.SetDoneFunc(func(key tcell.Key) {
						if key == tcell.KeyEnter {
							ans := strings.TrimSpace(strings.ToLower(portfolioInput.GetText()))
							if ans == "y" {
								delete(portfolio, ticker)
								portfolioEditMessage.SetText(fmt.Sprintf("$%s removed from portfolio", ticker))
								updatePortfolioStockTable()
							} else {
								portfolioEditMessage.SetText("Delete canceled")
							}
							portfolioInput.SetText("")
							portfolioInput.SetDoneFunc(portfolioInputDoneFunc) // Restore original
						}
					})
				} else {
					portfolioEditMessage.SetText("Invalid command or format")
				}
			}
			portfolioInput.SetText("")
		}
	}

	portfolioInput.SetDoneFunc(portfolioInputDoneFunc)

	searchStocksInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := strings.TrimSpace(strings.ToLower(searchStocksInput.GetText()))
			switch cmd {
			case "search $aapl", "search apple":
				showMore = false
				renderStocksTable()
			case "show-more":
				showMore = true
				renderStocksTable()
			case "add-stocks":
				pages.SwitchToPage("addStocks")
				app.SetFocus(addStocksInput)
			case "portfolio":
				portfolioEditMessage.SetText("")
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				PromptQuit(app, searchStocksLayout, searchStocksCommands, searchStocksInput, searchStocksCommandsText)
			default:
			}
			searchStocksInput.SetText("")
		}
	})

	addStocksInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			addStocksMessage.SetText("")
			cmd := strings.TrimSpace(strings.ToLower(addStocksInput.GetText()))
			switch cmd {
			case "search-stocks":
				pages.SwitchToPage("searchStocks")
				app.SetFocus(searchStocksInput)
			case "portfolio":
				portfolioEditMessage.SetText("")
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				PromptQuit(app, addStocksLayout, addStocksCommands, addStocksInput, addStocksCommandsText)
			default:
				// Check for: add $TICKER <number>
				parts := strings.Split(cmd, " ")
				if len(parts) == 3 && strings.ToLower(parts[0]) == "add" &&
					strings.HasPrefix(parts[1], "$") {
					ticker := strings.TrimPrefix(parts[1], "$")
					count, err := strconv.Atoi(parts[2])
					if err == nil && count > 0 {
						portfolio[ticker] += count
						addStocksMessage.SetText(fmt.Sprintf("%d shares of '$%s' successfully added to portfolio", count, ticker))
						updatePortfolioStockTable()
					} else {
						addStocksMessage.SetText("Invalid number of shares")
					}
				} else {
					addStocksMessage.SetText("Invalid command or format. Use: add $TICKER <number>")
				}
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
				portfolioEditMessage.SetText("")
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				PromptQuit(app, searchCryptoLayout, searchCryptoCommands, searchCryptoInput, searchCryptoCommandsText)
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
				portfolioEditMessage.SetText("")
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				PromptQuit(app, addCryptoLayout, addCryptoCommands, addCryptoInput, addCryptoCommandsText)
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
				portfolioEditMessage.SetText("")
				pages.SwitchToPage("portfolio")
				app.SetFocus(portfolioInput)
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				PromptQuit(app, addIncomeLayout, addIncomeCommands, addIncomeInput, addIncomeCommandsText)
			default:
			}
			addIncomeInput.SetText("")
		}
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func PromptQuit(app *tview.Application, layout *tview.Flex, commandsView *tview.TextView, inputField *tview.InputField, originalText string) {
	quitConfirmInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(3)

	commandsView.SetText("Are you sure you want to quit? (y/n)")
	layout.RemoveItem(inputField)
	layout.AddItem(quitConfirmInput, 1, 1, true)
	app.SetFocus(quitConfirmInput)

	quitConfirmInput.SetText("")
	quitConfirmInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			answer := quitConfirmInput.GetText()
			if answer == "y" || answer == "Y" {
				app.Stop()
				return
			}
			// Restore layout
			commandsView.SetText(originalText)
			layout.RemoveItem(quitConfirmInput)
			layout.AddItem(inputField, 1, 1, true)
			app.SetFocus(inputField)
		}
	})
}
