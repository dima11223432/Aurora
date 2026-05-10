FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN useradd -m -u 1000 appuser
COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

RUN chown -R appuser:appuser /app

RUN mkdir -p /app/logs && chown appuser:appuser /app/logs

USER appuser

EXPOSE 44043

CMD ["python3", "-m", "command.analytic.main"]
