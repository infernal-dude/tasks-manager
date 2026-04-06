-- +goose Up
INSERT INTO tasks (title, description) 
VALUES ('Учиться, учиться', 'И еще раз учиться');

-- +goose Down
DELETE FROM tasks;



