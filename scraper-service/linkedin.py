import asyncio
import json
import random
from datetime import datetime
from typing import Optional

from playwright.async_api import async_playwright, Browser, Page, Playwright


class LinkedInScraper:
    def __init__(self, cookies_path: str = "cookies.json"):
        self.cookies_path = cookies_path
        self.playwright: Optional[Playwright] = None
        self.browser: Optional[Browser] = None

    async def __aenter__(self):
        self.playwright = await async_playwright().start()
        self.browser = await self.playwright.chromium.launch(headless=True)
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self.browser:
            await self.browser.close()
        if self.playwright:
            await self.playwright.stop()

    async def _load_cookies(self, page: Page):
        try:
            with open(self.cookies_path, "r") as f:
                cookies = json.load(f)
            await page.context.add_cookies(cookies)
        except FileNotFoundError:
            print(f"Warning: Cookies file {self.cookies_path} not found")

    def _random_delay(self, min_sec: float = 3.0, max_sec: float = 5.0):
        return random.uniform(min_sec, max_sec)

    async def _extract_job_data(self, page: Page) -> list[dict]:
        jobs = []
        await asyncio.sleep(2)

        job_cards = await page.locator(".job-card-container").all()

        for card in job_cards:
            try:
                job_data = {}

                title_elem = card.locator(".job-card-list__title")
                job_data["title"] = await title_elem.text_content() if await title_elem.count() > 0 else ""

                company_elem = card.locator(".job-card-container__company-name")
                job_data["company"] = await company_elem.text_content() if await company_elem.count() > 0 else ""

                location_elem = card.locator(".job-card-container__metadata-item")
                job_data["location"] = await location_elem.text_content() if await location_elem.count() > 0 else ""

                url_elem = card.locator(".job-card-list__title")
                if await url_elem.count() > 0:
                    href = await url_elem.get_attribute("href")
                    job_data["url"] = f"https://www.linkedin.com{href}" if href else ""

                posted_elem = card.locator(".job-card-container__listed-time")
                job_data["posted_at"] = await posted_elem.text_content() if await posted_elem.count() > 0 else ""

                if job_data.get("url"):
                    jobs.append(job_data)

            except Exception as e:
                print(f"Error extracting job card: {e}")
                continue

        return jobs

    async def search_jobs(self, keyword: str, location: str = "India") -> list[dict]:
        search_query = f"{keyword} jobs {location}"
        encoded_query = search_query.replace(" ", "%20")

        url = f"https://www.linkedin.com/jobs/search/?keywords={encoded_query}&location={location.replace(' ', '%20')}"

        context = await self.browser.new_context(
            viewport={"width": 1920, "height": 1080},
            user_agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        )

        page = await context.new_page()
        await self._load_cookies(page)

        await page.goto(url, wait_until="networkidle")

        await asyncio.sleep(self._random_delay())

        jobs = await self._extract_job_data(page)

        await context.close()

        return jobs