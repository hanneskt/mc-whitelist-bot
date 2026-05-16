-- name: GetPlayerByDiscordUuid :one
SELECT * FROM players
WHERE discord_uuid = ?;

-- name: CreatePlayer :one
INSERT INTO players (
    mc_uuid, mc_username, discord_uuid, country, invited_by, whitelisted
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;
