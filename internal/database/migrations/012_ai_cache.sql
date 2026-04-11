-- +goose Up
-- Table to store Gemini context caching metadata
CREATE TABLE prompt_caches (
   cache_key TEXT PRIMARY KEY,           -- SHA256(canonical_payload)
   model_name TEXT NOT NULL,             -- e.g., 'gemini-2.5-flash-001'
   gemini_cache_name TEXT NOT NULL,      -- The 'name' returned by Gemini API (e.g., 'cachedContents/...')
   expires_at TIMESTAMPTZ NOT NULL,      -- Expiration time from Gemini
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index to quickly find non-expired caches
CREATE INDEX idx_prompt_caches_expires_at ON prompt_caches (expires_at);

-- +goose Down
DROP TABLE IF EXISTS prompt_caches;