package model

import "time"

type EquipmentCategory string

const (
	CategoryCamera     EquipmentCategory = "camera"
	CategoryLens       EquipmentCategory = "lens"
	CategoryFlash      EquipmentCategory = "flash"
	CategoryBackground EquipmentCategory = "background_stand"
)

type EquipmentStatus string

const (
	StatusAvailable EquipmentStatus = "available"
	StatusBorrowed  EquipmentStatus = "borrowed"
	StatusDamaged   EquipmentStatus = "damaged"
	StatusRepairing EquipmentStatus = "repairing"
)

type Equipment struct {
	ID             string            `json:"id"`
	Category       EquipmentCategory `json:"category"`
	Brand          string            `json:"brand"`
	Model          string            `json:"model"`
	LensModel      string            `json:"lens_model,omitempty"`
	FlashPower     int               `json:"flash_power,omitempty"`
	Status         EquipmentStatus   `json:"status"`
	PreBorrowPhoto string            `json:"pre_borrow_photo,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type BorrowStatus string

const (
	BorrowActive       BorrowStatus = "active"
	BorrowReturned     BorrowStatus = "returned"
	BorrowAppealed     BorrowStatus = "appealed"
	BorrowDepositFrozen BorrowStatus = "deposit_frozen"
	BorrowClosed       BorrowStatus = "closed"
)

type PreConditionItem struct {
	Location    string `json:"location"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Photo       string `json:"photo,omitempty"`
}

type BorrowRecord struct {
	ID                 string            `json:"id"`
	EquipmentID        string            `json:"equipment_id"`
	CustomerName       string            `json:"customer_name"`
	CustomerPhone      string            `json:"customer_phone"`
	StudioPosition     string            `json:"studio_position"`
	Deposit            float64           `json:"deposit"`
	PreBorrowPhotos    []string          `json:"pre_borrow_photos"`
	PreConditionItems  []PreConditionItem `json:"pre_condition_items,omitempty"`
	PreConditionNote   string            `json:"pre_condition_note,omitempty"`
	BorrowTime         time.Time         `json:"borrow_time"`
	ReturnTime         *time.Time        `json:"return_time,omitempty"`
	ReturnPhotos       []string          `json:"return_photos,omitempty"`
	Status             BorrowStatus      `json:"status"`
	DepositFrozen      bool              `json:"deposit_frozen"`
	DepositFrozenAt    *time.Time        `json:"deposit_frozen_at,omitempty"`
	DepositFrozenReason string           `json:"deposit_frozen_reason,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
}

type ResponsibilityType string

const (
	NormalWear       ResponsibilityType = "normal_wear"
	CustomerDamage   ResponsibilityType = "customer_damage"
	PreviousRemnant  ResponsibilityType = "previous_remnant"
	TransportImpact  ResponsibilityType = "transport_impact"
	AccessoryMissing ResponsibilityType = "accessory_missing"
	Undetermined     ResponsibilityType = "undetermined"
)

type DamageReport struct {
	ID                 string             `json:"id"`
	BorrowRecordID     string             `json:"borrow_record_id"`
	EquipmentID        string             `json:"equipment_id"`
	FaultPoints        []FaultPoint       `json:"fault_points"`
	ReturnPhotos       []string           `json:"return_photos"`
	Responsibility     ResponsibilityType `json:"responsibility"`
	ResponsibilityNote string             `json:"responsibility_note"`
	CreatedAt          time.Time          `json:"created_at"`
}

type FaultPoint struct {
	Location    string `json:"location"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type AccessoryPrice struct {
	ID          string  `json:"id"`
	EquipmentID string  `json:"equipment_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
}

type AccessoryDeductionItem struct {
	AccessoryName string  `json:"accessory_name"`
	Price         float64 `json:"price"`
}

type DeductionRecord struct {
	ID               string                  `json:"id"`
	BorrowRecordID   string                  `json:"borrow_record_id"`
	RepairQuoteID    string                  `json:"repair_quote_id,omitempty"`
	DeductAmount     float64                 `json:"deduct_amount"`
	RefundAmount     float64                 `json:"refund_amount"`
	AccessoryItems   []AccessoryDeductionItem `json:"accessory_items,omitempty"`
	Note             string                  `json:"note"`
	CreatedAt        time.Time               `json:"created_at"`
}

type AppealStatus string

const (
	AppealPending  AppealStatus = "pending"
	AppealAccepted AppealStatus = "accepted"
	AppealRejected AppealStatus = "rejected"
)

type Appeal struct {
	ID             string       `json:"id"`
	BorrowRecordID string       `json:"borrow_record_id"`
	CustomerName   string       `json:"customer_name"`
	Reason         string       `json:"reason"`
	Evidence       []string     `json:"evidence"`
	Status         AppealStatus `json:"status"`
	ReviewNote     string       `json:"review_note,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	ReviewedAt     *time.Time   `json:"reviewed_at,omitempty"`
}

type EvidenceType string

const (
	EvidenceSurveillance  EvidenceType = "surveillance"
	EvidenceHandoverPhoto EvidenceType = "handover_photo"
	EvidencePreCheck      EvidenceType = "pre_check_record"
	EvidencePostCheck     EvidenceType = "post_check_record"
	EvidenceOther         EvidenceType = "other"
)

type SupplementalEvidence struct {
	ID             string       `json:"id"`
	AppealID       string       `json:"appeal_id"`
	BorrowRecordID string       `json:"borrow_record_id"`
	OperatorName   string       `json:"operator_name"`
	EvidenceType   EvidenceType `json:"evidence_type"`
	Description    string       `json:"description"`
	Attachments    []string     `json:"attachments"`
	CreatedAt      time.Time    `json:"created_at"`
}

type RepairQuote struct {
	ID             string     `json:"id"`
	DamageReportID string     `json:"damage_report_id"`
	RepairCost     float64    `json:"repair_cost"`
	LaborCost      float64    `json:"labor_cost"`
	TotalCost      float64    `json:"total_cost"`
	Description    string     `json:"description"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	IsUpdated      bool       `json:"is_updated"`
	UpdateNote     string     `json:"update_note,omitempty"`
}

type AdditionalCompensation struct {
	ID                string     `json:"id"`
	BorrowRecordID    string     `json:"borrow_record_id"`
	RepairQuoteID     string     `json:"repair_quote_id"`
	OriginalDeduct    float64    `json:"original_deduct"`
	OriginalRefund    float64    `json:"original_refund"`
	NewTotalCost      float64    `json:"new_total_cost"`
	DepositAmount     float64    `json:"deposit_amount"`
	AdditionalAmount  float64    `json:"additional_amount"`
	Status            string     `json:"status"`
	Note              string     `json:"note"`
	CreatedAt         time.Time  `json:"created_at"`
	CollectedAt       *time.Time `json:"collected_at,omitempty"`
}
