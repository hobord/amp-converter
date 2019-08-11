package web

type Request struct {
	BaseUrl  string
	document string
}

type Response struct {
	document     string
	dependencies []Dependency
}

type Dependency struct {
	name      string
	headertag string
}
