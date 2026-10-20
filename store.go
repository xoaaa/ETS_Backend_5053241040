package main

import (
	"net/http"
	"sync"
)

type Store struct {
	mu               sync.Mutex
	sumberDaya       []SumberDaya
	peminjaman       []Peminjaman
	nextSumberDayaID int
	nextPeminjamanID int
}

func NewHandler() http.Handler {
	store := &Store{
		sumberDaya:       make([]SumberDaya, 0),
		peminjaman:       make([]Peminjaman, 0),
		nextSumberDayaID: 1,
		nextPeminjamanID: 1,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /sumber-daya", store.createSumberDaya)
	mux.HandleFunc("GET /sumber-daya", store.listSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", store.getSumberDaya)
	mux.HandleFunc("POST /peminjaman", store.createPeminjaman)
	mux.HandleFunc("GET /peminjaman", store.listPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", store.getPeminjaman)
	return mux
}
func (s *Store) findSumberDayaLocked(id int) (SumberDaya, bool) {
	for _, item := range s.sumberDaya {
		if item.ID == id {
			return item, true
		}
	}
	return SumberDaya{}, false
}

func (s *Store) findPeminjamanLocked(id int) (Peminjaman, bool) {
	for _, item := range s.peminjaman {
		if item.ID == id {
			return item, true
		}
	}
	return Peminjaman{}, false
}
