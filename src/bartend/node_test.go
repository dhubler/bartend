package bartend

import (
	"testing"

	"github.com/freeconf/gconf/meta"
	"github.com/freeconf/gconf/meta/yang"
	"github.com/freeconf/gconf/node"
	"github.com/freeconf/gconf/nodes"
)

func TestNode(t *testing.T) {
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
				"id": 0,
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
				"id": 1,
				"name": "straight gin",
				"ingredient": [
					{
						"liquid": "gin",
						"amount": 0.5
					}
				]
			}
		]
	}`
	ypath := &meta.FileStreamSource{Root: "../../etc/yang"}
	m := yang.RequireModule(ypath, "bartend")
	app := NewBartend()
	b := node.NewBrowser(m, Node(app))
	root := b.Root()
	if err := root.InsertFrom(nodes.ReadJSON(data)).LastErr; err != nil {
		panic(err)
	}

	in := nodes.ReadJSON(`{"multiplier":10}`)
	update := make(chan bool)
	var unsub node.NotifyCloser
	var err error
	if err := root.Find("recipe=0/make").Action(in).LastErr; err != nil {
		panic(err)
	}
	unsub, err = root.Find("current/update").Notifications(func(msg node.Selection) {
		update <- true
		defer unsub()
	})
	if err != nil {
		t.Fatal(err)
	}
	<-update
}
