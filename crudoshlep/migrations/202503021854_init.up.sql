CREATE TABLE "event"
(
    "id"         varchar(255) UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    "name"       varchar(255) NOT NULL,
    "data"       jsonb,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);
