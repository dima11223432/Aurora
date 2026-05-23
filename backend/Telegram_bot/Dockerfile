FROM python:3.11-slim

USER root

RUN echo 'Acquire::ForceIPv4 "true";' > /etc/apt/apt.conf.d/99force-ipv4

RUN sed -i 's/deb.debian.org/mirror.yandex.ru/g' /etc/apt/sources.list.d/debian.sources || \
    sed -i 's/deb.debian.org/mirror.yandex.ru/g' /etc/apt/sources.list

WORKDIR /app

RUN useradd -m -u 1000 appuser

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 44042

CMD sh -c "xray version && xray run -c /app/vpn.json > /dev/null 2>&1 & sleep 3 && python3 src/main.py"
