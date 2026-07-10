-- name: get-segments
-- Returns the saved segments (a JSON array) stored under the 'segments' settings key.
SELECT value FROM settings WHERE key = 'segments';

-- name: save-segments
-- Upserts the full segments array (raw JSON) into the settings table.
INSERT INTO settings (key, value) VALUES ('segments', $1::JSONB)
    ON CONFLICT (key) DO UPDATE SET value = $1::JSONB, updated_at = NOW();
