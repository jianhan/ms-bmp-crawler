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

	// get suppliers from crawler
	suppliers := crawler.Products()
	if len(suppliers) == 0 {
		return fmt.Errorf("empty suppliers when output to service with %s", crawler.Name())
	}

	// start syncing
	if err := o.syncCategories(ctx, categories); err != nil {
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
			if categories[k].Url == v.Url {
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
			if products[k].Url == v.Url {
				products[k].ID = v.ID
			}
		}
	}

	// update via RPC
	_, err = o.productsServiceClient.UpsertProducts(ctx, &pproducts.UpsertProductsReq{Products: products})
	if err != nil {
		return err
	}

	return nil
}

func (o *service) syncSupplier(ctx context.Context, suppliers []*psuppliers.Supplier) error {
	// get existing
	rsp, err := o.suppliersServiceClient.Suppliers(ctx, &psuppliers.SuppliersReq{})
	if err != nil {
		return err
	}

	// assign IDs
	for _, v := range rsp.Suppliers {
		for k := range suppliers {
			if suppliers[k].HomePageUrl == v.HomePageUrl {
				suppliers[k].ID = v.ID
			}
		}
	}

	// update via RPC
	_, err = o.suppliersServiceClient.UpsertSuppliers(ctx, &psuppliers.UpsertSuppliersReq{Suppliers: suppliers})
	if err != nil {
		return err
	}

	return nil
}
