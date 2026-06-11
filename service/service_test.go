package service

import (
	"testing"

	"studio-equipment-manager/model"
	"studio-equipment-manager/store"
)

func TestPreviousRemnantViaNormalWear(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryLens,
		Brand:          "Canon",
		Model:          "EF 50mm",
		LensModel:      "EF 50mm f/1.4",
		PreBorrowPhoto: "before_001.jpg",
	})

	br1, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Alice",
		CustomerPhone:   "111",
		StudioPosition:  "A-1",
		Deposit:         1000,
		PreBorrowPhotos: []string{"borrow1_before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br1.ID,
		ReturnPhotos:   []string{"return1.jpg"},
	})

	dm1, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br1.ID,
		FaultPoints:    []model.FaultPoint{{Location: "barrel", Description: "light scratch", Severity: "minor"}},
		ReturnPhotos:   []string{"damage1.jpg"},
	})

	if dm1.Responsibility != model.NormalWear {
		t.Fatalf("first damage should be normal_wear, got %s", dm1.Responsibility)
	}

	eqAfter, _ := svc.GetEquipment(eq.ID)
	if eqAfter.Status != model.StatusAvailable {
		t.Fatalf("after normal_wear, equipment should be available, got %s", eqAfter.Status)
	}

	br2, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Bob",
		CustomerPhone:   "222",
		StudioPosition:  "A-2",
		Deposit:         1000,
		PreBorrowPhotos: []string{"borrow2_before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br2.ID,
		ReturnPhotos:   []string{"return2.jpg"},
	})

	dm2, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br2.ID,
		FaultPoints:    []model.FaultPoint{{Location: "barrel", Description: "light scratch", Severity: "minor"}},
		ReturnPhotos:   []string{"damage2.jpg"},
	})

	if dm2.Responsibility != model.PreviousRemnant {
		t.Fatalf("second damage should be previous_remnant, got %s (note: %s)", dm2.Responsibility, dm2.ResponsibilityNote)
	}

	eqAfter2, _ := svc.GetEquipment(eq.ID)
	if eqAfter2.Status != model.StatusAvailable {
		t.Fatalf("after previous_remnant, equipment should be available, got %s", eqAfter2.Status)
	}
}

func TestPreviousRemnantViaUndetermined(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category: model.CategoryCamera,
		Brand:    "Sony",
		Model:    "A7III",
	})

	br1, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:   eq.ID,
		CustomerName:  "Alice",
		CustomerPhone: "111",
		StudioPosition: "B-1",
		Deposit:       2000,
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br1.ID,
		ReturnPhotos:   []string{"return1.jpg"},
	})

	dm1, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br1.ID,
		FaultPoints:    []model.FaultPoint{{Location: "body", Description: "dent", Severity: "severe"}},
		ReturnPhotos:   []string{"damage1.jpg"},
	})

	if dm1.Responsibility != model.Undetermined {
		t.Fatalf("first damage without pre-borrow photos should be undetermined, got %s", dm1.Responsibility)
	}

	eqAfter, _ := svc.GetEquipment(eq.ID)
	if eqAfter.Status != model.StatusAvailable {
		t.Fatalf("after undetermined, equipment should be available, got %s", eqAfter.Status)
	}

	br2, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Bob",
		CustomerPhone:   "222",
		StudioPosition:  "B-2",
		Deposit:         2000,
		PreBorrowPhotos: []string{"borrow2_before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br2.ID,
		ReturnPhotos:   []string{"return2.jpg"},
	})

	dm2, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br2.ID,
		FaultPoints:    []model.FaultPoint{{Location: "body", Description: "dent found again", Severity: "severe"}},
		ReturnPhotos:   []string{"damage2.jpg"},
	})

	if dm2.Responsibility != model.PreviousRemnant {
		t.Fatalf("second damage with prior undetermined record should be previous_remnant, got %s (note: %s)", dm2.Responsibility, dm2.ResponsibilityNote)
	}
}

