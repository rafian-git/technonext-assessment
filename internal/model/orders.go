package model

import (
	"time"
)

type Order struct {
	tableName          struct{}  `pg:"orders"`
	ID                 int64     `pg:"id,pk"`
	ConsignmentID      string    `pg:"consignment_id,unique,notnull"`
	MerchantOrderID    string    `pg:"merchant_order_id"`
	RecipientName      string    `pg:"recipient_name"`
	RecipientPhone     string    `pg:"recipient_phone"`
	RecipientAddress   string    `pg:"recipient_address"`
	RecipientCity      int32     `pg:"recipient_city"`
	RecipientZone      int32     `pg:"recipient_zone"`
	RecipientArea      int32     `pg:"recipient_area"`
	DeliveryType       int32     `pg:"delivery_type"`
	ItemType           int32     `pg:"item_type"`
	SpecialInstruction string    `pg:"special_instruction"`
	ItemQuantity       int32     `pg:"item_quantity"`
	ItemWeight         float64   `pg:"item_weight"`
	AmountToCollect    float64   `pg:"amount_to_collect"`
	ItemDescription    string    `pg:"item_description"`
	OrderStatus        string    `pg:"order_status"` // Pending, Cancelled
	DeliveryFee        float64   `pg:"delivery_fee"`
	CodFee             float64   `pg:"cod_fee"`
	PromoDiscount      float64   `pg:"promo_discount"`
	Discount           float64   `pg:"discount"`
	TotalFee           float64   `pg:"total_fee"`
	OrderType          string    `pg:"order_type"`    // Delivery
	ItemTypeStr        string    `pg:"item_type_str"` // Parcel
	OrderTypeID        int32     `pg:"order_type_id"`
	CreatedAt          time.Time `pg:"created_at,default:now()"`
}
