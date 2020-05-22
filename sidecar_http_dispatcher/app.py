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
    new_headers, default_destination = {}, config.destination
    matched_header = request.headers.get(config.key)
    if matched_header in config.rewrites:
        rules = config.rewrites[matched_header]
        logger.info(f"patching headers: {rules['patch']}")
        for rule in rules["patch"]:
            new_headers[rule["key"]] = rule["val"]
        destination = rules.get("destination", default_destination)
    request.headers.update(new_headers)
    # TODO pass all properties to request
    return await pass_request(destination=f"{destination}/{path}", request=request)


async def pass_request(*, destination: str, request: Quart.request_class) -> Response:
    """Pass patched request to downstream services."""
    # cannot use **request. type(Quart.request) != type(session.request)
    async with aiohttp.ClientSession() as session:
        data = await request.get_data()
        logger.info(f"passing {request.method} to {destination}")
        async with session.request(
            request.method,
            destination,
            headers=request.headers,
            data=data,
            cookies=request.cookies,
            params=request.args,
        ) as response:
            resp_text, resp_status, resp_headers = (
                await response.text(),
                response.status,
                dict(response.headers),
            )
    # TODO pass all properties to response
    return Response(resp_text, status=resp_status, headers=resp_headers)
