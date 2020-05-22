from dataclasses import dataclass

import pytest

from sidecar_http_dispatcher.app import app, config, pass_request
from sidecar_http_dispatcher.config import Config, read_config


@pytest.fixture
def destination():
    return "https://example.com"


@pytest.fixture
def dev_headers():
    return {"environment": "dev"}


@pytest.fixture
def testapp():
    return app


@dataclass
class DummyRequest:
    headers: tuple = (("Host", "example.com"), ("foo", "Bar"))
    data: str = ""
    cookies: str = ""


def dummy_request_factory(*, method, data):
    class Req(DummyRequest):
        def __init__(self, method, data=""):
            self.method: str = method
            self.data = data
            self.args = None

        async def get_data(self):
            return self.data

    return Req(method)


@pytest.mark.parametrize("path,code", (("/", 200), ("/foo", 404)))
@pytest.mark.asyncio
async def test_app_dev_headers(path, code, testapp, dev_headers):
    client = testapp.test_client()
    response = await client.get(path, headers=dev_headers)
    assert response.status_code == code


# default path is pointing to not reachable url
@pytest.mark.parametrize("path,code", (("/", 500),))
@pytest.mark.asyncio
async def test_app_no_headers(path, code, testapp):
    client = testapp.test_client()
    response = await client.get(path)
    assert response.status_code == code


@pytest.mark.asyncio
@pytest.mark.parametrize(
    "method,data",
    (
        ("GET", None),
        ("HEAD", None),
        ("POST", "Foo"),
        ("OPTIONS", None),
        # ("DELETE", None)
        # ("PUT", "Foo"), not supported by example.com
    ),
)
async def test_pass_request(method, data, destination):
    resp = await pass_request(
        request=dummy_request_factory(method=method, data=data), destination=destination
    )
    assert 200 == resp.status_code


def test_read_config():
    config_dict = read_config()
    assert config_dict
    assert isinstance(config_dict, dict)


def test_config_object():
    assert config
    assert isinstance(config, Config)
    config_dict = read_config()
    for key, value in config_dict.items():
        assert getattr(config, key) == value
