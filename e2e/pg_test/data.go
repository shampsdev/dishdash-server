package pg_test

import (
	"context"
	"log"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo/pg"
	"dishdash.ru/internal/usecase"
	"dishdash.ru/pkg/filter"
)

var Tags = []*domain.Tag{
	{Name: "bar", Icon: "bar.png"},
	{Name: "cafe", Icon: "cafe.png"},
	{Name: "coffee", Icon: "coffee.png"},
	{Name: "food", Icon: "food.png"},
}

var Places = []*domain.Place{
	{
		Title:            "Пловная",
		ShortDescription: "Ресторан узбекской кухни",
		Description:      "«Пловная» — это ресторан узбекской кухни, расположенный недалеко от Сытного рынка на Петроградской стороне. По словам посетителей, это одно из лучших мест в городе для тех, кто хочет попробовать настоящий узбекский плов.",
		Images:           []string{"https://avatars.mds.yandex.net/get-altay/10814540/2a0000018b5c49e5aa64362e66f0e493e1a9/XXXL"},
		Location:         domain.Coordinate{Lon: 30.2956, Lat: 59.9505},
		Address:          "Сытнинская ул., 4П, Санкт-Петербург",
		PriceAvg:         450,
		ReviewRating:     4.5,
		ReviewCount:      1000,
		Tags:             []*domain.Tag{Tags[1], Tags[3]},
	},
	{
		Title:            "Zoomer Coffee",
		ShortDescription: "Лучший кофе у ИТМО",
		Description:      "Лучший кофе у ИТМО от Стаса. Вкусные зумер-сендвичи и зумер-доги.",
		Images:           []string{"https://avatars.mds.yandex.net/get-altay/777564/2a0000018789d73e0d0447c04671d4f00971/XXXL"},
		Location:         domain.Coordinate{Lon: 30.2871, Lat: 59.9606},
		Address:          "Сытнинская ул., 4П, Санкт-Петербург",
		PriceAvg:         120,
		ReviewRating:     4.6,
		ReviewCount:      1000,
		Tags:             []*domain.Tag{Tags[1], Tags[2], Tags[3]},
	},
	{
		Title:            "ЛюдиЛюбят",
		ShortDescription: "Пекарня, Кондитерская",
		Description:      "ЛюдиЛюбят - отличная пекарня с классными комбо на обед.",
		Images:           []string{"https://avatars.mds.yandex.net/get-altay/9753788/2a0000018a2cf1ef9ad2f7880c260e4379ef/XXXL"},
		Location:         domain.Coordinate{Lon: 30.2658, Lat: 59.9556},
		Address:          "Саблинская ул., 12, Санкт-Петербург",
		PriceAvg:         300,
		ReviewRating:     4.7,
		ReviewCount:      1000,
		Tags:             []*domain.Tag{Tags[1], Tags[2], Tags[3]},
	},
	{
		Title:            "Cous-Cous",
		ShortDescription: "Лучшая шаверма у ИТМО",
		Description:      "Красиво, вкусно и недорого",
		Images:           []string{"https://avatars.mds.yandex.net/get-altay/922263/2a00000185bf82020b53f168abfbae5e3e17/XXXL"},
		Location:         domain.Coordinate{Lon: 30.2555, Lat: 59.9522},
		Address:          "ул. Кропоткина, 19/8, Санкт-Петербург",
		PriceAvg:         320,
		ReviewRating:     4.8,
		ReviewCount:      1000,
		Tags:             []*domain.Tag{Tags[0], Tags[1], Tags[3]},
	},
	{
		Title:            "Шавафель",
		ShortDescription: "Вкусная шаверма",
		Description:      "Вкусная сырная шаверма и не только",
		Images:           []string{"https://avatars.mds.yandex.net/get-altay/9691438/2a0000018bc3142dd6dbd390b6c8ba5736de/XXXL"},
		Location:         domain.Coordinate{Lon: 30.3124, Lat: 59.9615},
		Address:          "Кронверкский проспект, 27",
		PriceAvg:         320,
		ReviewRating:     4.5,
		ReviewCount:      1000,
		Tags:             []*domain.Tag{Tags[1], Tags[3]},
	},
}

func ResetData(cases usecase.Cases) error {
	err := pg.MigrateDown()
	if err != nil {
		return err
	}
	err = pg.MigrateUp()
	if err != nil {
		return err
	}
	ctx := context.Background()
	for i := range Tags {
		var err error
		Tags[i], err = cases.Tag.SaveTag(ctx, Tags[i])
		if err != nil {
			return err
		}
	}
	for i, p := range Places {
		var err error
		placeInput := usecase.SavePlaceInput{
			Title:            p.Title,
			ShortDescription: p.ShortDescription,
			Description:      p.Description,
			Location:         p.Location,
			Address:          p.Address,
			PriceAvg:         p.PriceAvg,
			ReviewRating:     p.ReviewRating,
			ReviewCount:      p.ReviewCount,
			Images:           p.Images,
			Tags: filter.Map(p.Tags, func(t *domain.Tag) int64 {
				return t.ID
			}),
		}
		Places[i], err = cases.Place.SavePlace(ctx, placeInput)
		if err != nil {
			return err
		}
	}

	log.Printf("setup data")
	return nil
}
