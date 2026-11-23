\set ON_ERROR_STOP on
CREATE DATABASE todoapp_demo;
\connect todoapp_demo
\i /migrations/000001_init.up.sql
\i /migrations/000002_task_due_date.up.sql

-- Assignments and titles are fictional. Relative dates keep the first run useful.
INSERT INTO todoapp.users (full_name)
VALUES ('Анна Орлова'), ('Михаил Волков'), ('Ирина Соколова'), ('Денис Морозов');

WITH demo AS (
    SELECT n,
           (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date - (n % 30 + 2) AS opened,
           n % 3 <> 0 AND n % 7 <> 0 AS completed,
           CASE WHEN n % 9 = 0 THEN NULL
                ELSE (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date - (n % 30 + 2) + (n % 9 + 1)
           END AS due
    FROM generate_series(1, 64) AS series(n)
)
INSERT INTO todoapp.tasks (title, completed, created_at, completed_at, author_user_id, due_date)
SELECT (ARRAY['Проверить интеграцию', 'Разобрать обращение', 'Обновить инструкцию',
              'Проверить квитанцию', 'Уточнить сценарий', 'Подготовить отчёт',
              'Проверить настройку'])[1 + (n % 7)] || ' · ' || n,
       completed,
       (opened + time '09:00') AT TIME ZONE 'UTC',
       CASE WHEN completed THEN
           LEAST((opened + (n % 8))::date + time '14:00',
                 ((CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date - 1) + time '14:00')
           AT TIME ZONE 'UTC'
       END,
       1 + n % 4,
       due
FROM demo;

GRANT CONNECT ON DATABASE todoapp_demo TO redash;
\getenv reporting_password REPORTING_PASSWORD
CREATE ROLE todoapp_reporting LOGIN PASSWORD :'reporting_password';
GRANT CONNECT ON DATABASE todoapp_demo TO todoapp_reporting;

CREATE SCHEMA analytics;
CREATE VIEW analytics.tasks AS
SELECT t.id,
       t.title,
       u.id AS user_id,
       u.full_name AS task_user,
       t.completed,
       (t.created_at AT TIME ZONE 'UTC')::date AS opened_on,
       (t.completed_at AT TIME ZONE 'UTC')::date AS completed_on,
       t.due_date,
       t.completed_at - t.created_at AS time_to_complete,
       NOT t.completed AND t.due_date < (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date AS overdue,
       (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date - t.due_date AS days_after_due
FROM todoapp.tasks AS t
JOIN todoapp.users AS u ON u.id = t.author_user_id;

GRANT USAGE ON SCHEMA analytics TO todoapp_reporting;
GRANT SELECT ON analytics.tasks TO todoapp_reporting;
REVOKE CONNECT ON DATABASE redash FROM PUBLIC;
