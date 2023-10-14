package common

// for defining ENUM stasues
type OrderStatusType string

// payment types
type PaymentType string

const (
	// order status
	StatusPaymentPending  OrderStatusType = "payment pending"
	StatusOrderPlaced     OrderStatusType = "order placed"
	StatusOrderCancelled  OrderStatusType = "order cancelled"
	StatusOrderDelivered  OrderStatusType = "order delivered"
	StatusReturnRequested OrderStatusType = "return requested"
	StatusReturnApproved  OrderStatusType = "return approved"
	StatusReturnCancelled OrderStatusType = "return cancelled"
	StatusOrderReturned   OrderStatusType = "order returned"

	// payment type
	RazopayPayment        PaymentType = "razor pay"
	RazorPayMaximumAmount             = 50000 // this is only for initial admin can later change this
	CodPayment            PaymentType = "cod"
	CodMaximumAmount                  = 20000
	StripePayment         PaymentType = "stripe"
	StripeMaximumAmount               = 50000
)
