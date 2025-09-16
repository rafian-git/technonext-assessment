package service

import (
	"fmt"
	ordersv1 "gitlab.com/sample_projects/technonext-assessment/gen/orders/v1"
	"gitlab.com/sample_projects/technonext-assessment/internal/model"
	"gitlab.com/sample_projects/technonext-assessment/internal/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"regexp"
	"strings"
	"time"
)

type OrderService struct{ repo *repo.OrderRepo }

func NewOrderService(or *repo.OrderRepo) *OrderService { return &OrderService{repo: or} }

// ---- helpers ----
var bdPhone = regexp.MustCompile(`^(01)[3-9]{1}[0-9]{8}$`)

type valErrors map[string][]string

func (ve valErrors) add(field, msg string) { ve[field] = append(ve[field], msg) }

func generateConsignmentID() string {
	now := time.Now().UTC()
	return fmt.Sprintf("DA%s%05d", now.Format("060102"), now.Nanosecond()%100000)
}

func (s *OrderService) CreateOrder(req *ordersv1.CreateOrderRequest) (*ordersv1.CreateOrderResponse, error) {
	// validations → collect field-wise errors to match spec
	ve := valErrors{}

	// Fixed fields (exact values required by the task specification)
	if req.GetStoreId() == 0 {
		ve.add("store_id", "The store field is required")
	}
	if req.GetStoreId() != 131172 {
		ve.add("store_id", "Wrong Store selected")
	}
	if req.GetDeliveryType() == 0 {
		ve.add("delivery_type", "The delivery type field is required.")
	}
	if req.GetDeliveryType() != 48 {
		ve.add("delivery_type", "Invalid delivery type.")
	}
	if req.GetItemType() == 0 {
		ve.add("item_type", "The item type field is required.")
	}
	if req.GetItemType() != 2 {
		ve.add("item_type", "Invalid item type.")
	}

	// Required inputs
	if strings.TrimSpace(req.GetRecipientName()) == "" {
		ve.add("recipient_name", "The recipient name field is required.")
	}
	if strings.TrimSpace(req.GetRecipientPhone()) == "" {
		ve.add("recipient_phone", "The recipient phone field is required.")
	} else if !bdPhone.MatchString(req.GetRecipientPhone()) {
		ve.add("recipient_phone", "Invalid Bangladeshi phone number.")
	}
	if strings.TrimSpace(req.GetRecipientAddress()) == "" {
		ve.add("recipient_address", "The recipient address field is required.")
	}
	if req.GetAmountToCollect() <= 0 {
		ve.add("amount_to_collect", "The amount to collect field is required.")
	}
	if req.GetItemQuantity() < 1 {
		ve.add("item_quantity", "The item quantity field is required.")
	}
	if req.GetItemWeight() <= 0 {
		ve.add("item_weight", "The item weight field is required.")
	}

	// Address tokens rule (exact phrase components required by spec)
	addr := strings.ToLower(req.GetRecipientAddress())
	if !(strings.Contains(addr, "banani") && strings.Contains(addr, "gulshan 2") && strings.Contains(addr, "dhaka") && strings.Contains(addr, "bangladesh")) {
		ve.add("recipient_address", "Address must include: banani, gulshan 2, dhaka, bangladesh")
	}

	// If any validation errors found, return
	if len(ve) > 0 {
		errs := map[string]*ordersv1.ErrorList{}
		for k, list := range ve {
			errs[k] = &ordersv1.ErrorList{Messages: list}
		}
		return &ordersv1.CreateOrderResponse{
			Message: "Please fix the given errors",
			Type:    "error",
			Code:    422,
			Errors:  errs,
		}, nil
	}

	// Fees calculation
	delivery := s.CalcDeliveryFee(req.GetRecipientCity(), req.GetItemWeight())
	cod := math.Round(req.GetAmountToCollect()*0.01*100) / 100.0
	total := delivery + cod

	consignmentID := generateConsignmentID()

	order := &model.Order{
		ConsignmentID:      consignmentID,
		MerchantOrderID:    req.GetMerchantOrderId(),
		RecipientName:      req.GetRecipientName(),
		RecipientPhone:     req.GetRecipientPhone(),
		RecipientAddress:   req.GetRecipientAddress(),
		RecipientCity:      req.GetRecipientCity(),
		RecipientZone:      req.GetRecipientZone(),
		RecipientArea:      req.GetRecipientArea(),
		DeliveryType:       req.GetDeliveryType(),
		ItemType:           req.GetItemType(),
		SpecialInstruction: req.GetSpecialInstruction(),
		ItemQuantity:       req.GetItemQuantity(),
		ItemWeight:         req.GetItemWeight(),
		AmountToCollect:    req.GetAmountToCollect(),
		ItemDescription:    req.GetItemDescription(),
		OrderStatus:        "Pending",
		DeliveryFee:        delivery,
		CodFee:             cod,
		PromoDiscount:      0,
		Discount:           0,
		TotalFee:           total,
		OrderType:          "Delivery",
		ItemTypeStr:        "Parcel",
		OrderTypeID:        1,
	}

	if err := s.Insert(order); err != nil {
		return nil, status.Errorf(codes.Internal, "db: %v", err)
	}

	return &ordersv1.CreateOrderResponse{
		Message: "Order Created Successfully",
		Type:    "success",
		Code:    200,
		Data: &ordersv1.OrderData{
			ConsignmentId:   consignmentID,
			MerchantOrderId: order.MerchantOrderID,
			OrderStatus:     order.OrderStatus,
			DeliveryFee:     order.DeliveryFee,
		},
	}, nil
}

