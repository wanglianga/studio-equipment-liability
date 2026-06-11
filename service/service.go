package service

import (
	"fmt"
	"time"

	"studio-equipment-manager/model"
	"studio-equipment-manager/store"
)

type Service struct {
	store *store.Store
}

func New(s *store.Store) *Service {
	return &Service{store: s}
}

type CreateEquipmentInput struct {
	Category       model.EquipmentCategory `json:"category"`
	Brand          string                  `json:"brand"`
	Model          string                  `json:"model"`
	LensModel      string                  `json:"lens_model,omitempty"`
	FlashPower     int                     `json:"flash_power,omitempty"`
	PreBorrowPhoto string                  `json:"pre_borrow_photo,omitempty"`
}

func (svc *Service) CreateEquipment(input CreateEquipmentInput) (*model.Equipment, error) {
	now := time.Now()
	eq := &model.Equipment{
		ID:             svc.store.NextEquipID(),
		Category:       input.Category,
		Brand:          input.Brand,
		Model:          input.Model,
		LensModel:      input.LensModel,
		FlashPower:     input.FlashPower,
		Status:         model.StatusAvailable,
		PreBorrowPhoto: input.PreBorrowPhoto,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	svc.store.SaveEquipment(eq)
	return eq, nil
}

func (svc *Service) GetEquipment(id string) (*model.Equipment, error) {
	eq, ok := svc.store.GetEquipment(id)
	if !ok {
		return nil, fmt.Errorf("equipment %s not found", id)
	}
	return eq, nil
}

func (svc *Service) ListEquipments() []*model.Equipment {
	return svc.store.ListEquipments()
}

type BorrowInput struct {
	EquipmentID      string                `json:"equipment_id"`
	CustomerName     string                `json:"customer_name"`
	CustomerPhone    string                `json:"customer_phone"`
	StudioPosition   string                `json:"studio_position"`
	Deposit          float64               `json:"deposit"`
	PreBorrowPhotos  []string              `json:"pre_borrow_photos"`
	PreConditionItems []model.PreConditionItem `json:"pre_condition_items"`
	PreConditionNote string                `json:"pre_condition_note"`
}

func (svc *Service) BorrowEquipment(input BorrowInput) (*model.BorrowRecord, error) {
	eq, ok := svc.store.GetEquipment(input.EquipmentID)
	if !ok {
		return nil, fmt.Errorf("equipment %s not found", input.EquipmentID)
	}
	if eq.Status != model.StatusAvailable {
		return nil, fmt.Errorf("equipment %s is not available, current status: %s", input.EquipmentID, eq.Status)
	}

	now := time.Now()
	record := &model.BorrowRecord{
		ID:                svc.store.NextBorrowID(),
		EquipmentID:       input.EquipmentID,
		CustomerName:      input.CustomerName,
		CustomerPhone:     input.CustomerPhone,
		StudioPosition:    input.StudioPosition,
		Deposit:           input.Deposit,
		PreBorrowPhotos:   input.PreBorrowPhotos,
		PreConditionItems: input.PreConditionItems,
		PreConditionNote:  input.PreConditionNote,
		BorrowTime:        now,
		Status:            model.BorrowActive,
		CreatedAt:         now,
	}

	eq.Status = model.StatusBorrowed
	eq.UpdatedAt = now
	svc.store.SaveEquipment(eq)
	svc.store.SaveBorrowRecord(record)

	return record, nil
}

type ReturnInspectionInput struct {
	BorrowRecordID string   `json:"borrow_record_id"`
	ReturnPhotos   []string `json:"return_photos"`
}

func (svc *Service) ReturnInspection(input ReturnInspectionInput) (*model.BorrowRecord, error) {
	record, ok := svc.store.GetBorrowRecord(input.BorrowRecordID)
	if !ok {
		return nil, fmt.Errorf("borrow record %s not found", input.BorrowRecordID)
	}
	if record.Status != model.BorrowActive {
		return nil, fmt.Errorf("borrow record %s is not active, current status: %s", input.BorrowRecordID, record.Status)
	}

	now := time.Now()
	record.ReturnTime = &now
	record.ReturnPhotos = input.ReturnPhotos
	record.Status = model.BorrowReturned

	eq, ok := svc.store.GetEquipment(record.EquipmentID)
	if ok {
		eq.Status = model.StatusAvailable
		eq.UpdatedAt = now
		svc.store.SaveEquipment(eq)
	}

	svc.store.SaveBorrowRecord(record)
	return record, nil
}

type RegisterDamageInput struct {
	BorrowRecordID string             `json:"borrow_record_id"`
	FaultPoints    []model.FaultPoint `json:"fault_points"`
	ReturnPhotos   []string           `json:"return_photos"`
}

func (svc *Service) RegisterDamage(input RegisterDamageInput) (*model.DamageReport, error) {
	record, ok := svc.store.GetBorrowRecord(input.BorrowRecordID)
	if !ok {
		return nil, fmt.Errorf("borrow record %s not found", input.BorrowRecordID)
	}

	respType, note := svc.determineResponsibility(record, input.FaultPoints)

	report := &model.DamageReport{
		ID:                 svc.store.NextDamageID(),
		BorrowRecordID:     input.BorrowRecordID,
		EquipmentID:        record.EquipmentID,
		FaultPoints:        input.FaultPoints,
		ReturnPhotos:       input.ReturnPhotos,
		Responsibility:     respType,
		ResponsibilityNote: note,
		CreatedAt:          time.Now(),
	}

	eq, ok := svc.store.GetEquipment(record.EquipmentID)
	if ok {
		switch respType {
		case model.NormalWear, model.Undetermined, model.PreviousRemnant:
			eq.Status = model.StatusAvailable
		default:
			eq.Status = model.StatusDamaged
		}
		eq.UpdatedAt = time.Now()
		svc.store.SaveEquipment(eq)
	}

	svc.store.SaveDamageReport(report)
	return report, nil
}

func (svc *Service) determineResponsibility(record *model.BorrowRecord, faults []model.FaultPoint) (model.ResponsibilityType, string) {
	hasMissingAccessory := false
	hasPhysicalDamage := false
	hasMinorScratch := false

	for _, fp := range faults {
		switch fp.Severity {
		case "missing":
			hasMissingAccessory = true
		case "severe":
			hasPhysicalDamage = true
		case "minor":
			hasMinorScratch = true
		}
	}

	if hasMissingAccessory {
		return model.AccessoryMissing, "归还时发现配件缺失，判定为客户责任"
	}

	hasPrevDamage := svc.store.HasPreviousDamage(record.EquipmentID, record.BorrowTime)
	hasPrevUndetermined := svc.store.HasPreviousUndetermined(record.EquipmentID, record.BorrowTime)

	if hasPrevUndetermined {
		return model.PreviousRemnant, "该器材在此前借出期间已有未判定损坏记录，判定为前序遗留"
	}

	preConditionMatchCount := 0
	nonPreConditionFaults := make([]model.FaultPoint, 0)
	for _, fp := range faults {
		if fp.Severity == "missing" {
			nonPreConditionFaults = append(nonPreConditionFaults, fp)
			continue
		}
		matched := false
		for _, pc := range record.PreConditionItems {
			if fp.Location == pc.Location && (fp.Severity == pc.Severity || fp.Severity == "minor") {
				matched = true
				preConditionMatchCount++
				break
			}
		}
		if !matched {
			nonPreConditionFaults = append(nonPreConditionFaults, fp)
		}
	}

	if preConditionMatchCount > 0 && len(nonPreConditionFaults) == 0 {
		return model.PreviousRemnant, fmt.Sprintf("所有故障点均与借前已记录状态一致（%d处匹配），判定为前序遗留", preConditionMatchCount)
	}

	if hasPhysicalDamage && len(record.PreBorrowPhotos) == 0 && preConditionMatchCount == 0 {
		return model.Undetermined, "缺乏借前照片作为对比依据，无法确定损坏发生时间，判定为无法判定"
	}

	if hasPrevDamage && hasMinorScratch && !hasPhysicalDamage {
		return model.PreviousRemnant, "仅存在轻微磨损且该器材有历史损坏记录，判定为前序遗留"
	}

	if hasPhysicalDamage {
		if preConditionMatchCount > 0 {
			return model.PreviousRemnant, fmt.Sprintf("部分损坏与借前状态一致（%d处匹配），判定为前序遗留", preConditionMatchCount)
		}
		eq, ok := svc.store.GetEquipment(record.EquipmentID)
		if ok && (eq.Category == model.CategoryFlash || eq.Category == model.CategoryBackground) && record.StudioPosition != "" {
			return model.TransportImpact, "灯具或背景架在棚位间移动中可能出现碰撞，判定为运输碰撞"
		}
		return model.CustomerDamage, "借前照片显示器材完好，归还时发现明显损坏，判定为客户损坏"
	}

	if hasMinorScratch && !hasPhysicalDamage {
		return model.NormalWear, "仅存在轻微划痕或磨损，属于正常使用损耗，判定为正常磨损"
	}

	return model.Undetermined, "现有证据不足以判定责任归属，判定为无法判定"
}

type CreateRepairQuoteInput struct {
	DamageReportID string  `json:"damage_report_id"`
	RepairCost     float64 `json:"repair_cost"`
	LaborCost      float64 `json:"labor_cost"`
	Description    string  `json:"description"`
}

func (svc *Service) CreateRepairQuote(input CreateRepairQuoteInput) (*model.RepairQuote, error) {
	report, ok := svc.store.GetDamageReport(input.DamageReportID)
	if !ok {
		return nil, fmt.Errorf("damage report %s not found", input.DamageReportID)
	}

	quote := &model.RepairQuote{
		ID:             svc.store.NextRepairID(),
		DamageReportID: input.DamageReportID,
		RepairCost:     input.RepairCost,
		LaborCost:      input.LaborCost,
		TotalCost:      input.RepairCost + input.LaborCost,
		Description:    input.Description,
		CreatedAt:      time.Now(),
	}

	eq, ok := svc.store.GetEquipment(report.EquipmentID)
	if ok {
		eq.Status = model.StatusRepairing
		eq.UpdatedAt = time.Now()
		svc.store.SaveEquipment(eq)
	}

	svc.store.SaveRepairQuote(quote)
	return quote, nil
}

type DeductDepositInput struct {
	BorrowRecordID string `json:"borrow_record_id"`
	RepairQuoteID  string `json:"repair_quote_id"`
	Note           string `json:"note,omitempty"`
}

func (svc *Service) DeductDeposit(input DeductDepositInput) (*model.DeductionRecord, error) {
	record, ok := svc.store.GetBorrowRecord(input.BorrowRecordID)
	if !ok {
		return nil, fmt.Errorf("borrow record %s not found", input.BorrowRecordID)
	}

	if record.DepositFrozen {
		appeal := svc.store.FindAppealByBorrow(input.BorrowRecordID)
		if appeal != nil && appeal.Status == model.AppealPending {
			return nil, fmt.Errorf("押金已冻结：客户申诉待审核，暂停押金结算。原因：%s", record.DepositFrozenReason)
		}
	}

	quote, ok := svc.store.GetRepairQuote(input.RepairQuoteID)
	if !ok {
		return nil, fmt.Errorf("repair quote %s not found", input.RepairQuoteID)
	}

	damageReport, ok := svc.store.GetDamageReport(quote.DamageReportID)
	if !ok {
		return nil, fmt.Errorf("damage report for quote not found")
	}

	var deductAmount float64
	var refundAmount float64
	var note string

	switch damageReport.Responsibility {
	case model.NormalWear, model.PreviousRemnant:
		deductAmount = 0
		refundAmount = record.Deposit
		note = "责任判定为" + string(damageReport.Responsibility) + "，不扣除押金，全额退还"
	case model.CustomerDamage, model.AccessoryMissing, model.TransportImpact:
		if quote.TotalCost >= record.Deposit {
			deductAmount = record.Deposit
			refundAmount = 0
			note = "维修费用 %.2f 超出押金 %.2f，全额扣除押金，需追加赔付 %.2f"
			note = fmt.Sprintf(note, quote.TotalCost, record.Deposit, quote.TotalCost-record.Deposit)
		} else {
			deductAmount = quote.TotalCost
			refundAmount = record.Deposit - quote.TotalCost
			note = "扣除维修费用 %.2f，退还剩余押金 %.2f"
			note = fmt.Sprintf(note, deductAmount, refundAmount)
		}
	default:
		deductAmount = 0
		refundAmount = record.Deposit
		note = "责任未确定，暂不扣除押金，全额退还；待申诉或补充证据后重新判定"
	}

	if input.Note != "" {
		note = input.Note + "；" + note
	}

	deduction := &model.DeductionRecord{
		ID:             svc.store.NextDeductionID(),
		BorrowRecordID: input.BorrowRecordID,
		RepairQuoteID:  input.RepairQuoteID,
		DeductAmount:   deductAmount,
		RefundAmount:   refundAmount,
		Note:           note,
		CreatedAt:      time.Now(),
	}

	svc.store.SaveDeductionRecord(deduction)

	if damageReport.Responsibility == model.CustomerDamage ||
		damageReport.Responsibility == model.AccessoryMissing ||
		damageReport.Responsibility == model.TransportImpact {
		if quote.TotalCost > record.Deposit {
			existing := svc.store.FindAdditionalCompensationByRepairQuote(input.RepairQuoteID)
			if existing == nil {
				svc.createAdditionalCompensation(record, quote, deduction)
			}
		}
	}

	return deduction, nil
}

type CreateAppealInput struct {
	BorrowRecordID string   `json:"borrow_record_id"`
	CustomerName   string   `json:"customer_name"`
	Reason         string   `json:"reason"`
	Evidence       []string `json:"evidence"`
}

func (svc *Service) CreateAppeal(input CreateAppealInput) (*model.Appeal, error) {
	record, ok := svc.store.GetBorrowRecord(input.BorrowRecordID)
	if !ok {
		return nil, fmt.Errorf("borrow record %s not found", input.BorrowRecordID)
	}

	existingAppeal := svc.store.FindAppealByBorrow(input.BorrowRecordID)
	if existingAppeal != nil && existingAppeal.Status == model.AppealPending {
		return nil, fmt.Errorf("borrow record %s already has a pending appeal", input.BorrowRecordID)
	}

	damageReport := svc.store.FindDamageReportByBorrow(input.BorrowRecordID)
	if damageReport == nil {
		return nil, fmt.Errorf("no damage report found for borrow record %s", input.BorrowRecordID)
	}

	now := time.Now()
	appeal := &model.Appeal{
		ID:             svc.store.NextAppealID(),
		BorrowRecordID: input.BorrowRecordID,
		CustomerName:   input.CustomerName,
		Reason:         input.Reason,
		Evidence:       input.Evidence,
		Status:         model.AppealPending,
		CreatedAt:      now,
	}

	record.Status = model.BorrowDepositFrozen
	record.DepositFrozen = true
	record.DepositFrozenAt = &now
	record.DepositFrozenReason = "客户申诉：不同意责任结论(" + string(damageReport.Responsibility) + ")，暂停押金结算"
	svc.store.SaveBorrowRecord(record)
	svc.store.SaveAppeal(appeal)

	return appeal, nil
}

type ReviewAppealInput struct {
	AppealID   string `json:"appeal_id"`
	Accepted   bool   `json:"accepted"`
	ReviewNote string `json:"review_note"`
}

func (svc *Service) ReviewAppeal(input ReviewAppealInput) (*model.Appeal, error) {
	appeal, ok := svc.store.GetAppeal(input.AppealID)
	if !ok {
		return nil, fmt.Errorf("appeal %s not found", input.AppealID)
	}

	if appeal.Status != model.AppealPending {
		return nil, fmt.Errorf("appeal %s is not pending, current status: %s", input.AppealID, appeal.Status)
	}

	now := time.Now()
	appeal.ReviewedAt = &now
	appeal.ReviewNote = input.ReviewNote

	if input.Accepted {
		appeal.Status = model.AppealAccepted

		damageReport := svc.store.FindDamageReportByBorrow(appeal.BorrowRecordID)
		if damageReport != nil {
			damageReport.Responsibility = model.Undetermined
			damageReport.ResponsibilityNote = "客户申诉通过，原判定撤销，重新判定为无法判定；" + input.ReviewNote
			svc.store.SaveDamageReport(damageReport)

			deduction := svc.store.FindDeductionByBorrow(appeal.BorrowRecordID)
			if deduction != nil {
				record, _ := svc.store.GetBorrowRecord(appeal.BorrowRecordID)
				if record != nil {
					deduction.DeductAmount = 0
					deduction.RefundAmount = record.Deposit
					deduction.Note = "申诉通过，撤销押金扣除，全额退还押金"
					svc.store.SaveDeductionRecord(deduction)
				}
			}
		}
	} else {
		appeal.Status = model.AppealRejected
	}

	record, _ := svc.store.GetBorrowRecord(appeal.BorrowRecordID)
	if record != nil {
		record.Status = model.BorrowClosed
		svc.store.SaveBorrowRecord(record)
	}

	svc.store.SaveAppeal(appeal)
	return appeal, nil
}

func (svc *Service) ListBorrowRecords() []*model.BorrowRecord {
	return svc.store.ListBorrowRecords()
}

func (svc *Service) ListDamageReports() []*model.DamageReport {
	return svc.store.ListDamageReports()
}

func (svc *Service) ListAppeals() []*model.Appeal {
	return svc.store.ListAppeals()
}

func (svc *Service) GetDamageReport(id string) (*model.DamageReport, error) {
	r, ok := svc.store.GetDamageReport(id)
	if !ok {
		return nil, fmt.Errorf("damage report %s not found", id)
	}
	return r, nil
}

func (svc *Service) CompleteRepair(equipmentID string) (*model.Equipment, error) {
	eq, ok := svc.store.GetEquipment(equipmentID)
	if !ok {
		return nil, fmt.Errorf("equipment %s not found", equipmentID)
	}
	if eq.Status != model.StatusRepairing {
		return nil, fmt.Errorf("equipment %s is not under repair, current status: %s", equipmentID, eq.Status)
	}

	eq.Status = model.StatusAvailable
	eq.UpdatedAt = time.Now()
	svc.store.SaveEquipment(eq)
	return eq, nil
}

func (svc *Service) GetBorrowRecord(id string) (*model.BorrowRecord, error) {
	r, ok := svc.store.GetBorrowRecord(id)
	if !ok {
		return nil, fmt.Errorf("borrow record %s not found", id)
	}
	return r, nil
}

type AddAccessoryPriceInput struct {
	EquipmentID string  `json:"equipment_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
}

func (svc *Service) AddAccessoryPrice(input AddAccessoryPriceInput) (*model.AccessoryPrice, error) {
	_, ok := svc.store.GetEquipment(input.EquipmentID)
	if !ok {
		return nil, fmt.Errorf("equipment %s not found", input.EquipmentID)
	}
	if input.Name == "" {
		return nil, fmt.Errorf("accessory name is required")
	}
	if input.Price <= 0 {
		return nil, fmt.Errorf("accessory price must be positive")
	}

	existing := svc.store.FindAccessoryPriceByEquipAndName(input.EquipmentID, input.Name)
	if existing != nil {
		existing.Price = input.Price
		svc.store.SaveAccessoryPrice(existing)
		return existing, nil
	}

	ap := &model.AccessoryPrice{
		ID:          svc.store.NextAccessoryID(),
		EquipmentID: input.EquipmentID,
		Name:        input.Name,
		Price:       input.Price,
	}
	svc.store.SaveAccessoryPrice(ap)
	return ap, nil
}

func (svc *Service) ListAccessoryPrices(equipmentID string) []*model.AccessoryPrice {
	if equipmentID != "" {
		return svc.store.ListAccessoryPricesByEquipment(equipmentID)
	}
	return svc.store.ListAccessoryPrices()
}

type DeductAccessoryInput struct {
	BorrowRecordID string   `json:"borrow_record_id"`
	AccessoryNames []string `json:"accessory_names"`
	Note           string   `json:"note,omitempty"`
}

func (svc *Service) DeductAccessory(input DeductAccessoryInput) (*model.DeductionRecord, error) {
	record, ok := svc.store.GetBorrowRecord(input.BorrowRecordID)
	if !ok {
		return nil, fmt.Errorf("borrow record %s not found", input.BorrowRecordID)
	}

	if len(input.AccessoryNames) == 0 {
		return nil, fmt.Errorf("at least one accessory name is required")
	}

	accessoryItems := make([]model.AccessoryDeductionItem, 0, len(input.AccessoryNames))
	var totalDeduct float64
	var missingNames []string

	for _, name := range input.AccessoryNames {
		ap := svc.store.FindAccessoryPriceByEquipAndName(record.EquipmentID, name)
		if ap == nil {
			missingNames = append(missingNames, name)
			continue
		}
		accessoryItems = append(accessoryItems, model.AccessoryDeductionItem{
			AccessoryName: name,
			Price:         ap.Price,
		})
		totalDeduct += ap.Price
	}

	if len(missingNames) > 0 && len(accessoryItems) == 0 {
		return nil, fmt.Errorf("no accessory prices found for: %v; please add accessory prices first", missingNames)
	}

	var deductAmount, refundAmount float64
	if totalDeduct >= record.Deposit {
		deductAmount = record.Deposit
		refundAmount = 0
	} else {
		deductAmount = totalDeduct
		refundAmount = record.Deposit - totalDeduct
	}

	note := fmt.Sprintf("配件缺失扣款：共 %d 项配件，合计 %.2f", len(accessoryItems), totalDeduct)
	if len(missingNames) > 0 {
		note += fmt.Sprintf("；以下配件未找到价目：%v", missingNames)
	}
	if input.Note != "" {
		note = input.Note + "；" + note
	}

	deduction := &model.DeductionRecord{
		ID:             svc.store.NextDeductionID(),
		BorrowRecordID: input.BorrowRecordID,
		DeductAmount:   deductAmount,
		RefundAmount:   refundAmount,
		AccessoryItems: accessoryItems,
		Note:           note,
		CreatedAt:      time.Now(),
	}

	svc.store.SaveDeductionRecord(deduction)

	if totalDeduct > record.Deposit {
		existing := svc.store.FindAdditionalCompensationByBorrow(input.BorrowRecordID)
		if len(existing) == 0 {
			quote := &model.RepairQuote{
				ID:        "ACC-DEDUCT",
				TotalCost: totalDeduct,
			}
			svc.createAdditionalCompensation(record, quote, deduction)
		}
	}

	return deduction, nil
}

type AddSupplementalEvidenceInput struct {
	AppealID       string             `json:"appeal_id"`
	BorrowRecordID string             `json:"borrow_record_id"`
	OperatorName   string             `json:"operator_name"`
	EvidenceType   model.EvidenceType `json:"evidence_type"`
	Description    string             `json:"description"`
	Attachments    []string           `json:"attachments"`
}

func (svc *Service) AddSupplementalEvidence(input AddSupplementalEvidenceInput) (*model.SupplementalEvidence, error) {
	if input.BorrowRecordID == "" {
		return nil, fmt.Errorf("borrow_record_id is required")
	}
	if input.OperatorName == "" {
		return nil, fmt.Errorf("operator_name is required")
	}
	if input.EvidenceType == "" {
		return nil, fmt.Errorf("evidence_type is required")
	}
	if len(input.Attachments) == 0 {
		return nil, fmt.Errorf("at least one attachment is required")
	}

	record, ok := svc.store.GetBorrowRecord(input.BorrowRecordID)
	if !ok {
		return nil, fmt.Errorf("borrow record %s not found", input.BorrowRecordID)
	}

	if input.AppealID != "" {
		appeal, ok := svc.store.GetAppeal(input.AppealID)
		if !ok {
			return nil, fmt.Errorf("appeal %s not found", input.AppealID)
		}
		if appeal.BorrowRecordID != input.BorrowRecordID {
			return nil, fmt.Errorf("appeal %s does not belong to borrow record %s", input.AppealID, input.BorrowRecordID)
		}
	}

	evidence := &model.SupplementalEvidence{
		ID:             svc.store.NextSupplementalEvidenceID(),
		AppealID:       input.AppealID,
		BorrowRecordID: input.BorrowRecordID,
		OperatorName:   input.OperatorName,
		EvidenceType:   input.EvidenceType,
		Description:    input.Description,
		Attachments:    input.Attachments,
		CreatedAt:      time.Now(),
	}

	svc.store.SaveSupplementalEvidence(evidence)
	_ = record
	return evidence, nil
}

type UpdateRepairQuoteInput struct {
	RepairQuoteID string  `json:"repair_quote_id"`
	RepairCost    float64 `json:"repair_cost"`
	LaborCost     float64 `json:"labor_cost"`
	Description   string  `json:"description"`
	UpdateNote    string  `json:"update_note"`
}

func (svc *Service) UpdateRepairQuote(input UpdateRepairQuoteInput) (*model.RepairQuote, *model.AdditionalCompensation, error) {
	if input.RepairQuoteID == "" {
		return nil, nil, fmt.Errorf("repair_quote_id is required")
	}
	if input.RepairCost < 0 || input.LaborCost < 0 {
		return nil, nil, fmt.Errorf("repair_cost and labor_cost must not be negative")
	}

	quote, ok := svc.store.GetRepairQuote(input.RepairQuoteID)
	if !ok {
		return nil, nil, fmt.Errorf("repair quote %s not found", input.RepairQuoteID)
	}

	now := time.Now()
	oldTotalCost := quote.TotalCost
	quote.RepairCost = input.RepairCost
	quote.LaborCost = input.LaborCost
	quote.TotalCost = input.RepairCost + input.LaborCost
	quote.Description = input.Description
	quote.UpdatedAt = &now
	quote.IsUpdated = true
	quote.UpdateNote = input.UpdateNote
	svc.store.UpdateRepairQuote(quote)

	damageReport, _ := svc.store.GetDamageReport(quote.DamageReportID)
	var additionalComp *model.AdditionalCompensation

	if damageReport != nil {
		if damageReport.Responsibility == model.CustomerDamage ||
			damageReport.Responsibility == model.AccessoryMissing ||
			damageReport.Responsibility == model.TransportImpact {

			record, _ := svc.store.GetBorrowRecord(damageReport.BorrowRecordID)
			if record != nil && quote.TotalCost > record.Deposit {
				existingDeduction := svc.store.FindDeductionByBorrow(record.ID)
				if existingDeduction != nil {
					existingComp := svc.store.FindAdditionalCompensationByRepairQuote(input.RepairQuoteID)
					if existingComp != nil {
						existingComp.NewTotalCost = quote.TotalCost
						existingComp.AdditionalAmount = quote.TotalCost - record.Deposit
						existingComp.Note = fmt.Sprintf("报价更新：原总费用%.2f→新总费用%.2f，追加赔付调整为%.2f",
							oldTotalCost, quote.TotalCost, existingComp.AdditionalAmount)
						svc.store.SaveAdditionalCompensation(existingComp)
						additionalComp = existingComp
					} else {
						additionalComp = svc.createAdditionalCompensation(record, quote, existingDeduction)
					}

					existingDeduction.DeductAmount = record.Deposit
					existingDeduction.RefundAmount = 0
					existingDeduction.Note = fmt.Sprintf("维修报价更新(%.2f→%.2f)：维修费用超出押金 %.2f，全额扣除押金，需追加赔付 %.2f",
						oldTotalCost, quote.TotalCost, record.Deposit, quote.TotalCost-record.Deposit)
					svc.store.SaveDeductionRecord(existingDeduction)
				}
			}
		}
	}

	return quote, additionalComp, nil
}

func (svc *Service) createAdditionalCompensation(
	record *model.BorrowRecord,
	quote *model.RepairQuote,
	deduction *model.DeductionRecord,
) *model.AdditionalCompensation {
	additionalAmount := quote.TotalCost - record.Deposit
	if additionalAmount <= 0 {
		return nil
	}

	comp := &model.AdditionalCompensation{
		ID:               svc.store.NextAdditionalCompensationID(),
		BorrowRecordID:   record.ID,
		RepairQuoteID:    quote.ID,
		OriginalDeduct:   deduction.DeductAmount,
		OriginalRefund:   deduction.RefundAmount,
		NewTotalCost:     quote.TotalCost,
		DepositAmount:    record.Deposit,
		AdditionalAmount: additionalAmount,
		Status:           "pending",
		Note:             fmt.Sprintf("维修总费用%.2f超过押金%.2f，客户需追加赔付%.2f", quote.TotalCost, record.Deposit, additionalAmount),
		CreatedAt:        time.Now(),
	}
	svc.store.SaveAdditionalCompensation(comp)
	return comp
}

func (svc *Service) ListSupplementalEvidences(appealID, borrowRecordID string) []*model.SupplementalEvidence {
	if appealID != "" {
		return svc.store.FindSupplementalEvidenceByAppeal(appealID)
	}
	if borrowRecordID != "" {
		return svc.store.FindSupplementalEvidenceByBorrow(borrowRecordID)
	}
	return svc.store.ListSupplementalEvidences()
}

func (svc *Service) ListAdditionalCompensations(borrowRecordID string) []*model.AdditionalCompensation {
	if borrowRecordID != "" {
		return svc.store.FindAdditionalCompensationByBorrow(borrowRecordID)
	}
	return svc.store.ListAdditionalCompensations()
}

type CollectAdditionalCompensationInput struct {
	CompensationID string `json:"compensation_id"`
}

func (svc *Service) CollectAdditionalCompensation(input CollectAdditionalCompensationInput) (*model.AdditionalCompensation, error) {
	if input.CompensationID == "" {
		return nil, fmt.Errorf("compensation_id is required")
	}

	comp, ok := svc.store.GetAdditionalCompensation(input.CompensationID)
	if !ok {
		return nil, fmt.Errorf("additional compensation %s not found", input.CompensationID)
	}

	if comp.Status == "collected" {
		return nil, fmt.Errorf("additional compensation %s already collected", input.CompensationID)
	}

	now := time.Now()
	comp.Status = "collected"
	comp.CollectedAt = &now
	svc.store.SaveAdditionalCompensation(comp)
	return comp, nil
}
