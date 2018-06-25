package server

//响应码
const (
	RequestOK        = 200
	RequestRefuse    = 400
	RequestFail      = 500
	RequestTimeout   = 600
	RequestNoAuth    = 900
	InvalidRequest   = 3000
	EmptyRequest     = 3001
	InvalidSign      = 3002
	InvalidParameter = 3003
)

//结果状态码
const (
	ResultOK = iota + 1000
	BussinessIdErr
	ProductIdErr
	CountNotEnoughErr
	PurchaseTimeErr
	CouponValidTimeErr
	CouponCartIdErr
	CouponPasswordErr
	CouponUseTimeErr
	CouponStatusErr
	CouponChannelErr
	CouponIdErr
	OrderInfoErr
	OrderInfoUpdateErr
	SignErr
	BussinessNotApproveErr
	RequestUrlErr
	ProductInfoErr
	CouponInfoErr
	OrderStatusErr
	OrderIdErr
)

//订单状态
const (
	OrderCanceled = iota - 1
	OrderWait
	OrderFinished
)
