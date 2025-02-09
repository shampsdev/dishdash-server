ALTER TABLE "lobby" ADD COLUMN "price_avg" bigint;
ALTER TABLE "lobby" ADD COLUMN "location" geography;

CREATE TABLE "lobby_tag"
(
    "lobby_id" varchar(255) NOT NULL,
    "tag_id"   serial       NOT NULL,
    PRIMARY KEY ("lobby_id", "tag_id"),
    FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("tag_id") REFERENCES "tag" ("id") ON DELETE CASCADE
);

UPDATE "lobby"
SET
    "price_avg" = ("settings"->'classicPlaces'->>'priceAvg')::numeric,
    "location" = ST_SetSRID(ST_MakePoint(
        ("settings"->'classicPlaces'->'location'->>'lon')::float,
        ("settings"->'classicPlaces'->'location'->>'lat')::float
    ), 4326);

INSERT INTO "lobby_tag" ("lobby_id", "tag_id")
SELECT l.id, jsonb_array_elements_text(l."settings"->'classicPlaces'->'tags')::int
FROM "lobby" l
WHERE l."settings" IS NOT NULL;

ALTER TABLE "lobby" DROP COLUMN "settings";
ALTER TABLE "lobby" DROP COLUMN "type";   
