package main

import (
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"
)

func (s *Store) createSumberDaya(w http.ResponseWriter, r *http.Request) {
	var input sumberDayaInput
	if !decodeJSON(w, r, &input) {
		return
	}
	if input.Nama == nil || input.Jenis == nil || input.Stok == nil {
		writeError(w, http.StatusBadRequest, "nama, jenis, dan stok wajib diisi.")
		return
	}
	nama := strings.TrimSpace(*input.Nama)
	jenis := strings.TrimSpace(*input.Jenis)
	stok := *input.Stok
	if nama == "" {
		writeError(w, http.StatusBadRequest, "nama tidak boleh kosong.")
		return
	}
	if utf8.RuneCountInString(nama) > 50 {
		writeError(w, http.StatusBadRequest, "nama maksimal 50 karakter.")
		return
	}
	if jenis != "ruang" && jenis != "alat" {
		writeError(w, http.StatusBadRequest, "jenis harus ruang atau alat.")
		return
	}
	if stok < 1 {
		writeError(w, http.StatusBadRequest, "stok minimal 1.")
		return
	}
	s.mu.Lock()
	item := SumberDaya{
		ID:    s.nextSumberDayaID,
		Nama:  nama,
		Jenis: jenis,
		Stok:  stok,
	}
	s.nextSumberDayaID++
	s.sumberDaya = append(s.sumberDaya, item)
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, item)
}

func (s *Store) listSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	items := make([]SumberDaya, len(s.sumberDaya))
	copy(items, s.sumberDaya)
	s.mu.Unlock()
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
	writeJSON(w, http.StatusOK, items)
}

func (s *Store) getSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	item, found := s.findSumberDayaLocked(id)
	s.mu.Unlock()
	if !found {
		writeError(w, http.StatusNotFound, "Sumber daya tidak ditemukan.")
		return
	}
	writeJSON(w, http.StatusOK, item)
}
