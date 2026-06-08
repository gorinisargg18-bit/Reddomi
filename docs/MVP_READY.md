# Reddomi MVP readiness

This document captures the minimum path to make the current repo feel launch-ready for the Reddomi MVP.

## What is already in place
- Local infrastructure stack for Postgres, Redis, and pgweb
- Go backend worker and portal bootstrap flow
- Frontend portal app for dashboard and onboarding
- Subscription/automation plan logic for usage limits
- Reddit automation paths via browser cookie-based login for MVP only

## MVP launch blockers to confirm
1. Reddit credentials or browser automation are configured for customer demos
2. AI generation keys (OpenAI/LiteLLM) are enabled for DM and comment suggestions
3. Auth0/portal environment values are set for real sign-in
4. Basic billing, trial, and usage limit paths are wired for the first customer
5. Product copy, onboarding, and support docs use the Reddomi name consistently

## Recommended launch path
- Keep the MVP narrow: cookie-based DM automation, lead discovery, plan limits, and onboarding
- Use the cookie-based Reddit flow for the first test because it does not require official Reddit API credentials
- Avoid overbuilding before validating the first real customer workflow
- Use this repo as the core SaaS foundation and add only the features customers explicitly request

## Verification checklist
- `devel/up.sh` starts local infrastructure
- `go run ./cmd/doota start portal-api redora-spooler` boots the backend services
- `pnpm --dir frontend dev:portal` starts the portal UI
- The portal page loads with the updated Reddomi branding
