package shortener

// Redirect info
type Redirect struct {
	Code      string `bson:"code" json:"code" msgpack:"code"`
	URL       string `bson:"url" json:"url" msgpack:"url" validate:"format=url"`
	CreatedAt int64  `bson:"created_at" msgpack:"created_at" json:"created_at"`
}
