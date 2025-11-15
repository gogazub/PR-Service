--- Основные таблицы ---
CREATE TABLE users (
    user_id UUID PRIMARY KEY, 
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,

    -- Добавим метки со временем на будущее
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);
 
CREATE TABLE teams (
    team_name TEXT PRIMARY KEY,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TYPE pull_request_status AS ENUM (
    'OPEN',
    'MERGED'
);
CREATE TABLE pull_requests (
    pr_id UUID PRIMARY KEY,
    name TEXT,
    -- Если автор удален, то оставляем pr без автора
    author_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    status pull_request_status NOT NULL DEFAULT 'OPEN',
    
    --first_reviewers_id UUID REFERENCES users(user_id) DEFAULT NULL,
    --second_reviewers_id UUID REFERENCES users(user_id) DEFAULT NULL

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);


--- Отношения между основными таблицами ---

-- Допустим, что юзер может находится сразу в двух командах. Тогда имеем one to many.
CREATE TABLE user_teams (
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    team_name UUID REFERENCES teams(team_name) ON DELETE CASCADE,
    PRIMARY KEY (user_id, team_name)
);

-- Отдельная таблица для связи pr и reviewrs. Это таблицу стоит вынести, потому в дальнейшем можно будет легко увеличить кол-во ревьюеров  
CREATE TABLE pr_reviewers (
    pr_id UUID REFERENCES pull_requests(pr_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (pr_id, user_id)
    -- assigned_at
);


--- Функции ---

-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер обновления updated_at при апдейте в users
CREATE TRIGGER users_updated_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Триггер обновления updated_at при апдейте в teams
CREATE TRIGGER teams_updated_timestamp
BEFORE UPDATE ON teams
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Триггер обновления updated_at при апдейте в pull_requests
CREATE TRIGGER pull_requests_updated_timestamp
BEFORE UPDATE ON pull_requests
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();