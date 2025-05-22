package jwt_metadata

type SerializebleBase64 interface {
	Header | Claims
}

type Validator interface {
	Validate() error
}
