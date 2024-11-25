-- +goose Up
CREATE TABLE IF NOT EXISTS events (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        title VARCHAR(255) NOT NULL,
                        event_datetime TIMESTAMP NOT NULL,
                        duration INTERVAL,
                        description TEXT,
                        user_id UUID NOT NULL,
                        notify_before INTERVAL,
                        created_at TIMESTAMP DEFAULT NOW(),
                        updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW()
);

CREATE INDEX idx_events_user_id_event_datetime ON events (user_id, event_datetime);

-- +goose Down
DROP TABLE IF EXISTS events;
DROP INDEX  IF EXISTS idx_events_user_id_event_datetime;