FROM python:buster

RUN pip install --upgrade pip
RUN pip install poetry
WORKDIR /app
COPY . .
