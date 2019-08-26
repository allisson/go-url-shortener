package shortener

// RedirectService ...
type RedirectService interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
