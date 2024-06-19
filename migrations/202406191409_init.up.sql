CREATE TABLE "card"
(
    "id"                serial       NOT NULL,
    "title"             varchar(255) NOT NULL,
    "short_description" varchar(255) NOT NULL,
    "description"       text         NOT NULL,
    "image"             varchar(255) NOT NULL,
    "location"          geometry     NOT NULL,
    "address"           varchar(255) NOT NULL,
    "price"             decimal      NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "tag"
(
    "id"   serial       NOT NULL,
    "name" varchar(255) NOT NULL,
    "icon" varchar(255),
    PRIMARY KEY ("id")
);

CREATE TABLE "card_tag"
(
    "card_id" int NOT NULL,
    "tag_id"  int NOT NULL
);

CREATE TABLE "lobby"
(
    "id"         varchar(255) NOT NULL,
    "location"   geometry     NOT NULL,
    "created_at" timestamp    NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "temp_user"
(
    "id"         varchar(255) NOT NULL,
    "created_at" timestamp    NOT NULL,
    PRIMARY KEY ("id")
);


CREATE TABLE "swipe"
(
    "card_id"  int          NOT NULL,
    "lobby_id" varchar(255) NOT NULL,
    "user_id"  varchar(255) NOT NULL,
    "type"     varchar(255) NOT NULL
);

ALTER TABLE "card_tag"
    ADD FOREIGN KEY ("tag_id") REFERENCES "tag" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "card_tag"
    ADD FOREIGN KEY ("card_id") REFERENCES "card" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "swipe"
    ADD FOREIGN KEY ("card_id") REFERENCES "card" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "swipe"
    ADD FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "swipe"
    ADD FOREIGN KEY ("user_id") REFERENCES "temp_user" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE;