package repository

// Pagination is used as pagination when listing data
// Why using lastID instead of offset? Because in mysql, using limit X, Y will read
// X + Y rows, instead of only Y rows. Mysql reads all X + Y rows, then drop rows
// before X. Read this for more explanation:
// https://stackoverflow.com/questions/3799193/mysql-data-best-way-to-implement-paging#comment4027585_3799223
type Pagination struct {
	LastID int
	Limit  int
}
