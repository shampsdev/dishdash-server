ALTER TABLE
    "tag"
ADD
    COLUMN "visible" boolean NOT NULL DEFAULT true;

ALTER TABLE
    "tag"
ADD
    COLUMN "order" int NOT NULL DEFAULT 0;
