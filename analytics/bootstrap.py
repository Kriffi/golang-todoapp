"""Create and verify the TodoApp dashboard in Redash 26.9."""

import os
import time
from pathlib import Path

from redash.handlers.setup import create_org
from redash.models import Organization, User
from redash.wsgi import app


HERE = Path(__file__).resolve().parent
HOST = os.environ.get("REDASH_HOST", "http://localhost:5001")
ADMIN_EMAIL = os.environ.get("REDASH_ADMIN_EMAIL", "todoapp@example.com")
DASHBOARD = "TodoApp · Аналитика задач"
DATASOURCE = "TodoApp · PostgreSQL demo"
USER_PARAM = {
    "name": "user_id",
    "title": "Пользователь (0 — все, 1–4 — ID)",
    "type": "number",
    "value": 0,
    "global": True,
}
REPORTS = (
    ("overview", "Сводка", "TABLE", None, (), (0, 0, 12, 4)),
    ("workload", "Задачи по пользователям", "CHART", "Пользователь",
     ("Открыто без просрочки", "Просрочено", "Выполнено"), (0, 4, 7, 7)),
    ("aging", "Давность просрочки", "CHART", "Просрочка", ("Задачи",), (7, 4, 5, 7)),
    ("throughput", "Создано и закрыто за 30 дней", "CHART", "День",
     ("Создано", "Закрыто"), (0, 11, 12, 7)),
    ("overdue", "Просроченные задачи", "TABLE", None, (), (0, 18, 12, 9)),
)


def read_admin_key():
    with app.app_context():
        org = Organization.query.filter_by(slug="default").one_or_none()
        if org is None:
            _, user = create_org("TodoApp", "TodoApp Analytics", ADMIN_EMAIL, os.environ["REDASH_ADMIN_PASSWORD"])
        else:
            user = User.query.filter_by(org_id=org.id, email=ADMIN_EMAIL).one()
        return user.api_key


def request(client, path, body=None):
    if body is None:
        response = client.get(path, base_url=HOST)
    else:
        response = client.post(path, json=body, base_url=HOST)
    if response.status_code >= 300:
        if path == "/api/data_sources":
            raise RuntimeError(f"{path}: HTTP {response.status_code}")
        raise RuntimeError(f"{path}: HTTP {response.status_code}: {response.get_data(as_text=True)[:300]}")
    return response.get_json()


def result_rows(client, query_id, user_id):
    result = request(client, f"/api/queries/{query_id}/results", {"parameters": {"user_id": user_id}, "max_age": 0})
    if "query_result" in result:
        return result["query_result"]["data"]["rows"]
    job = result["job"]
    for _ in range(60):
        job = request(client, f"/api/jobs/{job['id']}")["job"]
        if job["status"] == 3:
            result = request(client, f"/api/query_results/{job['query_result_id']}.json")
            return result["query_result"]["data"]["rows"]
        if job["status"] in (4, 5):
            raise RuntimeError(job.get("error") or "Redash query failed")
        time.sleep(0.5)
    raise TimeoutError(f"Redash did not finish query {query_id}")


def display_options(x, series, key):
    if x is None:
        return {}
    chart_type = "line" if key == "throughput" else "column"
    return {
        "globalSeriesType": chart_type,
        "sortX": key == "throughput",
        "legend": {"enabled": len(series) > 1},
        "xAxis": {"type": "datetime" if key == "throughput" else "category"},
        "yAxis": [{"type": "linear"}, {"type": "linear", "opposite": True}],
        "columnMapping": {x: "x", **{name: "y" for name in series}},
    }


def verify_default_results(rows):
    overview = rows["overview"][0]
    overdue = rows["overdue"]
    assert overview["Всего"] == 64
    assert overview["Открыто"] + overview["Выполнено"] == overview["Всего"]
    assert overview["Просрочено"] == len(overdue) > 0
    assert sum(row["Просрочено"] for row in rows["workload"]) == len(overdue)
    assert sum(row["Задачи"] for row in rows["aging"]) == len(overdue)
    assert len(rows["throughput"]) == 30
    assert all(row["Дней просрочки"] > 0 for row in overdue)


def main():
    client = app.test_client()
    client.environ_base["HTTP_AUTHORIZATION"] = f"Key {read_admin_key()}"
    sources = request(client, "/api/data_sources")
    source = next((item for item in sources if item["name"] == DATASOURCE), None)
    if source is None:
        source = request(client, "/api/data_sources", {
            "name": DATASOURCE,
            "type": "pg",
            "options": {
                "host": os.environ.get("REPORTING_HOST", "postgres"),
                "port": int(os.environ.get("REPORTING_PORT", "5432")),
                "dbname": os.environ.get("REPORTING_DB", "todoapp_demo"),
                "sslmode": os.environ.get("REPORTING_SSLMODE", "prefer"),
                "user": "todoapp_reporting",
                "password": os.environ["REPORTING_PASSWORD"],
            },
        })
    dashboards = request(client, "/api/dashboards/my")["results"]
    dashboard = next((item for item in dashboards if item["name"] == DASHBOARD), None)
    if dashboard is None:
        dashboard = request(client, "/api/dashboards", {"name": DASHBOARD})
    dashboard_path = f"/api/dashboards/{dashboard['id']}"
    dashboard = request(client, dashboard_path)
    widgets = {widget["visualization"]["name"]: widget
               for widget in dashboard["widgets"] if widget.get("visualization")}
    existing_queries = request(client, "/api/queries/my?page_size=100")["results"]
    results = {}
    for key, name, visual_type, x, series, (col, row, width, height) in REPORTS:
        query_name = f"TodoApp · {name}"
        query = next((item for item in existing_queries if item["name"] == query_name), None)
        if query is None:
            query = request(client, "/api/queries", {
                "name": query_name,
                "query": (HERE / "queries" / f"{key}.sql").read_text(encoding="utf-8"),
                "data_source_id": source["id"],
                "options": {"parameters": [USER_PARAM]},
            })
        query_id = query["id"]
        results[key] = result_rows(client, query_id, 0)
        if key == "overview":
            assert result_rows(client, query_id, 1)[0]["Всего"] == 16
            assert result_rows(client, query_id, -1)[0]["Всего"] == 0
        if name in widgets:
            continue
        query = request(client, f"/api/queries/{query_id}")
        visual = next((item for item in query["visualizations"] if item["name"] == name), None)
        if visual is None:
            visual = request(client, "/api/visualizations", {
                "query_id": query_id,
                "name": name,
                "type": visual_type,
                "description": "",
                "options": display_options(x, series, key),
            })
        request(client, "/api/widgets", {
            "dashboard_id": dashboard["id"],
            "visualization_id": visual["id"],
            "text": "",
            "options": {"position": {"col": col, "row": row, "sizeX": width, "sizeY": height}},
            "width": 1,
        })
        request(client, f"/api/queries/{query_id}", {"is_draft": False})
    verify_default_results(results)
    dashboard = request(client, dashboard_path, {"is_draft": False})
    print(f"{HOST}/dashboards/{dashboard['slug']}")
    if os.environ.get("REDASH_PUBLIC_DEMO") == "1":
        public_url = dashboard.get("public_url") or request(client, f"{dashboard_path}/share", {})["public_url"]
        if target := os.environ.get("REDASH_DEMO_REDIRECT_FILE"):
            from urllib.parse import urlsplit

            url = urlsplit(public_url)
            Path(target).write_text(f"{url.path}?{url.query}", encoding="utf-8")
        else:
            print(public_url)


if __name__ == "__main__":
    main()
