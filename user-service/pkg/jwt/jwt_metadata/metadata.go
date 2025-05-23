package jwt_metadata

type SerializebleBase64 interface {
	Header | Payload
}
