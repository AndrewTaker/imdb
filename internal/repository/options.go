package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginationOptions struct {
	Limit, Offset int
}

type SortOptions struct {
	SortBy SortableVar
	Order  OrderVar
}

type FilterOptions struct {
	FilterBy FilterableVar
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
		filterBy = append(filterBy, bson.E{Key: string(f.FilterBy), Value: f.Value})
	}

	return filterBy
}

func CreateQueryOptions(pag PaginationOptions, sort []SortOptions) *options.FindOptions {
	opts := options.Find()
	opts.SetLimit(int64(pag.Limit))
	opts.SetSkip(int64(pag.Offset))

	sortBy := bson.D{}
	for _, s := range sort {
		sortBy = append(sortBy, bson.E{Key: string(s.SortBy), Value: s.Order})
	}
	opts.SetSort(sortBy)

	return opts
}
