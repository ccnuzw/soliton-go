package shippingapp

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/soliton-go/application/internal/domain/shipping"
)

// ShippingService 处理跨领域的业务逻辑编排。
type ShippingService struct {
	repo shipping.ShippingRepository
}

// NewShippingService 创建 ShippingService 实例。
func NewShippingService(
	repo shipping.ShippingRepository,
) *ShippingService {
	return &ShippingService{
		repo: repo,
	}
}

// CreateShipment 实现 CreateShipment 用例。
func (s *ShippingService) CreateShipment(ctx context.Context, req CreateShipmentServiceRequest) (*CreateShipmentServiceResponse, error) {
	if req.OrderId == "" {
		return nil, fmt.Errorf("order_id is required")
	}
	if req.Carrier == "" {
		return nil, fmt.Errorf("carrier is required")
	}
	if req.ShippingMethod == "" {
		return nil, fmt.Errorf("shipping_method is required")
	}
	if req.ReceiverName == "" || req.ReceiverPhone == "" || req.ReceiverAddress == "" {
		return nil, fmt.Errorf("receiver_name, receiver_phone, receiver_address are required")
	}

	entity := shipping.NewShipping(
		uuid.New().String(),
		req.OrderId,
		req.Carrier,
		shipping.ShippingShippingMethod(req.ShippingMethod),
		"",
		shipping.ShippingStatusPending,
		nil,
		nil,
		req.ReceiverName,
		req.ReceiverPhone,
		req.ReceiverAddress,
		req.ReceiverCity,
		req.ReceiverState,
		req.ReceiverCountry,
		req.ReceiverPostalCode,
		req.Notes,
	)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &CreateShipmentServiceResponse{
		Success:        true,
		Message:        "created",
		ShipmentId:     string(entity.ID),
		TrackingNumber: entity.TrackingNumber,
		Status:         string(entity.Status),
		ShippedAt:      entity.ShippedAt,
	}, nil
}

// UpdateTracking 实现 UpdateTracking 用例。
func (s *ShippingService) UpdateTracking(ctx context.Context, req UpdateTrackingServiceRequest) (*UpdateTrackingServiceResponse, error) {
	if req.ShipmentId == "" {
		return nil, fmt.Errorf("shipment_id is required")
	}
	if req.TrackingNumber == "" {
		return nil, fmt.Errorf("tracking_number is required")
	}

	entity, err := s.repo.Find(ctx, shipping.ShippingID(req.ShipmentId))
	if err != nil {
		return nil, err
	}
	if entity.Status == shipping.ShippingStatusDelivered || entity.Status == shipping.ShippingStatusCancelled {
		return nil, fmt.Errorf("shipment cannot be updated in current status")
	}

	var status *shipping.ShippingStatus
	if req.Status != "" {
		parsed := shipping.ShippingStatus(req.Status)
		status = &parsed
	}
	var shippedAt *time.Time
	if status != nil && (*status == shipping.ShippingStatusLabelCreated || *status == shipping.ShippingStatusInTransit) {
		if req.UpdatedAt != nil {
			shippedAt = req.UpdatedAt
		} else {
			now := time.Now()
			shippedAt = &now
		}
	}

	entity.Update(nil, nil, nil, &req.TrackingNumber, status, shippedAt, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &UpdateTrackingServiceResponse{
		Success:        true,
		Message:        "updated",
		ShipmentId:     string(entity.ID),
		TrackingNumber: entity.TrackingNumber,
		Status:         string(entity.Status),
	}, nil
}

// MarkDelivered 实现 MarkDelivered 用例。
func (s *ShippingService) MarkDelivered(ctx context.Context, req MarkDeliveredServiceRequest) (*MarkDeliveredServiceResponse, error) {
	if req.ShipmentId == "" {
		return nil, fmt.Errorf("shipment_id is required")
	}

	entity, err := s.repo.Find(ctx, shipping.ShippingID(req.ShipmentId))
	if err != nil {
		return nil, err
	}
	if entity.Status == shipping.ShippingStatusCancelled {
		return nil, fmt.Errorf("shipment cannot be delivered after cancellation")
	}
	if entity.Status == shipping.ShippingStatusDelivered {
		return &MarkDeliveredServiceResponse{
			Success:     true,
			Message:     "already delivered",
			ShipmentId:  string(entity.ID),
			Status:      string(entity.Status),
			DeliveredAt: entity.DeliveredAt,
		}, nil
	}

	status := shipping.ShippingStatusDelivered
	deliveredAt := req.DeliveredAt
	if deliveredAt == nil {
		now := time.Now()
		deliveredAt = &now
	}

	entity.Update(nil, nil, nil, nil, &status, nil, deliveredAt, nil, nil, nil, nil, nil, nil, nil, nil)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &MarkDeliveredServiceResponse{
		Success:     true,
		Message:     "delivered",
		ShipmentId:  string(entity.ID),
		Status:      string(entity.Status),
		DeliveredAt: entity.DeliveredAt,
	}, nil
}

// CancelShipment 实现 CancelShipment 用例。
func (s *ShippingService) CancelShipment(ctx context.Context, req CancelShipmentServiceRequest) (*CancelShipmentServiceResponse, error) {
	if req.ShipmentId == "" {
		return nil, fmt.Errorf("shipment_id is required")
	}

	entity, err := s.repo.Find(ctx, shipping.ShippingID(req.ShipmentId))
	if err != nil {
		return nil, err
	}
	if entity.Status == shipping.ShippingStatusDelivered {
		return nil, fmt.Errorf("shipment cannot be cancelled after delivery")
	}
	if entity.Status == shipping.ShippingStatusCancelled {
		return &CancelShipmentServiceResponse{
			Success:    true,
			Message:    "already cancelled",
			ShipmentId: string(entity.ID),
			Status:     string(entity.Status),
		}, nil
	}

	status := shipping.ShippingStatusCancelled
	var notes *string
	if req.Reason != "" {
		notes = &req.Reason
	}

	entity.Update(nil, nil, nil, nil, &status, nil, nil, nil, nil, nil, nil, nil, nil, nil, notes)
	if err := s.repo.Save(ctx, entity); err != nil {
		return nil, err
	}

	return &CancelShipmentServiceResponse{
		Success:    true,
		Message:    "cancelled",
		ShipmentId: string(entity.ID),
		Status:     string(entity.Status),
	}, nil
}
