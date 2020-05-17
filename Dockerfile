FROM python:buster
LABEL maintainer=marcin.niemria@gmail.com
LABEL author="Marcin Niemira <n0npax>"
RUN pip install --upgrade pip
RUN pip install poetry
WORKDIR /app
COPY pyproject.toml .
RUN poetry env use system
RUN poetry config virtualenvs.create false --local
RUN poetry install --no-dev
COPY . .
ENTRYPOINT ./app.py