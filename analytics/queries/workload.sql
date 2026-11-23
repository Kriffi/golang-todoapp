SELECT task_user AS "Пользователь",
       count(*) FILTER (WHERE NOT completed AND overdue IS NOT TRUE) AS "Открыто без просрочки",
       count(*) FILTER (WHERE overdue) AS "Просрочено",
       count(*) FILTER (WHERE completed) AS "Выполнено"
FROM analytics.tasks
WHERE ({{ user_id }} = 0 OR user_id = {{ user_id }})
GROUP BY user_id, task_user
ORDER BY count(*) FILTER (WHERE overdue) DESC, task_user;
