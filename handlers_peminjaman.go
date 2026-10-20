package main

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (s *Store) createPeminjaman(w http.ResponseWriter, r *http.Request) {
	var input peminjamanInput
	if !decodeJSON(w, r, &input) {
		return
	}
	if input.SumberDayaID == nil ||
		input.Peminjam == nil ||
		input.Tanggal == nil ||
		input.SlotMulai == nil ||
		input.SlotSelesai == nil ||
		input.Jumlah == nil {
		writeError(w, http.StatusBadRequest, "Semua field peminjaman wajib diisi.")
		return
	}
	sumberDayaID := *input.SumberDayaID
	peminjam := strings.TrimSpace(*input.Peminjam)
	tanggal := strings.TrimSpace(*input.Tanggal)
	slotMulai := *input.SlotMulai
	slotSelesai := *input.SlotSelesai
	jumlah := *input.Jumlah
	if peminjam == "" {
		writeError(w, http.StatusBadRequest, "peminjam tidak boleh kosong.")
		return
	}
	if utf8.RuneCountInString(peminjam) > 50 {
		writeError(w, http.StatusBadRequest, "peminjam maksimal 50 karakter.")
		return
	}
	if !validDate(tanggal) {
		writeError(w, http.StatusBadRequest, "tanggal harus valid dengan format YYYY-MM-DD.")
		return
	}
	if slotMulai < 1 || slotMulai > 8 || slotSelesai < 1 || slotSelesai > 8 {
		writeError(w, http.StatusBadRequest, "slot harus berada di antara 1 dan 8.")
		return
	}
	if slotMulai > slotSelesai {
		writeError(w, http.StatusBadRequest, "slot_mulai tidak boleh lebih besar dari slot_selesai.")
		return
	}

	s.mu.Lock()
	sumberDaya, found := s.findSumberDayaLocked(sumberDayaID)
	if !found {
		s.mu.Unlock()
		writeError(w, http.StatusBadRequest, "sumber_daya_id tidak menunjuk sumber daya yang ada.")
		return
	}
	if jumlah < 1 || jumlah > sumberDaya.Stok {
		s.mu.Unlock()
		writeError(w, http.StatusBadRequest, "jumlah harus antara 1 dan stok sumber daya.")
		return
	}
	item := Peminjaman{
		ID:           s.nextPeminjamanID,
		SumberDayaID: sumberDayaID,
		Peminjam:     peminjam,
		Tanggal:      tanggal,
		SlotMulai:    slotMulai,
		SlotSelesai:  slotSelesai,
		Jumlah:       jumlah,
		Status:       "aktif",
	}
	s.nextPeminjamanID++
	s.peminjaman = append(s.peminjaman, item)
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, item)
}

func (s *Store) listPeminjaman(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var sumberDayaFilter *int
	if values, exists := query["sumber_daya_id"]; exists {
		rawID := ""
		if len(values) > 0 {
			rawID = values[0]
		}
		id, err := strconv.Atoi(rawID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id harus berupa angka.")
			return
		}
		sumberDayaFilter = &id
	}
	tanggalFilter := strings.TrimSpace(query.Get("tanggal"))
	statusFilter := strings.TrimSpace(query.Get("status"))
	s.mu.Lock()
	items := make([]Peminjaman, 0)
	for _, item := range s.peminjaman {
		if sumberDayaFilter != nil && item.SumberDayaID != *sumberDayaFilter {
			continue
		}
		if tanggalFilter != "" && item.Tanggal != tanggalFilter {
			continue
		}
		if statusFilter != "" && item.Status != statusFilter {
			continue
		}
		items = append(items, item)
	}
	s.mu.Unlock()
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
	writeJSON(w, http.StatusOK, items)
}

func (s *Store) getPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	item, found := s.findPeminjamanLocked(id)
	s.mu.Unlock()
	if !found {
		writeError(w, http.StatusNotFound, "Peminjaman tidak ditemukan.")
		return
	}
	writeJSON(w, http.StatusOK, item)
}
