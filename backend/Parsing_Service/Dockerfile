FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc libpq-dev curl \
    && rm -rf /var/lib/apt/lists/*

RUN curl -Ls https://raw.githubusercontent.com/XTLS/Xray-install/main/install-release.sh | bash

WORKDIR /app

RUN useradd -m -u 1000 appuser

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
USER root

RUN mkdir -p /app/logs && chown -R appuser:appuser /app

USER appuser

EXPOSE 44042

CMD sh -c "xray run -c /app/vpn.json > /app/logs/xray.log 2>&1 & sleep 3 && alembic upgrade head && python3 -m cmd.parser.main"
