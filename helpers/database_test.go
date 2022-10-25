package helpers_test

import (
	"testing"

	"github.com/go-catupiry/catu/helpers"
	"github.com/stretchr/testify/assert"
)

func TestParseUrlQueryOrder(t *testing.T) {
	type args struct {
		order         string
		sort          string
		sortDirection string
	}

	type result struct {
		field   string
		isDesc  bool
		isValid bool
	}

	tests := []struct {
		name   string
		args   args
		result result
	}{
		{
			name: "success order desc",
			args: args{
				order: "createdAt DESC",
			},
			result: result{
				field:   "createdAt",
				isDesc:  true,
				isValid: true,
			},
		},
		{
			name: "success order asc",
			args: args{
				order: "createdAt ASC",
			},
			result: result{
				field:   "createdAt",
				isDesc:  false,
				isValid: true,
			},
		},
		{
			name: "invalid order without sort",
			args: args{
				order: "createdAt",
			},
			result: result{
				field:   "",
				isDesc:  true,
				isValid: false,
			},
		},
		{
			name: "success with sort and sortDirection",
			args: args{
				sort:          "createdAt",
				sortDirection: "DESC",
			},
			result: result{
				field:   "createdAt",
				isDesc:  true,
				isValid: true,
			},
		},
		{
			name: "success with sort only",
			args: args{
				sort: "publishedAt",
			},
			result: result{
				field:   "publishedAt",
				isDesc:  true,
				isValid: true,
			},
		},
		{
			name: "success with sort and sortDirection ASC",
			args: args{
				sort:          "id",
				sortDirection: "ASC",
			},
			result: result{
				field:   "id",
				isDesc:  false,
				isValid: true,
			},
		},
		{
			name: "invalid with no params",
			args: args{},
			result: result{
				field:   "",
				isDesc:  true,
				isValid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, isDesc, isValid := helpers.ParseUrlQueryOrder(tt.args.order, tt.args.sort, tt.args.sortDirection)
			assert.Equal(t, tt.result, result{field, isDesc, isValid})
		})
	}

}
