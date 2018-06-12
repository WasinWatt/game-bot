package mongo

// New creates new mongo Repository
func New() *Repository {
	return &Repository{}
}

// Repository contains all functions connecting mongo
type Repository struct{}
