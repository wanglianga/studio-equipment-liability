package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"studio-equipment-manager/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	var input service.CreateEquipmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.Category == "" {
		writeError(w, http.StatusBadRequest, "category is required")
		return
	}
	if input.Brand == "" || input.Model == "" {
		writeError(w, http.StatusBadRequest, "brand and model are required")
		return
	}

	eq, err := h.svc.CreateEquipment(input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, eq)
}

func (h *Handler) GetEquipment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "equipment id is required")
		return
	}

	eq, err := h.svc.GetEquipment(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, eq)
}

func (h *Handler) ListEquipments(w http.ResponseWriter, r *http.Request) {
	eqs := h.svc.ListEquipments()
	writeJSON(w, http.StatusOK, eqs)
}

func (h *Handler) BorrowEquipment(w http.ResponseWriter, r *http.Request) {
	var input service.BorrowInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.EquipmentID == "" {
		writeError(w, http.StatusBadRequest, "equipment_id is required")
		return
	}
	if input.CustomerName == "" {
		writeError(w, http.StatusBadRequest, "customer_name is required")
		return
	}
	if input.Deposit <= 0 {
		writeError(w, http.StatusBadRequest, "deposit must be positive")
		return
	}

	record, err := h.svc.BorrowEquipment(input)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (h *Handler) ReturnInspection(w http.ResponseWriter, r *http.Request) {
	var input service.ReturnInspectionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.BorrowRecordID == "" {
		writeError(w, http.StatusBadRequest, "borrow_record_id is required")
		return
	}

	record, err := h.svc.ReturnInspection(input)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (h *Handler) RegisterDamage(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterDamageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.BorrowRecordID == "" {
		writeError(w, http.StatusBadRequest, "borrow_record_id is required")
		return
	}
	if len(input.FaultPoints) == 0 {
		writeError(w, http.StatusBadRequest, "fault_points must not be empty")
		return
	}

	report, err := h.svc.RegisterDamage(input)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, report)
}

func (h *Handler) CreateRepairQuote(w http.ResponseWriter, r *http.Request) {
	var input service.CreateRepairQuoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.DamageReportID == "" {
		writeError(w, http.StatusBadRequest, "damage_report_id is required")
		return
	}
	if input.RepairCost < 0 || input.LaborCost < 0 {
		writeError(w, http.StatusBadRequest, "repair_cost and labor_cost must not be negative")
		return
	}

	quote, err := h.svc.CreateRepairQuote(input)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, quote)
}

func (h *Handler) DeductDeposit(w http.ResponseWriter, r *http.Request) {
	var input service.DeductDepositInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.BorrowRecordID == "" || input.RepairQuoteID == "" {
		writeError(w, http.StatusBadRequest, "borrow_record_id and repair_quote_id are required")
		return
	}

	deduction, err := h.svc.DeductDeposit(input)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, deduction)
}

func (h *Handler) CreateAppeal(w http.ResponseWriter, r *http.Request) {
	var input service.CreateAppealInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.BorrowRecordID == "" {
		writeError(w, http.StatusBadRequest, "borrow_record_id is required")
		return
	}
	if input.CustomerName == "" {
		writeError(w, http.StatusBadRequest, "customer_name is required")
		return
	}
	if input.Reason == "" {
		writeError(w, http.StatusBadRequest, "reason is required")
		return
	}

	appeal, err := h.svc.CreateAppeal(input)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, appeal)
}

func (h *Handler) ReviewAppeal(w http.ResponseWriter, r *http.Request) {
	var input service.ReviewAppealInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.AppealID == "" {
		writeError(w, http.StatusBadRequest, "appeal_id is required")
		return
	}

	appeal, err := h.svc.ReviewAppeal(input)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, appeal)
}

func (h *Handler) ListBorrowRecords(w http.ResponseWriter, r *http.Request) {
	records := h.svc.ListBorrowRecords()
	writeJSON(w, http.StatusOK, records)
}

func (h *Handler) GetBorrowRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "borrow record id is required")
		return
	}

	record, err := h.svc.GetBorrowRecord(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (h *Handler) ListDamageReports(w http.ResponseWriter, r *http.Request) {
	reports := h.svc.ListDamageReports()
	writeJSON(w, http.StatusOK, reports)
}

func (h *Handler) GetDamageReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "damage report id is required")
		return
	}

	report, err := h.svc.GetDamageReport(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) CompleteRepair(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "equipment id is required")
		return
	}

	eq, err := h.svc.CompleteRepair(id)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, eq)
}

func (h *Handler) ListAppeals(w http.ResponseWriter, r *http.Request) {
	appeals := h.svc.ListAppeals()
	writeJSON(w, http.StatusOK, appeals)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/equipments", h.CreateEquipment)
	mux.HandleFunc("GET /api/equipments", h.ListEquipments)
	mux.HandleFunc("GET /api/equipments/{id}", h.GetEquipment)
	mux.HandleFunc("POST /api/borrow", h.BorrowEquipment)
	mux.HandleFunc("POST /api/borrow/return", h.ReturnInspection)
	mux.HandleFunc("GET /api/borrow", h.ListBorrowRecords)
	mux.HandleFunc("GET /api/borrow/{id}", h.GetBorrowRecord)
	mux.HandleFunc("POST /api/damage", h.RegisterDamage)
	mux.HandleFunc("GET /api/damage", h.ListDamageReports)
	mux.HandleFunc("GET /api/damage/{id}", h.GetDamageReport)
	mux.HandleFunc("POST /api/repair-quote", h.CreateRepairQuote)
	mux.HandleFunc("POST /api/repair-complete/{id}", h.CompleteRepair)
	mux.HandleFunc("POST /api/deduction", h.DeductDeposit)
	mux.HandleFunc("POST /api/appeal", h.CreateAppeal)
	mux.HandleFunc("POST /api/appeal/review", h.ReviewAppeal)
	mux.HandleFunc("GET /api/appeal", h.ListAppeals)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "studio-equipment-manager",
	})
	fmt.Println("[HealthCheck] service is running")
}
