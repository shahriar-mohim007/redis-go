CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(255) UNIQUE NOT NULL,
    original_url TEXT NOT NULL
    );

INSERT INTO urls (short_url, original_url) VALUES
    ('42v8nK', 'https://example.com/my-original-long-url')
    ON CONFLICT (short_url) DO NOTHING;
