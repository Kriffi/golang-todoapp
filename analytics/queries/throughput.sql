WITH events AS (
    SELECT opened_on AS day, 1 AS opened, 0 AS closed
    FROM analytics.tasks WHERE ({{ user_id }} = 0 OR user_id = {{ user_id }})
    UNION ALL
    SELECT completed_on, 0, 1
    FROM analytics.tasks
    WHERE completed AND ({{ user_id }} = 0 OR user_id = {{ user_id }})
),
days AS (
    SELECT ((CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date - offset_days) AS day
    FROM generate_series(0, 29) AS offsets(offset_days)
)
SELECT days.day AS "День",
       coalesce(sum(events.opened), 0) AS "Создано",
       coalesce(sum(events.closed), 0) AS "Закрыто"
FROM days
LEFT JOIN events ON events.day = days.day
GROUP BY days.day
ORDER BY days.day;
