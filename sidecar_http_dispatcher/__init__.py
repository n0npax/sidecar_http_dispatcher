import logging
import sys

APP_NAME = "sidecar http dispatcher"

logging.basicConfig(stream=sys.stdout, level=logging.INFO)
logger = logging.getLogger(APP_NAME)
