# Event Manager

Предлагает слои для работы с паттерном Transaction Outbox.

## Getting Started
Для работы необходимо заранее промигрировать таблицу:
```sql
CREATE TYPE event_status AS ENUM ('WAIT', 'DONE');
CREATE TABLE IF NOT EXISTS event
(
    id          bigserial PRIMARY KEY,
    topic       TEXT NOT NULL,
    payload     JSONB NOT NULL,
    status      event_status NOT NULL DEFAULT 'WAIT',
    created_at  timestamp NOT NULL DEFAULT NOW(),
    reserved_to timestamp DEFAULT NULL
);
CREATE INDEX idx_created_at ON "event"(created_at);
```

Слои представлены в events_manager.go