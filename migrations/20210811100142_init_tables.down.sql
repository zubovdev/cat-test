-- task table
ALTER TABLE task DROP CONSTRAINT fk_task_user;
DROP TABLE task;

-- user table
DROP TABLE "user";