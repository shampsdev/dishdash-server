-- 21.06 dump

DELETE
from card_tag;
DELETE
from tag;
DELETE
from card;

INSERT INTO "card" ("id", "title", "short_description", "description", "image", "location", "address", "price_min",
                    "price_max")
VALUES (1, 'Пловная', 'Ресторан узбекской кухни',
        '«Пловная» — это ресторан узбекской кухни, расположенный недалеко от Сытного рынка на Петроградской стороне. По словам посетителей, это одно из лучших мест в городе для тех, кто хочет попробовать настоящий узбекский плов.',
        'https://avatars.mds.yandex.net/get-altay/10814540/2a0000018b5c49e5aa64362e66f0e493e1a9/XXXL',
        '0101000020E6100000CFA3E2FF8EFA4D40DEE522BE134F3E40', 'Сытнинская ул., 4П, Санкт-Петербург', 450, 700),
       (2, 'Zoomer Coffee', 'Лучший кофе у ИТМО', 'Лучший кофе у ИТМО от Стаса. Вкусные зумер-сендвичи и зумер-доги.',
        'https://avatars.mds.yandex.net/get-altay/777564/2a0000018789d73e0d0447c04671d4f00971/XXXL',
        '0101000020E6100000D190F12895FA4D4062BD512B4C4F3E40', 'Сытнинская ул., 4П, Санкт-Петербург', 90, 300),
       (3, 'ЛюдиЛюбят', 'Пекарня, Кондитерская', 'ЛюдиЛюбят - отличная пекарня с классными комбо на обед.',
        'https://avatars.mds.yandex.net/get-altay/9753788/2a0000018a2cf1ef9ad2f7880c260e4379ef/XXXL',
        '0101000020E6100000BB9C121093FA4D405FEFFE78AF4E3E40', 'Саблинская ул., 12, Санкт-Петербург', 150, 450),
       (4, 'Cous-Cous', 'Лучшая шаверма у ИТМО', 'Красиво, вкусно и недорого',
        'https://avatars.mds.yandex.net/get-altay/922263/2a00000185bf82020b53f168abfbae5e3e17/XXXL',
        '0101000020E6100000D07F0F5EBBFA4D400056478E744E3E40', 'ул. Кропоткина, 19/8, Санкт-Петербург', 250, 350),
       (5, 'Шавафель', 'Вкусная шаверма', 'Вкусная сырная шаверма и не только',
        'https://avatars.mds.yandex.net/get-altay/9691438/2a0000018bc3142dd6dbd390b6c8ba5736de/XXXL',
        '0101000020E610000048533D997FFA4D40EB353D2828513E40', 'Кронверкский проспект, 27', 250, 500);

INSERT INTO "tag" ("id", "name", "icon")
VALUES (1, 'bar', 'bar.png'),
       (2, 'cafe', 'cafe.png'),
       (3, 'coffee', 'coffee.png');

INSERT INTO "card_tag" ("card_id", "tag_id")
VALUES (1, 2),
       (2, 2),
       (2, 3),
       (3, 2),
       (3, 3),
       (4, 1),
       (4, 2),
       (5, 2);
