package kamemaru

type List struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

func (k *kamemaru) Route() {
	// front end
	//k.Echo.File("/", "assets/index.html")
	// backend api
	k.Echo.POST("/api/v1/list", k.List)
}
