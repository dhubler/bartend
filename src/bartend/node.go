package bartend

import (
	"log"
	"reflect"

	"github.com/c2stack/c2g/node"
)

// Node is management for bartend app
func Node(app *Bartend) node.Node {
	return &node.MyNode{
		OnChild: func(r node.ChildRequest) (node.Node, error) {
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
					available := AutomaticRecipes(liquids, app.Recipes)
					return drinksNode(app, available), nil
				}
			case "recipe":
				if r.New {
					app.Recipes = make(map[int]*Recipe, 0)
				}
				if app.Recipes != nil {
					return recipesNode(app, app.Recipes), nil
				}
			case "current":
				if app.Current != nil {
					return currentDrinkNode(app), nil
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

func currentDrinkNode(app *Bartend) node.Node {
	return &node.Extend{
		Node: node.ReflectNode(app.Current),
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "auto":
				if len(app.Current.Automatic) > 0 {
					return autoNodes(app.Current.Automatic), nil
				}
			case "manual":
				if len(app.Current.Manual) > 0 {
					return manualNodes(app.Current.Manual), nil
				}
			default:
				return p.Child(r)
			}
			return nil, nil
		},
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.GetIdent() {
			case "percentComplete":
				hnd.Val = &node.Value{Int: app.Current.PercentComplete()}
			case "complete":
				hnd.Val = &node.Value{Bool: app.Current.Complete()}
			default:
				return p.Field(r, hnd)
			}
			return nil
		},
		OnNotify: func(p node.Node, r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.GetIdent() {
			case "update":
				sub := app.OnDrinkUpdate(func(d *Drink) {
					r.Stream.Notify(r.Meta, r.Selection.Path, currentDrinkNode(app))
				})
				return sub.Close, nil
			}
			return nil, nil
		},
		OnAction: func(p node.Node, r node.ActionRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "stop":
				app.Current.Stop()
			}
			return nil, nil
		},
	}
}

func manualNodes(steps []*ManualStep) node.Node {
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var step *ManualStep
			key := r.Key
			var id int
			if key != nil {
				id := key[0].Int
				if id < len(steps) {
					step = steps[id]
				}
			} else if r.Row < len(steps) {
				step = steps[id]
				key = node.SetValues(r.Meta.KeyMeta(), id)
			}
			if step != nil {
				return stepNode(node.ReflectNode(step), id, step.Ingredient), key, nil
			}
			return nil, nil, nil
		},
	}
}

func autoNodes(steps []*AutoStep) node.Node {
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var step *AutoStep
			key := r.Key
			var id int
			if key != nil {
				id = key[0].Int
				if id < len(steps) {
					step = steps[id]
				}
			} else if r.Row < len(steps) {
				id = r.Row
				step = steps[id]
				key = node.SetValues(r.Meta.KeyMeta(), id)
			}
			if step != nil {
				return stepNode(node.ReflectNode(step), id, step.Ingredient), key, nil
			}
			return nil, nil, nil
		},
	}
}

func stepNode(base node.Node, id int, ingredient *Ingredient) node.Node {
	return &node.Extend{
		Node: base,
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.GetIdent() {
			case "id":
				hnd.Val = &node.Value{Int: id}
			default:
				return p.Field(r, hnd)
			}
			return nil
		},
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "ingredient":
				return ingredientNode(ingredient), nil
			}
			return p.Child(r)
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

func recipesNode(app *Bartend, recipes map[int]*Recipe) node.Node {
	index := node.NewIndex(recipes)
	index.Sort(func(a, b reflect.Value) bool {
		return a.Int() < b.Int()
	})
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var recipe *Recipe
			key := r.Key
			if r.New {
				recipe = &Recipe{Id: r.Key[0].Int}
				recipes[recipe.Id] = recipe
			} else if key != nil {
				recipe = recipes[key[0].Int]
			} else {
				v := index.NextKey(r.Row)
				if v != node.NO_VALUE {
					id := int(v.Int())
					key = node.SetValues(r.Meta.KeyMeta(), id)
					recipe = recipes[id]
				}
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
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
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
				return ingredientNode(ingredient), key, nil
			}
			return nil, nil, nil
		},
	}
}

func ingredientNode(ingredient *Ingredient) node.Node {
	return &node.Extend{
		Node: node.ReflectNode(ingredient),
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.GetIdent() {
			case "weight":
				hnd.Val = &node.Value{Int: ingredient.Weight()}
			default:
				return p.Field(r, hnd)
			}
			return nil
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
				log.Printf("pump=%d, on=%v", pump.GpioPin, on)
				return nil, pump.Enable(on)
			}
			return nil, nil
		},
	}
}
