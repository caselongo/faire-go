package faire_go

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"time"
)

type Product struct {
	Id                           string                     `json:"id,omitempty"`
	CreatedAt                    *time.Time                 `json:"created_at,omitempty"`
	UpdatedAt                    *time.Time                 `json:"updated_at,omitempty"`
	BrandId                      string                     `json:"brand_id,omitempty"`
	Name                         string                     `json:"name,omitempty"`
	Description                  string                     `json:"description,omitempty"`
	ShortDescription             string                     `json:"short_description,omitempty"`
	LifecycleState               string                     `json:"lifecycle_state,omitempty"`
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
	Id                string     `json:"id,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	ProductId         string     `json:"product_id,omitempty"`
	Name              string     `json:"name,omitempty"`
	SaleState         string     `json:"sale_state,omitempty"`
	LifecycleState    string     `json:"lifecycle_state,omitempty"`
	Sku               string     `json:"sku,omitempty"`
	AvailableQuantity *int       `json:"available_quantity,omitempty"`
	TariffCode        string     `json:"tariff_code,omitempty"`
	Measurements      *struct {
		MassUnit     string  `json:"mass_unit"`
		Weight       float64 `json:"weight"`
		DistanceUnit string  `json:"distance_unit"`
		Length       float64 `json:"length"`
		Width        float64 `json:"width"`
		Height       float64 `json:"height"`
	} `json:"measurements,omitempty"`
	Gtin             string `json:"gtin,omitempty"`
	OrderabilityType string `json:"orderability_type,omitempty"`
	CaseMeasurements *struct {
		MassUnit     string  `json:"mass_unit"`
		Weight       float64 `json:"weight"`
		DistanceUnit string  `json:"distance_unit"`
		Length       float64 `json:"length"`
		Width        float64 `json:"width"`
		Height       float64 `json:"height"`
	} `json:"case_measurements,omitempty"`
	Images  []*ProductImage `json:"images,omitempty"`
	Options []*struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"options,omitempty"`
	Prices []*struct {
		GeoConstraint struct {
			CountryGroup string `json:"country_group"`
		} `json:"geo_constraint"`
		WholesalePrice struct {
			AmountMinor int    `json:"amount_minor"`
			Currency    string `json:"currency"`
		} `json:"wholesale_price"`
		RetailPrice struct {
			AmountMinor int    `json:"amount_minor"`
			Currency    string `json:"currency"`
		} `json:"retail_price"`
	} `json:"prices,omitempty"`
	IdempotenceToken string `json:"idempotence_token,omitempty"`
}

type ProductVariantOptionSet struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type ProductImage struct {
	Id          string   `json:"id"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Sequence    int      `json:"sequence"`
	Url         string   `json:"url"`
	OriginalUrl string   `json:"original_url"`
	Tags        []string `json:"tags"`
}

type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TaxonomyType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
