# Nusantara Analytics

Nusantara Analytics adalah proyek monorepo yang terdiri dari dua bagian utama:

## Struktur Proyek

- **[backend/](file:///d:/kuliah/Project-After-Lulus/nusantara-analytics/backend)**: Layanan backend API yang dibangun menggunakan **Golang** dan framework **Fiber**.
- **[frontend/](file:///d:/kuliah/Project-After-Lulus/nusantara-analytics/frontend)**: Aplikasi frontend yang akan dibangun menggunakan **Nuxt 3**.

## Memulai Proyek

### Backend
Untuk menjalankan backend secara lokal, masuk ke direktori `backend/`:
```bash
cd backend
go run main.go
```

### Frontend
Untuk menjalankan frontend secara lokal, masuk ke direktori `frontend/`:
```bash
cd frontend
pnpm install
pnpm dev
```

## Panduan Deployment

### 1. Backend (Railway)
- Push proyek monorepo ini ke GitHub.
- Hubungkan/connect **Railway** ke repositori tersebut dan arahkan ke folder `backend/`.
- Railway akan otomatis membaca konfigurasi `railway.toml` dan `Dockerfile` untuk mendeploy backend Go.

### 2. Frontend (Vercel)
- Push proyek monorepo ini ke GitHub.
- Hubungkan/connect **Vercel** ke repositori tersebut dan arahkan ke folder `frontend/`.
- Atur Environment Variable di dashboard Vercel:
  - **`NUXT_PUBLIC_API_BASE`**: URL backend Railway (contoh: `https://nusantara-analytics-production.up.railway.app`)
- Vercel akan otomatis membaca file `vercel.json` untuk konfigurasi build dan deploy.
