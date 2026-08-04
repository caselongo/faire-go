package faire_go

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"time"
)

type Order struct {
	Id        string     `json:"id"`
	DisplayId string     `json:"display_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	State     string     `json:"state"`
	Items     []*struct {
		Id               string     `json:"id"`
		CreatedAt        *time.Time `json:"created_at"`
		UpdatedAt        *time.Time `json:"updated_at"`
		OrderId          string     `json:"order_id"`
		ProductId        string     `json:"product_id"`
		VariantId        string     `json:"variant_id"`
		Quantity         *int       `json:"quantity"`
		Sku              string     `json:"sku"`
		PriceCents       *int       `json:"price_cents"`
		ProductName      string     `json:"product_name"`
		VariantName      string     `json:"variant_name"`
		IncludesTester   *bool      `json:"includes_tester"`
		TesterPriceCents *int       `json:"tester_price_cents"`
		Customizations   []*struct {
			Token string `json:"token"`
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"customizations"`
		Price       *Amount `json:"price"`
		TesterPrice *Amount `json:"tester_price"`
		Discounts   []*struct {
			Id                   string  `json:"id"`
			Code                 string  `json:"code"`
			DiscountType         string  `json:"discount_type"`
			DiscountAmountCents  *int    `json:"discount_amount_cents"`
			DiscountPercentage   *int    `json:"discount_percentage"`
			IncludesFreeShipping *bool   `json:"includes_free_shipping"`
			DiscountAmount       *Amount `json:"discount_amount"`
		} `json:"discounts"`
		State string `json:"state"`
	} `json:"items"`
	Shipments []*struct {
		Id               string     `json:"id"`
		CreatedAt        *time.Time `json:"created_at"`
		UpdatedAt        *time.Time `json:"updated_at"`
		OrderId          string     `json:"order_id"`
		MakerCostCents   *int       `json:"maker_cost_cents"`
		Carrier          string     `json:"carrier"`
		TrackingCode     string     `json:"tracking_code"`
		MakerCost        *Amount    `json:"maker_cost"`
		ShippingType     string     `json:"shipping_type"`
		ShippingLabelUrl string     `json:"shipping_label_url"`
	} `json:"shipments"`
	Address *struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Address1    string `json:"address1"`
		Address2    string `json:"address2"`
		PostalCode  string `json:"postal_code"`
		City        string `json:"city"`
		State       string `json:"state"`
		StateCode   string `json:"state_code"`
		PhoneNumber string `json:"phone_number"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		CompanyName string `json:"company_name"`
		AddressType string `json:"address_type"`
	} `json:"address"`
	ShipAfter   *time.Time `json:"ship_after"`
	PayoutCosts *struct {
		PayoutFeeCents         *int    `json:"payout_fee_cents"`
		PayoutFeeBps           *int    `json:"payout_fee_bps"`
		PayoutFlatFee          *Amount `json:"payout_flat_fee"`
		CommissionCents        *int    `json:"commission_cents"`
		CommissionBps          *int    `json:"commission_bps"`
		CommissionFlatFee      *Amount `json:"commission_flat_fee"`
		PayoutFee              *Amount `json:"payout_fee"`
		Commission             *Amount `json:"commission"`
		TotalPayout            *Amount `json:"total_payout"`
		PayoutProtectionFee    *Amount `json:"payout_protection_fee"`
		DamagedAndMissingItems *Amount `json:"damaged_and_missing_items"`
		NetTax                 *Amount `json:"net_tax"`
		ShippingSubsidy        *Amount `json:"shipping_subsidy"`
		Taxes                  []*struct {
			Value           *Amount `json:"value"`
			TaxableItemType string  `json:"taxable_item_type"`
			TaxType         string  `json:"tax_type"`
			Effect          string  `json:"effect"`
		} `json:"taxes"`
		SubtotalAfterBrandDiscounts *Amount `json:"subtotal_after_brand_discounts"`
		TotalBrandDiscounts         *Amount `json:"total_brand_discounts"`
	} `json:"payout_costs"`
	PaymentInitiatedAt *time.Time `json:"payment_initiated_at"`
	OriginalOrderId    string     `json:"original_order_id"`
	RetailerId         string     `json:"retailer_id"`
	Source             string     `json:"source"`
	ExpectedShipDate   string     `json:"expected_ship_date"`
	Customer           *struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"customer"`
	BrandDiscounts []*struct {
		Id                   string  `json:"id"`
		Code                 string  `json:"code"`
		DiscountType         string  `json:"discount_type"`
		DiscountAmountCents  *int    `json:"discount_amount_cents"`
		DiscountPercentage   *int    `json:"discount_percentage"`
		IncludesFreeShipping *bool   `json:"includes_free_shipping"`
		DiscountAmount       *Amount `json:"discount_amount"`
	} `json:"brand_discounts"`
	RequestedShipDate                     string     `json:"requested_ship_date"`
	ProcessingAt                          string     `json:"processing_at"`
	IsFreeShipping                        *bool      `json:"is_free_shipping"`
	FreeShippingReason                    string     `json:"free_shipping_reason"`
	FaireCoveredShippingCost              *Amount    `json:"faire_covered_shipping_cost"`
	EstimatedPayoutAt                     *time.Time `json:"estimated_payout_at"`
	IsFulfilledByFaire                    bool       `json:"is_fulfilled_by_faire"`
	PurchaseOrderNumber                   string     `json:"purchase_order_number"`
	Notes                                 string     `json:"notes"`
	HasPendingRetailerCancellationRequest bool       `json:"has_pending_retailer_cancellation_request"`
	SalesRepName                          string     `json:"sales_rep_name"`
}

