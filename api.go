package bartend

import (
	"log"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

// Node is management for bartend app
func Node(app *Bartend) node.Node {
	return &nodeutil.Basic{
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
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
					available := Recipes(liquids, app.Recipes)
					if len(available) > 0 {
						return recipesNode(app, available), nil
					}
				}
			case "recipe":
				if r.New {
					app.Recipes = make([]*Recipe, 0)
				}
				if app.Recipes != nil {
					return recipesNode(app, app.Recipes), nil
				}
			case "drink":
				return drinkNode(app), nil
			}
			return nil, nil
		},
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "liquids":
				hnd.Val = val.StringList(DistinctLiquids(app.Recipes))
			}
			return nil
		},
	}
}

func drinkNode(app *Bartend) node.Node {
	return &nodeutil.Basic{
		OnNotify: func(r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.Ident() {
			case "update":
				sub := app.OnDrinkUpdate(func(d *Drink) {
					r.Send(drinkUpdateNode(app))
				})
				return sub.Close, nil
			}
			return nil, nil
		},
		OnAction: func(r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "stop":
				app.Current.Stop()
			}
			return nil, nil
		},
	}
}

func drinkUpdateNode(app *Bartend) node.Node {
	drink := app.Current
	if drink == nil {
		drink = &Drink{}
	}
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(drink),
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "pour":
				if len(app.Current.Pour) > 0 {
					return pourNodes(app.Current.Pour), nil
				}
			default:
				return p.Child(r)
			}
			return nil, nil
		},
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "percentComplete":
				hnd.Val = val.Int32(app.Current.PercentComplete())
			case "complete":
				hnd.Val = val.Bool(app.Current.Complete())
			default:
				return p.Field(r, hnd)
			}
			return nil
		},
	}
}

func pourNodes(pours []*Step) node.Node {
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var pour *Step
			key := r.Key
			var id int
			if key != nil {
				id = key[0].Value().(int)
				if id < len(pours) {
					pour = pours[id]
				}
			} else if r.Row < len(pours) {
				id = r.Row
				pour = pours[id]
				var err error
				if key, err = node.NewValues(r.Meta.KeyMeta(), id); err != nil {
					return nil, nil, err
				}
			}
			if pour != nil {
				return pourNode(pour, pour.Ingredient), key, nil
			}
			return nil, nil, nil
		},
	}
}

func pourNode(pour *Step, ingredient *Ingredient) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(pour),
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "pumpId":
				hnd.Val = val.Int32(pour.pump.Id)
			case "liquid":
				hnd.Val = val.String(ingredient.Liquid)
			case "amount":
				hnd.Val = val.Decimal64(ingredient.Amount)
			default:
				return p.Field(r, hnd)
			}
			return nil
		},
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "ingredient":
				return ingredientNode(ingredient), nil
			}
			return p.Child(r)
		},
	}
}

func pumpsNode(app *Bartend) node.Node {
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var p *Pump
			row := int(r.Row)
			key := r.Key
			if r.New {
				row = len(app.Pumps)
				p = &Pump{}
				app.Pumps = append(app.Pumps, p)
			} else if key != nil {
				if row = key[0].Value().(int); row < len(app.Pumps) {
					p = app.Pumps[row]
				}
			} else if row < len(app.Pumps) {
				p = app.Pumps[row]
			}
			if p != nil {
				if key == nil {
					if k, err := node.NewValues(r.Meta.KeyMeta(), row); err != nil {
						return nil, nil, err
					} else {
						key = k
					}
				}
				return pumpNode(p), key, nil
			}
			return nil, nil, nil
		},
	}
}

func recipesNode(app *Bartend, recipes []*Recipe) node.Node {
	return &nodeutil.Basic{
		OnNextItem: func(r node.ListRequest) nodeutil.BasicNextItem {
			var recipe *Recipe
			var id int
			return nodeutil.BasicNextItem{
				New: func() error {
					recipe = &Recipe{}
					id = len(app.Recipes)
					app.Recipes = append(app.Recipes, recipe)
					return nil
				},
				GetByKey: func() error {
					if r.Key[0] != nil {
						name := r.Key[0].String()
						for _, candidate := range recipes {
							if candidate.Name == name {
								recipe = candidate
								break
							}
						}
					}
					return nil
				},
				GetByRow: func() ([]val.Value, error) {
					if r.Row < len(recipes) {
						recipe = recipes[r.Row]
						id = r.Row
						return []val.Value{val.Int32(r.Row)}, nil
					}
					return nil, nil
				},
				Node: func() (node.Node, error) {
					if recipe != nil {
						return recipeNode(app, id, recipe), nil
					}
					return nil, nil
				},
			}
		},
	}
}

func recipeNode(app *Bartend, id int, recipe *Recipe) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(recipe),
		OnField: func(parent node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "id":
				hnd.Val = val.Int32(id)
			default:
				return parent.Field(r, hnd)
			}
			return nil
		},
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
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
			switch r.Meta.Ident() {
			case "make":
				scale := 1.0
				if !r.Input.IsNil() {
					if scaleVal, err := r.Input.GetValue("multiplier"); err != nil {
						return nil, err
					} else {
						scale = scaleVal.Value().(float64)
					}
				}
				if err := app.MakeDrink(recipe, scale); err != nil {
					return nil, err
				}
				return nil, nil
			}
			return nil, nil
		},
	}
}

func ingredientsNode(recipe *Recipe) node.Node {
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var ingredient *Ingredient
			var key = r.Key
			if r.New {
				ingredient = &Ingredient{}
				recipe.Ingredients = append(recipe.Ingredients, ingredient)
			} else if key != nil {
				for _, candidate := range recipe.Ingredients {
					if candidate.Liquid == key[0].String() {
						ingredient = candidate
						break
					}
				}
			} else if r.Row < len(recipe.Ingredients) {
				ingredient = recipe.Ingredients[r.Row]
				var err error
				if key, err = node.NewValues(r.Meta.KeyMeta(), ingredient.Liquid); err != nil {
					return nil, nil, err
				}
			}
			if ingredient != nil {
				return ingredientNode(ingredient), key, nil
			}
			return nil, nil, nil
		},
	}
}

func ingredientNode(ingredient *Ingredient) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(ingredient),
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "weight":
				hnd.Val = val.Int32(ingredient.Weight())
			default:
				return p.Field(r, hnd)
			}
			return nil
		},
	}
}

func pumpNode(pump *Pump) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(pump),
		OnAction: func(p node.Node, r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "on", "off":
				on := r.Meta.Ident() == "on"
				log.Printf("pump=%d, on=%v", pump.GpioPin, on)
				return nil, pump.Enable(on)
			}
			return nil, nil
		},
	}
}
