FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc libpq-dev wget unzip \
    && rm -rf /var/lib/apt/lists/* \
    && wget -q https://github.com/XTLS/Xray-core/releases/download/v26.5.3/Xray-linux-64.zip -O /tmp/xray.zip \
    && unzip -j /tmp/xray.zip -d /usr/local/bin/ \
    && chmod +x /usr/local/bin/xray \
    && rm /tmp/xray.zip

WORKDIR /app

RUN useradd -m -u 1000 appuser

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

USER root

COPY . .

RUN mkdir -p /app/logs && chown appuser:appuser /app/logs

USER appuser

EXPOSE 44042

WORKDIR /app

CMD sh -c "xray version && xray run -c /app/vpn.json > /app/logs/xray.log 2>&1 & sleep 3 && alembic upgrade head && python3 -m cmd.parser.main"
