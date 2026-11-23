"""Expose a stable link to the regenerated public demo dashboard."""

import os
from pathlib import Path

from werkzeug.wrappers import Request, Response

from redash.wsgi import app as redash


@Request.application
def app(request):
    if request.path in ("/", "/demo"):
        url = Path(os.environ["REDASH_DEMO_REDIRECT_FILE"]).read_text(encoding="utf-8")
        return Response("", status=302, headers={"Location": url, "Cache-Control": "no-store"})
    response = Response.from_app(redash, request.environ)
    if request.path.startswith("/public/dashboards/"):
        response.headers["Cache-Control"] = "no-store"
    return response
