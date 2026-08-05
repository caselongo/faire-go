package faire_go

import (
	"encoding/base64"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"io"
	"net/http"
	"net/url"
	"time"
)

type SaleState string

const (
	SaleStateForSale    SaleState = "FOR_SALE"
	SaleStateSalePaused SaleState = "SALE_PAUSED"
)

type LifecycleState string

const (
	LifecycleStateDraft       LifecycleState = "DRAFT"
	LifecycleStatePublished   LifecycleState = "PUBLISHED"
	LifecycleStateUnpublished LifecycleState = "UNPUBLISHED"
	LifecycleStateDeleted     LifecycleState = "DELETED"
)

type OrderabilityType string

const (
	OrderabilityTypeImmediate OrderabilityType = "IMMEDIATE"
	OrderabilityTypePreorder  OrderabilityType = "PREORDER"
)

type Product struct {
	Id                           string                     `json:"id,omitempty"`
	CreatedAt                    *time.Time                 `json:"created_at,omitempty"`
	UpdatedAt                    *time.Time                 `json:"updated_at,omitempty"`
	BrandId                      string                     `json:"brand_id,omitempty"`
	Name                         string                     `json:"name,omitempty"`
	Description                  string                     `json:"description,omitempty"`
	ShortDescription             string                     `json:"short_description,omitempty"`
	LifecycleState               LifecycleState             `json:"lifecycle_state,omitempty"`
	UnitMultiplier               *int                       `json:"unit_multiplier,omitempty"`
	MinimumOrderQuantity         *int                       `json:"minimum_order_quantity,omitempty"`
	PerStyleMinimumOrderQuantity *int                       `json:"per_style_minimum_order_quantity,omitempty"`
	AllowSalesWhenOutOfStock     *bool                      `json:"allow_sales_when_out_of_stock,omitempty"`
	TaxonomyType                 *TaxonomyType              `json:"taxonomy_type,omitempty"`
	Preorderable                 *bool                      `json:"preorderable,omitempty"`
	MadeInCountry                string                     `json:"made_in_country,omitempty"`
	Variants                     []*ProductVariant          `json:"variants,omitempty"`
	Images                       []*ProductImage            `json:"images,omitempty"`
	VariantOptionSets            []*ProductVariantOptionSet `json:"variant_option_sets,omitempty"`
	ProductAttributes            []*ProductAttribute        `json:"product_attributes,omitempty"`
	IdempotenceToken             string                     `json:"idempotence_token,omitempty"`
}

type ProductVariant struct {
	Id                string               `json:"id,omitempty"`
	CreatedAt         *time.Time           `json:"created_at,omitempty"`
	UpdatedAt         *time.Time           `json:"updated_at,omitempty"`
	ProductId         string               `json:"product_id,omitempty"`
	Name              string               `json:"name,omitempty"`
	SaleState         SaleState            `json:"sale_state,omitempty"`
	LifecycleState    LifecycleState       `json:"lifecycle_state,omitempty"`
	Sku               string               `json:"sku,omitempty"`
	AvailableQuantity *int                 `json:"available_quantity,omitempty"`
	TariffCode        string               `json:"tariff_code,omitempty"`
	Measurements      *ProductMeasurements `json:"measurements,omitempty"`
	Gtin              string               `json:"gtin,omitempty"`
	OrderabilityType  OrderabilityType     `json:"orderability_type,omitempty"`
	CaseMeasurements  *struct {
		MassUnit     string   `json:"mass_unit,omitempty"`
		Weight       *float64 `json:"weight,omitempty"`
		DistanceUnit string   `json:"distance_unit,omitempty"`
		Length       *float64 `json:"length,omitempty"`
		Width        *float64 `json:"width,omitempty"`
		Height       *float64 `json:"height,omitempty"`
	} `json:"case_measurements,omitempty"`
	Images           []*ProductImage `json:"images,omitempty"`
	Options          []*Option       `json:"options,omitempty"`
	Prices           []*ProductPrice `json:"prices,omitempty"`
	IdempotenceToken string          `json:"idempotence_token,omitempty"`
}

