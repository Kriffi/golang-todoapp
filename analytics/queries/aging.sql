WITH buckets AS (
    SELECT CASE
               WHEN days_after_due <= 3 THEN 1
               WHEN days_after_due <= 7 THEN 2
               ELSE 3
           END AS bucket
    FROM analytics.tasks
    WHERE overdue AND ({{ user_id }} = 0 OR user_id = {{ user_id }})
)
SELECT labels.age AS "Просрочка", count(buckets.bucket) AS "Задачи"
FROM (VALUES (1, '1–3 дня'), (2, '4–7 дней'), (3, '8 дней и более')) AS labels(id, age)
LEFT JOIN buckets ON buckets.bucket = labels.id
GROUP BY labels.id, labels.age
ORDER BY labels.id;
