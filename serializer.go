package shortener

// RedirectSerializer represents a contract to convert bytes to pointer of Redirect and vice-versa
type RedirectSerializer interface {
	Decode(input []byte) (*Redirect, error)
	Encode(input *Redirect) ([]byte, error)
}
