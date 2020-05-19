import aiohttp
from quart import Quart, Response, request

from sidecar_http_dispatcher import logger
from sidecar_http_dispatcher.config import Config

ALL_HTTP_METHODS = ("GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS")


app = Quart(__name__)


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
