package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/shank318/doota/browser_automation"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()
	token := os.Getenv("DOOTA_START_COMMON_STEEL_API_KEY")
	if token == "" {
		fmt.Println("DOOTA_START_COMMON_STEEL_API_KEY not set")
		os.Exit(1)
	}

	steel := browser_automation.NewSteelBrowser(token, logger)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	info, err := steel.GetCDPInfo(ctx, browser_automation.CDPInput{StartURL: "https://www.reddit.com", UseProxy: false, LiveURL: false})
	if err != nil {
		fmt.Println("GetCDPInfo error:", err)
		os.Exit(2)
	}

	fmt.Printf("WSEndpoint: %s\n", info.WSEndpoint)

	// Install Playwright driver (skip browser install) to ensure playwright.Run works
	if err := playwright.Install(&playwright.RunOptions{SkipInstallBrowsers: true}); err != nil {
		fmt.Println("playwright.Install error:", err)
		os.Exit(3)
	}

	pw, err := playwright.Run()
	if err != nil {
		fmt.Println("playwright.Run error:", err)
		os.Exit(4)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.ConnectOverCDP(info.WSEndpoint)
	if err != nil {
		fmt.Println("ConnectOverCDP error:", err)
		os.Exit(4)
	}
	defer browser.Close()

	pageContext := browser.Contexts()[0]
	page := pageContext.Pages()[0]

	fmt.Println("Connected, current page URL:", page.URL())

	fmt.Println("Success")
}
