package entity

type Route struct {
	Name    string
	Path    string
	Methods []Method
}

type Method struct {
	Type        string
	Name        string
	Summary     string
	Description string
	Request     Schema
	Response    Schema
}