type ProductMeasurements struct {
	MassUnit     string   `json:"mass_unit,omitempty"`
	Weight       *float64 `json:"weight,omitempty"`
	DistanceUnit string   `json:"distance_unit,omitempty"`
	Length       *float64 `json:"length,omitempty"`
	Width        *float64 `json:"width,omitempty"`
	Height       *float64 `json:"height,omitempty"`
}

type ProductPrice struct {
	GeoConstraint  *GeoConstraint `json:"geo_constraint,omitempty"`
	WholesalePrice *Amount        `json:"wholesale_price,omitempty"`
	RetailPrice    *Amount        `json:"retail_price,omitempty"`
}

type ProductVariantOptionSet struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type ProductImage struct {
	Id          string   `json:"id,omitempty"`
	Width       *int     `json:"width,omitempty"`
	Height      *int     `json:"height,omitempty"`
	Sequence    *int     `json:"sequence,omitempty"`
	Url         string   `json:"url"`
	OriginalUrl string   `json:"original_url,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TaxonomyType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetProductsResponse struct {
	Products     []*Product `json:"products"`
	Page         int        `json:"page"`
	Limit        int        `json:"limit"`
	UpdatedAtMin *time.Time `json:"updated_at_min"`
	SortBy       string     `json:"sort_by"`
	Cursor       string     `json:"cursor"`
}

type GetProductsOptions struct {
	Cursor         string
	IncludeDeleted string
	Limit          *int
	Page           *int
	Sku            string
	UpdatedAtMin   *time.Time
}

func (service *Service) GetProducts(options *GetProductsOptions) ([]*Product, *errortools.Error) {
	var products []*Product

	values := url.Values{}
	if options != nil {
		if options.Cursor != "" {
			values.Add("cursor", options.Cursor)
		}
		if options.IncludeDeleted != "" {
			values.Add("include_deleted", options.IncludeDeleted)
		}
		if options.Limit != nil {
			values.Add("limit", fmt.Sprint(*options.Limit))
		}
		if options.Page != nil {
			values.Add("page", fmt.Sprint(*options.Page))
		}
		if options.Sku != "" {
			values.Add("sku", options.Sku)
		}
		if options.UpdatedAtMin != nil {
			values.Add("updated_at_min", options.UpdatedAtMin.Format(time.RFC3339))
		}
	}

	for {
		var getProductsResponse GetProductsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("/products?%s", values.Encode())),
			ResponseModel: &getProductsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			//fmt.Println(*service.errorResponse)
			return nil, e
		}

		if getProductsResponse.Products == nil || len(getProductsResponse.Products) == 0 {
			break
		}

		products = append(products, getProductsResponse.Products...)

		if getProductsResponse.Cursor == "" {
			break
		}
		values.Set("cursor", getProductsResponse.Cursor)
	}

	return products, nil
}

func (service *Service) GetProduct(productId string) (*Product, *errortools.Error) {
	var product Product

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("products/%s", productId)),
		ResponseModel: &product,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		//fmt.Println(*service.errorResponse)
		return nil, e
	}

	return &product, nil
}

func (service *Service) CreateProduct(product *Product) (*Product, *errortools.Error) {
	if product == nil {
		return nil, errortools.ErrorMessage("product must not be nil")
	}

	var productCreated Product

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url("products"),
		BodyModel:     product,
		ResponseModel: &productCreated,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		//fmt.Println(*service.errorResponse)
		return nil, e
	}

	return &productCreated, nil
}

type UploadImageResponse struct {
	Url string `json:"url"`
}

func (service *Service) UploadImage(url string) (string, *errortools.Error) {
	// Download image
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", errortools.ErrorMessage(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errortools.ErrorMessagef("unexpected status: %s", resp.Status)
	}

	// Read image bytes
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errortools.ErrorMessage(err)
	}

	// Convert to Base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Create API request
	uploadImageResponse := UploadImageResponse{}
	requestConfig := go_http.RequestConfig{
		Method: http.MethodPost,
		Url:    service.url("products/upload-image"),
		BodyModel: struct {
			Attachment string `json:"attachment"`
		}{
			Attachment: imageBase64,
		},
		ResponseModel: &uploadImageResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		//fmt.Println(*service.errorResponse)
		return "", e
	}

	return uploadImageResponse.Url, nil
}
