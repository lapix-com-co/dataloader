package pagination

// Input represent the pagination options to retrieve data.
type Input struct {
	First  *uint16
	Last   *uint16
	Before *string
	After  *string
}

// Output represent the PageInfo.
type Output struct {
	Total           uint32
	HasNextPage     bool
	HasPreviousPage bool
	StartCursor     *string
	EndCursor       *string
}
