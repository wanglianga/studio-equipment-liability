package store

import (
	"fmt"
	"sync"
	"time"

	"studio-equipment-manager/model"
)

type Store struct {
	mu                      sync.RWMutex
	equipments              map[string]*model.Equipment
	borrowRecords           map[string]*model.BorrowRecord
	damageReports           map[string]*model.DamageReport
	repairQuotes            map[string]*model.RepairQuote
	deductionRecords        map[string]*model.DeductionRecord
	appeals                 map[string]*model.Appeal
	accessoryPrices         map[string]*model.AccessoryPrice
	supplementalEvidences   map[string]*model.SupplementalEvidence
	additionalCompensations map[string]*model.AdditionalCompensation
	equipCounter            int
	borrowCounter           int
	damageCounter           int
	repairCounter           int
	deductCounter           int
	appealCounter           int
	accessoryCounter        int
	supplementalCounter     int
	additionalCompCounter   int
}

func New() *Store {
	return &Store{
		equipments:              make(map[string]*model.Equipment),
		borrowRecords:           make(map[string]*model.BorrowRecord),
		damageReports:           make(map[string]*model.DamageReport),
		repairQuotes:            make(map[string]*model.RepairQuote),
		deductionRecords:        make(map[string]*model.DeductionRecord),
		appeals:                 make(map[string]*model.Appeal),
		accessoryPrices:         make(map[string]*model.AccessoryPrice),
		supplementalEvidences:   make(map[string]*model.SupplementalEvidence),
		additionalCompensations: make(map[string]*model.AdditionalCompensation),
	}
}

func (s *Store) SaveEquipment(eq *model.Equipment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.equipments[eq.ID] = eq
}

func (s *Store) GetEquipment(id string) (*model.Equipment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	eq, ok := s.equipments[id]
	return eq, ok
}

func (s *Store) ListEquipments() []*model.Equipment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.Equipment, 0, len(s.equipments))
	for _, eq := range s.equipments {
		result = append(result, eq)
	}
	return result
}

func (s *Store) NextEquipID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.equipCounter++
	return fmt.Sprintf("EQ-%04d", s.equipCounter)
}

func (s *Store) SaveBorrowRecord(r *model.BorrowRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.borrowRecords[r.ID] = r
}

func (s *Store) GetBorrowRecord(id string) (*model.BorrowRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.borrowRecords[id]
	return r, ok
}

func (s *Store) ListBorrowRecords() []*model.BorrowRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.BorrowRecord, 0, len(s.borrowRecords))
	for _, r := range s.borrowRecords {
		result = append(result, r)
	}
	return result
}

func (s *Store) FindActiveBorrowByEquipment(equipID string) *model.BorrowRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.borrowRecords {
		if r.EquipmentID == equipID && r.Status == model.BorrowActive {
			return r
		}
	}
	return nil
}

func (s *Store) NextBorrowID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.borrowCounter++
	return fmt.Sprintf("BR-%04d", s.borrowCounter)
}

func (s *Store) SaveDamageReport(r *model.DamageReport) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.damageReports[r.ID] = r
}

func (s *Store) GetDamageReport(id string) (*model.DamageReport, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.damageReports[id]
	return r, ok
}

func (s *Store) ListDamageReports() []*model.DamageReport {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.DamageReport, 0, len(s.damageReports))
	for _, r := range s.damageReports {
		result = append(result, r)
	}
	return result
}

func (s *Store) FindDamageReportByBorrow(borrowID string) *model.DamageReport {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.damageReports {
		if r.BorrowRecordID == borrowID {
			return r
		}
	}
	return nil
}

func (s *Store) NextDamageID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.damageCounter++
	return fmt.Sprintf("DM-%04d", s.damageCounter)
}

func (s *Store) SaveRepairQuote(q *model.RepairQuote) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repairQuotes[q.ID] = q
}

func (s *Store) GetRepairQuote(id string) (*model.RepairQuote, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.repairQuotes[id]
	return q, ok
}

func (s *Store) FindRepairQuoteByDamage(damageID string) *model.RepairQuote {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, q := range s.repairQuotes {
		if q.DamageReportID == damageID {
			return q
		}
	}
	return nil
}

func (s *Store) NextRepairID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repairCounter++
	return fmt.Sprintf("RQ-%04d", s.repairCounter)
}

func (s *Store) SaveDeductionRecord(d *model.DeductionRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deductionRecords[d.ID] = d
}

func (s *Store) GetDeductionRecord(id string) (*model.DeductionRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.deductionRecords[id]
	return d, ok
}

func (s *Store) FindDeductionByBorrow(borrowID string) *model.DeductionRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, d := range s.deductionRecords {
		if d.BorrowRecordID == borrowID {
			return d
		}
	}
	return nil
}

func (s *Store) NextDeductionID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deductCounter++
	return fmt.Sprintf("DD-%04d", s.deductCounter)
}

func (s *Store) SaveAppeal(a *model.Appeal) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.appeals[a.ID] = a
}

func (s *Store) GetAppeal(id string) (*model.Appeal, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.appeals[id]
	return a, ok
}