func (s *OrderService) CalcDeliveryFee(city int32, weight float64) float64 {
	if city == 1 {
		if weight <= 0.5 {
			return 60
		}
		if weight <= 1.0 {
			return 70
		}
		return 70 + math.Ceil(weight-1.0)*15
	}
	if weight <= 1.0 {
		return 100
	}
	return 100 + math.Ceil(weight-1.0)*15
}

func (s *OrderService) ValidAddress(addr string) bool {
	a := strings.ToLower(addr)
	return strings.Contains(a, "banani") && strings.Contains(a, "gulshan 2") &&
		strings.Contains(a, "dhaka") && strings.Contains(a, "bangladesh")
}

func (s *OrderService) ListOrders(req *ordersv1.ListOrdersRequest) (*ordersv1.ListOrdersResponse, error) {
	if req.GetTransferStatus() != 1 || req.GetArchive() != 0 {
		return &ordersv1.ListOrdersResponse{
			Message: "Orders successfully fetched.",
			Type:    "success",
			Code:    200,
			Data: &ordersv1.PagedOrders{
				Data:        []*ordersv1.Order{},
				Total:       0,
				CurrentPage: 1,
				PerPage:     0,
				TotalInPage: 0,
				LastPage:    1,
			},
		}, nil
	}

	limit := req.GetLimit()
	if limit < 1 {
		limit = 10
	}
	page := req.GetPage()
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	orders, total, err := s.List(int(limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db: %v", err)
	}

	items := make([]*ordersv1.Order, 0, len(orders))
	for _, o := range orders {
		items = append(items, &ordersv1.Order{
			OrderConsignmentId: o.ConsignmentID,
			OrderCreatedAt:     o.CreatedAt.Format("2006-01-02 15:04:05"),
			OrderDescription:   o.ItemDescription,
			MerchantOrderId:    o.MerchantOrderID,
			RecipientName:      o.RecipientName,
			RecipientAddress:   o.RecipientAddress,
			RecipientPhone:     o.RecipientPhone,
			OrderAmount:        o.AmountToCollect,
			TotalFee:           o.TotalFee,
			Instruction:        o.SpecialInstruction,
			OrderTypeId:        o.OrderTypeID,
			CodFee:             o.CodFee,
			PromoDiscount:      o.PromoDiscount,
			Discount:           o.Discount,
			DeliveryFee:        o.DeliveryFee,
			OrderStatus:        o.OrderStatus,
			OrderType:          o.OrderType,
			ItemType:           o.ItemTypeStr,
		})
	}

	perPage := limit
	lastPage := int32((total + int(perPage) - 1) / int(perPage))

	return &ordersv1.ListOrdersResponse{
		Message: "Orders successfully fetched.",
		Type:    "success",
		Code:    200,
		Data: &ordersv1.PagedOrders{
			Data:        items,
			Total:       int32(total),
			CurrentPage: page,
			PerPage:     perPage,
			TotalInPage: int32(len(items)),
			LastPage:    lastPage,
		},
	}, nil
}

func (s *OrderService) CancelOrder(req *ordersv1.CancelOrderRequest) (*ordersv1.GenericResponse, error) {
	cons := strings.TrimSpace(req.GetConsignmentId())
	if cons == "" {
		return nil, status.Error(codes.InvalidArgument, `{"message":"Please contact cx to cancel order","type":"error","code":400}`)
	}

	o, err := s.FindByConsignmentID(cons)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, `{"message":"Please contact cx to cancel order","type":"error","code":400}`)
	}

	if o.OrderStatus != "Cancelled" {
		o.OrderStatus = "Cancelled"
		if err := s.Update(o); err != nil {
			return nil, status.Errorf(codes.Internal, "db: %v", err)
		}
	}

	return &ordersv1.GenericResponse{
		Message: "Order Cancelled Successfully",
		Type:    "success",
		Code:    200,
	}, nil
}

func (s *OrderService) Insert(o *model.Order) error { return s.repo.Insert(o) }
func (s *OrderService) FindByConsignmentID(id string) (*model.Order, error) {
	return s.repo.FindByConsignmentID(id)
}
func (s *OrderService) Update(o *model.Order) error { return s.repo.Update(o) }
func (s *OrderService) List(limit, offset int) ([]model.Order, int, error) {
	return s.repo.List(limit, offset)
}
