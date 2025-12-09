package domain

// --------------------------------------------------------
// ENUMS
// --------------------------------------------------------

type UserRole string

const (
	RoleAdmin    UserRole = "ADMIN"
	RoleManager  UserRole = "MANAGER"
	RoleStaff    UserRole = "STAFF"
	RoleCustomer UserRole = "CUSTOMER"
	RoleSupplier UserRole = "SUPPLIER"
)

type PurchaseOrderStatus string

const (
	PODraft             PurchaseOrderStatus = "DRAFT"
	POSent              PurchaseOrderStatus = "SENT"
	POPartiallyReceived PurchaseOrderStatus = "PARTIALLY_RECEIVED"
	POCompleted         PurchaseOrderStatus = "COMPLETED"
	POCancelled         PurchaseOrderStatus = "CANCELLED"
)

type MovementReason string

const (
	ReasonSale                MovementReason = "SALE"
	ReasonPurchase            MovementReason = "PURCHASE"
	ReasonAssemblyConsumption MovementReason = "ASSEMBLY_CONSUMPTION"
	ReasonAssemblyOutput      MovementReason = "ASSEMBLY_OUTPUT"
	ReasonTransfer            MovementReason = "TRANSFER"
	ReasonAdjustment          MovementReason = "ADJUSTMENT"
	ReasonReturn              MovementReason = "RETURN"
)

type POSSessionStatus string

const (
	SessionOpeningControl POSSessionStatus = "OPENING_CONTROL"
	SessionOpened         POSSessionStatus = "OPENED"
	SessionClosingControl POSSessionStatus = "CLOSING_CONTROL"
	SessionClosed         POSSessionStatus = "CLOSED"
)

type OrderStatus string

const (
	OrderDraft     OrderStatus = "DRAFT"
	OrderQuotation OrderStatus = "QUOTATION"
	OrderConfirmed OrderStatus = "CONFIRMED"
	OrderCompleted OrderStatus = "COMPLETED"
	OrderCancelled OrderStatus = "CANCELLED"
)

type ShipmentStatus string

const (
	ShipmentPending     ShipmentStatus = "PENDING"
	ShipmentReadyToPack ShipmentStatus = "READY_TO_PACK"
	ShipmentShipped     ShipmentStatus = "SHIPPED"
	ShipmentDelivered   ShipmentStatus = "DELIVERED"
	ShipmentReturned    ShipmentStatus = "RETURNED"
)

type PaymentStatus string

const (
	PaymentUnpaid            PaymentStatus = "UNPAID"
	PaymentPaid              PaymentStatus = "PAID"
	PaymentPartiallyRefunded PaymentStatus = "PARTIALLY_REFUNDED"
	PaymentRefunded          PaymentStatus = "REFUNDED"
)

type ProductCondition string

const (
	ConditionNew         ProductCondition = "NEW"
	ConditionRefurbished ProductCondition = "REFURBISHED"
	ConditionDamaged     ProductCondition = "DAMAGED"
)

type TransactionType string

const (
	TransCredit TransactionType = "CREDIT"
	TransDebit  TransactionType = "DEBIT"
)

type OrderChannel string

const (
	ChannelWeb         OrderChannel = "WEB"
	ChannelPOS         OrderChannel = "POS"
	ChannelMarketplace OrderChannel = "MARKETPLACE"
	ChannelManual      OrderChannel = "MANUAL"
)

type InvoiceType string

const (
	InvoiceCustomer InvoiceType = "CUSTOMER_INVOICE"
	InvoiceSupplier InvoiceType = "SUPPLIER_BILL"
)

type InvoiceStatus string

const (
	InvoiceDraft  InvoiceStatus = "DRAFT"
	InvoicePosted InvoiceStatus = "POSTED"
	InvoicePaid   InvoiceStatus = "PAID"
	InvoiceVoid   InvoiceStatus = "VOID"
)

type CashMoveType string

const (
	CashAdd  CashMoveType = "ADD"
	CashDrop CashMoveType = "DROP"
)
