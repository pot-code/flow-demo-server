package hello

type PostHelloDto struct {
	Foo *string `json:"foo" validate:"required"`
	Bar *int32  `json:"bar" validate:"required"`
}
