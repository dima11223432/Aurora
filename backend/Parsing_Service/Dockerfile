FROM python:3.11-slim

RUN echo 'Acquire::ForceIPv4 "true";' > /etc/apt/apt.conf.d/99force-ipv4

RUN sed -i 's/deb.debian.org/mirror.yandex.ru/g' /etc/apt/sources.list.d/debian.sources || \
    sed -i 's/deb.debian.org/mirror.yandex.ru/g' /etc/apt/sources.list

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

COPY config/config.env /app/config/config.env
COPY vpn.json /app/vpn.json

COPY . .

RUN chown -R appuser:appuser /app

RUN mkdir -p /app/logs && chown appuser:appuser /app/logs

USER appuser

COPY --chown=appuser:appuser pars_session.session /app/pars_session.session

EXPOSE 44042

WORKDIR /app

CMD sh -c "xray version && xray run -c /app/vpn.json > /dev/null 2>&1 & sleep 3 && alembic upgrade head && python3 -m cmd.parser.main"
