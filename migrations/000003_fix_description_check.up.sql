-- description по домену допускается до 1000 символов,
-- в начальной миграции CHECK ошибочно ограничивал 100.
ALTER TABLE todoapp.tasks
    DROP CONSTRAINT tasks_description_check;

ALTER TABLE todoapp.tasks
    ADD CONSTRAINT tasks_description_check
    CHECK (description IS NULL OR char_length(description) BETWEEN 1 AND 1000);
