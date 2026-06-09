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
