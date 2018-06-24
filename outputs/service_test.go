package outputs

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
)

func TestRemoveDuplicatesByName(t *testing.T) {
	products := []*pproducts.Product{
		{
			ID:   "test1",
			Name: "test1",
		},
		{
			ID:   "test2",
			Name: "test2",
		},
		{
			ID:   "test3",
			Name: "test3",
		},
		{
			ID:   "test1",
			Name: "test1",
		},
	}
	spew.Dump(removeDuplicatesByName(products))

}
