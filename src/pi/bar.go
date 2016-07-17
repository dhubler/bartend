package pi

import (
	"sort"
)

func AvailableLiquids(pumps []*Pump) []string {
	liquids := make([]string, len(pumps))
	for i, pump := range pumps {
		liquids[i] = pump.Liquid
	}
	sort.Strings(liquids)
	return liquids
}

func FindString(sorted []string, a string) bool {
	index := sort.SearchStrings(sorted, a)
	if index >= len(sorted) || sorted[index] != a {
		return false
	}
	return true
}

func AvailableRecipes(liquids []string, recipes map[string]*Recipe) []*Recipe {
	available := make([]*Recipe, 0, len(recipes))
	for _, recipe := range recipes {
		var found bool
		for _, ingredient := range recipe.Ingredients {
			if found = FindString(liquids, ingredient.Liquid); !found {
				break
			}
		}
		if found {
			available = append(available, recipe)
		}
	}
	return available
}

func DistinctLiquids(recipes map[string]*Recipe) []string {
	distinct := make(map[string]struct{}, 10)
	for _, recipe := range recipes {
		for _, ingredient := range recipe.Ingredients {
			distinct[ingredient.Liquid] = struct{}{}
		}
	}
	liquids := make([]string, len(distinct))
	var i int
	for liquid, _ := range distinct {
		liquids[i] = liquid
		i++
	}
	sort.Strings(liquids)
	return liquids
}

type Liquid string

type Ingredient struct {
	Liquid string
	Amount float64
}

type Recipe struct {
	Name        string
	Ingredients []*Ingredient
}

type Pump struct {
	Id                int
	Liquid            string
	TimeToVolumeRatio float64
}

type Isaac struct {
	Pumps   []*Pump
	Recipes map[string]*Recipe
}
