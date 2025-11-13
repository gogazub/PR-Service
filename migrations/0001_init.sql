CREATE TABLE users (
    user_id UUID PRIMARY KEY, 
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE

    -- Добавим метки со временем на будущее
    created_at TIMESTAMPTZ DEFAULT NOW()
    updated_at TIMESTAMPTZ DEFAULT NOW()
    deleted_at TIMESTAMPTZ NULL
);

-- Хотя и можно использовать team_name как естественный ключ, но тогда может быть две проблемы.
-- 1. Если team_name длинное, то будет много копий длинного названия. 
-- 2. Если захотим изменить team_name - придется каскадно менять во всех строках. 
-- Так что для гибкости добавим id 
CREATE TABLE teams (
    team_id UUID PRIMARY KEY,
    team_name TEXT

    created_at TIMESTAMPTZ DEFAULT NOW()
    updated_at TIMESTAMPTZ DEFAULT NOW()
    deleted_at TIMESTAMPTZ NULL
);

-- Допустим, что юзер может находится сразу в двух командах. Тогда имеем one to many.
CREATE TABLE user_teams (
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(team_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, team_id)
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

    created_at TIMESTAMPTZ DEFAULT NOW()
    updated_at TIMESTAMPTZ DEFAULT NOW()
    deleted_at TIMESTAMPTZ NULL
);

-- Отдельная таблица для связи pr и reviewrs. Это таблицу стоит вынести, потому в дальнейшем можно будет легко увеличить кол-во ревьюеров  
CREATE TABLE pr_reviewers (
    pr_id UUID REFERENCES pull_requests(pr_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (pr_id, user_id)
);
