-- Up Migration

-- Create lobbysettings table
CREATE TABLE "lobbysettings"
(
    "id"           serial       NOT NULL,
    "lobby_id"     varchar(255) NOT NULL,
    "price_min"    decimal      NOT NULL DEFAULT 0,
    "price_max"    decimal      NOT NULL DEFAULT 1000000,
    "max_distance" decimal      NOT NULL DEFAULT 1000000,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create lobbysettings_tag table
CREATE TABLE "lobbysettings_tag"
(
    "id"               serial NOT NULL,
    "tag_id"           int    NOT NULL,
    "lobbysettings_id" int    NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("tag_id") REFERENCES "tag" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE,
    FOREIGN KEY ("lobbysettings_id") REFERENCES "lobbysettings" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create match table
CREATE TABLE "match"
(
    "id"       serial       NOT NULL,
    "lobby_id" varchar(255) NOT NULL,
    "card_id"  int          NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("card_id") REFERENCES "card" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE,
    FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create final_vote table
CREATE TABLE "final_vote"
(
    "id"       serial       NOT NULL,
    "lobby_id" varchar(255) NOT NULL,
    "card_id"  int          NOT NULL,
    "user_id"  varchar(255) NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE,
    FOREIGN KEY ("card_id") REFERENCES "card" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Create lobby_card table
CREATE TABLE "lobby_card"
(
    "id"       serial       NOT NULL,
    "lobby_id" varchar(255) NOT NULL,
    "card_id"  int          NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE,
    FOREIGN KEY ("card_id") REFERENCES "card" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE
);
