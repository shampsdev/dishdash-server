CREATE TABLE "collection"
(
    "id"          serial       NOT NULL,
    "name"        varchar(255) NOT NULL,
    "description" TEXT         NOT NULL,
    PRIMARY KEY ("id")

);

CREATE TABLE "collection_place"
(
    "collection_id" int NOT NULL,
    "place_id"      int NOT NULL,

    FOREIGN KEY ("collection_id") REFERENCES "collection" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("place_id") REFERENCES "place" ("id") ON DELETE CASCADE,
    UNIQUE ("collection_id", "place_id")
);
