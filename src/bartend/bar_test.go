package bartend

import (
	"testing"
	"github.com/c2g/c2"
)

var ab = &Recipe{
	Name: "ab",
	Ingredients: []*Ingredient{
		{
			Liquid: "a",
		},
		{
			Liquid: "b",
		},
	},
}
var bc = &Recipe{
	Name: "bc",
	Ingredients: []*Ingredient{
		{
			Liquid: "b",
		},
		{
			Liquid: "c",
		},
	},
}
var cd = &Recipe{
	Name: "cd",
	Ingredients: []*Ingredient{
		{
			Liquid: "c",
		},
		{
			Liquid: "d",
		},
	},
}

func TestAvailableRecipes(t *testing.T) {
	tests := []struct {
		recipes  map[string]*Recipe
		expected []*Recipe
	}{
		{
			recipes: map[string]*Recipe{
				ab.Name : ab,
				bc.Name : bc,
				cd.Name : cd,
			},
			expected: []*Recipe{
				bc,
				cd,
			},
		},
	}
	for _, test := range tests {
		available := []string{"b","c","d"}
		if err := c2.CheckEqual(test.expected, AvailableRecipes(available, test.recipes)); err != nil {
			t.Error(err)
		}
	}
}

func TestDistinctLiquids(t *testing.T) {
	recipes := map[string]*Recipe{
		ab.Name : ab,
		bc.Name : bc,
		cd.Name : cd,
	}
	actual := DistinctLiquids(recipes)
	expected := []string{"a", "b", "c", "d"}
	if notEqual := c2.CheckEqual(expected, actual); notEqual != nil {
		t.Error(notEqual)
	}
}
