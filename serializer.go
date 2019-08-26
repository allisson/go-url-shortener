package shortener

// RedirectSerializer represents a contract to convert bytes to pointer of Redirect and vice-versa
type RedirectSerializer interface {
	Encode(input []byte) (*Redirect, error)
	Decode(input *Redirect) ([]byte, error)
}
