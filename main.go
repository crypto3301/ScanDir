package main

import (
	"fmt"
	"log"
	"os"
	"scandir/backend"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

type Command struct {
	Name        string
	Description string
	Handler     func(g *gocui.Gui, args []string) error
}

var commands = make(map[string]Command)

func registerCommand(name, description string, handler func(g *gocui.Gui, args []string) error) {
	commands[name] = Command{
		Name:        name,
		Description: description,
		Handler:     handler,
	}
}

var (
	filesList   []string
	selectedIdx int
	viewArr     = []string{"window1", "window2", "window3"}
	activeView  = 0
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("window3", gocui.KeyEnter, gocui.ModNone, processInput); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("window1", gocui.KeyArrowUp, gocui.ModNone, scrollUpContent); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("window1", gocui.KeyArrowDown, gocui.ModNone, scrollDownContent); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("window1", gocui.KeyEnter, gocui.ModNone, openSelectedFile); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("window1", gocui.KeySpace, gocui.ModNone, updateFilesList); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("window2", gocui.KeyArrowUp, gocui.ModNone, scrollUpContent); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("window2", gocui.KeyArrowDown, gocui.ModNone, scrollDownContent); err != nil {
		log.Panicln(err)
	}

	g.Cursor = true

	registerCommand("help", "Show available commands", helpCommand)
	registerCommand("echo", "Echo the input", echoCommand)
	registerCommand("auto", "auto mode", autoCommand)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("window1", 0, 0, maxX/2-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "INFO"
		v.Wrap = true
		v.Highlight = true
		v.Autoscroll = false

		if err := loadFilesList(); err != nil {
			return err
		}
		if err := updateFilesList(g, v); err != nil {
			return err
		}
	}

	if v, err := g.SetView("window2", maxX/2, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "LOG"
		v.Wrap = true
		v.Autoscroll = false
		v.Highlight = true
	}

	if v, err := g.SetView("window3", 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Terminal"
		v.Editable = true
		v.Wrap = false
	}

	if _, err := g.SetCurrentView(viewArr[activeView]); err != nil {
		return err
	}

	updateActiveViewColor(g)

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (activeView + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}

	activeView = nextIndex

	updateActiveViewColor(g)

	return nil
}

func processInput(g *gocui.Gui, v *gocui.View) error {
	input := v.Buffer()
	input = strings.TrimSpace(input)

	v.Clear()
	v.SetCursor(0, 0)

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	cmdName := parts[0]
	args := parts[1:]

	if cmd, exists := commands[cmdName]; exists {
		if err := cmd.Handler(g, args); err != nil {
			showErrorModal(g, fmt.Sprintf("Error executing command: %v", err))
		} else {
			g.DeleteView("errorModal")
		}
	} else {
		showErrorModal(g, fmt.Sprintf("Unknown command: %s", cmdName))
	}

	return nil
}

func scrollUpContent(g *gocui.Gui, v *gocui.View) error {
	v.Autoscroll = false
	ox, oy := v.Origin()
	if oy > 0 {
		v.SetOrigin(ox, oy-1)
	}
	return nil
}

func scrollDownContent(g *gocui.Gui, v *gocui.View) error {
	v.Autoscroll = false
	ox, oy := v.Origin()
	_, sy := v.Size()
	lines := len(v.BufferLines())
	if oy+sy < lines {
		v.SetOrigin(ox, oy+1)
	}
	return nil
}

func updateActiveViewColor(g *gocui.Gui) {
	for _, name := range viewArr {
		v, _ := g.View(name)
		if name == viewArr[activeView] {
			v.FgColor = gocui.ColorWhite
			v.BgColor = gocui.ColorBlue
		} else {
			v.FgColor = gocui.ColorBlack
			v.BgColor = gocui.ColorWhite
		}
	}
}

func helpCommand(g *gocui.Gui, args []string) error {
	logView, err := g.View("window2")
	if err != nil {
		return err
	}

	fmt.Fprintln(logView, time.Now().Format(time.RFC850), " > ", "Available commands:")
	for _, cmd := range commands {
		fmt.Fprintf(logView, "  %s - %s\n", cmd.Name, cmd.Description)
	}

	return nil
}

func echoCommand(g *gocui.Gui, args []string) error {
	logView, err := g.View("window2")
	if err != nil {
		return err
	}

	fmt.Fprintln(logView, time.Now().Format(time.RFC850), " > ", strings.Join(args, " "))
	return nil
}

func autoCommand(g *gocui.Gui, args []string) error {
	logView, err := g.View("window2")
	if err != nil {
		return err
	}

	if len(args) < 1 {
		showErrorModal(g, "Choose: auto <dir> [ext...]")
		return nil
	}

	dir := args[0]
	extensions := args[1:]

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		showErrorModal(g, fmt.Sprintf("No dir: %s", dir))
		return nil
	}

	files, err := backend.ParseDirectory(dir, extensions)
	if err != nil {
		showErrorModal(g, fmt.Sprintf("Error of scan: %v", err))
		return nil
	}

	outputPath := "snapshot_with_hash.json"
	err = backend.WriteToJSON(files, outputPath)
	if err != nil {
		showErrorModal(g, fmt.Sprintf("Error save in JSON: %v", err))
		return nil
	}

	fmt.Fprintf(logView, time.Now().Format(time.RFC850), " > ", "A snapshot of the file system is saved in %s\n", outputPath)

	if err := loadFilesList(); err != nil {
		showErrorModal(g, fmt.Sprintf("File list update error: %v", err))
		return nil
	}
	if err := updateFilesList(g, logView); err != nil {
		showErrorModal(g, fmt.Sprintf("Window update error: %v", err))
		return nil
	}

	g.DeleteView("errorModal")
	return nil
}

func loadFilesList() error {
	files, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	filesList = nil
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			filesList = append(filesList, file.Name())
		}
	}
	return nil
}

func updateFilesList(g *gocui.Gui, v *gocui.View) error {
	v, err := g.View("window1")
	if err != nil {
		return err
	}

	v.Clear()
	for i, file := range filesList {
		if i == selectedIdx {
			fmt.Fprintf(v, "> %s\n", file)
		} else {
			fmt.Fprintf(v, "  %s\n", file)
		}
	}
	return nil
}

func openSelectedFile(g *gocui.Gui, v *gocui.View) error {
	if selectedIdx < 0 || selectedIdx >= len(filesList) {
		return nil
	}

	filePath := filesList[selectedIdx]
	content, err := os.ReadFile(filePath)
	if err != nil {
		showErrorModal(g, fmt.Sprintf("File opening error: %v", err))
		return nil
	}

	v, err = g.View("window1")
	if err != nil {
		return err
	}

	v.Clear()
	fmt.Fprintf(v, "File Contents %s:\n%s\n", filePath, string(content))
	return nil
}

func showErrorModal(g *gocui.Gui, message string) {
	maxX, maxY := g.Size()

	g.DeleteView("errorModal")

	if v, err := g.SetView("errorModal", 0, maxY/2-1, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return
		}
		v.Title = "Error"
		v.Wrap = true
		v.FgColor = gocui.ColorRed
		v.BgColor = gocui.ColorBlack
		fmt.Fprintln(v, message)

		if err := g.SetKeybinding("errorModal", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			g.DeleteView("errorModal")
			return nil
		}); err != nil {
			log.Println("Error setting keybinding:", err)
		}

		if _, err := g.SetCurrentView("errorModal"); err != nil {
			log.Println("Error setting current view:", err)
		}
	}
}
