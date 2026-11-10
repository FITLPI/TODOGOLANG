ALTER TABLE todoapp.tasks
    DROP CONSTRAINT tasks_description_check;

ALTER TABLE todoapp.tasks
    ADD CONSTRAINT tasks_description_check
    CHECK (char_length(description) BETWEEN 1 AND 100);
