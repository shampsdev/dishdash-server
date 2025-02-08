ALTER TABLE place_lobby ADD COLUMN "order" int NOT NULL DEFAULT 0;

WITH ranked_places AS (
    SELECT
        pl.lobby_id,
        pl.place_id,
        ROW_NUMBER() OVER (
            PARTITION BY pl.lobby_id
            ORDER BY (
                (1.0 * ST_Distance(p.location, l.location)) +
                (0.1 * ABS(p.price_avg - l.price_avg))
            ) /
            CASE
                WHEN p.boost IS NOT NULL AND p.boost_radius IS NOT NULL AND 
                    ST_Distance(p.location, l.location) <= p.boost_radius THEN
                    p.boost
                ELSE
                    1
            END
        ) AS rn
    FROM place_lobby pl
    JOIN place p ON pl.place_id = p.id
    JOIN lobby l ON pl.lobby_id = l.id
)
UPDATE place_lobby
SET "order" = ranked_places.rn
FROM ranked_places
WHERE place_lobby.lobby_id = ranked_places.lobby_id
AND place_lobby.place_id = ranked_places.place_id;

ALTER TABLE place_lobby
ALTER COLUMN "order" DROP DEFAULT;
