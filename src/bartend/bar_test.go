package bartend

import (
	"testing"
	"time"

	"github.com/c2stack/c2g/c2"
)

var ab = &Recipe{
	Id:   0,
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
	Id:   1,
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
	Id:   2,
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
		recipes  map[int]*Recipe
		expected []*Recipe
	}{
		{
			recipes: map[int]*Recipe{
				ab.Id: ab,
				bc.Id: bc,
				cd.Id: cd,
			},
			expected: []*Recipe{
				bc,
				cd,
			},
		},
	}
	for _, test := range tests {
		available := []string{"b", "c", "d"}
		if err := c2.CheckEqual(test.expected, AutomaticRecipes(available, test.recipes)); err != nil {
			t.Error(err)
		}
	}
}

func TestDistinctLiquids(t *testing.T) {
	recipes := map[int]*Recipe{
		ab.Id: ab,
		bc.Id: bc,
		cd.Id: cd,
	}
	actual := DistinctLiquids(recipes)
	expected := []string{"a", "b", "c", "d"}
	if notEqual := c2.CheckEqual(expected, actual); notEqual != nil {
		t.Error(notEqual)
	}
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
	a := &AutoStep{PourTime: 10}
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
		t.Error(err)
	}
	if b.Current == nil {
		t.Error("no job in progress")
	}
	assertEqual(t, 1, len(b.Current.Manual))
	assertEqual(t, 1, len(b.Current.Automatic))
	assertEqual(t, false, b.Current.Complete())
	if err := b.MakeDrink(r, 1); err == nil {
		t.Error("supposed to get error that drink is in progress")
	}
	t.Log(b.Current.Automatic[0].PourTime)
	b.Current.Manual[0].Complete = true
	done := make(chan bool)
	go func() {
		<-time.After(b.Current.Automatic[0].PourTime + (1 * time.Second))
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
}
