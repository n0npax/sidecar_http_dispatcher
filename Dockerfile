FROM python:buster
LABEL maintainer=marcin.niemria@gmail.com
LABEL author="Marcin Niemira <n0npax>"
RUN pip install --upgrade pip
RUN pip install poetry
WORKDIR /app
COPY pyproject.toml .
RUN poetry install
COPY . .
RUN python app.py
ENTRYPOINT ./app.py