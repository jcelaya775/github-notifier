package main

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/services/badge"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/github-light.png assets/github-dark.png
var iconFS embed.FS

const APP_TITLE = "GitHub Notifier"

var windowIsVisible = false

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// System tray notification count badge
	badgeService := badge.New()
	if err := badgeService.SetBadge("3"); err != nil {
		log.Fatal(err)
	}

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
			application.NewService(badgeService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	systray := app.SystemTray.New()

	// System tray icons
	lightModeIconBytes, _ := iconFS.ReadFile("assets/github-light.png")
	darkModeIconBytes, _ := iconFS.ReadFile("assets/github-dark.png")
	systray.SetIcon(lightModeIconBytes)
	systray.SetDarkModeIcon(darkModeIconBytes)

	// System tray tooltip (Windows) and label (macOS)
	systray.SetTooltip(APP_TITLE)
	systray.SetLabel(APP_TITLE)

	// System tray menu
	menu := application.NewMenu()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systray.SetMenu(menu)

	// System tray window
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  APP_TITLE,
		Width:  200,
		Height: 100,
	})
	systray.AttachWindow(window)

	systray.OnClick(func() {
		fmt.Println("Clicked on tray icon")
		window.Show()
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	//app.Window.NewWithOptions(application.WebviewWindowOptions{
	//	Title: "Window 1",
	//	Mac: application.MacWindow{
	//		InvisibleTitleBarHeight: 50,
	//		Backdrop:                application.MacBackdropTranslucent,
	//		TitleBar:                application.MacTitleBarHiddenInset,
	//	},
	//	Linux:            application.LinuxWindow{},
	//	BackgroundColour: application.NewRGB(27, 38, 54),
	//	URL:              "/",
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
