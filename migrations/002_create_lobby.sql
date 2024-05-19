CREATE TABLE "lobby"
(
    "id"         serial       NOT NULL,
    "location"   varchar(255) NOT NULL,
    "created_at" timestamp    NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);
