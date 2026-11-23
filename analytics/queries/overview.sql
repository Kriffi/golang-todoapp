SELECT count(*) AS "Всего",
       count(*) FILTER (WHERE completed) AS "Выполнено",
       count(*) FILTER (WHERE NOT completed) AS "Открыто",
       count(*) FILTER (WHERE overdue) AS "Просрочено",
       count(*) FILTER (WHERE NOT completed AND due_date BETWEEN
           (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date AND
           (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date + 7) AS "Срок в ближайшие 7 дней",
       round(100.0 * count(*) FILTER (WHERE completed) / NULLIF(count(*), 0), 1) AS "Выполнено, %",
       round((avg(extract(epoch FROM time_to_complete)) / 3600)::numeric, 1) AS "Среднее время выполнения, ч"
FROM analytics.tasks
WHERE ({{ user_id }} = 0 OR user_id = {{ user_id }});
