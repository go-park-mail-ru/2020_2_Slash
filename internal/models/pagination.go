package models

type Pagination struct {
	From  uint64 `query:"from"`
	Count uint64 `query:"count"`
}
