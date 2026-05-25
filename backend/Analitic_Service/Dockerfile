FROM python:3.11-slim

USER root

RUN echo 'Acquire::ForceIPv4 "true";' > /etc/apt/apt.conf.d/99force-ipv4 && \
    (sed -i 's/deb.debian.org/mirror.yandex.ru/g' /etc/apt/sources.list.d/debian.sources || \
     sed -i 's/deb.debian.org/mirror.yandex.ru/g' /etc/apt/sources.list)

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

CMD ["/bin/sh", "-c", "sleep 5 && python3 -m command.analytic.main"]
