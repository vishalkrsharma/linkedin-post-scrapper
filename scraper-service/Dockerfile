FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN python -m venv /app/.venv && /app/.venv/bin/pip install -r requirements.txt
RUN /app/.venv/bin/pip install playwright && /app/.venv/bin/playwright install chromium --with-deps

COPY . .

EXPOSE 8000

CMD ["/app/.venv/bin/python", "main.py"]