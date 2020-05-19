import os
import sys

import yaml

from sidecar_http_dispatcher import logger

CONFIG_FILE = os.environ.get("SIDECAR_CONFIG", "config.yaml")
PORT = int(os.environ.get("SIDECAR_PORT", 8192))


def read_config():
    """Read config from yaml."""
    with open(CONFIG_FILE, "r") as f:
        return yaml.safe_load(f.read())


class ConfigMeta(type):

    """ConfigMeta creates config class base on yaml definition."""

    def __new__(
        config_metaclass, future_class_name, future_class_parents, future_class_attr
    ):
        new_attrs = {}
        try:
            config_dict = read_config()
        except FileNotFoundError as e:
            logger.critical(f"cannot open config file: {e}")
            sys.exit(1)
        for name, val in config_dict.items():
            new_attrs[name] = val
        return type(future_class_name, future_class_parents, new_attrs)


class Config(metaclass=ConfigMeta):

    """Create config from meta."""
