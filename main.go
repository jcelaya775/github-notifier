package main

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/events"
	"log"
	"os/exec"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/*
var iconFS embed.FS

const AppTitle = "GitHub Notifier"

var frontMostAppId string

// get frontmost app bundle ID
func getFrontmostAppId() (string, error) {
	out, err := exec.Command("osascript", "-e", `tell application "System Events" to get bundle identifier of first application process whose frontmost is true`).Output()
	return string(out), err
}

// re-focus previously active app
func focusApp(bundleID string) error {
	return exec.Command("open", "-b", bundleID).Run()
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "github-notifier",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: false,
		},
		Linux: application.LinuxOptions{
			DisableQuitOnLastWindowClosed: false,
		},
	})

	systray := app.SystemTray.New()

	// System tray icons
	lightModeIconBytes, _ := iconFS.ReadFile("assets/github-light.png")
	darkModeIconBytes, _ := iconFS.ReadFile("assets/github-dark.png")
	systray.SetIcon(lightModeIconBytes)
	systray.SetDarkModeIcon(darkModeIconBytes)

	// System tray tooltip and label
	systray.SetTooltip(AppTitle) // Windows

	// System tray window
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:         AppTitle,
		Frameless:     true,
		Hidden:        true,
		Width:         500,
		Height:        400,
		DisableResize: true,
	})

	app.Event.On("escape-pressed", func(event *application.CustomEvent) {
		fmt.Println("Escape pressed!")
		window.Hide()
	})

	window.RegisterHook(events.Common.WindowLostFocus, func(event *application.WindowEvent) {
		window.Hide()
	})

	//window.RegisterHook(events.Common.WindowFocus, func(event *application.WindowEvent) {
	//	fmt.Println("Window shown!")
	//	//
	//	//	appId, err := getFrontmostAppId()
	//	//	if err != nil {
	//	//		fmt.Printf("Error getting frontmost app: %v\n", err)
	//	//	}
	//	//	fmt.Printf("Frontmost app before showing: %s\n", appId)
	//	//	frontMostAppId = appId
	//})
	//
	//window.RegisterHook(events.Common.WindowLostFocus, func(event *application.WindowEvent) {
	//	fmt.Println("Window unfocusing!")
	//
	//	//fmt.Printf("Frontmost app before refocus: %s\n", frontMostAppId)
	//	//if frontMostAppId != "" {
	//	//	if err := focusApp(frontMostAppId); err != nil {
	//	//		fmt.Printf("Error refocusing app: %v\n", err)
	//	//	}
	//	//}
	//})

	systray.AttachWindow(window)

	//systray.OnClick(func() {
	//	if !window.IsVisible() {
	//		window.SetAlwaysOnTop(true)
	//		window.Show()
	//		window.Focus()
	//		window.SetAlwaysOnTop(false)
	//	} else {
	//		window.Hide()
	//	}
	//})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
