package bartend

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/c2stack/c2g/c2"
)

const (
	Id int = iota
	DrinkName
	Type
	Glass
	Garnish
	Occasion
	Ingredient1
	Ingredient2
	Ingredient3
	Ingredient4
	Ingredient5
	Ingredient6
	Ingredient7
	Ingredient8
	Ingredient9
	Ingredient10
	Ingredient11
	Ingredient12
)

type Reader struct {
	Csv *csv.Reader
}

func NewReader(in io.Reader) *Reader {
	return &Reader{Csv: csv.NewReader(in)}
}

func (self *Reader) Read() (*Recipe, error) {
	l1, err := self.Csv.Read()
	if err != nil {
		return nil, err
	}
	r := &Recipe{
		Name: l1[DrinkName],
	}
	if r.Id, err = strconv.Atoi(l1[Id]); err != nil {
		return nil, err
	}
	for i := Ingredient1; i <= Ingredient12; i++ {
		if l1[i] == "" {
			r.Ingredients = self.Ingredients(r.Id, l1[Ingredient1:i])
			break
		}
	}
	c2.Debug.Print(r.Name, l1[Ingredient1], len(l1))
	return r, nil
}

func (self *Reader) Ingredients(id int, s []string) []*Ingredient {
	ingreds := make([]*Ingredient, len(s))
	for i, data := range s {
		ingreds[i] = &Ingredient{}
		ingreds[i].Amount, ingreds[i].Liquid = self.parseIngredient(data)
	}
	return ingreds
}

func (self *Reader) parseIngredient(s string) (float64, string) {
	return 0, s
}
