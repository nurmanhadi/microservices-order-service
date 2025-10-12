package enum

type STATUS string

var (
	STATUS_PENDING  STATUS = "pending"
	STATUS_PAID     STATUS = "paid"
	STATUS_SHIPPED  STATUS = "shipped"
	STATUS_CENCELED STATUS = "cenceled"
)