func TestPreviousRemnantViaRepairComplete(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryLens,
		Brand:          "Nikon",
		Model:          "24-70mm",
		PreBorrowPhoto: "before_001.jpg",
	})

	br1, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Alice",
		CustomerPhone:   "111",
		StudioPosition:  "C-1",
		Deposit:         3000,
		PreBorrowPhotos: []string{"borrow1_before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br1.ID,
		ReturnPhotos:   []string{"return1.jpg"},
	})

	dm1, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br1.ID,
		FaultPoints:    []model.FaultPoint{{Location: "front element", Description: "crack", Severity: "severe"}},
		ReturnPhotos:   []string{"damage1.jpg"},
	})

	if dm1.Responsibility != model.CustomerDamage {
		t.Fatalf("first damage should be customer_damage, got %s", dm1.Responsibility)
	}

	eqAfter, _ := svc.GetEquipment(eq.ID)
	if eqAfter.Status != model.StatusDamaged {
		t.Fatalf("after customer_damage, equipment should be damaged, got %s", eqAfter.Status)
	}

	svc.CreateRepairQuote(CreateRepairQuoteInput{
		DamageReportID: dm1.ID,
		RepairCost:     500,
		LaborCost:      100,
		Description:    "replace front element",
	})

	eqRepairing, _ := svc.GetEquipment(eq.ID)
	if eqRepairing.Status != model.StatusRepairing {
		t.Fatalf("after repair quote, equipment should be repairing, got %s", eqRepairing.Status)
	}

	eqFixed, _ := svc.CompleteRepair(eq.ID)
	if eqFixed.Status != model.StatusAvailable {
		t.Fatalf("after repair complete, equipment should be available, got %s", eqFixed.Status)
	}

	br2, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Bob",
		CustomerPhone:   "222",
		StudioPosition:  "C-2",
		Deposit:         3000,
		PreBorrowPhotos: []string{"borrow2_before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br2.ID,
		ReturnPhotos:   []string{"return2.jpg"},
	})

	dm2, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br2.ID,
		FaultPoints:    []model.FaultPoint{{Location: "barrel", Description: "minor scratch", Severity: "minor"}},
		ReturnPhotos:   []string{"damage2.jpg"},
	})

	if dm2.Responsibility != model.PreviousRemnant {
		t.Fatalf("second damage with prior damage record + minor scratch should be previous_remnant, got %s (note: %s)", dm2.Responsibility, dm2.ResponsibilityNote)
	}
}

func TestCustomerDamageSetsDamaged(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryLens,
		Brand:          "Canon",
		Model:          "EF 85mm",
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Alice",
		CustomerPhone:   "111",
		StudioPosition:  "D-1",
		Deposit:         1000,
		PreBorrowPhotos: []string{"before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints:    []model.FaultPoint{{Location: "glass", Description: "crack", Severity: "severe"}},
		ReturnPhotos:   []string{"damage.jpg"},
	})

	if dm.Responsibility != model.CustomerDamage {
		t.Fatalf("expected customer_damage, got %s", dm.Responsibility)
	}

	eqAfter, _ := svc.GetEquipment(eq.ID)
	if eqAfter.Status != model.StatusDamaged {
		t.Fatalf("after customer_damage, equipment should be damaged, got %s", eqAfter.Status)
	}
}

func TestTransportImpactSetsDamaged(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryFlash,
		Brand:          "Godox",
		Model:          "AD600",
		FlashPower:     600,
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Alice",
		CustomerPhone:   "111",
		StudioPosition:  "E-1",
		Deposit:         2000,
		PreBorrowPhotos: []string{"before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints:    []model.FaultPoint{{Location: "body", Description: "dent", Severity: "severe"}},
		ReturnPhotos:   []string{"damage.jpg"},
	})

	if dm.Responsibility != model.TransportImpact {
		t.Fatalf("expected transport_impact, got %s", dm.Responsibility)
	}

	eqAfter, _ := svc.GetEquipment(eq.ID)
	if eqAfter.Status != model.StatusDamaged {
		t.Fatalf("after transport_impact, equipment should be damaged, got %s", eqAfter.Status)
	}
}

func TestAccessoryMissingSetsDamaged(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryCamera,
		Brand:          "Sony",
		Model:          "A7IV",
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Alice",
		CustomerPhone:   "111",
		StudioPosition:  "F-1",
		Deposit:         3000,
		PreBorrowPhotos: []string{"before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints:    []model.FaultPoint{{Location: "body cap", Description: "body cap missing", Severity: "missing"}},
		ReturnPhotos:   []string{"damage.jpg"},
	})

	if dm.Responsibility != model.AccessoryMissing {
		t.Fatalf("expected accessory_missing, got %s", dm.Responsibility)
	}

	eqAfter, _ := svc.GetEquipment(eq.ID)
	if eqAfter.Status != model.StatusDamaged {
		t.Fatalf("after accessory_missing, equipment should be damaged, got %s", eqAfter.Status)
	}
}

