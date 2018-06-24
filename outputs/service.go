package outputs

import (
	"context"
	"fmt"

	"github.com/jianhan/ms-bmp-crawler/crawlers"
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
	psuppliers "github.com/jianhan/ms-bmp-products/proto/suppliers"
)

type service struct {
	categoriesServiceClient pcategories.CategoriesServiceClient
	productsServiceClient   pproducts.ProductsServiceClient
	suppliersServiceClient  psuppliers.SuppliersServiceClient
}

func NewService(cs pcategories.CategoriesServiceClient, ps pproducts.ProductsServiceClient, ss psuppliers.SuppliersServiceClient) OutputWriter {
	return &service{categoriesServiceClient: cs, productsServiceClient: ps, suppliersServiceClient: ss}
}

func (o *service) Output(ctx context.Context, crawler crawlers.Crawler) error {
	// get categories from crawler
	categories := crawler.Categories()
	if len(categories) == 0 {
		return fmt.Errorf("empty categories when output to service with %s", crawler.Name())
	}

	// get products from crawler
	products := crawler.Products()
	if len(products) == 0 {
		return fmt.Errorf("empty products when output to service with %s", crawler.Name())
	}

	// start syncing
	if err := o.syncSupplier(ctx, crawler.Supplier()); err != nil {
		return err
	}
	if err := o.syncCategories(ctx, categories); err != nil {
		return err
	}

	if err := o.syncProducts(ctx, products); err != nil {
		return err
	}

	return nil
}

func (o *service) syncCategories(ctx context.Context, categories []*pcategories.Category) error {
	// get existing categories
	rsp, err := o.categoriesServiceClient.Categories(ctx, &pcategories.CategoriesReq{})
	if err != nil {
		return err
	}

	// assign IDs
	for _, v := range rsp.Categories {
		for k := range categories {
			if categories[k].Url == v.Url || categories[k].Name == v.Name {
				categories[k].ID = v.ID
			}
		}
	}

	// update categories via RPC
	_, err = o.categoriesServiceClient.UpsertCategories(ctx, &pcategories.UpsertCategoriesReq{Categories: categories})
	if err != nil {
		return err
	}

	return nil
}

func (o *service) syncProducts(ctx context.Context, products []*pproducts.Product) error {
	// get existing
	rsp, err := o.productsServiceClient.Products(ctx, &pproducts.ProductsReq{})
	if err != nil {
		return err
	}

	// assign IDs
	for _, v := range rsp.Products {
		for k := range products {
			if products[k].Url == v.Url || products[k].Name == v.Name {
				products[k].ID = v.ID
			}
		}
	}

	filterdProducts := removeDuplicatesByName(products)
	// update via RPC
	_, err = o.productsServiceClient.UpsertProducts(ctx, &pproducts.UpsertProductsReq{Products: filterdProducts})
	if err != nil {
		return err
	}

	return nil
}

func (o *service) syncSupplier(ctx context.Context, supplier *psuppliers.Supplier) error {
	// get existing
	rsp, err := o.suppliersServiceClient.Suppliers(ctx, &psuppliers.SuppliersReq{})
	if err != nil {
		return err
	}

	// assign IDs
	for k := range rsp.Suppliers {
		if supplier.HomePageUrl == rsp.Suppliers[k].HomePageUrl || supplier.Name == rsp.Suppliers[k].Name {
			supplier.ID = rsp.Suppliers[k].ID
		}
	}

	// update via RPC
	_, err = o.suppliersServiceClient.UpsertSuppliers(ctx, &psuppliers.UpsertSuppliersReq{Suppliers: []*psuppliers.Supplier{supplier}})
	if err != nil {
		return err
	}

	return nil
}

func removeDuplicatesByName(products []*pproducts.Product) []*pproducts.Product {
	encountered := map[string]bool{}
	result := []*pproducts.Product{}

	for v := range products {
		if encountered[products[v].Name] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[products[v].Name] = true
			// Append to result slice.
			result = append(result, products[v])
		}
	}
	// Return the new slice.
	return result
}
