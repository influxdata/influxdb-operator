package foo

type fooCodec struct {
}

func (fooCodec) Marshal(v interface{}) ([]byte, error) {
	switch t := v.(type) {
	case *Request:
		return []byte(t.In), nil
	case *Response:
		return []byte(t.Out), nil
	default:
		panic("invalid type")
	}
}

func (fooCodec) Unmarshal(data []byte, v interface{}) error {
	str := string(data)
	switch t := v.(type) {
	case *Request:
		t.In = str
	case *Response:
		t.Out = str
	default:
		panic("invalid type")
	}
	return nil
}

func (fooCodec) String() string {
	return "FooCodec"
}