func TestCompleteRepairOnlyFromRepairing(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category: model.CategoryCamera,
		Brand:    "Sony",
		Model:    "A7III",
	})

	_, err := svc.CompleteRepair(eq.ID)
	if err == nil {
		t.Fatal("should not be able to complete repair on available equipment")
	}
}

func TestPreConditionItemsMatchFaultPoints(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryCamera,
		Brand:          "Canon",
		Model:          "R5",
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:   eq.ID,
		CustomerName:  "Alice",
		CustomerPhone: "111",
		StudioPosition: "G-1",
		Deposit:       5000,
		PreBorrowPhotos: []string{"before.jpg"},
		PreConditionItems: []model.PreConditionItem{
			{Location: "body top", Description: "scratch near hot shoe", Severity: "minor", Photo: "pre_scratch1.jpg"},
			{Location: "lens mount", Description: "slight looseness", Severity: "minor", Photo: "pre_loose1.jpg"},
		},
		PreConditionNote: "器材出借前已有热靴附近划痕和卡口轻微松动",
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints: []model.FaultPoint{
			{Location: "body top", Description: "scratch near hot shoe", Severity: "minor"},
			{Location: "lens mount", Description: "slight looseness", Severity: "minor"},
		},
		ReturnPhotos: []string{"damage.jpg"},
	})

	if dm.Responsibility != model.PreviousRemnant {
		t.Fatalf("expected previous_remnant when all faults match pre-condition items, got %s (note: %s)", dm.Responsibility, dm.ResponsibilityNote)
	}

	eqAfter, _ := svc.GetEquipment(eq.ID)
	if eqAfter.Status != model.StatusAvailable {
		t.Fatalf("after previous_remnant, equipment should be available, got %s", eqAfter.Status)
	}
}

func TestPreConditionItemsPartialMatch(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryCamera,
		Brand:          "Canon",
		Model:          "R5",
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:   eq.ID,
		CustomerName:  "Bob",
		CustomerPhone: "222",
		StudioPosition: "G-2",
		Deposit:       5000,
		PreBorrowPhotos: []string{"before.jpg"},
		PreConditionItems: []model.PreConditionItem{
			{Location: "body top", Description: "scratch near hot shoe", Severity: "minor", Photo: "pre_scratch1.jpg"},
		},
		PreConditionNote: "器材出借前已有热靴附近划痕",
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints: []model.FaultPoint{
			{Location: "body top", Description: "scratch near hot shoe", Severity: "minor"},
			{Location: "screen", Description: "cracked LCD", Severity: "severe"},
		},
		ReturnPhotos: []string{"damage.jpg"},
	})

	if dm.Responsibility != model.PreviousRemnant {
		t.Fatalf("expected previous_remnant when partial faults match pre-condition items with preConditionMatchCount>0, got %s (note: %s)", dm.Responsibility, dm.ResponsibilityNote)
	}
}

func TestPreConditionItemsNoMatch(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryCamera,
		Brand:          "Canon",
		Model:          "R5",
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:   eq.ID,
		CustomerName:  "Charlie",
		CustomerPhone: "333",
		StudioPosition: "G-3",
		Deposit:       5000,
		PreBorrowPhotos: []string{"before.jpg"},
		PreConditionItems: []model.PreConditionItem{
			{Location: "body top", Description: "scratch near hot shoe", Severity: "minor", Photo: "pre_scratch1.jpg"},
		},
		PreConditionNote: "器材出借前已有热靴附近划痕",
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints: []model.FaultPoint{
			{Location: "screen", Description: "cracked LCD", Severity: "severe"},
		},
		ReturnPhotos: []string{"damage.jpg"},
	})

	if dm.Responsibility != model.CustomerDamage {
		t.Fatalf("expected customer_damage when faults don't match pre-condition items, got %s (note: %s)", dm.Responsibility, dm.ResponsibilityNote)
	}
}

