package bartend

import (
	"testing"
	"time"

	"github.com/freeconf/yang/fc"
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
		recipes  []*Recipe
		expected []*Recipe
	}{
		{
			recipes: []*Recipe{
				ab,
				bc,
				cd,
			},
			expected: []*Recipe{
				bc,
				cd,
			},
		},
	}
	for _, test := range tests {
		available := []string{"b", "c", "d"}
		fc.AssertEqual(t, test.expected, Recipes(available, test.recipes))
	}
}

func TestDistinctLiquids(t *testing.T) {
	recipes := []*Recipe{
		ab,
		bc,
		cd,
	}
	actual := DistinctLiquids(recipes)
	expected := []string{"a", "b", "c", "d"}
	fc.AssertEqual(t, expected, actual)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("%s != %s", a, b)
	}
}

func TestPourTime(t *testing.T) {
	p := Pump{
		TimeToVolumeRatioMs: 2,
	}
	assertEqual(t, time.Duration(10*time.Millisecond), p.calculatePourTime(5))
}

func TestPercentComplete(t *testing.T) {
	a := &Step{PourTime: 10}
	a.calculatePercentageDone(1)
	assertEqual(t, 10, a.PercentComplete)
	a.calculatePercentageDone(5)
	assertEqual(t, 50, a.PercentComplete)
}

func TestMakeDrink(t *testing.T) {
	b := NewBartend()
	b.Pumps = []*Pump{
		{
			Liquid:              "OJ",
			TimeToVolumeRatioMs: 1,
		},
		{
			Liquid:              "Vodka",
			TimeToVolumeRatioMs: 1,
		},
	}
	r := &Recipe{
		Ingredients: []*Ingredient{
			{
				Amount: 1,
				Liquid: "Vodka",
			},
			{
				Amount: 2,
				Liquid: "OJ",
			},
		},
	}

	if err := b.MakeDrink(r, 1); err != nil {
		t.Fatal(err)
	}
	if b.Current == nil {
		t.Fatal("no job in progress")
	}
	assertEqual(t, 2, len(b.Current.Pour))
	assertEqual(t, false, b.Current.Complete())
	if err := b.MakeDrink(r, 1); err == nil {
		t.Fatal("supposed to get error that drink is in progress")
	}
	t.Log(b.Current.Pour[0].PourTime)
	b.Current.Pour[0].Complete = true
	done := make(chan bool)
	go func() {
		<-time.After(b.Current.Pour[0].PourTime + time.Minute)
		close(done)
	}()
	b.OnDrinkUpdate(func(d *Drink) {
		if d.Complete() {
			done <- true
		}
	})
	if _, ok := <-done; !ok {
		t.Error("timeout")
	}

	rBad := &Recipe{
		Name: "Gimlet",
		Ingredients: []*Ingredient{
			{
				Amount: 1,
				Liquid: "Gin",
			},
			{
				Amount: 1,
				Liquid: "Lime juice",
			},
		},
	}
	if err := b.MakeDrink(rBad, 1); err == nil {
		t.Error("expected error")
	}
	if err := b.MakeDrink(r, 1); err != nil {
		t.Fatal(err)
	}
}
