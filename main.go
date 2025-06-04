package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type IndexData struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
	Name   string
	Ticker string
}

type StockData struct {
	Date   string
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Ticker string
}

type BudgetCategory struct {
	Name       string
	Percentage int
}

var (
	totalBudget         int
	remainingPercentage = 100
	budgetCategories    = make(map[string]int)
)

func main() {
	app := tview.NewApplication()

	// COMMAND TEXTS
	mainCommandsText := (`COMMANDS
	summary			Get a summary of the three major stock indices
	budget			Enter a budget
	search-stocks   Search for stocks
	search-crypto   Search for cryptocurrencies
	quit            Quit the application`)

	summaryCommandsText := (`COMMANDS
	main		Go to main screen
	quit		Quit the application`)

	budgetCommandsText := (`COMMANDS
	main		Go to main screen
	quit		Quit the application`)

	searchStocksCommandsText := (`COMMANDS
	search $TICKER		Search for company
	show-more			Show additional price details
	main				Go to main screen
	quit            	Quit the application`)

	searchCryptoCommandsText := (`COMMANDS
	search COIN			Search for cryptocurrency
	main				Go to main screen
	quit            	Quit the application`)

	// MAIN PAGE
	mainTitle := tview.NewTextView().
		SetText("Go Finance TUI")

	mainDescription := tview.NewTextView().
		SetText("Set a budget and search prices of stocks and cryptocurrencies")

	mainCommands := tview.NewTextView().SetText(mainCommandsText)

	mainInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	mainLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainTitle, 3, 1, false).
		AddItem(mainDescription, 3, 1, false).
		AddItem(mainCommands, 8, 1, false).
		AddItem(mainInput, 1, 1, true)

	// SUMMARY PAGE
	summaryWaiting := tview.NewTextView().SetText("Waiting for data...")

	summaryTitle := tview.NewTextView().
		SetText("Major Stock Indices")

	summaryDescription := tview.NewTextView().
		SetText(`Summary of the Dow Jones, S&P 500, and NASDAQ stock indices`)

	summaryCommands := tview.NewTextView().SetText(summaryCommandsText)

	summaryInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	indicesTable := tview.NewTable().SetBorders(true)

	summaryLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(summaryTitle, 3, 1, false).
		AddItem(summaryDescription, 2, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(summaryCommands, 5, 1, false).
		AddItem(summaryInput, 1, 1, true).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(summaryWaiting, 1, 1, false).
		AddItem(indicesTable, 0, 1, true)

	// BUDGET PAGE
	budgetTitle := tview.NewTextView().
		SetText("Budget Calculator")

	budgetDescription := tview.NewTextView().
		SetText(`Enter budget total, categories, and percentages
		
		The calculated values of the total budget will then be displayed`)

	budgetCommands := tview.NewTextView().SetText(budgetCommandsText)

	budgetMainInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	budgetInput := tview.NewInputField().SetLabel("Enter Total Budget: ").SetFieldWidth(30)
	categoryInput := tview.NewInputField().SetLabel("Category Name: ").SetFieldWidth(30)
	percentageInput := tview.NewInputField().SetLabel("Percentage: ").SetFieldWidth(30)
	budgetMessage := tview.NewTextView().SetText("100% remaining to allocate.")
	categoryTable := tview.NewTable().SetBorders(true)
	budgetTable := tview.NewTable().SetBorders(true)

	budgetInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			budgetStr := strings.TrimSpace(budgetInput.GetText())
			budgetVal, err := strconv.Atoi(budgetStr)
			if err != nil || budgetVal <= 0 {
				budgetMessage.SetText("Please enter a positive number.")
				return
			}
			totalBudget = budgetVal
			budgetMessage.SetText(fmt.Sprintf("Success! 100%% remaining to allocate."))
			app.SetFocus(categoryInput)
		}
	})

	categoryInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			category := strings.TrimSpace(categoryInput.GetText())
			if _, exists := budgetCategories[category]; exists {
				budgetMessage.SetText("Category already exists.")
				return
			}
			app.SetFocus(percentageInput)
		}
	})

	percentageInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			percentageStr := strings.TrimSpace(percentageInput.GetText())
			percentage, err := strconv.Atoi(percentageStr)
			if err != nil || percentage <= 0 || percentage > remainingPercentage {
				budgetMessage.SetText(fmt.Sprintf("Enter value between 1 and %d", remainingPercentage))
				return
			}

			category := strings.TrimSpace(categoryInput.GetText())
			budgetCategories[category] = percentage
			remainingPercentage -= percentage
			budgetMessage.SetText(fmt.Sprintf("Success! Remaining: %d%%.", remainingPercentage))

			// Update the table
			row := len(budgetCategories)
			categoryTable.SetCell(row, 0, tview.NewTableCell(category))
			categoryTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d%%", percentage)))

			// Clear inputs
			categoryInput.SetText("")
			percentageInput.SetText("")

			if remainingPercentage == 0 {
				budgetMessage.SetText("Budget allocated: Writing to file...")
				err = saveBudgetToFile("../sprint3/microservice-a/input.json", totalBudget, budgetCategories)
				if err != nil {
					budgetMessage.SetText("Failed to write file.")
				} else {
					budgetMessage.SetText("Budget saved successfully.")
				}
			}

			app.SetFocus(categoryInput)

			go waitForBudgetOutput(app, "../sprint3/microservice-a/output.json", budgetTable)
		}
	})

	budgetLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(budgetTitle, 3, 1, false).
		AddItem(budgetDescription, 2, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(budgetCommands, 5, 1, false).
		AddItem(budgetMainInput, 3, 0, false).
		AddItem(budgetInput, 2, 0, true).
		AddItem(categoryInput, 2, 0, false).
		AddItem(percentageInput, 2, 0, false).
		AddItem(budgetMessage, 2, 0, false).
		AddItem(categoryTable, 0, 1, false).
		AddItem(budgetTable, 0, 1, false)

	// SEARCH STOCKS PAGE
	searchStocksWaiting := tview.NewTextView().
		SetText("Enter stock ticker")

	searchStocksTitle := tview.NewTextView().
		SetText("Search Stocks")

	searchStocksDescription := tview.NewTextView().
		SetText(`Search for stocks using the ticker symbol

Example: search $AAPL`)

	searchStocksCommands := tview.NewTextView().SetText(searchStocksCommandsText)

	searchStocksInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	searchStocksTable := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 0)

	var showMore bool = false

	searchStocksLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchStocksTitle, 3, 1, false).
		AddItem(searchStocksDescription, 5, 1, false).
		AddItem(searchStocksCommands, 7, 1, false).
		AddItem(searchStocksInput, 1, 1, true).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(searchStocksWaiting, 1, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(searchStocksTable, 0, 1, false)

	// SEARCH CRYPTO PAGE
	searchCryptoWaiting := tview.NewTextView().
		SetText("Waiting for crytpocurrency...")

	searchCryptoTitle := tview.NewTextView().
		SetText("Search Cryptocurrencies")

	searchCryptoDescription := tview.NewTextView().
		SetText(`Search for cryptocurrencies using their name
		
Example: search bitcoin`)

	searchCryptoCommands := tview.NewTextView().SetText(searchCryptoCommandsText)

	searchCryptoInput := tview.NewInputField().
		SetLabel("→ ").
		SetFieldWidth(30)

	searchCryptoTable := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 0)

	searchCryptoLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchCryptoTitle, 3, 1, false).
		AddItem(searchCryptoDescription, 5, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(searchCryptoCommands, 7, 1, false).
		AddItem(searchCryptoInput, 1, 1, true).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(searchCryptoWaiting, 1, 1, false).
		AddItem(tview.NewTextView().SetText(""), 1, 0, false).
		AddItem(searchCryptoTable, 0, 1, false)

	// PAGE ROUTES
	pages := tview.NewPages().
		AddPage("main", mainLayout, true, true).
		AddPage("summary", summaryLayout, true, false).
		AddPage("budget", budgetLayout, true, false).
		AddPage("searchStocks", searchStocksLayout, true, false).
		AddPage("searchCrypto", searchCryptoLayout, true, false)

	// INPUTS
	mainInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := mainInput.GetText()
			switch cmd {
			case "summary":
				os.Remove("../sprint3/microservice-b/output_summary.json")
				err := writeSummaryInput("../sprint3/microservice-b/input_summary.json")
				if err != nil {
					log.Printf("Error writing summary trigger: %v", err)
				} else {
					summaryWaiting.SetText("Waiting for data...")
					indicesTable.Clear()
					summaryLayout.RemoveItem(indicesTable)

					waitForSummaryData(app, "../sprint3/microservice-b/output_summary.json", func(indices []IndexData) {
						summaryLayout.AddItem(tview.NewTextView().SetText("Market Summary"), 1, 1, false)
						summaryLayout.AddItem(indicesTable, 0, 1, false)
						renderSummaryTable(indicesTable, indices)
						summaryWaiting.SetText("")
					})
				}
				pages.SwitchToPage("summary")
				app.SetFocus(summaryInput)
			case "budget":
				err := os.Remove("../sprint3/microservice-a/output.json")
				if err != nil {
				}
				pages.SwitchToPage("budget")
			case "search-stocks":
				pages.SwitchToPage("searchStocks")
				app.SetFocus(searchStocksInput)
			case "search-crypto":
				pages.SwitchToPage("searchCrypto")
				app.SetFocus(searchCryptoInput)
			case "quit":
				PromptQuit(app, mainLayout, mainCommands, mainInput, mainCommandsText)
			default:
			}
			mainInput.SetText("")
		}
	})

	summaryInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := summaryInput.GetText()
			switch cmd {
			case "main":
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				PromptQuit(app, summaryLayout, summaryCommands, summaryInput, summaryCommandsText)
			default:
			}
		}
		summaryInput.SetText("")
	})

	budgetMainInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := budgetMainInput.GetText()
			switch cmd {
			case "main":
				budgetMainInput.SetText("")
				pages.SwitchToPage("main")
				app.SetFocus(mainInput)
			case "quit":
				PromptQuit(app, budgetLayout, budgetCommands, budgetMainInput, budgetCommandsText)
			default:
			}
		}
		summaryInput.SetText("")
	})

	searchStocksInput.SetDoneFunc(func(key tcell.Key) {
		cmd := strings.TrimSpace(searchStocksInput.GetText())
		if key == tcell.KeyEnter {
			if strings.HasPrefix(strings.ToLower(cmd), "search $") {
				searchStocksWaiting.SetText("Waiting for stock data...")
				os.Remove("../sprint3/microservice-c/output_stock.json")
				ticker := strings.TrimPrefix(cmd[7:], "$")
				err := writeTickerInput("../sprint3/microservice-c/input_stock.json", ticker)
				if err != nil {
					searchStocksWaiting.SetText("Failed to write input.")
				} else {
					searchStocksTable.Clear()
					searchStocksLayout.RemoveItem(searchStocksTable)
					searchStocksLayout.AddItem(searchStocksTable, 0, 1, false)
					showMore = false
					go waitForStockDataAndRender(app, "../sprint3/microservice-c/output_stock.json", searchStocksTable, showMore, searchStocksWaiting)
				}
			} else if strings.ToLower(cmd) == "show-more" {
				showMore = true
				go waitForStockDataAndRender(app, "../sprint3/microservice-c/output_stock.json", searchStocksTable, showMore, searchStocksWaiting)
				searchStocksWaiting.SetText("")
			}
		}
		switch cmd {
		case "main":
			pages.SwitchToPage("main")
			app.SetFocus(mainInput)
		case "quit":
			PromptQuit(app, searchStocksLayout, searchStocksCommands, searchStocksInput, searchStocksCommandsText)
		default:
		}
		searchStocksInput.SetText("")
	})

	searchCryptoInput.SetDoneFunc(func(key tcell.Key) {
		cmd := searchCryptoInput.GetText()
		if key == tcell.KeyEnter {
			if strings.HasPrefix(cmd, "search ") {
				searchCryptoWaiting.SetText("Waiting for cryptocurrency data...")
				os.Remove("../sprint3/microservice-d/output_crypto.json")
				coin := strings.TrimPrefix(cmd, "search ")
				err := writeCryptoInput("../sprint3/microservice-d/input_crypto.json", coin)
				if err != nil {
					searchCryptoWaiting.SetText("Failed to write input.")
				} else {
					searchCryptoTable.Clear()
					searchCryptoLayout.RemoveItem(searchCryptoTable)
					searchCryptoLayout.AddItem(searchCryptoTable, 0, 1, false)
					go waitForCryptoDataAndRender(app, "../sprint3/microservice-d/output_crypto.json", searchCryptoTable, searchCryptoWaiting)
				}
			}
			switch cmd {
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

// SUMMARY FUNCTIONS
func loadIndexDataFromFile(path string) ([]IndexData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var indices []IndexData
	if err := json.Unmarshal(data, &indices); err != nil {
		return nil, err
	}

	return indices, nil
}

func writeSummaryInput(filePath string) error {
	data := map[string]int{"summary": 1}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func waitForSummaryData(app *tview.Application, path string, onLoaded func([]IndexData)) {
	go func() {
		for {
			if _, err := os.Stat(path); err == nil {
				data, err := loadIndexDataFromFile(path)
				if err == nil {
					app.QueueUpdateDraw(func() {
						onLoaded(data)
					})
					break
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func renderSummaryTable(summaryTable *tview.Table, indices []IndexData) {
	summaryTable.Clear()

	headers := []string{"Name", "Ticker", "Date", "Open", "High", "Low", "Close", "Volume"}
	for col, h := range headers {
		summaryTable.SetCell(0, col, tview.NewTableCell(h).SetAlign(tview.AlignCenter).SetSelectable(false))
	}

	for row, index := range indices {
		summaryTable.SetCell(row+1, 0, tview.NewTableCell(index.Name))
		summaryTable.SetCell(row+1, 1, tview.NewTableCell(index.Ticker))
		summaryTable.SetCell(row+1, 2, tview.NewTableCell(index.Date))
		summaryTable.SetCell(row+1, 3, tview.NewTableCell(fmt.Sprintf("%.2f", index.Open)).SetAlign(tview.AlignRight))
		summaryTable.SetCell(row+1, 4, tview.NewTableCell(fmt.Sprintf("%.2f", index.High)).SetAlign(tview.AlignRight))
		summaryTable.SetCell(row+1, 5, tview.NewTableCell(fmt.Sprintf("%.2f", index.Low)).SetAlign(tview.AlignRight))
		summaryTable.SetCell(row+1, 6, tview.NewTableCell(fmt.Sprintf("%.2f", index.Close)).SetAlign(tview.AlignRight))
		summaryTable.SetCell(row+1, 7, tview.NewTableCell(fmt.Sprintf("%d", index.Volume)).SetAlign(tview.AlignRight))
	}
}

// STOCK FUNCTIONS
func loadStockFromFile(path string) (StockData, error) {
	var data StockData
	bytes, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(bytes, &data)
	return data, err
}

func writeTickerInput(filePath string, ticker string) error {
	input := map[string]string{"ticker": ticker}
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func waitForStockDataAndRender(app *tview.Application, path string, stockTable *tview.Table, showMore bool, message *tview.TextView) {
	for {
		if _, err := os.Stat(path); err == nil {
			data, err := loadStockFromFile(path)
			if err != nil {
				log.Printf("Failed to parse stock data: %v", err)
				return
			}

			app.QueueUpdateDraw(func() {
				renderStockTable(stockTable, data, showMore)
				message.SetText("")
			})

			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func renderStockTable(stockTable *tview.Table, stockData StockData, showMore bool) {
	stockTable.Clear()

	if showMore {
		// Full view
		headers := []string{"Ticker", "Date", "Open", "High", "Low", "Close"}

		for col, h := range headers {
			stockTable.SetCell(0, col,
				tview.NewTableCell(h).
					SetAlign(tview.AlignCenter).
					SetSelectable(false))
		}

		stockTable.SetCell(1, 0, tview.NewTableCell(stockData.Ticker))
		stockTable.SetCell(1, 1, tview.NewTableCell(stockData.Date))
		stockTable.SetCell(1, 2, tview.NewTableCell(fmt.Sprintf("%.2f", stockData.Open)).SetAlign(tview.AlignRight))
		stockTable.SetCell(1, 3, tview.NewTableCell(fmt.Sprintf("%.2f", stockData.High)).SetAlign(tview.AlignRight))
		stockTable.SetCell(1, 4, tview.NewTableCell(fmt.Sprintf("%.2f", stockData.Low)).SetAlign(tview.AlignRight))
		stockTable.SetCell(1, 5, tview.NewTableCell(fmt.Sprintf("%.2f", stockData.Close)).SetAlign(tview.AlignRight))
	} else {
		// Compact view
		headers := []string{"Ticker", "Date", "Close"}

		for col, h := range headers {
			stockTable.SetCell(0, col,
				tview.NewTableCell(h).
					SetAlign(tview.AlignCenter).
					SetSelectable(false))
		}

		stockTable.SetCell(1, 0, tview.NewTableCell(stockData.Ticker))
		stockTable.SetCell(1, 1, tview.NewTableCell(stockData.Date))
		stockTable.SetCell(1, 2, tview.NewTableCell(fmt.Sprintf("%.2f", stockData.Close)).SetAlign(tview.AlignRight))
	}
}

// CRYPTO FUNCTIONS
func writeCryptoInput(path string, coin string) error {
	data := map[string]string{"coin": coin}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonData, 0644)
}

type OrderedPair struct {
	Key   string
	Value interface{}
}

func waitForCryptoDataAndRender(app *tview.Application, path string, table *tview.Table, message *tview.TextView) {
	for {
		if _, err := os.Stat(path); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		app.QueueUpdateDraw(func() {
			table.Clear()
			table.SetCell(0, 0, tview.NewTableCell("Failed to read crypto data"))
		})
		return
	}

	// Use a decoder to preserve order
	var raw map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&raw); err != nil {
		app.QueueUpdateDraw(func() {
			table.Clear()
			table.SetCell(0, 0, tview.NewTableCell("Invalid crypto JSON"))
		})
		return
	}

	// Convert to ordered pairs
	var ordered []OrderedPair
	for k, v := range raw {
		ordered = append(ordered, OrderedPair{Key: k, Value: v})
	}

	app.QueueUpdateDraw(func() {
		renderCryptoTable(table, ordered)
		message.SetText("")
	})
}

func renderCryptoTable(table *tview.Table, data []OrderedPair) {
	table.Clear()

	table.SetCell(0, 0,
		tview.NewTableCell("Coin").
			SetAlign(tview.AlignCenter).
			SetSelectable(false).
			SetAttributes(tcell.AttrBold))

	table.SetCell(0, 1,
		tview.NewTableCell("Price").
			SetAlign(tview.AlignCenter).
			SetSelectable(false).
			SetAttributes(tcell.AttrBold))

	for i, pair := range data {
		table.SetCell(i+1, 0, tview.NewTableCell(pair.Key))
		table.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("%v", pair.Value)))
	}
}

// BUDGET FUNCTIONS
func saveBudgetToFile(filename string, total int, categories map[string]int) error {
	// Build JSON manually: start with total
	builder := make(map[string]interface{})
	builder["total"] = total

	// Serialize separately
	totalJSON, err := json.MarshalIndent(builder, "", "  ")
	if err != nil {
		return err
	}

	// Remove the final } so we can append to it
	totalStr := string(totalJSON[:len(totalJSON)-2]) // Cut off last "}"

	// Append the categories
	categoryJSON := ""
	for k, v := range categories {
		catEntry, err := json.MarshalIndent(map[string]int{k: v}, "", "  ")
		if err != nil {
			return err
		}
		// Trim opening { and newline
		line := string(catEntry[2 : len(catEntry)-2])
		categoryJSON += ",\n  " + line
	}

	final := totalStr + categoryJSON + "\n}\n"

	return os.WriteFile(filename, []byte(final), 0644)
}

func waitForBudgetOutput(app *tview.Application, path string, table *tview.Table) {
	for {
		if _, err := os.Stat(path); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		app.QueueUpdateDraw(func() {
			table.SetCell(0, 0, tview.NewTableCell("Failed to read output file"))
		})
		return
	}

	var budget map[string]int
	if err := json.Unmarshal(data, &budget); err != nil {
		app.QueueUpdateDraw(func() {
			table.SetCell(0, 0, tview.NewTableCell("Invalid JSON in output file"))
		})
		return
	}

	app.QueueUpdateDraw(func() {
		renderBudgetTable(table, budget)
	})
}

func renderBudgetTable(table *tview.Table, budget map[string]int) {
	table.Clear()

	row := 0
	// Add headers
	table.SetCell(row, 0, tview.NewTableCell("Category").SetAlign(tview.AlignCenter).SetSelectable(false))
	table.SetCell(row, 1, tview.NewTableCell("Amount").SetAlign(tview.AlignCenter).SetSelectable(false))
	row++

	for k, v := range budget {
		if k == "total" {
			continue
		}
		table.SetCell(row, 0, tview.NewTableCell(k))
		table.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d", v)))
		row++
	}

	// Optionally show total last or at the top
	if total, ok := budget["total"]; ok {
		table.SetCell(row, 0, tview.NewTableCell("Total"))
		table.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d", total)))
	}
}
