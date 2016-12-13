package bartend

import (
	"container/list"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/c2stack/c2g/c2"
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

// AutomaticRecipes is list of drinks that can be made completely automatically
func AutomaticRecipes(liquids []string, recipes map[int]*Recipe) []*Recipe {
	available := make([]*Recipe, 0, len(recipes))
	for _, recipe := range recipes {
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

func DistinctLiquids(recipes map[int]*Recipe) []string {
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

func (self *Ingredient) Weight() int {
	return int(self.Amount * LiquidToGrams)
}

type Recipe struct {
	Id          int
	Name        string
	Description string
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

func (self *Pump) Enable(on bool) error {
	var v int
	// not sure why, but  1 - off,  0 - on
	if on {
		v = embd.Low
	} else {
		v = embd.High
	}
	p, err := GetPin(self.GpioPin)
	if err != nil {
		log.Printf("Err pin %d - %s", self.GpioPin, err)
		return err
	}
	return p.Write(v)
}

func (self *Pump) calculatePourTime(amount float64) time.Duration {
	oneUnit := time.Millisecond * time.Duration(self.TimeToVolumeRatioMs)
	return time.Duration(amount * float64(oneUnit))
}

type Bartend struct {
	Current   *Drink
	Pumps     []*Pump
	Recipes   map[int]*Recipe
	listeners *list.List
}

func NewBartend() *Bartend {
	return &Bartend{
		listeners: list.New(),
	}
}

var DrinkInProgress = c2.NewErrC("Drink in progress", 400)

type Drink struct {
	Automatic []*AutoStep
	Manual    []*ManualStep
	ticker    *time.Ticker
	Aborted   bool
}

func (self *Drink) Complete() bool {
	if self.Aborted == true {
		return true
	}
	for _, a := range self.Automatic {
		if !a.Complete {
			return false
		}
	}
	for _, m := range self.Manual {
		if !m.Complete {
			return false
		}
	}
	return true
}

func (self *Drink) Stop() {
	self.ticker.Stop()
	self.allPumpsOn(false)
	self.Aborted = true
}

func (self *Drink) PercentComplete() int {
	n := len(self.Automatic) + len(self.Manual)
	var pct int
	for _, a := range self.Automatic {
		pct += (a.PercentComplete / n)
	}
	for _, m := range self.Manual {
		if m.Complete {
			pct += 100 / n
		}
	}
	return pct
}

type ManualStep struct {
	Ingredient *Ingredient
	Complete   bool
}

type AutoStep struct {
	pump            *Pump
	Ingredient      *Ingredient
	PourTime        time.Duration
	PercentComplete int
	Complete        bool
}

func (self *AutoStep) pumpOn(on bool) error {
	return self.pump.Enable(on)
}

func (self *AutoStep) calculatePercentageDone(t time.Duration) {
	if t > self.PourTime {
		self.PercentComplete = 100
	} else if t == 0 {
		self.PercentComplete = 0
	} else {
		self.PercentComplete = int(100 * (1 - (float32(self.PourTime-t) / float32(self.PourTime))))
	}
}

func (self *AutoStep) update(t time.Duration) error {
	self.calculatePercentageDone(t)
	complete := t > self.PourTime
	if complete != self.Complete {
		if err := self.pumpOn(!complete); err != nil {
			return err
		}
		self.Complete = complete
		c2.Debug.Printf("complete %v", self.Complete)
	}
	return nil
}

func (self *Drink) allPumpsOn(on bool) error {
	var err error
	for _, step := range self.Automatic {
		if e := step.pumpOn(on); e != nil {
			err = e
		}
	}
	return err
}

func (self *Bartend) OnDrinkUpdate(l DrinkProgressListener) c2.Subscription {
	return c2.NewSubscription(self.listeners, self.listeners.PushBack(l))
}

func (self *Bartend) updateJob(job *Drink) {
	e := self.listeners.Front()
	for e != nil {
		e.Value.(DrinkProgressListener)(job)
		e = e.Next()
	}
}

type DrinkProgressListener func(job *Drink)

func (self *Drink) Start(l DrinkProgressListener) {
	timeStep := time.Millisecond * 100
	self.ticker = time.NewTicker(timeStep)
	var t time.Duration
	self.allPumpsOn(true)
	defer func() {
		// shouldn't be nec. unless error happened
		self.allPumpsOn(false)
	}()
	for {
		var incomplete bool
		for _, step := range self.Automatic {
			if err := step.update(t); err != nil {
				log.Printf("Cannot update pump : %s", err)
				break
			}
			if !step.Complete {
				incomplete = true
			}
		}
		l(self)
		if !incomplete {
			break
		}
		if _, more := <-self.ticker.C; !more {
			break
		}
		t += timeStep
	}
}

func (self *Bartend) MakeDrink(recipe *Recipe) error {
	if self.Current != nil && !self.Current.Complete() {
		return DrinkInProgress
	}
	var drink Drink
	self.Current = &drink
	for _, ingredient := range recipe.Ingredients {
		p := FindPumpByLiquid(self.Pumps, ingredient.Liquid)
		if p == nil {
			drink.Manual = append(drink.Manual, &ManualStep{
				Ingredient: ingredient,
			})
		} else {
			drink.Automatic = append(drink.Automatic, &AutoStep{
				Ingredient: ingredient,
				pump:       p,
				PourTime:   p.calculatePourTime(ingredient.Amount),
			})
		}
	}
	go drink.Start(self.updateJob)
	// drink responsibly
	return nil
}
