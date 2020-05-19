import logging
import os
import sys

APP_NAME = "sidecar http dispatcher"
SIDECAR_LOGGING_LEVEL = os.environ.get("SIDECAR_LOGGING_LEVEL", "").upper()


logging.basicConfig(
    stream=sys.stdout, level=getattr(logging, SIDECAR_LOGGING_LEVEL, logging.INFO),
)
logger = logging.getLogger(APP_NAME)
