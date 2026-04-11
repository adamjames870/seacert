-- name: GetPromptCache :one
SELECT cache_key, model_name, gemini_cache_name, expires_at, created_at
FROM prompt_caches
WHERE cache_key = $1 AND expires_at > NOW()
LIMIT 1;

-- name: UpsertPromptCache :exec
INSERT INTO prompt_caches (
    cache_key, model_name, gemini_cache_name, expires_at, created_at
) VALUES (
    $1, $2, $3, $4, NOW()
)
ON CONFLICT (cache_key) DO UPDATE SET
    gemini_cache_name = EXCLUDED.gemini_cache_name,
    expires_at = EXCLUDED.expires_at,
    created_at = NOW();

-- name: DeleteExpiredPromptCaches :exec
DELETE FROM prompt_caches
WHERE expires_at <= NOW();
