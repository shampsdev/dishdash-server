INSERT INTO
    "tag" (name, icon)
VALUES
    ('bar', 'bar.png'),
    ('cafe', 'cafe.png'),
    ('coffee', 'coffee.png'),
    ('food', 'food.png');

INSERT INTO
    "place" (
        title,
        short_description,
        description,
        images,
        location,
        address,
        price_avg,
        review_rating,
        review_count,
        updated_at
    )
VALUES
    (
        'Пловная',
        'Ресторан узбекской кухни',
        '«Пловная» — это ресторан узбекской кухни, расположенный недалеко от Сытного рынка на Петроградской стороне. По словам посетителей, это одно из лучших мест в городе для тех, кто хочет попробовать настоящий узбекский плов.',
        'https://avatars.mds.yandex.net/get-altay/10814540/2a0000018b5c49e5aa64362e66f0e493e1a9/XXXL',
        ST_SetSRID(ST_MakePoint(30.2956, 59.9505), 4326),
        'Сытнинская ул., 4П, Санкт-Петербург',
        450,
        4.5,
        1000,
        NOW()
    ),
    (
        'Zoomer Coffee',
        'Лучший кофе у ИТМО',
        'Лучший кофе у ИТМО от Стаса. Вкусные зумер-сендвичи и зумер-доги.',
        'https://avatars.mds.yandex.net/get-altay/777564/2a0000018789d73e0d0447c04671d4f00971/XXXL',
        ST_SetSRID(ST_MakePoint(30.2871, 59.9606), 4326),
        'Сытнинская ул., 4П, Санкт-Петербург',
        120,
        4.6,
        1000,
        NOW()
    ),
    (
        'ЛюдиЛюбят',
        'Пекарня, Кондитерская',
        'ЛюдиЛюбят - отличная пекарня с классными комбо на обед.',
        'https://avatars.mds.yandex.net/get-altay/9753788/2a0000018a2cf1ef9ad2f7880c260e4379ef/XXXL',
        ST_SetSRID(ST_MakePoint(30.2658, 59.9556), 4326),
        'Саблинская ул., 12, Санкт-Петербург',
        300,
        4.7,
        1000,
        NOW()
    ),
    (
        'Cous-Cous',
        'Лучшая шаверма у ИТМО',
        'Красиво, вкусно и недорого',
        'https://avatars.mds.yandex.net/get-altay/922263/2a00000185bf82020b53f168abfbae5e3e17/XXXL',
        ST_SetSRID(ST_MakePoint(30.2555, 59.9522), 4326),
        'ул. Кропоткина, 19/8, Санкт-Петербург',
        320,
        4.8,
        1000,
        NOW()
    ),
    (
        'Шавафель',
        'Вкусная шаверма',
        'Вкусная сырная шаверма и не только',
        'https://avatars.mds.yandex.net/get-altay/9691438/2a0000018bc3142dd6dbd390b6c8ba5736de/XXXL',
        ST_SetSRID(ST_MakePoint(30.3124, 59.9615), 4326),
        'Кронверкский проспект, 27',
        320,
        4.5,
        1000,
        NOW()
    );

INSERT INTO "place_tag" (place_id, tag_id) VALUES
((SELECT id FROM place WHERE title = 'Пловная'), (SELECT id FROM tag WHERE name = 'cafe')),
((SELECT id FROM place WHERE title = 'Пловная'), (SELECT id FROM tag WHERE name = 'food'));

INSERT INTO "place_tag" (place_id, tag_id) VALUES
((SELECT id FROM place WHERE title = 'Zoomer Coffee'), (SELECT id FROM tag WHERE name = 'cafe')),
((SELECT id FROM place WHERE title = 'Zoomer Coffee'), (SELECT id FROM tag WHERE name = 'coffee')),
((SELECT id FROM place WHERE title = 'Zoomer Coffee'), (SELECT id FROM tag WHERE name = 'food'));

INSERT INTO "place_tag" (place_id, tag_id) VALUES
((SELECT id FROM place WHERE title = 'ЛюдиЛюбят'), (SELECT id FROM tag WHERE name = 'cafe')),
((SELECT id FROM place WHERE title = 'ЛюдиЛюбят'), (SELECT id FROM tag WHERE name = 'coffee')),
((SELECT id FROM place WHERE title = 'ЛюдиЛюбят'), (SELECT id FROM tag WHERE name = 'food'));

INSERT INTO "place_tag" (place_id, tag_id) VALUES
((SELECT id FROM place WHERE title = 'Cous-Cous'), (SELECT id FROM tag WHERE name = 'bar')),
((SELECT id FROM place WHERE title = 'Cous-Cous'), (SELECT id FROM tag WHERE name = 'cafe')),
((SELECT id FROM place WHERE title = 'Cous-Cous'), (SELECT id FROM tag WHERE name = 'food'));

INSERT INTO "place_tag" (place_id, tag_id) VALUES
((SELECT id FROM place WHERE title = 'Шавафель'), (SELECT id FROM tag WHERE name = 'cafe')),
((SELECT id FROM place WHERE title = 'Шавафель'), (SELECT id FROM tag WHERE name = 'food'));