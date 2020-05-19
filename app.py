#!/usr/bin/env python3

from sidecar_http_dispatcher.app import app
from sidecar_http_dispatcher.config import PORT

if __name__ == "__main__":
    app.run(port=PORT)
