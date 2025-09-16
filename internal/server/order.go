package server

import (
	"context"
	ordersv1 "gitlab.com/sample_projects/technonext-assessment/gen/orders/v1"
	"gitlab.com/sample_projects/technonext-assessment/internal/service"
)

type OrderServer struct {
	ordersv1.UnimplementedOrderServiceServer
	svc   *service.OrderService
	authz *Authz
}

func NewOrderServer(svc *service.OrderService, authz *Authz) *OrderServer {
	return &OrderServer{svc: svc, authz: authz}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *ordersv1.CreateOrderRequest) (*ordersv1.CreateOrderResponse, error) {
	// auth check
	if _, err := s.authz.Require(ctx); err != nil {
		return nil, err
	}
	res, err := s.svc.CreateOrder(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *OrderServer) ListOrders(ctx context.Context, req *ordersv1.ListOrdersRequest) (*ordersv1.ListOrdersResponse, error) {
	if _, err := s.authz.Require(ctx); err != nil {
		return nil, err
	}

	res, err := s.svc.ListOrders(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *ordersv1.CancelOrderRequest) (*ordersv1.GenericResponse, error) {

	if _, err := s.authz.Require(ctx); err != nil {
		return nil, err
	}

	res, err := s.svc.CancelOrder(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
