FROM python:buster
LABEL maintainer=marcin.niemria@gmail.com
LABEL author="Marcin Niemira <n0npax>"
RUN apt-get update
RUN apt-get install nginx -y
RUN pip install --upgrade pip==20.1.1
RUN pip install poetry==1.0.5 uvicorn==0.11.5
WORKDIR /app
RUN mkdir /app/config
COPY pyproject.toml .
RUN poetry env use system
RUN poetry config virtualenvs.create false --local
RUN poetry install --no-dev
COPY sidecar_http_dispatcher sidecar_http_dispatcher
COPY app.py app.py
COPY _nginx.conf /etc/nginx/nginx.conf
ENV SIDECAR_PORT 5000
ENTRYPOINT sed -i "s/SIDECAR_PORT/$SIDECAR_PORT/g" /etc/nginx/nginx.conf && nginx && uvicorn app:app --uds /tmp/sidecar.sock
