package main

import (
	"context"
	"fmt"
	"os"
	"time"

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
	info, err := steel.GetCDPInfo(ctx, browser_automation.CDPInput{StartURL: "https://www.reddit.com", UseProxy: false})
	if err != nil {
		fmt.Println("GetCDPInfo error:", err)
		os.Exit(2)
	}

	fmt.Printf("WSEndpoint: %s\n", info.WSEndpoint)
	fmt.Printf("LiveURL: %s\n", info.LiveURL)
}
