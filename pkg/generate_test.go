package pkg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"
)

type test struct {
	name   string
	args   LoaderInput
	input  []byte
	output []byte
}

var tests = []*test{
	{
		name: "pets",
		args: LoaderInput{
			Type:           "Pet",
			CreatePageType: true,
			Expose:         true,
		},
		input:  []byte(testLoaderInput),
		output: []byte(testLoaderOutput),
	},
}

var testLoaderInput = `package generate

import "context"

type Pet struct {
	ID        uint
	Name      string
}

type petDataLoader interface {
	Create(context.Context, *Pet) error
	Update(context.Context, *Pet) error
	Delete(context.Context, *Pet) error
	Find(context.Context, []uint) ([]*Pet, error)
	FindOne(context.Context, uint) (*Pet, error)
	All(context.Context, dataloader.Input) (*PetPage, error)

	// Dataloader does not know about this method, it will be skiped.
	AnyOtherMethod(context.Context, []string) error
}
`

var testLoaderOutput = `// Code generated by github.com/lapix-com-co/dataloader; do not edit!

package generate

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/lapix-com-co/dataloader/pkg"
	"github.com/lapix-com-co/dataloader/pkg/slice"
	"github.com/lapix-com-co/dataloader/pkg/pagination"
	local "github.com/lapix-com-co/dataloader/pkg/gorm"
)

type PetPage struct {
	Nodes    []*Pet
	PageInfo pagination.Output
}

type petMySQLDataLoader struct {
	db        *gorm.DB
	tableName string
	sortKey   string
}

func NewPetMySQLDataLoader(db *gorm.DB) *petMySQLDataLoader {
	return &petMySQLDataLoader{
		db:        db,
		tableName: "pets",
		sortKey:   "id",
	}
}

func (r *petMySQLDataLoader) Create(ctx context.Context, pet *Pet) error {
	if err := local.Create(ctx, r.queryTable(), pet); err != nil {
		return fmt.Errorf("could not create the item: %w", err)
	}

	return nil
}

func (r *petMySQLDataLoader) Update(ctx context.Context, pet *Pet) error {
	if pet.ID == 0 {
		return errors.New("could not update the item because it does not have a valid ID")
	}

	if err := local.Save(ctx, r.queryTable(), pet); err != nil {
		return fmt.Errorf("could not update the item: %w", err)
	}

	return nil
}

func (r *petMySQLDataLoader) Delete(ctx context.Context, pet *Pet) error {
	if err := local.Delete(ctx, r.queryTable(), pet); err != nil {
		return fmt.Errorf("could not delete the item: %w", err)
	}

	return nil
}

func (r *petMySQLDataLoader) Find(ctx context.Context, i []uint) ([]*Pet, error) {
	elements := make([]*Pet, 0)
	if len(i) == 0 {
		return elements, nil
	}

	if err := local.Find(ctx, r.queryTable().Where("id IN (?)", i), &elements); err != nil {
		if errors.Is(err, pkg.ErrRecordNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("could not find the item: %w", err)
	}

	return elements, nil
}

func (r *petMySQLDataLoader) FindOne(ctx context.Context, i uint) (*Pet, error) {
	elements, err := r.Find(ctx, []uint{i})
	if err != nil {
		return nil, err
	}
	if len(elements) == 0 {
		return nil, fmt.Errorf("%w: '%d'", pkg.ErrNotFound, i)
	}
	return elements[0], nil
}

func (r *petMySQLDataLoader) All(ctx context.Context, i pagination.Input) (*PetPage, error) {
	return r.paginateElements(ctx, r.queryTable(), i)
}

func (r *petMySQLDataLoader) paginateElements(ctx context.Context, query *gorm.DB, i pagination.Input) (*PetPage, error) {
	elements := make([]*Pet, 0)
	total, err := local.Count(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not count the total items: %w", err)
	}

	if total == 0 {
		return createPetPageInfo(elements, 0, i), nil
	}

	if err = local.Find(ctx, r.applyPaginationToQuery(query, i), &elements); err != nil {
		return nil, fmt.Errorf("could not query the items: %w", err)
	}

	if i.Last != nil {
		slice.Reverse(elements)
	}

	return createPetPageInfo(elements, total, i), nil
}

func (r *petMySQLDataLoader) applyPaginationToQuery(query *gorm.DB, input pagination.Input) *gorm.DB {
	if input.After != nil {
		query = query.Where(fmt.Sprintf("%s > ?", r.sortKey), *input.After)
	} else if input.Before != nil {
		query = query.Where(fmt.Sprintf("%s < ?", r.sortKey), *input.Before)
	}

	if input.Last != nil {
		query = query.Order(fmt.Sprintf("%s DESC", r.sortKey)).Limit(*input.Last)
	} else if input.First != nil {
		query = query.Order(fmt.Sprintf("%s ASC", r.sortKey)).Limit(*input.First)
	}

	return query
}

func createPetPageInfo(items []*Pet, totalItems uint32, input pagination.Input) *PetPage {
	var page = &PetPage{Nodes: items}
	var pageLength = len(items)
	var pageInfo = pagination.Output{
		Total:           totalItems,
		HasNextPage:     true,
		HasPreviousPage: true,
	}

	if pageLength > 0 {
		startCursorID := strconv.Itoa(int(items[0].ID))
		endCursorID := strconv.Itoa(int(items[pageLength-1].ID))

		pageInfo.StartCursor = &startCursorID
		pageInfo.EndCursor = &endCursorID
	}

	if input.Last != nil {
		pageInfo.HasPreviousPage = uint16(pageLength) == *input.Last
		pageInfo.HasNextPage = false
	} else if input.First != nil {
		pageInfo.HasPreviousPage = false
		pageInfo.HasNextPage = uint16(pageLength) == *input.First
	}

	page.PageInfo = pageInfo

	return page
}

func (r *petMySQLDataLoader) queryTable() *gorm.DB {
	return r.db.New().Table(r.tableName)
}
`

func TestBuildLoader(t *testing.T) {
	folder, err := ioutil.TempDir("", "generate")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(folder)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.name + ".go"
			absPath := filepath.Join(folder, filename)
			if err := ioutil.WriteFile(absPath, tt.input, 0644); err != nil {
				t.Error(err)
			}

			tt.args.Pattern = []string{absPath}

			got, err := BuildLoader(tt.args)

			require.NoError(t, err)
			assertNoDiff(t, string(tt.output), string(got))
		})
	}
}

func assertNoDiff(t *testing.T, expected, current string) {
	t.Helper()

	diff := difflib.ContextDiff{
		A:        difflib.SplitLines(expected),
		FromFile: "Expected",
		B:        difflib.SplitLines(current),
		ToFile:   "Current",
		Eol:      "\n",
		Context:  0,
	}
	result, _ := difflib.GetContextDiffString(diff)
	require.Empty(t, result)
}
