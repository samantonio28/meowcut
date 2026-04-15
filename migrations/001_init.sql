CREATE TABLE IF NOT EXISTS links (
    short_id VARCHAR(10) PRIMARY KEY,
    original_url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_links_original_url ON links(original_url);