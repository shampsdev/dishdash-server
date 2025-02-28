CREATE TABLE IF NOT EXISTS  "collection" 
(
    "id"          varchar(255) NOT NULL,
    "name"        varchar(255) NOT NULL,
    "description" TEXT         NOT NULL,
    "avatar"      varchar(255) NOT NULL,
    "created_at"  timestamp    NOT NULL DEFAULT NOW(),
    "updated_at"  timestamp    NOT NULL DEFAULT NOW(),
    "visible"     boolean      NOT NULL DEFAULT true,
    "order"       int          NOT NULL DEFAULT 0,
PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "collection_place" 
(
    "collection_id" varchar(255) NOT NULL,
    "place_id"      int NOT NULL,

    FOREIGN KEY ("collection_id") REFERENCES "collection" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("place_id") REFERENCES "place" ("id") ON DELETE CASCADE,
    UNIQUE ("collection_id", "place_id")
);