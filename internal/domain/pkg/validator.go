package pkg

type Validator interface {
	Struct(s interface{}) error
}
