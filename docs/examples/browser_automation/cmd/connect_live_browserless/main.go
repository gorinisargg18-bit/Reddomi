package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/shank318/doota/browser_automation"
	"github.com/streamingfast/dstore"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	browserlessToken := os.Getenv("DOOTA_START_COMMON_BROWSERLESS_API_KEY")
	browserlessWarmup := os.Getenv("DOOTA_START_COMMON_BROWSERLESS_WARMUP_API_KEY")
	if browserlessWarmup == "" {
		browserlessWarmup = "2SIxpPBYG6XJqLj5ec45cd436c170abdbec8713fd1bbaffe4"
	}

	if browserlessToken == "" {
		fmt.Println("No Browserless token found. Set DOOTA_START_COMMON_BROWSERLESS_API_KEY")
		os.Exit(1)
	}

	browserless := browser_automation.NewBrowserLessBrowser(browserlessToken, browserlessWarmup, logger)
	provider := browserless

	debugStore, _ := dstore.NewStore("data/debugstore", "", "", true)
	redditBA := browser_automation.NewRedditBrowserAutomation(provider, logger, debugStore)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	cdp, err := redditBA.StartLogin(ctx, "")
	if err != nil {
		fmt.Println("StartLogin error:", err)
		os.Exit(2)
	}

	fmt.Println("LiveURL:", cdp.LiveURL)
	fmt.Println("WSEndpoint:", cdp.WSEndpoint)
	fmt.Println("Open the LiveURL in your browser to perform the Reddit login. Waiting for login...")

	// Ensure Playwright driver installed
	if err := playwright.Install(&playwright.RunOptions{SkipInstallBrowsers: true}); err != nil {
		fmt.Println("playwright.Install error:", err)
		os.Exit(3)
	}

	cfg, err := redditBA.WaitAndGetCookies(ctx, cdp)
	if err != nil {
		fmt.Println("WaitAndGetCookies error:", err)
		os.Exit(4)
	}

	fmt.Println("Login detected!")
	fmt.Println("Username:", cfg.Username)
	fmt.Println("Cookies length:", len(cfg.Cookies))
	if len(cfg.Cookies) > 1000 {
		fmt.Println(cfg.Cookies[:1000])
	} else {
		fmt.Println(cfg.Cookies)
	}
}