func (s *Store) ListAppeals() []*model.Appeal {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.Appeal, 0, len(s.appeals))
	for _, a := range s.appeals {
		result = append(result, a)
	}
	return result
}

func (s *Store) FindAppealByBorrow(borrowID string) *model.Appeal {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.appeals {
		if a.BorrowRecordID == borrowID {
			return a
		}
	}
	return nil
}

func (s *Store) NextAppealID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.appealCounter++
	return fmt.Sprintf("AP-%04d", s.appealCounter)
}

func (s *Store) HasPreviousDamage(equipID string, beforeTime time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.damageReports {
		if r.EquipmentID == equipID && !r.CreatedAt.After(beforeTime) {
			return true
		}
	}
	return false
}

func (s *Store) HasPreviousUndetermined(equipID string, beforeTime time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.damageReports {
		if r.EquipmentID == equipID && !r.CreatedAt.After(beforeTime) && r.Responsibility == model.Undetermined {
			return true
		}
	}
	return false
}

func (s *Store) CountBorrowByEquipment(equipID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, r := range s.borrowRecords {
		if r.EquipmentID == equipID {
			count++
		}
	}
	return count
}

func (s *Store) SaveAccessoryPrice(ap *model.AccessoryPrice) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accessoryPrices[ap.ID] = ap
}

func (s *Store) GetAccessoryPrice(id string) (*model.AccessoryPrice, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ap, ok := s.accessoryPrices[id]
	return ap, ok
}

func (s *Store) FindAccessoryPriceByEquipAndName(equipID, name string) *model.AccessoryPrice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ap := range s.accessoryPrices {
		if ap.EquipmentID == equipID && ap.Name == name {
			return ap
		}
	}
	return nil
}

func (s *Store) ListAccessoryPricesByEquipment(equipID string) []*model.AccessoryPrice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.AccessoryPrice, 0)
	for _, ap := range s.accessoryPrices {
		if ap.EquipmentID == equipID {
			result = append(result, ap)
		}
	}
	return result
}

func (s *Store) ListAccessoryPrices() []*model.AccessoryPrice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.AccessoryPrice, 0, len(s.accessoryPrices))
	for _, ap := range s.accessoryPrices {
		result = append(result, ap)
	}
	return result
}

func (s *Store) NextAccessoryID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accessoryCounter++
	return fmt.Sprintf("AC-%04d", s.accessoryCounter)
}

func (s *Store) SaveSupplementalEvidence(se *model.SupplementalEvidence) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.supplementalEvidences[se.ID] = se
}

func (s *Store) GetSupplementalEvidence(id string) (*model.SupplementalEvidence, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	se, ok := s.supplementalEvidences[id]
	return se, ok
}

func (s *Store) ListSupplementalEvidences() []*model.SupplementalEvidence {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.SupplementalEvidence, 0, len(s.supplementalEvidences))
	for _, se := range s.supplementalEvidences {
		result = append(result, se)
	}
	return result
}

func (s *Store) FindSupplementalEvidenceByAppeal(appealID string) []*model.SupplementalEvidence {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.SupplementalEvidence, 0)
	for _, se := range s.supplementalEvidences {
		if se.AppealID == appealID {
			result = append(result, se)
		}
	}
	return result
}

func (s *Store) FindSupplementalEvidenceByBorrow(borrowID string) []*model.SupplementalEvidence {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.SupplementalEvidence, 0)
	for _, se := range s.supplementalEvidences {
		if se.BorrowRecordID == borrowID {
			result = append(result, se)
		}
	}
	return result
}

func (s *Store) NextSupplementalEvidenceID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.supplementalCounter++
	return fmt.Sprintf("SE-%04d", s.supplementalCounter)
}

func (s *Store) SaveAdditionalCompensation(ac *model.AdditionalCompensation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.additionalCompensations[ac.ID] = ac
}

func (s *Store) GetAdditionalCompensation(id string) (*model.AdditionalCompensation, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ac, ok := s.additionalCompensations[id]
	return ac, ok
}

func (s *Store) ListAdditionalCompensations() []*model.AdditionalCompensation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.AdditionalCompensation, 0, len(s.additionalCompensations))
	for _, ac := range s.additionalCompensations {
		result = append(result, ac)
	}
	return result
}

func (s *Store) FindAdditionalCompensationByBorrow(borrowID string) []*model.AdditionalCompensation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.AdditionalCompensation, 0)
	for _, ac := range s.additionalCompensations {
		if ac.BorrowRecordID == borrowID {
			result = append(result, ac)
		}
	}
	return result
}

func (s *Store) FindAdditionalCompensationByRepairQuote(repairQuoteID string) *model.AdditionalCompensation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ac := range s.additionalCompensations {
		if ac.RepairQuoteID == repairQuoteID {
			return ac
		}
	}
	return nil
}

func (s *Store) NextAdditionalCompensationID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.additionalCompCounter++
	return fmt.Sprintf("XC-%04d", s.additionalCompCounter)
}

func (s *Store) UpdateRepairQuote(q *model.RepairQuote) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repairQuotes[q.ID] = q
}