type GetOrdersResponse struct {
	Orders       []*Order   `json:"orders"`
	Page         int        `json:"page"`
	Limit        int        `json:"limit"`
	UpdatedAtMin *time.Time `json:"updated_at_min"`
	SortBy       string     `json:"sort_by"`
	Cursor       string     `json:"cursor"`
}

type GetOrdersOptions struct {
	CreatedAtMin    *time.Time
	Cursor          string
	ExcludedStates  string
	Limit           *int
	OriginalOrderId string
	Page            *int
	ShipAfterMax    *time.Time
	SortBy          string
	UpdatedAtMin    *time.Time
}

func (service *Service) GetOrders(options *GetOrdersOptions) ([]*Order, *errortools.Error) {
	var orders []*Order

	values := url.Values{}
	if options != nil {
		if options.CreatedAtMin != nil {
			values.Add("created_at_min", options.CreatedAtMin.Format(time.RFC3339))
		}
		if options.Cursor != "" {
			values.Add("cursor", options.Cursor)
		}
		if options.ExcludedStates != "" {
			values.Add("excluded_states", options.ExcludedStates)
		}
		if options.Limit != nil {
			values.Add("limit", fmt.Sprint(*options.Limit))
		}
		if options.OriginalOrderId != "" {
			values.Add("original_order_id", options.OriginalOrderId)
		}
		if options.Page != nil {
			values.Add("page", fmt.Sprint(*options.Page))
		}
		if options.ShipAfterMax != nil {
			values.Add("ship_after_max", options.CreatedAtMin.Format(time.RFC3339))
		}
		if options.SortBy != "" {
			values.Add("sort_by", options.SortBy)
		}
		if options.UpdatedAtMin != nil {
			values.Add("updated_at_min", options.UpdatedAtMin.Format(time.RFC3339))
		}
	}

	for {
		var getOrdersResponse GetOrdersResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("/orders?%s", values.Encode())),
			ResponseModel: &getOrdersResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			//fmt.Println(*service.errorResponse)
			return nil, e
		}

		if getOrdersResponse.Orders == nil || len(getOrdersResponse.Orders) == 0 {
			break
		}

		orders = append(orders, getOrdersResponse.Orders...)

		if getOrdersResponse.Cursor == "" {
			break
		}
		values.Set("cursor", getOrdersResponse.Cursor)
	}

	return orders, nil
}

func (service *Service) GetOrder(orderId string) (*Order, *errortools.Error) {
	var order Order

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("orders/%s", orderId)),
		ResponseModel: &order,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		//fmt.Println(*service.errorResponse)
		return nil, e
	}

	return &order, nil
}
