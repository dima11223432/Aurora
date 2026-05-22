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

CMD ["python", "src/main.py"]
