package main

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/events"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/*
var iconFS embed.FS

const AppTitle = "GitHub Notifier"

var frontMostAppId string

func getFrontmostAppId() (string, error) {
	out, err := exec.Command("osascript", "-e", `tell application "System Events" to get bundle identifier of first application process whose frontmost is true`).Output()
	return string(out), err
}

func focusApp(bundleID string) error {
	return exec.Command("open", "-b", bundleID).Run()
}

func main() {
	app := application.New(application.Options{
		Name:        AppTitle,
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GitHubAPIService{
				httpClient: http.DefaultClient,
				token:      os.Getenv("GITHUB_TOKEN"), // TODO(later): Support multiple authentication methods
			}),
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

	// System tray
	systray := app.SystemTray.New()

	lightModeIconBytes, _ := iconFS.ReadFile("assets/github-light.png")
	darkModeIconBytes, _ := iconFS.ReadFile("assets/github-dark.png")
	systray.SetIcon(lightModeIconBytes)
	systray.SetDarkModeIcon(darkModeIconBytes)

	systray.SetTooltip(AppTitle) // Windows

	// System tray window
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:           AppTitle,
		Name:            AppTitle,
		Frameless:       true,
		Hidden:          true,
		Width:           500,
		Height:          400,
		DisableResize:   true,
		DevToolsEnabled: true,
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

	app.Event.On("escape-pressed", func(event *application.CustomEvent) {
		fmt.Println("Escape pressed!")
		window.Hide()
	})

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
