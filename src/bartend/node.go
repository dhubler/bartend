package bartend

import (
	"reflect"

	"github.com/c2stack/c2g/node"
)

// Node is management for bartend app
func Node(app *Bartend) node.Node {
	return &node.MyNode{
		OnSelect: func(r node.ContainerRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "pump":
				if r.New {
					app.Pumps = make([]*Pump, 0)
				}
				if app.Pumps != nil {
					return pumpsNode(app), nil
				}
			case "available":
				if app.Recipes != nil {
					liquids := AvailableLiquids(app.Pumps)
					available := AvailableRecipes(liquids, app.Recipes)
					return drinksNode(app, available), nil
				}
			case "recipe":
				if r.New {
					app.Recipes = make(map[string]*Recipe, 0)
				}
				if app.Recipes != nil {
					return recipesNode(app, app.Recipes), nil
				}
			}
			return nil, nil
		},
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.GetIdent() {
			case "liquids":
				hnd.Val = &node.Value{Strlist: DistinctLiquids(app.Recipes)}
			}
			return nil
		},
	}
}

func pumpsNode(app *Bartend) node.Node {
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var p *Pump
			row := int(r.Row)
			key := r.Key
			if r.New {
				row = len(app.Pumps)
				p = &Pump{}
				app.Pumps = append(app.Pumps, p)
			} else if key != nil {
				if row = key[0].Int; row < len(app.Pumps) {
					p = app.Pumps[row]
				}
			} else if row < len(app.Pumps) {
				p = app.Pumps[row]
			}
			if p != nil {
				if key == nil {
					key = node.SetValues(r.Meta.KeyMeta(), row)
				}
				return pumpNode(p), key, nil
			}
			return nil, nil, nil
		},
	}
}

func recipesNode(app *Bartend, recipes map[string]*Recipe) node.Node {
	index := node.NewIndex(recipes)
	index.Sort(func(a, b reflect.Value) bool {
		return a.String() < b.String()
	})
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var recipe *Recipe
			key := r.Key
			if r.New {
				recipe = &Recipe{Name: r.Key[0].Str}
				recipes[recipe.Name] = recipe
			} else if key != nil {
				recipe = recipes[key[0].Str]
			} else {
				name := index.NextKey(r.Row).String()
				key = node.SetValues(r.Meta.KeyMeta(), name)
				recipe = recipes[name]
			}
			if recipe != nil {
				return recipeNode(app, recipe), key, nil
			}
			return nil, nil, nil
		},
	}
}

func drinksNode(app *Bartend, drinks []*Recipe) node.Node {
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var recipe *Recipe
			row := int(r.Row)
			key := r.Key
			if key != nil {
				if row = key[0].Int; row < len(drinks) {
					recipe = drinks[row]
				}
			} else if row < len(drinks) {
				recipe = drinks[row]
				key = node.SetValues(r.Meta.KeyMeta(), row)
			}
			if recipe != nil {
				return recipeNode(app, recipe), key, nil
			}
			return nil, nil, nil
		},
	}
}

func recipeNode(app *Bartend, recipe *Recipe) node.Node {
	return &node.Extend{
		Node: node.ReflectNode(recipe),
		OnSelect: func(p node.Node, r node.ContainerRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "ingredient":
				if r.New {
					recipe.Ingredients = make([]*Ingredient, 0)
				}
				if recipe.Ingredients != nil {
					return ingredientsNode(recipe), nil
				}
			}
			return nil, nil
		},
		OnAction: func(p node.Node, r node.ActionRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "make":
				if err := app.MakeDrink(recipe); err != nil {
					return nil, err
				}
				return nil, nil
			}
			return nil, nil
		},
	}
}

func ingredientsNode(recipe *Recipe) node.Node {
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var ingredient *Ingredient
			var key = r.Key
			if r.New {
				ingredient = &Ingredient{}
				recipe.Ingredients = append(recipe.Ingredients, ingredient)
			} else if key != nil {
				for _, candidate := range recipe.Ingredients {
					if candidate.Liquid == key[0].Str {
						ingredient = candidate
						break
					}
				}
			} else if r.Row < len(recipe.Ingredients) {
				ingredient = recipe.Ingredients[r.Row]
				key = node.SetValues(r.Meta.KeyMeta(), ingredient.Liquid)
			}
			if ingredient != nil {
				return node.ReflectNode(ingredient), key, nil
			}
			return nil, nil, nil
		},
	}
}

func pumpNode(pump *Pump) node.Node {
	return &node.Extend{
		Node: node.ReflectNode(pump),
		OnAction: func(p node.Node, r node.ActionRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "on", "off":
				on := r.Meta.GetIdent() == "on"
				pump.Enable(on)
			}
			return nil, nil
		},
	}
}
