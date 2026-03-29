package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginationOptions struct {
	Limit, Offset int
}

type SortOptions struct {
	SortBy string
	Order  OrderVar
}

type FilterOptions struct {
	FilterBy string
	Value    any
}

type SortableVar string

type OrderVar int8

type FilterableVar string

var (
	ASC  OrderVar = 1
	DESC OrderVar = -1
)

func CreateFilterOptions(filters []FilterOptions) bson.D {
	filterBy := bson.D{}
	for _, f := range filters {
		if f.FilterBy != "" {
			filterBy = append(filterBy, bson.E{Key: f.FilterBy, Value: f.Value})
		}
	}

	return filterBy
}

func CreateQueryOptions(pag PaginationOptions, sort []SortOptions, filter []FilterOptions) *options.FindOptions {
	opts := options.Find()
	opts.SetLimit(int64(pag.Limit))
	opts.SetSkip(int64(pag.Offset))

	sortBy := bson.D{}
	for _, s := range sort {
		sortBy = append(sortBy, bson.E{Key: s.SortBy, Value: s.Order})
	}
	opts.SetSort(sortBy)

	for _, f := range filter {
		if f.FilterBy == "genres" {
			opts.SetCollation(&options.Collation{
				Locale:   "en",
				Strength: 2,
			})
			break
		}
	}

	return opts
}
