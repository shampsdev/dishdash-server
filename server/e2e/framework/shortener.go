package framework

import (
	"dishdash.ru/e2e/framework/session"
	"dishdash.ru/pkg/algo"
	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/usecase/event"
)

type H = map[string]any

var CardsShortener = session.EventShortener(func(ev event.Cards) any {
	cards := algo.Map(ev.Cards, func(c *domain.Place) any {
		return H{
			"id":    c.ID,
			"title": c.Title,
		}
	})

	return H{
		"cards": cards,
	}
})

var ResultsShortener = session.EventShortener(func(ev event.Results) any {
	results := algo.Map(ev.Top, func(r event.TopPosition) any {
		return H{
			"card": H{
				"id":    r.Card.ID,
				"title": r.Card.Title,
			},
			"likes": algo.Map(r.Likes, func(u *domain.User) any {
				return H{
					"id":   u.ID,
					"name": u.Name,
				}
			}),
		}
	})

	return H{
		"results": results,
	}
})
