package manta

// FindOptions represents options passed to all find methods with multiple results.
type FindOptions struct {
	Limit      int
	Offset     int
	SortBy     string
	Descending bool
}
