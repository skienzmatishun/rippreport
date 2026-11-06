"""One-page Playwright adapter. Go owns verification, persistence, and cleanup."""
import json
import sys
from playwright.sync_api import sync_playwright

request = json.load(sys.stdin)
with sync_playwright() as playwright:
    browser = playwright.chromium.launch(headless=True)
    try:
        page = browser.new_page(
            viewport={"width": 1920, "height": 1080},
            user_agent="Mozilla/5.0 (compatible; ArchiveVerifier/2.0)",
        )
        response = page.goto(
            request["url"], wait_until="domcontentloaded", timeout=request["timeout_ms"]
        )
        page.wait_for_timeout(request["wait_ms"])
        screenshot_error = ""
        if request.get("screenshot"):
            try:
                page.screenshot(path=request["screenshot"], full_page=True)
            except Exception as error:
                screenshot_error = str(error)
        print(json.dumps({
            "html": page.content(), "url": page.url,
            "status": response.status if response else 0,
            "screenshot_error": screenshot_error,
        }))
    finally:
        browser.close()
