SELECT id AS "ID", title AS "Задача", task_user AS "Пользователь",
       due_date AS "Срок (UTC)", days_after_due AS "Дней просрочки"
FROM analytics.tasks
WHERE overdue AND ({{ user_id }} = 0 OR user_id = {{ user_id }})
ORDER BY days_after_due DESC, id;
