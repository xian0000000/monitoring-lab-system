package models

import "time"

// Lab adalah struktur data untuk tabel labs
type Lab struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Capacity int    `json:"capacity"`
}

// Device adalah struktur data untuk tabel devices
type Device struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	IPAddress string     `json:"ip_address"`
	LabID     int        `json:"lab_id"`
	Status    string     `json:"status"`
	LastSeen  *time.Time `json:"last_seen"`
}