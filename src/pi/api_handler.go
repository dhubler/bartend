package pi

import (
	"github.com/c2g/node"
	"github.com/c2g/meta/yang"
	"reflect"
)

type ApiHandler struct {
}

func (self ApiHandler) Manage(root *node.Browser, app *Isaac) node.Node {
	return &node.MyNode{
		OnSelect: func(r node.ContainerRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "pump":
				if r.New {
					app.Pumps = make([]*Pump, 0)
				}
				if app.Pumps != nil {
					return self.Pumps(app), nil
				}
			case "available":
				if app.Recipes != nil {
					liquids := AvailableLiquids(app.Pumps)
					available := AvailableRecipes(liquids, app.Recipes)
					return self.Drinks(available), nil
				}
			case "recipe":
				if r.New {
					app.Recipes = make(map[string]*Recipe, 0)
				}
				if app.Recipes != nil {
					return self.Recipes(app.Recipes), nil
				}
			}
			return nil, nil
		},
		OnRead: func(r node.FieldRequest) (*node.Value, error) {
			switch r.Meta.GetIdent() {
			case "liquids":
				return &node.Value{Strlist:DistinctLiquids(app.Recipes)}, nil
			}
			return nil, nil
		},
	}
}


func (self ApiHandler) Pumps(app *Isaac) node.Node {
	return &node.MyNode{
		OnNext:func(r node.ListRequest) (node.Node, []*node.Value, error) {
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
				return self.Pump(p), key, nil
			}
			return nil, nil, nil
		},
	}
}

func (self ApiHandler) Recipes(recipes map[string]*Recipe) node.Node {
	index := node.NewIndex(recipes)
	index.Sort(func(a, b reflect.Value) bool {
		return a.String() < b.String()
	} )
	return &node.MyNode{
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var recipe *Recipe
			key := r.Key
			if r.New {
				recipe = &Recipe{Name:r.Key[0].Str}
				recipes[recipe.Name] = recipe
			} else if key != nil {
				recipe = recipes[key[0].Str]
			} else {
				name := index.NextKey(r.Row).String()
				key = node.SetValues(r.Meta.KeyMeta(), name)
				recipe = recipes[name]
			}
			if recipe != nil {
				return self.Recipe(recipe), key, nil
			}
			return nil, nil, nil
		},
	}
}

func (self ApiHandler) Drinks(drinks []*Recipe) node.Node {
	return &node.MyNode{
		OnNext:func(r node.ListRequest) (node.Node, []*node.Value, error) {
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
				return self.Recipe(recipe), key, nil
			}
			return nil, nil, nil
		},
	}
}

func (self ApiHandler) Recipe(recipe *Recipe) node.Node {
	return &node.Extend{
		Node: node.MarshalContainer(recipe),
		OnSelect : func(p node.Node, r node.ContainerRequest) (node.Node, error) {
			switch r.Meta.GetIdent() {
			case "ingredient":
				if r.New {
					recipe.Ingredients = make([]*Ingredient, 0)
				}
				if recipe.Ingredients != nil {
					return self.Ingredients(recipe), nil
				}
			}
			return nil, nil
		},
	}
}

func (self ApiHandler) Ingredients(recipe *Recipe) node.Node {
	return &node.MyNode {
		OnNext: func(r node.ListRequest) (node.Node, []*node.Value, error) {
			var ingredient *Ingredient
			if r.New {
				ingredient = &Ingredient{}
				recipe.Ingredients = append(recipe.Ingredients, ingredient)
			} else {
				if r.Row < len(recipe.Ingredients) {
					ingredient = recipe.Ingredients[r.Row]
				}
			}
			if ingredient != nil {
				return node.MarshalContainer(ingredient), nil, nil
			}
			return nil, nil, nil
		},
	}
}

func (self ApiHandler) Pump(pump *Pump) node.Node {
	return node.MarshalContainer(pump)
}

func init() {
	yang.InternalYang()["isaac"] = `
module isaac {
  namespace "";
  prefix "";
  revision 0;

  list pump {
    key "id";
    leaf id {
      type int32;
    }
    leaf timeToVolumeRatio {
      description "Number of millisecs to turn on pump to pour one milliliter";
      type decimal64;
    }
    leaf liquid {
      type string;
    }
  }

  leaf-list liquids {
    config "false";
    type string;
  }

  grouping drink {
    leaf name {
      type string;
    }
    list ingredient {
      leaf liquid {
        type string;
      }
      leaf amount {
        type decimal64;
      }
    }
  }

  list available {
    key "name";
    uses drink;
    action make {
      input {}
    }
  }

  list recipe {
    key "name";
    uses drink;
  }

  action calibrateStart {
    input {
    }
  }

  action calibrateStop {
    input {
      leaf update {
        type boolean;
      }
    }
    output {
      leaf timeToVolumeRatio {
        type decimal64;
      }
    }
  }
}
`
}

