package responses

import commonConstant "online-shop-2N/pkg/common/constants"

type OrderPayment struct {
	PaymentType  commonConstant.PaymentType `json:"payment_type"`
	PaymentOrder any                        `json:"payment_order"`
}
