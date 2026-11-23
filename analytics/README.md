# Аналитика задач в Redash

Дашборд показывает нагрузку по пользователям, сроки и динамику выполнения задач. В PostgreSQL загружаются 64 демонстрационные задачи и четыре вымышленных пользователя. Redash читает представление `analytics.tasks` через отдельную роль с правом `SELECT`.

| Отчёт | Что показывает |
| --- | --- |
| [Сводка](queries/overview.sql) | Открытые, выполненные и просроченные задачи, долю выполнения, среднее время выполнения и сроки на ближайшие 7 дней |
| [Задачи по пользователям](queries/workload.sql) | Задачи по трём статусам для каждого пользователя, по убыванию числа просрочек |
| [Давность просрочки](queries/aging.sql) | Сколько задач просрочено на 1–3, 4–7 и 8+ дней |
| [Динамика за 30 дней](queries/throughput.sql) | Сколько задач создали и закрыли за каждый день, включая дни без событий |
| [Просроченные задачи](queries/overdue.sql) | Список задач с датой и днями просрочки для приоритизации |

Общий параметр `Пользователь` применяется ко всем отчётам. `0` означает всех пользователей, `1` — Анну Орлову, `2` — Михаила Волкова, `3` — Ирину Соколову, `4` — Дениса Морозова.

Просроченной считается открытая задача со сроком раньше текущей даты UTC. Задачи без срока, выполненные задачи и задачи со сроком сегодня в просрочку не входят. Это совпадает с логикой `GET /api/v1/tasks?overdue=true` в TodoApp.

## Запуск демо

Из корня репозитория с Docker в Linux-контейнерах:

```sh
docker build -f analytics/Dockerfile -t todoapp-redash .
docker run --rm --name todoapp-redash-demo --memory 512m --memory-swap 512m -p 127.0.0.1:5002:5000 -e REDASH_HOST=http://localhost:5002 todoapp-redash
```

Открой [дашборд](http://localhost:5002/demo) после запуска Gunicorn. При каждом старте контейнер создаёт демонстрационные данные и дашборд заново. В GitHub доступны исходные SQL-запросы.

## Локальный Redash с сохранением настроек

Создай `analytics/.env` и задай пять разных случайных значений из букв и цифр. Используй длинные значения, подходящие для паролей.

```dotenv
REDASH_DB_PASSWORD=
REPORTING_PASSWORD=
REDASH_COOKIE_SECRET=
REDASH_SECRET_KEY=
REDASH_ADMIN_PASSWORD=
```

Выполни из корня репозитория:

```sh
docker compose --env-file analytics/.env -f analytics/compose.yaml run --rm initialize
docker compose --env-file analytics/.env -f analytics/compose.yaml up -d
docker compose --env-file analytics/.env -f analytics/compose.yaml exec -T -e PYTHONPATH=/app server python /analytics/bootstrap.py
```

Открой [Redash](http://localhost:5001/) и войди как `todoapp@example.com` с паролем `REDASH_ADMIN_PASSWORD`. Скрипт создаёт отчёты и сверяет итоги просрочки между таблицей, сводкой и графиками. Для следующих запусков используй `docker compose up` с тем же файлом Compose и сохранёнными значениями `.env`.
