package domain

import "dishdash.ru/internal/dto"

type Tag struct {
	ID   int64
	Name string
	Icon string
}

func (t *Tag) ToDto() dto.Tag {
	return dto.Tag{
		ID:   t.ID,
		Name: t.Name,
		Icon: t.Icon,
	}
}

func (t *Tag) ParseDto(tDto dto.Tag) {
	t.ID = tDto.ID
	t.Name = tDto.Name
	t.Icon = tDto.Icon
}
