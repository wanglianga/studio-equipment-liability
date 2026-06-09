package store

import (
	"fmt"
	"sync"
	"time"

	"studio-equipment-manager/model"
)

type Store struct {
	mu               sync.RWMutex
	equipments       map[string]*model.Equipment
	borrowRecords    map[string]*model.BorrowRecord
	damageReports    map[string]*model.DamageReport
	repairQuotes     map[string]*model.RepairQuote
	deductionRecords map[string]*model.DeductionRecord
	appeals          map[string]*model.Appeal
	equipCounter     int
	borrowCounter    int
	damageCounter    int
	repairCounter    int
	deductCounter    int
	appealCounter    int
}

func New() *Store {
	return &Store{
		equipments:       make(map[string]*model.Equipment),
		borrowRecords:    make(map[string]*model.BorrowRecord),
		damageReports:    make(map[string]*model.DamageReport),
		repairQuotes:     make(map[string]*model.RepairQuote),
		deductionRecords: make(map[string]*model.DeductionRecord),
		appeals:          make(map[string]*model.Appeal),
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
		if r.EquipmentID == equipID && r.CreatedAt.Before(beforeTime) {
			return true
		}
	}
	return false
}

func (s *Store) HasPreviousUndetermined(equipID string, beforeTime time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.damageReports {
		if r.EquipmentID == equipID && r.CreatedAt.Before(beforeTime) && r.Responsibility == model.Undetermined {
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
