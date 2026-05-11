from contextlib import asynccontextmanager
from typing import Optional

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import uvicorn

from linkedin import LinkedInScraper


class ScrapeRequest(BaseModel):
    keyword: str
    location: Optional[str] = "India"


class JobPost(BaseModel):
    title: str
    company: str
    location: str
    url: str
    posted_at: Optional[str] = None


scraper: Optional[LinkedInScraper] = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global scraper
    scraper = LinkedInScraper()
    async with scraper:
        yield


app = FastAPI(title="LinkedIn Scraper Service", lifespan=lifespan)


@app.get("/health")
async def health_check():
    return {"status": "healthy"}


@app.post("/scrape")
async def scrape_jobs(request: ScrapeRequest):
    if not scraper:
        raise HTTPException(status_code=500, detail="Scraper not initialized")

    try:
        async with scraper:
            jobs = await scraper.search_jobs(request.keyword, request.location or "India")
            return {"jobs": jobs, "count": len(jobs)}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Scraping failed: {str(e)}")


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)