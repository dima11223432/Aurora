FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc libpq-dev curl \
    && rm -rf /var/lib/apt/lists/*

RUN curl -L -o /usr/local/bin/xray https://github.com/XTLS/Xray/releases/download/v1.8.6/xray-linux-64 && \
    chmod +x /usr/local/bin/xray

WORKDIR /app

RUN useradd -m -u 1000 appuser

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
USER root

RUN mkdir -p /app/logs && chown -R appuser:appuser /app

USER appuser

EXPOSE 44042

CMD sh -c "alembic upgrade head && python3 -m cmd.parser.main"
