create table "tag"
(
    "id"   serial NOT NULL,
    "name" varchar(255),
    "icon" varchar(255),
    PRIMARY KEY ("id")
);

create table "tag_card"
(
    "card_id" serial NOT NULL,
    "tag_id"  serial NOT NULL
);

ALTER TABLE "tag_card"
    ADD FOREIGN KEY ("card_id") REFERENCES "card" ("id")
        ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "tag_card"
    ADD FOREIGN KEY ("tag_id") REFERENCES "tag" ("id")
        ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER table "card"
    drop column "type"