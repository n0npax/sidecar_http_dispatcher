#!/usr/bin/env python3
import yaml
import requests
import sys
import logging
from quart import Quart, request, Response

logging.basicConfig(stream=sys.stdout, level=logging.INFO)
app_name = "sidecar http dispatcher"
logger = logging.getLogger(app_name)

app = Quart(__name__)
CONFIG_FILE = "config.yaml"


def read_config():
    with open(CONFIG_FILE, "r") as f:
        return yaml.safe_load(f.read())


class ConfigMeta(type):
    def __new__(
        config_metaclass, future_class_name, future_class_parents, future_class_attr
    ):
        new_attrs = {}
        for name, val in read_config().items():
            new_attrs[name] = val
        return type(future_class_name, future_class_parents, new_attrs)


class Config(metaclass=ConfigMeta):
    pass


config = Config()


@app.route("/", methods=("GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"))
async def hello():
    new_headers, destination = {}, config.destination
    matched_header = request.headers.get(config.key)
    if matched_header in config.rewrites:
        rules = config.rewrites[matched_header]
        logger.info(f"patching headers: {rules['patch']}")
        for rule in rules["patch"]:
            new_headers[rule["key"]] = rule["val"]
        destination = rules.get(destination, destination)
    request.headers.update(new_headers)
    # cannot use **request. type(Quart.request) != (restests.restuest)
    resp = requests.request(
        request.method, destination, headers=request.headers, cookies=request.cookies
    )
    return Response(response=resp)


def main():
    app.run()


if __name__ == "__main__":
    main()