func TestDeductAccessorySeparate(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryCamera,
		Brand:          "Sony",
		Model:          "A7IV",
		PreBorrowPhoto: "before.jpg",
	})

	svc.AddAccessoryPrice(AddAccessoryPriceInput{
		EquipmentID: eq.ID,
		Name:        "lens_cap",
		Price:       50,
	})
	svc.AddAccessoryPrice(AddAccessoryPriceInput{
		EquipmentID: eq.ID,
		Name:        "battery",
		Price:       300,
	})
	svc.AddAccessoryPrice(AddAccessoryPriceInput{
		EquipmentID: eq.ID,
		Name:        "softbox",
		Price:       150,
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:   eq.ID,
		CustomerName:  "Dave",
		CustomerPhone: "444",
		StudioPosition: "H-1",
		Deposit:       1000,
		PreBorrowPhotos: []string{"before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	deduction, err := svc.DeductAccessory(DeductAccessoryInput{
		BorrowRecordID: br.ID,
		AccessoryNames: []string{"lens_cap", "battery"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if deduction.DeductAmount != 350 {
		t.Fatalf("expected deduct_amount 350 (50+300), got %.2f", deduction.DeductAmount)
	}
	if deduction.RefundAmount != 650 {
		t.Fatalf("expected refund_amount 650 (1000-350), got %.2f", deduction.RefundAmount)
	}
	if len(deduction.AccessoryItems) != 2 {
		t.Fatalf("expected 2 accessory items, got %d", len(deduction.AccessoryItems))
	}
	if deduction.RepairQuoteID != "" {
		t.Fatalf("expected empty repair_quote_id for accessory deduction, got %s", deduction.RepairQuoteID)
	}

	found := false
	for _, item := range deduction.AccessoryItems {
		if item.AccessoryName == "lens_cap" && item.Price == 50 {
			found = true
		}
	}
	if !found {
		t.Fatal("expected lens_cap accessory item with price 50")
	}
}

func TestDeductAccessoryExceedsDeposit(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryCamera,
		Brand:          "Sony",
		Model:          "A7IV",
		PreBorrowPhoto: "before.jpg",
	})

	svc.AddAccessoryPrice(AddAccessoryPriceInput{
		EquipmentID: eq.ID,
		Name:        "battery",
		Price:       800,
	})
	svc.AddAccessoryPrice(AddAccessoryPriceInput{
		EquipmentID: eq.ID,
		Name:        "carrying_bag",
		Price:       500,
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:   eq.ID,
		CustomerName:  "Eve",
		CustomerPhone: "555",
		StudioPosition: "H-2",
		Deposit:       1000,
		PreBorrowPhotos: []string{"before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	deduction, err := svc.DeductAccessory(DeductAccessoryInput{
		BorrowRecordID: br.ID,
		AccessoryNames: []string{"battery", "carrying_bag"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if deduction.DeductAmount != 1000 {
		t.Fatalf("expected deduct_amount 1000 (capped at deposit), got %.2f", deduction.DeductAmount)
	}
	if deduction.RefundAmount != 0 {
		t.Fatalf("expected refund_amount 0, got %.2f", deduction.RefundAmount)
	}
}

func TestAppealFreezesDeposit(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryLens,
		Brand:          "Canon",
		Model:          "EF 24-70mm",
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Frank",
		CustomerPhone:   "666",
		StudioPosition:  "I-1",
		Deposit:         2000,
		PreBorrowPhotos: []string{"before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints:    []model.FaultPoint{{Location: "glass", Description: "crack", Severity: "severe"}},
		ReturnPhotos:   []string{"damage.jpg"},
	})

	if dm.Responsibility != model.CustomerDamage {
		t.Fatalf("expected customer_damage, got %s", dm.Responsibility)
	}

	rq, _ := svc.CreateRepairQuote(CreateRepairQuoteInput{
		DamageReportID: dm.ID,
		RepairCost:     800,
		LaborCost:      200,
		Description:    "repair",
	})

	appeal, err := svc.CreateAppeal(CreateAppealInput{
		BorrowRecordID: br.ID,
		CustomerName:   "Frank",
		Reason:         "I disagree with the responsibility conclusion",
		Evidence:       []string{"my_evidence.jpg"},
	})
	if err != nil {
		t.Fatalf("unexpected error creating appeal: %v", err)
	}

	brAfter, _ := svc.GetBorrowRecord(br.ID)
	if !brAfter.DepositFrozen {
		t.Fatal("expected deposit to be frozen after appeal")
	}
	if brAfter.Status != model.BorrowDepositFrozen {
		t.Fatalf("expected borrow status to be deposit_frozen, got %s", brAfter.Status)
	}
	if brAfter.DepositFrozenReason == "" {
		t.Fatal("expected deposit_frozen_reason to be set")
	}
	if appeal.Status != model.AppealPending {
		t.Fatalf("expected appeal status pending, got %s", appeal.Status)
	}

	_, err = svc.DeductDeposit(DeductDepositInput{
		BorrowRecordID: br.ID,
		RepairQuoteID:  rq.ID,
	})
	if err == nil {
		t.Fatal("expected error when deducting frozen deposit")
	}
}

func TestSupplementalEvidence(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryCamera,
		Brand:          "Nikon",
		Model:          "Z6",
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Grace",
		CustomerPhone:   "777",
		StudioPosition:  "J-1",
		Deposit:         3000,
		PreBorrowPhotos: []string{"before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints:    []model.FaultPoint{{Location: "body", Description: "scratch", Severity: "severe"}},
		ReturnPhotos:   []string{"damage.jpg"},
	})
	_ = dm

	appeal, _ := svc.CreateAppeal(CreateAppealInput{
		BorrowRecordID: br.ID,
		CustomerName:   "Grace",
		Reason:         "disagree",
		Evidence:       []string{"cust_ev.jpg"},
	})

	se, err := svc.AddSupplementalEvidence(AddSupplementalEvidenceInput{
		AppealID:       appeal.ID,
		BorrowRecordID: br.ID,
		OperatorName:   "Operator_01",
		EvidenceType:   model.EvidenceSurveillance,
		Description:    "监控录像显示客户使用期间有摔落",
		Attachments:    []string{"surveillance_clip1.mp4", "surveillance_clip2.mp4"},
	})
	if err != nil {
		t.Fatalf("unexpected error adding supplemental evidence: %v", err)
	}
	if se.EvidenceType != model.EvidenceSurveillance {
		t.Fatalf("expected evidence_type surveillance, got %s", se.EvidenceType)
	}
	if len(se.Attachments) != 2 {
		t.Fatalf("expected 2 attachments, got %d", len(se.Attachments))
	}
	if se.OperatorName != "Operator_01" {
		t.Fatalf("expected operator_name Operator_01, got %s", se.OperatorName)
	}

	se2, _ := svc.AddSupplementalEvidence(AddSupplementalEvidenceInput{
		BorrowRecordID: br.ID,
		OperatorName:   "Operator_02",
		EvidenceType:   model.EvidenceHandoverPhoto,
		Description:    "交接照片显示器材完好",
		Attachments:    []string{"handover_1.jpg"},
	})
	if se2.AppealID != "" {
		t.Fatalf("expected empty appeal_id, got %s", se2.AppealID)
	}

	allByAppeal := svc.ListSupplementalEvidences(appeal.ID, "")
	if len(allByAppeal) != 1 {
		t.Fatalf("expected 1 supplemental evidence by appeal, got %d", len(allByAppeal))
	}

	allByBorrow := svc.ListSupplementalEvidences("", br.ID)
	if len(allByBorrow) != 2 {
		t.Fatalf("expected 2 supplemental evidences by borrow, got %d", len(allByBorrow))
	}
}

func TestUpdateRepairQuoteAndAdditionalCompensation(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryCamera,
		Brand:          "Sony",
		Model:          "A7S III",
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Henry",
		CustomerPhone:   "888",
		StudioPosition:  "K-1",
		Deposit:         2000,
		PreBorrowPhotos: []string{"before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints:    []model.FaultPoint{{Location: "sensor", Description: "damaged", Severity: "severe"}},
		ReturnPhotos:   []string{"damage.jpg"},
	})

	rq, _ := svc.CreateRepairQuote(CreateRepairQuoteInput{
		DamageReportID: dm.ID,
		RepairCost:     800,
		LaborCost:      200,
		Description:    "initial quote",
	})

	if rq.TotalCost != 1000 {
		t.Fatalf("expected initial total_cost 1000, got %.2f", rq.TotalCost)
	}

	_, err := svc.DeductDeposit(DeductDepositInput{
		BorrowRecordID: br.ID,
		RepairQuoteID:  rq.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error deducting deposit: %v", err)
	}

	updatedRq, comp, err := svc.UpdateRepairQuote(UpdateRepairQuoteInput{
		RepairQuoteID: rq.ID,
		RepairCost:    2500,
		LaborCost:     800,
		Description:   "shop callback: sensor replacement + board repair",
		UpdateNote:    "维修店回传最新报价，发现主板也需更换",
	})
	if err != nil {
		t.Fatalf("unexpected error updating repair quote: %v", err)
	}

	if updatedRq.TotalCost != 3300 {
		t.Fatalf("expected updated total_cost 3300, got %.2f", updatedRq.TotalCost)
	}
	if !updatedRq.IsUpdated {
		t.Fatal("expected is_updated to be true")
	}
	if updatedRq.UpdatedAt == nil {
		t.Fatal("expected updated_at to be set")
	}
	if updatedRq.UpdateNote == "" {
		t.Fatal("expected update_note to be set")
	}

	if comp == nil {
		t.Fatal("expected additional compensation to be created")
	}
	if comp.AdditionalAmount != 1300 {
		t.Fatalf("expected additional_amount 1300 (3300-2000), got %.2f", comp.AdditionalAmount)
	}
	if comp.DepositAmount != 2000 {
		t.Fatalf("expected deposit_amount 2000, got %.2f", comp.DepositAmount)
	}
	if comp.Status != "pending" {
		t.Fatalf("expected status pending, got %s", comp.Status)
	}

	comps := svc.ListAdditionalCompensations(br.ID)
	if len(comps) != 1 {
		t.Fatalf("expected 1 additional compensation by borrow, got %d", len(comps))
	}

	compCollected, err := svc.CollectAdditionalCompensation(CollectAdditionalCompensationInput{
		CompensationID: comp.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error collecting compensation: %v", err)
	}
	if compCollected.Status != "collected" {
		t.Fatalf("expected status collected, got %s", compCollected.Status)
	}
	if compCollected.CollectedAt == nil {
		t.Fatal("expected collected_at to be set")
	}

	_, err = svc.CollectAdditionalCompensation(CollectAdditionalCompensationInput{
		CompensationID: comp.ID,
	})
	if err == nil {
		t.Fatal("expected error when collecting already collected compensation")
	}
}

func TestRepairQuoteExceedsDepositCreatesAdditionalCompensation(t *testing.T) {
	s := store.New()
	svc := New(s)

	eq, _ := svc.CreateEquipment(CreateEquipmentInput{
		Category:       model.CategoryLens,
		Brand:          "Sigma",
		Model:          "35mm Art",
		PreBorrowPhoto: "before.jpg",
	})

	br, _ := svc.BorrowEquipment(BorrowInput{
		EquipmentID:     eq.ID,
		CustomerName:    "Ivy",
		CustomerPhone:   "999",
		StudioPosition:  "L-1",
		Deposit:         1500,
		PreBorrowPhotos: []string{"before.jpg"},
	})

	svc.ReturnInspection(ReturnInspectionInput{
		BorrowRecordID: br.ID,
		ReturnPhotos:   []string{"return.jpg"},
	})

	dm, _ := svc.RegisterDamage(RegisterDamageInput{
		BorrowRecordID: br.ID,
		FaultPoints:    []model.FaultPoint{{Location: "mount", Description: "broken", Severity: "severe"}},
		ReturnPhotos:   []string{"damage.jpg"},
	})

	rq, _ := svc.CreateRepairQuote(CreateRepairQuoteInput{
		DamageReportID: dm.ID,
		RepairCost:     1800,
		LaborCost:      400,
		Description:    "replace mount assembly",
	})

	if rq.TotalCost != 2200 {
		t.Fatalf("expected total_cost 2200, got %.2f", rq.TotalCost)
	}

	deduction, err := svc.DeductDeposit(DeductDepositInput{
		BorrowRecordID: br.ID,
		RepairQuoteID:  rq.ID,
	})
	if err != nil {
		t.Fatalf("unexpected error deducting deposit: %v", err)
	}

	if deduction.DeductAmount != 1500 {
		t.Fatalf("expected deduct_amount 1500 (full deposit), got %.2f", deduction.DeductAmount)
	}
	if deduction.RefundAmount != 0 {
		t.Fatalf("expected refund_amount 0, got %.2f", deduction.RefundAmount)
	}

	comp := svc.store.FindAdditionalCompensationByRepairQuote(rq.ID)
	if comp == nil {
		t.Fatal("expected additional compensation record to be created automatically")
	}
	if comp.AdditionalAmount != 700 {
		t.Fatalf("expected additional_amount 700 (2200-1500), got %.2f", comp.AdditionalAmount)
	}
	if comp.OriginalDeduct != 1500 {
		t.Fatalf("expected original_deduct 1500, got %.2f", comp.OriginalDeduct)
	}
	if comp.Status != "pending" {
		t.Fatalf("expected status pending, got %s", comp.Status)
	}
}
