package bartend

import (
	"flag"
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

var updateGoldFiles = flag.Bool("update", false, "update expected golden file(s)")

func TestApi(t *testing.T) {
	data := `
	{
		"pump": [
			{
				"id": 0,
				"gpioPin": 1,
				"liquid": "gin",
				"timeToVolumeRatioMs": 10
			},
			{
				"id": 1,
				"gpioPin": 2,
				"liquid": "juice",
				"timeToVolumeRatioMs": 10
			}
		],
		"recipe": [
			{
				"name": "snoop",
				"description": "good for sipping",
				"ingredient": [
					{
						"liquid": "gin",
						"amount": 1.5
					},
					{
						"liquid": "juice",
						"amount": 1
					}
				]
			},
			{
				"name": "straight gin",
				"ingredient": [
					{
						"liquid": "gin",
						"amount": 1
					}
				]
			},
			{
				"name": "screwdriver gin",
				"ingredient": [
					{
						"liquid": "vodka",
						"amount": 1
					},
					{
						"liquid": "juice",
						"amount": 1
					}
				]
			}
		]
	}`
	ypath := source.Dir("./etc/yang")
	m := parser.RequireModule(ypath, "bartend")
	app := NewBartend()
	b := node.NewBrowser(m, Node(app))
	root := b.Root()
	if err := root.InsertFrom(nodeutil.ReadJSON(data)).LastErr; err != nil {
		t.Fatal(err)
	}

	in := nodeutil.ReadJSON(`{"multiplier":10}`)
	update := make(chan bool)
	var unsub node.NotifyCloser
	var err error
	if err := root.Find("available=snoop/make").Action(in).LastErr; err != nil {
		t.Fatal(err)
	}
	unsub, err = root.Find("drink/update").Notifications(func(msg node.Notification) {
		update <- true
		defer unsub()
	})
	if err != nil {
		t.Fatal(err)
	}
	<-update

	actual, err := nodeutil.WritePrettyJSON(root)
	fc.AssertEqual(t, nil, err)
	fc.Gold(t, *updateGoldFiles, []byte(actual), "testdata/api.json")
}
