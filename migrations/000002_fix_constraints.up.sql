-- phone_number обязателен по NULL-ограничению и не соответствует домену
-- (PhoneNumber *string: nil = NULL в базе) -> делаем column nullable.
ALTER TABLE todoapp.users
    ALTER COLUMN phone_number DROP NOT NULL;

-- author_user_id не может быть "SET NULL" при удалении автора: колонка NOT NULL,
-- а домен требует AuthorUserID int (обязателен) -> каскадное удаление задач.
ALTER TABLE todoapp.tasks
    DROP CONSTRAINT tasks_author_user_id_fkey;

ALTER TABLE todoapp.tasks
    ADD CONSTRAINT tasks_author_user_id_fkey
    FOREIGN KEY (author_user_id) REFERENCES todoapp.users(id) ON DELETE CASCADE;
