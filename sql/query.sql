-- name: GetPlayerByDiscordUuid :one
SELECT * FROM players
WHERE discord_uuid = ?;

-- name: CreatePlayer :one
INSERT INTO players (
    mc_uuid, mc_username, discord_uuid, country, invited_by, whitelisted
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: CreateBirthday :one
INSERT INTO birthdays (
    discord_uuid, day, month
) VALUES (
    ?,?,?
) RETURNING *;

-- name: GetBirthdayById :one
SELECT * FROM birthdays WHERE discord_uuid = ?;

-- name: GetBirthdays :many
SELECT * FROM birthdays;

-- name: DeleteBirthday :exec
DELETE FROM birthdays WHERE discord_uuid = ?;

-- name: GetBirthdaysForDate :many
SELECT * FROM birthdays WHERE day = ? AND month = ?;

-- name: GetBirthdaysForMonth :many
SELECT * FROM birthdays WHERE month = ? ORDER BY day ASC;
