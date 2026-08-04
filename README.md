# Monitoring Lab System

Program monitoring status perangkat jaringan di laboratorium komputer berbasis **ping (ICMP)**. Sistem melakukan scan otomatis ke seluruh IP dalam satu subnet `/24` (maksimal 254 host) secara berkala, lalu menampilkan status *online/offline* tiap perangkat secara real-time di dashboard web.

## Fitur

- Scan otomatis subnet `/24` setiap 30 detik (background scheduler).
- Scan manual on-demand via endpoint API.
- Update status real-time ke frontend menggunakan **Server-Sent Events (SSE)**.
- Pencatatan waktu terakhir perangkat terlihat online (`last_seen`).
- Manajemen data lab (nama, lokasi, kapasitas).
- Manajemen data perangkat (tambah, lihat, hapus).
- Siap dijalankan dengan Docker Compose (database, backend, frontend).

## Arsitektur & Stack

| Komponen  | Teknologi                                  |
| --------- | ------------------------------------------- |
| Frontend  | Angular 22                                   |
| Backend   | Go + Gin, library ping `go-ping/ping`        |
| Database  | PostgreSQL 16                                |
| Deployment| Docker & Docker Compose                      |

Alur kerja singkat:

1. `scheduler` di backend memindai seluruh IP host dalam `NETWORK_CIDR` menggunakan ICMP ping (mode unprivileged).
2. Hasil setiap ping (`online`/`offline`) di-*upsert* ke tabel `devices` di PostgreSQL berdasarkan `ip_address`.
3. Hasil scan juga disiarkan (broadcast) ke semua client yang terhubung lewat SSE (`/api/scan/stream`).
4. Frontend Angular menampilkan data lab & perangkat, serta mendengarkan stream SSE untuk update status secara real-time.

## Struktur Proyek

```
monitoring-lab-system/
├── docker-compose.yml
├── lab-monitor-backend/     # Go + Gin API & ping scheduler
│   ├── db/                  # Koneksi & migrasi database
│   ├── handlers/            # HTTP handler (lab, device, scan)
│   ├── models/               # Struct data (Lab, Device)
│   ├── scheduler/            # Logic scan subnet & SSE hub
│   └── main.go
└── lab-monitor-frontend/    # Angular SPA
    └── src/app/
        ├── core/
        │   ├── models/        # Interface Device, Lab
        │   └── services/      # HTTP & SSE service
        └── features/
            └── dashboard/     # Halaman dashboard
```

## Cara Menjalankan

### 1. Menggunakan Docker Compose (disarankan)

Pastikan Docker & Docker Compose sudah terpasang, lalu sesuaikan `NETWORK_CIDR` pada `docker-compose.yml` dengan subnet host Docker kamu (cek dengan `ip addr`).

```bash
docker compose up --build
```

- Frontend: http://localhost:4200
- Backend API: http://localhost:8080/api

> Backend butuh capability `NET_RAW` dan `net.ipv4.ping_group_range` untuk bisa melakukan ping ICMP dari dalam container — sudah diatur di `docker-compose.yml`.

### 2. Menjalankan Manual (Development)

**Backend**

```bash
cd lab-monitor-backend
cp .env.example .env   # atau sesuaikan .env yang ada
go run main.go
```

Variabel environment (`.env`) yang dibutuhkan:

| Variabel      | Keterangan                                  | Contoh              |
| ------------- | -------------------------------------------- | -------------------- |
| `DB_HOST`     | Host database PostgreSQL                     | `localhost`           |
| `DB_PORT`     | Port database                                 | `5432`                |
| `DB_USER`     | Username database                             | `atminlab`             |
| `DB_PASSWORD` | Password database                             | `atmin123`             |
| `DB_NAME`     | Nama database                                 | `monitorlab`            |
| `PORT`        | Port server backend                           | `8080`                  |
| `NETWORK_CIDR`| Subnet `/24` yang akan discan                 | `192.168.18.0/24`        |

**Frontend**

```bash
cd lab-monitor-frontend
npm install
npm start
```

Frontend akan berjalan di `http://localhost:4200` dan memakai `proxy.conf.json` untuk meneruskan request `/api` ke backend.

## API Endpoints

| Method | Endpoint            | Keterangan                                   |
| ------ | -------------------- | --------------------------------------------- |
| GET    | `/api/labs`           | Ambil daftar lab                              |
| POST   | `/api/labs`            | Tambah lab baru                              |
| GET    | `/api/devices`        | Ambil daftar perangkat                        |
| POST   | `/api/devices`         | Tambah perangkat baru                        |
| DELETE | `/api/devices/:id`     | Hapus perangkat                              |
| POST   | `/api/scan`            | Trigger scan manual                          |
| GET    | `/api/scan/stream`     | Stream hasil scan real-time (SSE)            |

## Keterbatasan

- Jangkauan scan hanya sebatas satu subnet `/24` (maksimal 254 IP host).
- Deteksi status perangkat murni berdasarkan respons ICMP ping — perangkat yang memblokir ICMP akan terbaca `offline` meskipun sebenarnya aktif.
