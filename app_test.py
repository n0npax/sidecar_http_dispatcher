import asyncio
import pytest

from dataclasses import dataclass
from app import pass_request


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
