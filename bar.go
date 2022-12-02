package bartend

import (
	"container/list"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/nodeutil"
	"github.com/kidoman/embd"
)

func AvailableLiquids(pumps []*Pump) []string {
	liquids := make([]string, len(pumps))
	for i, pump := range pumps {
		liquids[i] = pump.Liquid
	}
	sort.Strings(liquids)
	return liquids
}

func findStringInSlice(sorted []string, a string) bool {
	index := sort.SearchStrings(sorted, a)
	if index >= len(sorted) || sorted[index] != a {
		return false
	}
	return true
}

type ByName []*Recipe

func (a ByName) Len() int {
	return len(a)
}
func (a ByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByName) Less(i, j int) bool {
	return strings.Compare(a[i].Name, a[j].Name) < 0
}

// Recipes is list of drinks that can be made completely automatically
func Recipes(liquids []string, all []*Recipe) []*Recipe {
	available := make([]*Recipe, 0, len(all))
	for _, recipe := range all {
		var found bool
		for _, ingredient := range recipe.Ingredients {
			if found = findStringInSlice(liquids, ingredient.Liquid); !found {
				break
			}
		}
		if found {
			available = append(available, recipe)
		}
	}
	sort.Sort(ByName(available))
	return available
}

func DistinctLiquids(recipes []*Recipe) []string {
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

func FindPumpByLiquid(pumps []*Pump, liquid string) *Pump {
	for _, pump := range pumps {
		if pump.Liquid == liquid {
			return pump
		}
	}
	return nil
}

type Liquid string

type Ingredient struct {
	Liquid string
	Amount float64
}

func (i *Ingredient) Scale(scale float64) *Ingredient {
	copy := *i
	copy.Amount = copy.Amount * scale
	return &copy
}

func (i *Ingredient) Weight() int {
	return int(i.Amount * LiquidToGrams)
}

type Recipe struct {
	Name        string
	Description string
	MadeCount   float64
	Ingredients []*Ingredient
}

// Standard volume to weight ratio for distilled water
const LiquidToGrams = 29.57

type Pump struct {
	Id                  int
	GpioPin             int
	Liquid              string
	TimeToVolumeRatioMs int
}

func (p *Pump) Enable(on bool) error {
	var v int
	// not sure why, but  1 - off,  0 - on
	if on {
		v = embd.Low
	} else {
		v = embd.High
	}
	pin, err := GetPin(p.GpioPin)
	if err != nil {
		log.Printf("Err pin %d - %s", p.GpioPin, err)
		return err
	}
	return pin.Write(v)
}

func (p *Pump) calculatePourTime(amount float64) time.Duration {
	oneUnit := time.Millisecond * time.Duration(p.TimeToVolumeRatioMs)
	return time.Duration(amount * float64(oneUnit))
}

type Bartend struct {
	Current   *Drink
	Pumps     []*Pump
	Recipes   []*Recipe
	listeners *list.List
}

func NewBartend() *Bartend {
	return &Bartend{
		listeners: list.New(),
	}
}

var ErrDrinkInProgress = fmt.Errorf("drink in progress. %w", fc.BadRequestError)

type Drink struct {
	Name    string
	Pour    []*Step
	ticker  *time.Ticker
	Aborted bool
}

func (d *Drink) Complete() bool {
	if d.Aborted {
		return true
	}
	for _, a := range d.Pour {
		if !a.Complete {
			return false
		}
	}
	return true
}

func (d *Drink) Stop() {
	d.ticker.Stop()
	d.allPumpsOn(false)
	d.Aborted = true
}

func (d *Drink) PercentComplete() int {
	n := len(d.Pour)
	var pct int
	for _, a := range d.Pour {
		pct += (a.PercentComplete / n)
	}
	return pct
}

type Step struct {
	pump            *Pump
	Ingredient      *Ingredient
	PourTime        time.Duration
	PercentComplete int
	Complete        bool
}

func (s *Step) pumpOn(on bool) error {
	return s.pump.Enable(on)
}

func (s *Step) calculatePercentageDone(t time.Duration) {
	if t > s.PourTime {
		s.PercentComplete = 100
	} else if t == 0 {
		s.PercentComplete = 0
	} else {
		s.PercentComplete = int(100 * (1 - (float32(s.PourTime-t) / float32(s.PourTime))))
	}
}

func (s *Step) update(t time.Duration) error {
	s.calculatePercentageDone(t)
	complete := t > s.PourTime
	if complete != s.Complete {
		if err := s.pumpOn(!complete); err != nil {
			return err
		}
		s.Complete = complete
	}
	return nil
}

func (b *Bartend) OnDrinkUpdate(l DrinkProgressListener) nodeutil.Subscription {
	return nodeutil.NewSubscription(b.listeners, b.listeners.PushBack(l))
}

func (b *Bartend) updateJob(job *Drink) {
	e := b.listeners.Front()
	for e != nil {
		e.Value.(DrinkProgressListener)(job)
		e = e.Next()
	}
}

type DrinkProgressListener func(job *Drink)

func (d *Drink) allPumpsOn(on bool) error {
	var err error
	for _, step := range d.Pour {
		if e := step.pumpOn(on); e != nil {
			err = e
		}
	}
	return err
}

func (d *Drink) Start(l DrinkProgressListener) {
	timeStep := time.Millisecond * 100
	d.ticker = time.NewTicker(timeStep)
	var t time.Duration
	d.allPumpsOn(true)
	defer func() {
		// shouldn't be nec. unless error happened
		d.allPumpsOn(false)
	}()
	for {
		var incomplete bool
		for _, step := range d.Pour {
			if err := step.update(t); err != nil {
				log.Printf("Cannot update pump : %s", err)
				break
			}
			if !step.Complete {
				incomplete = true
			}
		}
		l(d)
		if !incomplete {
			break
		}
		if _, more := <-d.ticker.C; !more {
			break
		}
		t += timeStep
	}
}

func (b *Bartend) MakeDrink(recipe *Recipe, scale float64) error {
	if b.Current != nil && !b.Current.Complete() {
		return ErrDrinkInProgress
	}
	drink := &Drink{Name: recipe.Name}
	b.Current = drink
	for _, ingredient := range recipe.Ingredients {
		scaled := ingredient.Scale(scale)
		p := FindPumpByLiquid(b.Pumps, ingredient.Liquid)
		if p == nil {
			return fmt.Errorf("%s is not available on any pump", ingredient.Liquid)
		} else {
			drink.Pour = append(drink.Pour, &Step{
				Ingredient: scaled,
				pump:       p,
				PourTime:   p.calculatePourTime(scaled.Amount),
			})
		}
	}
	recipe.MadeCount += scale
	go drink.Start(b.updateJob)
	// drink responsibly
	return nil
}
