-- Пользователи
CREATE TABLE IF NOT EXISTS users (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT         UNIQUE NOT NULL,
    role          TEXT         NOT NULL,
    password_hash TEXT         NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    );

-- -- Фиксированные пользователи для /dummyLogin
-- INSERT INTO users (id, email, role, password_hash) VALUES
--     ('00000000-0000-0000-0000-000000000001', 'admin@example.com', 'admin', ''),
--     ('00000000-0000-0000-0000-000000000002', 'user@example.com',  'user',  '')
--     ON CONFLICT (id) DO NOTHING;

-- Переговорки
CREATE TABLE IF NOT EXISTS rooms (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT         NOT NULL,
    description TEXT,
    capacity    INTEGER,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    );

-- Расписания
CREATE TABLE IF NOT EXISTS schedules (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id      UUID         NOT NULL REFERENCES rooms(id),
    days_of_week INTEGER[]    NOT NULL,
    start_time   TIME         NOT NULL,
    end_time     TIME         NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    );

-- Слоты
CREATE TABLE IF NOT EXISTS slots (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id     UUID         NOT NULL REFERENCES rooms(id),
    start_time  TIMESTAMPTZ  NOT NULL,
    end_time    TIMESTAMPTZ  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    );

-- Брони
CREATE TABLE IF NOT EXISTS bookings (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id         UUID         NOT NULL REFERENCES slots(id),
    user_id         UUID         NOT NULL REFERENCES users(id),
    status          TEXT         NOT NULL DEFAULT 'active',
    conference_link TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    );
