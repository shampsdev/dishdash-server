ALTER TABLE "swipe" ADD COLUMN "swiped_at" timestamp;

WITH "timed_swipe" AS (
    SELECT s.id, l.created_at
    FROM swipe s
    JOIN lobby l ON s.lobby_id = l.id
)
UPDATE "swipe" s
SET swiped_at = ts.created_at
FROM timed_swipe ts
WHERE s.id = ts.id;

ALTER TABLE "swipe" ALTER COLUMN "swiped_at" SET NOT NULL;
ALTER TABLE "swipe" ALTER COLUMN "swiped_at" SET DEFAULT NOW();
