ALTER TABLE "lobby" ADD COLUMN "settings" jsonb;

UPDATE "lobby" l
SET "settings" = jsonb_build_object(
    'type', 'classicPlaces'::text,
    'classicPlaces', jsonb_build_object(
        'location', jsonb_build_object(
            'lon', ST_X(l."location"::geometry)::float,
            'lat', ST_Y(l."location"::geometry)::float
        ),
        'tags', COALESCE((
            SELECT jsonb_agg(lt."tag_id")
            FROM "lobby_tag" lt
            WHERE lt."lobby_id" = l."id"
        ), '[]'::jsonb),
        'priceAvg', l."price_avg",
        'recommendation', null
    )
);

ALTER TABLE "lobby" ALTER COLUMN "settings" SET NOT NULL;

ALTER TABLE "lobby" ADD COLUMN "type" VARCHAR(255) DEFAULT 'classicPlaces';

ALTER TABLE "lobby" DROP COLUMN "price_avg";
ALTER TABLE "lobby" DROP COLUMN "location";
DROP TABLE "lobby_tag";
