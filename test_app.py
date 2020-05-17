from dataclasses import dataclass

import pytest

from app import Config, config, pass_request, read_config


@pytest.fixture
def destination():
    return "https://example.com"


@dataclass
class DummyRequest:
    headers: tuple = (("Host", "test"), ("foo", "Bar"))


def dummy_request_factory(method):
    class Req(DummyRequest):
        def __init__(self, method):
            self.method: str = method

    return Req(method)


@pytest.mark.asyncio
@pytest.mark.parametrize(
    "method", ("GET", "PUT", "HEAD", "POST", "OPTIONS"),
)
async def test_pass_request(method, destination):
    resp = await pass_request(
        request=dummy_request_factory(method), destination=destination
    )
    assert 200 <= resp.status_code < 500


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
