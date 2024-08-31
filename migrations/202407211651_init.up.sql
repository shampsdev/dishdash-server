CREATE TABLE "place"
(
    "id"                serial       NOT NULL,
    "title"             varchar(255) NOT NULL,
    "short_description" varchar(255) NOT NULL,
    "description"       text         NOT NULL,
    "images"            text         NOT NULL,
    "location"          geography    NOT NULL,
    "address"           varchar(255) NOT NULL,
    "price_avg"         decimal      NOT NULL,
    "review_rating"     float        NOT NULL,
    "review_count"      decimal      NOT NULL,
    "updated_at"        timestamp    NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "tag"
(
    "id"   serial       NOT NULL,
    "name" varchar(255) NOT NULL UNIQUE,
    "icon" varchar(255) NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "place_tag"
(
    "place_id" int NOT NULL,
    "tag_id"   int NOT NULL,
    FOREIGN KEY ("tag_id") REFERENCES "tag" ("id"),
    FOREIGN KEY ("place_id") REFERENCES "place" ("id")
);

CREATE TABLE "user"
(
    id         varchar(255) NOT NULL DEFAULT gen_random_uuid(),
    name       varchar(255) NOT NULL,
    avatar     varchar(255) NOT NULL,
    telegram   bigint       NULL,
    created_at timestamp    NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "lobby"
(
    id         varchar(255) NOT NULL,
    state      varchar(255) NOT NULL,
    price_avg  bigint       NOT NULL,
    location   geography    NOT NULL,
    created_at timestamp    NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "lobby_tag"
(
    lobby_id varchar(255) NOT NULL,
    tag_id   serial       NOT NULL,
    PRIMARY KEY ("lobby_id", "tag_id"),
    FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id"),
    FOREIGN KEY ("tag_id") REFERENCES "tag" ("id")
);

CREATE TABLE "place_lobby"
(
    lobby_id varchar(255) NOT NULL,
    place_id serial       NOT NULL,
    PRIMARY KEY ("lobby_id", "place_id"),
    FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id"),
    FOREIGN KEY ("place_id") REFERENCES "place" ("id")
);

CREATE TABLE "lobby_user"
(
    lobby_id varchar(255) NOT NULL,
    user_id  varchar(255) NOT NULL,
    PRIMARY KEY ("lobby_id", "user_id"),
    FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id"),
    FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);

CREATE TABLE "swipe"
(
    id       serial       NOT NULL,
    lobby_id varchar(255) NOT NULL,
    place_id serial       NOT NULL,
    user_id  varchar(255) NOT NULL,
    type     varchar(255) NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id"),
    FOREIGN KEY ("user_id") REFERENCES "user" ("id"),
    FOREIGN KEY ("place_id") REFERENCES "place" ("id")
);

