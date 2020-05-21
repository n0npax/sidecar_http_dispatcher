FROM python:buster
LABEL maintainer=marcin.niemria@gmail.com
LABEL author="Marcin Niemira <n0npax>"
RUN pip install --upgrade pip==20.1.1
RUN pip install poetry==1.0.5 uvicorn==0.11.5
WORKDIR /app
RUN mkdir /app/config
COPY pyproject.toml .
RUN poetry env use system
RUN poetry config virtualenvs.create false --local
RUN poetry install --no-dev
COPY . .
ENV SIDECAR_PORT=8192
ENTRYPOINT uvicorn app:app -b localhost.:${SIDECAR_PORT}