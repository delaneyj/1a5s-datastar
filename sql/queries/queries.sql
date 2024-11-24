-- name: RandomPokemon :many
SELECT * FROM pokemon ORDER BY random() LIMIT @limit;

-- name: UpvotePokemon :exec
UPDATE pokemon
SET
    up_votes = up_votes + 1,
    updated_at = @updated_at
WHERE id = @id;

-- name: DownvotePokemon :exec
UPDATE pokemon
SET
    down_votes = down_votes + 1,
    updated_at = @updated_at
WHERE id = @id;

-- name: AllIDs :many
SELECT id FROM pokemon
WHERE id < 1025;

-- name: Results :many
SELECT
    *,
    10000 * up_votes / total_votes as win_percentage
FROM (
SELECT
	name, id, dex_id,
	up_votes, down_votes,
	up_votes + down_votes as total_votes
FROM
	pokemon
WHERE id < 1025
) as subquery
ORDER BY up_votes DESC, win_percentage DESC;