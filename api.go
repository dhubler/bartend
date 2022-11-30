package bartend

import (
	"log"
	"reflect"

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
			switch r.Meta.Ident() {
			case "liquids":
				hnd.Val = val.StringList(DistinctLiquids(app.Recipes))
			}
			return nil
		},
	}
}

func currentDrinkNode(app *Bartend) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(app.Current),
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
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
		OnNotify: func(p node.Node, r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.Ident() {
			case "update":
				sub := app.OnDrinkUpdate(func(d *Drink) {
					r.Send(currentDrinkNode(app))
				})
				return sub.Close, nil
			}
			return nil, nil
		},
		OnAction: func(p node.Node, r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "stop":
				app.Current.Stop()
			}
			return nil, nil
		},
	}
}

func manualNodes(steps []*ManualStep) node.Node {
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var step *ManualStep
			key := r.Key
			var id int
			if key != nil {
				id := key[0].Value().(int)
				if id < len(steps) {
					step = steps[id]
				}
			} else if r.Row < len(steps) {
				id = r.Row
				step = steps[id]
				var err error
				if key, err = node.NewValues(r.Meta.KeyMeta(), id); err != nil {
					return nil, nil, err
				}
			}
			if step != nil {
				return stepNode(nodeutil.ReflectChild(step), id, step.Ingredient), key, nil
			}
			return nil, nil, nil
		},
	}
}

func autoNodes(steps []*AutoStep) node.Node {
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var step *AutoStep
			key := r.Key
			var id int
			if key != nil {
				id = key[0].Value().(int)
				if id < len(steps) {
					step = steps[id]
				}
			} else if r.Row < len(steps) {
				id = r.Row
				step = steps[id]
				var err error
				if key, err = node.NewValues(r.Meta.KeyMeta(), id); err != nil {
					return nil, nil, err
				}
			}
			if step != nil {
				return stepNode(nodeutil.ReflectChild(step), id, step.Ingredient), key, nil
			}
			return nil, nil, nil
		},
	}
}

func stepNode(base node.Node, id int, ingredient *Ingredient) node.Node {
	return &nodeutil.Extend{
		Base: base,
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "id":
				hnd.Val = val.Int32(id)
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

func recipesNode(app *Bartend, recipes map[int]*Recipe) node.Node {
	index := node.NewIndex(recipes)
	index.Sort(func(a, b reflect.Value) bool {
		return a.Int() < b.Int()
	})
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var recipe *Recipe
			key := r.Key
			if r.New {
				recipe = &Recipe{Id: r.Key[0].Value().(int)}
				recipes[recipe.Id] = recipe
			} else if key != nil {
				recipe = recipes[key[0].Value().(int)]
			} else {
				v := index.NextKey(r.Row)
				if v != node.NO_VALUE {
					id := int(v.Int())
					var err error
					if key, err = node.NewValues(r.Meta.KeyMeta(), id); err != nil {
						return nil, nil, err
					}
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
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var recipe *Recipe
			row := int(r.Row)
			key := r.Key
			if key != nil {
				if row = key[0].Value().(int); row < len(drinks) {
					recipe = drinks[row]
				}
			} else if row < len(drinks) {
				recipe = drinks[row]
				var err error
				if key, err = node.NewValues(r.Meta.KeyMeta(), row); err != nil {
					return nil, nil, err
				}
			}
			if recipe != nil {
				return recipeNode(app, recipe), key, nil
			}
			return nil, nil, nil
		},
	}
}

func recipeNode(app *Bartend, recipe *Recipe) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(recipe),
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
