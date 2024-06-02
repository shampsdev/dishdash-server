INSERT INTO "card" ("title", "short_description", "description", "image", "location", "address", "price")
VALUES ('Пловная', 'Ресторан узбекской кухни', '',
        'https://avatars.mds.yandex.net/get-altay/10814540/2a0000018b5c49e5aa64362e66f0e493e1a9/XXXL',
        '{"lat":59.957489,"lng":30.308895}', '', 1),
       ('ЛюдиЛюбят', 'Пекарня, Кондитерская', '',
        'https://avatars.mds.yandex.net/get-altay/9753788/2a0000018a2cf1ef9ad2f7880c260e4379ef/XXXL',
        '{"lat":59.957613,"lng":30.307365}', '', 1),
       ('Zoomer Coffee', 'Лучший кофе от Стаса', '',
        'https://avatars.mds.yandex.net/get-altay/777564/2a0000018789d73e0d0447c04671d4f00971/XXXL',
        '{"lat":59.957677,"lng":30.309756}', '', 1),
       ('Cous-Cous', 'Лучшая шаверма у ИТМО', '',
        'https://avatars.mds.yandex.net/get-altay/922263/2a00000185bf82020b53f168abfbae5e3e17/XXXL',
        '{"lat":59.958843,"lng":30.306465}', '', 1),
       ('Шавафель', 'Вкусная сырная шаверма и не только', '',
        'https://avatars.mds.yandex.net/get-altay/9691438/2a0000018bc3142dd6dbd390b6c8ba5736de/XXXL',
        '{"lat":59.957019,"lng":30.317019}', '', 1);

INSERT INTO "tag" ("name", "icon")
VALUES ('Bar', 'bar.svg'),
       ('Cafe', 'cafe.svg'),
       ('Restaurant', 'restaurant.svg');

INSERT INTO "tag_card" ("tag_id", "card_id")
SELECT 3, "id"
FROM "card"
WHERE "title" IN ('Пловная', 'Cous-Cous', 'Шавафель');

INSERT INTO "tag_card" ("tag_id", "card_id")
SELECT 2, "id"
FROM "card"
WHERE "title" IN ('ЛюдиЛюбят', 'Zoomer Coffee');
