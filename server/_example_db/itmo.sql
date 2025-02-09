-- 22.07 dump

DELETE
from place_tag;
DELETE
from tag;
DELETE
from place;

INSERT INTO "place" ("id",
                     "title",
                     "short_description",
                     "description",
                     "images",
                     "location",
                     "address",
                     "price_avg",
                     "review_rating",
                     "review_count",
                     "updated_at")
VALUES (1, 'Пловная', 'Ресторан узбекской кухни',
        '«Пловная» — это ресторан узбекской кухни, расположенный недалеко от Сытного рынка на Петроградской стороне. По словам посетителей, это одно из лучших мест в городе для тех, кто хочет попробовать настоящий узбекский плов.',
        'https://avatars.mds.yandex.net/get-altay/10814540/2a0000018b5c49e5aa64362e66f0e493e1a9/XXXL',
        '0101000020E6100000DEE522BE134F3E40CFA3E2FF8EFA4D40', 'Сытнинская ул., 4П, Санкт-Петербург', 450, 4.5, 1000,
        '2024-07-22 00:00:00'),
       (2, 'Zoomer Coffee', 'Лучший кофе у ИТМО', 'Лучший кофе у ИТМО от Стаса. Вкусные зумер-сендвичи и зумер-доги.',
        'https://avatars.mds.yandex.net/get-altay/777564/2a0000018789d73e0d0447c04671d4f00971/XXXL',
        '0101000020E610000062BD512B4C4F3E40D190F12895FA4D40', 'Сытнинская ул., 4П, Санкт-Петербург', 120, 4.6, 1000,
        '2024-07-22 00:00:00'),
       (3, 'ЛюдиЛюбят', 'Пекарня, Кондитерская', 'ЛюдиЛюбят - отличная пекарня с классными комбо на обед.',
        'https://avatars.mds.yandex.net/get-altay/9753788/2a0000018a2cf1ef9ad2f7880c260e4379ef/XXXL',
        '0101000020E61000005FEFFE78AF4E3E40BB9C121093FA4D40', 'Саблинская ул., 12, Санкт-Петербург', 300, 4.7, 1000,
        '2024-07-22 00:00:00'),
       (4, 'Cous-Cous', 'Лучшая шаверма у ИТМО', 'Красиво, вкусно и недорого',
        'https://avatars.mds.yandex.net/get-altay/922263/2a00000185bf82020b53f168abfbae5e3e17/XXXL',
        '0101000020E61000005F5E807D744E3E40D07F0F5EBBFA4D40', 'ул. Кропоткина, 19/8, Санкт-Петербург', 320, 4.8, 1000,
        '2024-07-22 00:00:00'),
       (5, 'Шавафель', 'Вкусная шаверма', 'Вкусная сырная шаверма и не только',
        'https://avatars.mds.yandex.net/get-altay/9691438/2a0000018bc3142dd6dbd390b6c8ba5736de/XXXL',
        '0101000020E6100000EB353D2828513E4048533D997FFA4D40', 'Кронверкский проспект, 27', 320, 4.5, 1000,
        '2024-07-22 00:00:00');

SELECT nextval('place_id_seq');
SELECT nextval('place_id_seq');
SELECT nextval('place_id_seq');
SELECT nextval('place_id_seq');
SELECT nextval('place_id_seq');

INSERT INTO "tag" ("id", "name", "icon")
VALUES (1, 'Бар', 'bar.png'),
       (2, 'Кафе', 'cafe.png'),
       (3, 'Кофейни', 'coffee.png');

SELECT nextval('tag_id_seq');
SELECT nextval('tag_id_seq');
SELECT nextval('tag_id_seq');

INSERT INTO "place_tag" ("place_id", "tag_id")
VALUES (1, 2),
       (2, 2),
       (2, 3),
       (3, 2),
       (3, 3),
       (4, 1),
       (4, 2),
       (5, 2);
