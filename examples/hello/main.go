package main

import (
	"log"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"gitlab.com/contextualcode/go-object-store/client"
	"gitlab.com/contextualcode/go-object-store/types"
)

func main() {
	client.URL = "http://localhost:31580"
	// create page view with object
	p := &PageView{
		Object: &types.Object{
			Data: map[string]interface{}{
				"name": "",
			},
		},
	}
	// set object
	res, err := client.Set([]*types.Object{p.Object}, "")
	if err != nil {
		panic(err)
	}
	p.Object.UID = res[0].UID
	log.Println("UID = " + p.Object.UID)
	// init web
	vecty.SetTitle("Object Store Get/Set Test")
	vecty.RenderBody(p)
}

// PageView is our main page component.
type PageView struct {
	vecty.Core
	Object *types.Object
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {
	name := "Unknown"
	if p.Object.Data["name"] != nil {
		name = p.Object.Data["name"].(string)
	}
	return elem.Body(
		vecty.Text("Hello "),
		elem.Div(
			elem.Input(
				vecty.Markup(
					vecty.Attribute("type", "text"),
					vecty.Attribute("value", name),
					event.Change(func(e *vecty.Event) {
						p.Object.Data["name"] = e.Target.Get("value").String()
					}),
				),
			),
			elem.Button(
				vecty.Text("Set"),
				vecty.Markup(
					event.Click(func(e *vecty.Event) {
						log.Println("Set")
						go func() {
							_, err := client.Set([]*types.Object{p.Object}, "")
							if err != nil {
								log.Println("ERROR: " + err.Error())
							}
						}()
					}),
				),
			),
			elem.Button(
				vecty.Text("Get"),
				vecty.Markup(
					event.Click(func(e *vecty.Event) {
						log.Println("Get")
						go func() {
							res, err := client.Get([]string{p.Object.UID}, "")
							if err != nil {
								log.Println("ERROR: " + err.Error())
								return
							}
							p.Object = res[0]
							vecty.Rerender(p)
						}()
					}),
				),
			),
		),
	)
}
