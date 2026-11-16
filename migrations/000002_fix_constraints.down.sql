ALTER TABLE todoapp.tasks
    DROP CONSTRAINT tasks_author_user_id_fkey;

ALTER TABLE todoapp.tasks
    ADD CONSTRAINT tasks_author_user_id_fkey
    FOREIGN KEY (author_user_id) REFERENCES todoapp.users(id) ON DELETE SET NULL;

ALTER TABLE todoapp.users
    ALTER COLUMN phone_number SET NOT NULL;
