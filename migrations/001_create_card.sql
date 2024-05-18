CREATE TYPE "type_t" AS ENUM ('BAR', 'CAFE', 'RESTAURANT');

CREATE TABLE "card"
(
    "id"                serial       NOT NULL,
    "title"             varchar(255) NOT NULL,
    "short_description" varchar(255) NOT NULL,
    "description"       text         NOT NULL,
    "image"             varchar(255) NOT NULL,
    "location"          varchar(255) NOT NULL,
    "address"           varchar(255) NOT NULL,
    "type"              type_t       NOT NULL,
    "price"             int          NOT NULL,
    PRIMARY KEY ("id")
);
