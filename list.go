package main

import (
	"context"
	"errors"
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
	isClose             bool
	rsyncUploader       RsyncUploader
}

func NewUIList(items []Cmd, selectedCommandChan chan<- Cmd) *SelectList {
	if len(items) == 0 {
		color.Red("Cmd list is empty, please fill in your configuration first.")
		os.Exit(1)
	}
	selectList := &SelectList{
		normalItems:         items,
		searchItems:         items,
		uiList:              widgets.NewList(),
		selectedMode:        NormalMode,
		selectedCommandChan: selectedCommandChan,
		normalTitle:         "Usage: (Search:</>) (Up/Down:<k>/<j>) (Exit:<C-c>/<Esc>) (Rsync:<C-r>)",
		searchTitle:         "Search: [%s](fg:red)  |  Usage: (Up/Down:<C-k>/<C-j>) (Exit:<C-c>/<Esc>) (Erase:<C-u>) (Rsync:<C-r>)",
		isClose:             false,
	}
	selectList.initUI()
	selectList.resizeUI()
	selectList.renderUI()
	return selectList
}

func (sl *SelectList) registerRsyncUploader(rsyncUploader RsyncUploader) {
	sl.rsyncUploader = rsyncUploader
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

//lint:file-ignore U1000 Ignore unused code, it will be used in the future, now for the simulation test.
func (sl *SelectList) listenEventsWithCancel(ctx context.Context) {
	uiEvents := ui.PollEvents()
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-uiEvents:
			switch sl.selectedMode {
			case NormalMode:
				sl.handleEventsAtNormalMode(e)
			case SearchMode:
				sl.handleEventsAtSearchMode(e)
			}
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
		sl.close()
		sl.selectedCommandChan <- Cmd{}
	case "<Enter>":
		sl.close()
		sl.selectedCommandChan <- sl.normalItems[sl.uiList.SelectedRow]
	case "<C-r>":
		sl.rsync()
	case "<Resize>":
		sl.resizeUI()
	case "/":
		sl.selectedMode = SearchMode
		sl.uiList.SelectedRow = 0
		sl.setSearchTitle()
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
	case "<C-u>":
		if len(sl.searchStr) != 0 {
			sl.searchStr = ""
			sl.setSearchTitle()
			sl.doSearch()
		}
	case "<Resize>":
		sl.resizeUI()
	case "<Enter>":
		if len(sl.searchItems) > 0 {
			sl.close()
			sl.selectedCommandChan <- sl.searchItems[sl.uiList.SelectedRow]
		}
	case "<C-r>":
		sl.rsync()
	case "<C-c>", "<Escape>":
		sl.selectedMode = NormalMode
		sl.searchStr = ""
		sl.uiList.Title = sl.normalTitle
		sl.uiList.SelectedRow = 0
		sl.searchItems = sl.normalItems
	case "<Backspace>":
		if len(sl.searchStr) > 0 {
			sl.searchStr = sl.searchStr[:len(sl.searchStr)-1]
			sl.setSearchTitle()
			sl.doSearch()
		}
	case "<Space>":
		sl.searchStr += " "
		sl.setSearchTitle()
		sl.doSearch()
	default:
		if len(e.ID) != 1 {
			return
		}
		sl.searchStr += e.ID
		sl.setSearchTitle()
		sl.doSearch()
	}
	sl.renderUI()
}

func (sl *SelectList) setSearchTitle() {
	sl.uiList.Title = fmt.Sprintf(sl.searchTitle, sl.searchStr)
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

func (sl *SelectList) close() {
	if sl.isClose {
		return
	}

	sl.isClose = true
	ui.Close()
}

func (sl *SelectList) rsync() {
	selectedCmd := Cmd{}
	if sl.selectedMode == NormalMode {
		selectedCmd = sl.normalItems[sl.uiList.SelectedRow]
	} else if sl.selectedMode == SearchMode {
		selectedCmd = sl.searchItems[sl.uiList.SelectedRow]
	}

	uploadCmd, err := sl.rsyncUploader.Upload(selectedCmd)
	if errors.Is(err, ErrRsUserCancel) || errors.Is(err, ErrRsNotSSHCmd) {
		debug("RsyncUpload get: %+v, then do nothing.", err)
		return
	} else if err != nil {
		sl.close()
		color.Red("RsyncUpload get: %v, will exit.", err)
		os.Exit(1)
	}

	sl.close()
	sl.selectedCommandChan <- Cmd{uploadCmd, fmt.Sprintf("Rsync %s", selectedCmd.Name), ""}
}
