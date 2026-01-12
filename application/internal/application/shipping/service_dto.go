package shippingapp

import "time"

// CreateShipmentServiceRequest 是 CreateShipment 方法的请求参数。
type CreateShipmentServiceRequest struct {
	OrderId            string `json:"order_id"`                       // 订单ID
	Carrier            string `json:"carrier"`                        // 物流承运商
	ShippingMethod     string `json:"shipping_method"`                // 配送方式
	ReceiverName       string `json:"receiver_name"`                  // 收件人姓名
	ReceiverPhone      string `json:"receiver_phone"`                 // 收件人电话
	ReceiverAddress    string `json:"receiver_address"`               // 收件人地址
	ReceiverCity       string `json:"receiver_city,omitempty"`        // 收件人城市
	ReceiverState      string `json:"receiver_state,omitempty"`       // 收件人省/州
	ReceiverCountry    string `json:"receiver_country,omitempty"`     // 收件人国家
	ReceiverPostalCode string `json:"receiver_postal_code,omitempty"` // 邮编
	Notes              string `json:"notes,omitempty"`                // 备注
}

// CreateShipmentServiceResponse 是 CreateShipment 方法的响应结果。
type CreateShipmentServiceResponse struct {
	Success        bool       `json:"success"`                   // 操作是否成功
	Message        string     `json:"message,omitempty"`         // 提示消息
	ShipmentId     string     `json:"shipment_id,omitempty"`     // 物流ID
	TrackingNumber string     `json:"tracking_number,omitempty"` // 物流单号
	Status         string     `json:"status,omitempty"`          // 物流状态
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`      // 发货时间
}

// UpdateTrackingServiceRequest 是 UpdateTracking 方法的请求参数。
type UpdateTrackingServiceRequest struct {
	ShipmentId     string     `json:"shipment_id"`          // 物流ID
	TrackingNumber string     `json:"tracking_number"`      // 物流单号
	Status         string     `json:"status,omitempty"`     // 物流状态
	UpdatedAt      *time.Time `json:"updated_at,omitempty"` // 更新时间
}

// UpdateTrackingServiceResponse 是 UpdateTracking 方法的响应结果。
type UpdateTrackingServiceResponse struct {
	Success        bool   `json:"success"`                   // 操作是否成功
	Message        string `json:"message,omitempty"`         // 提示消息
	ShipmentId     string `json:"shipment_id,omitempty"`     // 物流ID
	TrackingNumber string `json:"tracking_number,omitempty"` // 物流单号
	Status         string `json:"status,omitempty"`          // 物流状态
}

// MarkDeliveredServiceRequest 是 MarkDelivered 方法的请求参数。
type MarkDeliveredServiceRequest struct {
	ShipmentId        string     `json:"shipment_id"`                  // 物流ID
	DeliveredAt       *time.Time `json:"delivered_at,omitempty"`       // 签收时间
	ReceiverSignature string     `json:"receiver_signature,omitempty"` // 签收凭证
}

// MarkDeliveredServiceResponse 是 MarkDelivered 方法的响应结果。
type MarkDeliveredServiceResponse struct {
	Success     bool       `json:"success"`                // 操作是否成功
	Message     string     `json:"message,omitempty"`      // 提示消息
	ShipmentId  string     `json:"shipment_id,omitempty"`  // 物流ID
	Status      string     `json:"status,omitempty"`       // 物流状态
	DeliveredAt *time.Time `json:"delivered_at,omitempty"` // 签收时间
}

// CancelShipmentServiceRequest 是 CancelShipment 方法的请求参数。
type CancelShipmentServiceRequest struct {
	ShipmentId  string     `json:"shipment_id"`            // 物流ID
	Reason      string     `json:"reason,omitempty"`       // 取消原因
	CancelledAt *time.Time `json:"cancelled_at,omitempty"` // 取消时间
}

// CancelShipmentServiceResponse 是 CancelShipment 方法的响应结果。
type CancelShipmentServiceResponse struct {
	Success    bool   `json:"success"`               // 操作是否成功
	Message    string `json:"message,omitempty"`     // 提示消息
	ShipmentId string `json:"shipment_id,omitempty"` // 物流ID
	Status     string `json:"status,omitempty"`      // 物流状态
}
