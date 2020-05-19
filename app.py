#!/usr/bin/env python3
import logging
import os
import sys

import aiohttp
import yaml
from quart import Quart, Response, request

ALL_HTTP_METHODS = ("GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS")
CONFIG_FILE = os.environ.get("SIDECAR_CONFIG", "config.yaml")
APP_NAME = "sidecar http dispatcher"

logging.basicConfig(stream=sys.stdout, level=logging.INFO)
logger = logging.getLogger(APP_NAME)

app = Quart(__name__)


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
            os.exit(1)
        for name, val in config_dict.items():
            new_attrs[name] = val
        return type(future_class_name, future_class_parents, new_attrs)


class Config(metaclass=ConfigMeta):

    """Create config from meta."""


config = Config()


@app.route("/", defaults={"path": ""}, methods=ALL_HTTP_METHODS)
@app.route("/<path:path>", methods=ALL_HTTP_METHODS)
async def dispatch_and_pass(path) -> Response:
    """Patch request base on config and pass it to downstream."""
    new_headers, destination = {}, config.destination
    matched_header = request.headers.get(config.key)
    if matched_header in config.rewrites:
        rules = config.rewrites[matched_header]
        logger.info(f"patching headers: {rules['patch']}")
        for rule in rules["patch"]:
            new_headers[rule["key"]] = rule["val"]
        destination = rules.get("destination", destination)
    request.headers.update(new_headers)
    return await pass_request(destination=f"{destination}/{path}", request=request)


async def pass_request(*, destination: str, request: Quart.request_class) -> Response:
    """Pass patched request to downstream services."""
    # cannot use **request. type(Quart.request) != type(session.request)
    async with aiohttp.ClientSession() as session:
        data = await request.get_data()
        async with session.request(
            request.method,
            destination,
            headers=request.headers,
            data=data,
            cookies=request.cookies,
            params=request.args,
        ) as response:
            resp_text, resp_status = await response.text(), response.status
    return Response(resp_text, status=resp_status)


if __name__ == "__main__":
    app.run()
