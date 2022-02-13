package model

type Data struct {
	Action *string
	Type   string
	Data   map[string]interface{}
}

type MethodUrl struct {
	Method string
	Url    string
}
