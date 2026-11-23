WITH calendar AS (
    SELECT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date AS today
),
tasks AS (
    SELECT completed,
           (created_at AT TIME ZONE 'UTC')::date AS created_on,
           (completed_at AT TIME ZONE 'UTC')::date AS completed_on,
           completed_at - created_at AS completion_time,
           due_date,
           calendar.today,
           NOT completed AND due_date < calendar.today AS overdue
    FROM todoapp.tasks CROSS JOIN calendar
),
overview AS (
    SELECT count(*) AS total,
           count(*) FILTER (WHERE completed) AS completed,
           count(*) FILTER (WHERE NOT completed) AS open,
           count(*) FILTER (WHERE overdue) AS overdue,
           count(*) FILTER (WHERE NOT completed AND due_date BETWEEN today AND today + 7) AS due_soon,
           round(100.0 * count(*) FILTER (WHERE completed) / nullif(count(*), 0), 1) AS completion_rate,
           round((avg(extract(epoch FROM completion_time)) / 3600)::numeric, 1) AS average_hours
    FROM tasks
),
status AS (
    SELECT values_by_status.ordinal, values_by_status.label AS "Статус", values_by_status.amount AS "Задачи"
    FROM overview
    CROSS JOIN LATERAL (VALUES
        (1, 'Открыто без просрочки', open - overdue),
        (2, 'Просрочено', overdue),
        (3, 'Выполнено', completed)
    ) AS values_by_status(ordinal, label, amount)
),
aging AS (
    SELECT 1 AS ordinal, '1–3 дня' AS "Просрочка", count(*) FILTER (WHERE today - due_date <= 3) AS "Задачи"
    FROM tasks WHERE overdue
    UNION ALL
    SELECT 2, '4–7 дней', count(*) FILTER (WHERE today - due_date BETWEEN 4 AND 7)
    FROM tasks WHERE overdue
    UNION ALL
    SELECT 3, '8 дней и более', count(*) FILTER (WHERE today - due_date >= 8)
    FROM tasks WHERE overdue
),
events AS (
    SELECT created_on AS day, count(*) AS created, 0 AS completed FROM tasks GROUP BY created_on
    UNION ALL
    SELECT completed_on, 0, count(*) FROM tasks WHERE completed GROUP BY completed_on
),
throughput AS (
    SELECT today - offset_days AS "День",
           coalesce(sum(events.created), 0) AS "Создано",
           coalesce(sum(events.completed), 0) AS "Закрыто"
    FROM calendar CROSS JOIN generate_series(0, 29) AS days(offset_days)
    LEFT JOIN events ON events.day = today - offset_days
    GROUP BY today - offset_days
),
due_soon AS (
    SELECT calendar.today + offset_days AS "День", count(tasks.due_date) AS "Задачи"
    FROM calendar CROSS JOIN generate_series(0, 7) AS days(offset_days)
    LEFT JOIN tasks ON tasks.due_date = calendar.today + offset_days AND NOT tasks.completed
    GROUP BY calendar.today + offset_days
)
SELECT jsonb_build_object(
    'overview', jsonb_build_object(
        'Всего', total, 'Выполнено', completed, 'Открыто', open,
        'Просрочено', overdue, 'Срок в ближайшие 7 дней', due_soon,
        'Доля выполнения (%)', completion_rate, 'Среднее время выполнения (ч)', average_hours
    ),
    'status', (SELECT jsonb_agg(to_jsonb(status) - 'ordinal' ORDER BY ordinal) FROM status),
    'aging', (SELECT jsonb_agg(to_jsonb(aging) - 'ordinal' ORDER BY ordinal) FROM aging),
    'throughput', (SELECT jsonb_agg(to_jsonb(throughput) ORDER BY "День") FROM throughput),
    'due_soon', (SELECT jsonb_agg(to_jsonb(deadline) ORDER BY deadline."День") FROM due_soon AS deadline)
) FROM overview;
