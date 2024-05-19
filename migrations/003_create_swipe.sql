CREATE TABLE "swipe"
(
    "id"         serial NOT NULL,
    "lobby_id"   int,
    "card_id"    int,
    "user_id"    varchar(255),
    "swipe_type" varchar(255),
    PRIMARY KEY ("id")
);

ALTER TABLE "swipe"
    ADD FOREIGN KEY ("lobby_id") REFERENCES "lobby" ("id")
        ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "swipe"
    ADD FOREIGN KEY ("card_id") REFERENCES "card" ("id")
        ON UPDATE NO ACTION ON DELETE NO ACTION;