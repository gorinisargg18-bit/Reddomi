package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"encoding/json"
	"github.com/playwright-community/playwright-go"
	"github.com/shank318/doota/browser_automation"
	"github.com/streamingfast/dstore"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	steelToken := os.Getenv("DOOTA_START_COMMON_STEEL_API_KEY")
	browserlessToken := os.Getenv("DOOTA_START_COMMON_BROWSERLESS_API_KEY")
	browserlessWarmup := os.Getenv("DOOTA_START_COMMON_BROWSERLESS_WARMUP_API_KEY")
	if browserlessWarmup == "" {
		browserlessWarmup = "2SIxpPBYG6XJqLj5ec45cd436c170abdbec8713fd1bbaffe4"
	}

	if steelToken=="" && browserlessToken=="" {
		fmt.Println("no automation token configured; set DOOTA_START_COMMON_STEEL_API_KEY or DOOTA_START_COMMON_BROWSERLESS_API_KEY")
		os.Exit(1)
	}

	steel := browser_automation.NewSteelBrowser(steelToken, logger)
	browserless := browser_automation.NewBrowserLessBrowser(browserlessToken, browserlessWarmup, logger)
	provider := browser_automation.NewFallbackBrowserAutomation(steel, browserless, logger)

	debugStore, _ := dstore.NewStore("data/debugstore", "", "", true)
	redditBA := browser_automation.NewRedditBrowserAutomation(provider, logger, debugStore)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	cdp, err := redditBA.StartLogin(ctx, "")
	if err != nil {
		fmt.Println("StartLogin error:", err)
		os.Exit(2)
	}

	fmt.Println("LiveURL:", cdp.LiveURL)
	fmt.Println("WSEndpoint:", cdp.WSEndpoint)
	fmt.Println("Open the LiveURL in your browser to perform the Reddit login. I will poll for login for up to 30 minutes.")

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

	browser, err := pw.Chromium.ConnectOverCDP(cdp.WSEndpoint)
	if err != nil {
		fmt.Println("ConnectOverCDP error:", err)
		os.Exit(5)
	}
	defer browser.Close()

	pageContext := browser.Contexts()[0]
	page := pageContext.Pages()[0]

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("timeout waiting for login")
			os.Exit(6)
		case <-ticker.C:
			// Try to read display name; don't abort on errors
			displayName, _ := page.Locator("rs-current-user").GetAttribute("display-name")
			if displayName != "" {
				cookies, err := pageContext.Cookies()
				if err != nil {
					fmt.Println("failed to read cookies:", err)
					continue
				}
				if len(cookies) == 0 {
					fmt.Println("found display name but no cookies yet")
					continue
				}
				marshal, _ := json.Marshal(cookies)
				fmt.Println("Login detected: username=", displayName)
				fmt.Println("Cookies length bytes:", len(marshal))
				fmt.Println("Cookies (truncated 2000 chars):")
				if len(marshal) > 2000 {
					fmt.Println(string(marshal[:2000]))
				} else {
					fmt.Println(string(marshal))
				}
				// Keep session alive long enough for post processing, then exit
				os.Exit(0)
			}
			// otherwise continue polling
		}
	}
}
