-- +migrate Up
CREATE TABLE IF NOT EXISTS events (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        title VARCHAR(255) NOT NULL,
                        event_datetime TIMESTAMP NOT NULL,
                        duration INT,
                        description TEXT,
                        user_id UUID NOT NULL,
                        notify_before INT,
                        created_at TIMESTAMP DEFAULT NOW(),
                        updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_events_user_id_event_datetime ON events (user_id, event_datetime);

-- +migrate Down
DROP TABLE IF EXISTS events;
DROP INDEX  IF EXISTS idx_events_user_id_event_datetime;