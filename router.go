package kamemaru

import "fmt"

type List struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

func (k *kamemaru) Router() {
	// front end
	k.Echo.File("/", "public/index.html")

	// backend api
	k.Echo.POST("/api/v1/list", k.List)

	k.RunServe()
}

func (k *kamemaru) RunServe() {
	k.Echo.Start(fmt.Sprintf(":%d", k.Env.Port))
}
