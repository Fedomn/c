package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	ui "github.com/fedomn/termui/v3"
	"github.com/fedomn/termui/v3/widgets"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type listMode int

const (
	NormalMode listMode = iota
	SearchMode
)

type SelectList struct {
	normalItems         []Cmd
	searchItems         []Cmd
	uiList              *widgets.List
	selectedMode        listMode
	selectedCommandChan chan<- Cmd
	normalTitle         string
	searchTitle         string
	searchStr           string
}

func NewUIList(items []Cmd, selectedCommandChan chan<- Cmd) *SelectList {
	selectList := &SelectList{
		normalItems:         items,
		searchItems:         items,
		uiList:              widgets.NewList(),
		selectedMode:        NormalMode,
		selectedCommandChan: selectedCommandChan,
		normalTitle:         "Usage: (Search:</>) (Up/Down:<k>/<j>) (Exit:<C-c>/<Esc>)",
		searchTitle:         "Search: [%s](fg:red)  |  Usage: (Up/Down:<C-k>/<C-j>) (Exit:<C-c>/<Esc>) (Erase:<C-u>)",
	}
	selectList.initUI()
	selectList.resizeUI()
	selectList.renderUI()
	return selectList
}

func (sl *SelectList) initUI() {
	if err := ui.Init(); err != nil {
		color.Red("Failed to initialize termui: %v", err)
		os.Exit(1)
	}
	uiList := widgets.NewList()
	uiList.Title = sl.normalTitle
	uiList.TitleStyle = ui.NewStyle(ui.ColorBlue, ui.ColorClear, ui.ModifierBold)
	uiList.BorderStyle = ui.NewStyle(ui.ColorWhite)
	uiList.TextStyle = ui.NewStyle(ui.ColorCyan)
	uiList.WrapText = false

	sl.uiList = uiList
	debug("Init uiList successfully.")
}

func (sl *SelectList) resizeUI() {
	termWidth, termHeight := ui.TerminalDimensions()
	sl.uiList.SetRect(0, 0, termWidth, termHeight)
	debug("Resize uiList successfully.")
}

func (sl *SelectList) renderUI() {
	var rows []string
	var items []Cmd
	if sl.selectedMode == NormalMode {
		items = sl.normalItems
	} else if sl.selectedMode == SearchMode {
		items = sl.searchItems
	}
	for k, v := range items {
		if k == sl.uiList.SelectedRow {
			format := "[[%02d]](fg:green) [%s](fg:green,mod:underline) [-](fg:cyan,mod:bold) [%s](fg:green,mod:bold)"
			rows = append(rows, fmt.Sprintf(format, k, v.Name, v.Cmd))
		} else {
			rows = append(rows, fmt.Sprintf("[%02d] %s", k, v.Name))
		}
	}
	sl.uiList.Rows = rows
	ui.Render(sl.uiList)
	debug("Render uiList successfully. Selected Row Index: %v", sl.uiList.SelectedRow)
}

func (sl *SelectList) ListenEvents() {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch sl.selectedMode {
		case NormalMode:
			sl.handleEventsAtNormalMode(e)
		case SearchMode:
			sl.handleEventsAtSearchMode(e)
		}
	}
}

func (sl *SelectList) handleEventsAtNormalMode(e ui.Event) {
	debug("Normal Mode Event: %+v", e)
	switch e.ID {
	case "j", "<Down>":
		sl.uiList.ScrollDown()
	case "k", "<Up>":
		sl.uiList.ScrollUp()
	case "<C-d>":
		sl.uiList.ScrollHalfPageDown()
	case "<C-u>":
		sl.uiList.ScrollHalfPageUp()
	case "<C-f>":
		sl.uiList.ScrollPageDown()
	case "<C-b>":
		sl.uiList.ScrollPageUp()
	case "q", "<C-c>", "<Escape>":
		ui.Close()
		sl.selectedCommandChan <- Cmd{}
	case "<Enter>":
		ui.Close()
		sl.selectedCommandChan <- sl.normalItems[sl.uiList.SelectedRow]
	case "<Resize>":
		sl.resizeUI()
	case "/":
		sl.selectedMode = SearchMode
		sl.uiList.Title = fmt.Sprintf(sl.searchTitle, sl.searchStr)
		sl.uiList.SelectedRow = 0
	}
	sl.renderUI()
}

func (sl *SelectList) handleEventsAtSearchMode(e ui.Event) {
	debug("Search Mode Event: %+v", e)
	switch e.ID {
	case "<Down>", "<C-j>":
		if len(sl.searchItems) > 0 {
			sl.uiList.ScrollDown()
		}
	case "<Up>", "<C-k>":
		if len(sl.searchItems) > 0 {
			sl.uiList.ScrollUp()
		}
	case "<Resize>":
		sl.resizeUI()
	case "<Enter>":
		if len(sl.searchItems) > 0 {
			ui.Close()
			sl.selectedCommandChan <- sl.searchItems[sl.uiList.SelectedRow]
		}
	case "<C-c>", "<Escape>":
		sl.selectedMode = NormalMode
		sl.searchStr = ""
		sl.uiList.Title = sl.normalTitle
		sl.uiList.SelectedRow = 0
		sl.searchItems = sl.normalItems
	case "<Backspace>":
		if len(sl.searchStr) > 0 {
			sl.searchStr = sl.searchStr[:len(sl.searchStr)-1]
			sl.uiList.Title = fmt.Sprintf(sl.searchTitle, sl.searchStr)
			sl.doSearch()
		}
	case "<Space>":
		sl.searchStr += " "
		sl.uiList.Title = fmt.Sprintf(sl.searchTitle, sl.searchStr)
		sl.doSearch()
	default:
		if len(e.ID) != 1 {
			return
		}
		sl.searchStr += e.ID
		sl.uiList.Title = fmt.Sprintf(sl.searchTitle, sl.searchStr)
		sl.doSearch()
	}
	sl.renderUI()
}

func (sl *SelectList) doSearch() {
	var searchResult []Cmd
	for _, v := range sl.normalItems {
		if fuzzy.Match(sl.searchStr, v.Name) || fuzzy.Match(sl.searchStr, v.Cmd) {
			searchResult = append(searchResult, v)
		}
	}
	sl.uiList.SelectedRow = 0
	sl.searchItems = searchResult
}
