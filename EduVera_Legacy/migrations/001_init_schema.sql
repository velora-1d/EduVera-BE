-- 001_init_schema.sql

-- Enable UUID extension (biar ID-nya acak & aman, standar SaaS modern)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ENUMs untuk standarisasi (agar data tidak berantakan)
CREATE TYPE tenant_status AS ENUM ('ACTIVE', 'INACTIVE', 'SUSPENDED');
CREATE TYPE school_level AS ENUM ('PAUD', 'TK', 'SD', 'MI', 'SMP', 'MTS', 'SMA', 'MA', 'SMK', 'PESANTREN', 'OTHER');
CREATE TYPE user_role AS ENUM ('SUPER_ADMIN', 'ADMIN_SEKOLAH', 'GURU', 'STAFF', 'WALI', 'SISWA');

-- 1. Tabel TENANTS (Penyewa / Sekolah)
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,                 -- Nama Sekolah (Contoh: SMA 1 EduVera)
    subdomain VARCHAR(100) UNIQUE NOT NULL,     -- Alamat Web (Contoh: sma1.eduvera.id)
    npsn VARCHAR(20),                           -- Nomor Pokok Sekolah Nasional
    level school_level NOT NULL,                -- Jenjang utama

    -- MODUL SUBSCRIPTION (Bisa ambil dua-duanya)
    has_school_module BOOLEAN DEFAULT TRUE,     -- Langganan Fitur Sekolah (Akademik, Nilai)
    has_pesantren_module BOOLEAN DEFAULT FALSE, -- Langganan Fitur Pesantren (Asrama, Tahfidz)

    address TEXT,
    phone VARCHAR(50),
    status tenant_status DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Tabel USERS (Pengguna)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE, -- Link ke Sekolah. Jika NULL berarti Super Admin Global.
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'SISWA',
    phone VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes (Biar pencarian cepat)
CREATE INDEX idx_tenants_subdomain ON tenants(subdomain);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_tenant ON users(tenant_id);

-- ==========================================
-- MODUL PESANTREN (Data Asrama & Santri)
-- ==========================================

-- 3. Tabel ASRAMA (Gedung)
CREATE TABLE dorms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,    -- Nama Asrama (misal: Asrama Putra Al-Fatih)
    gender VARCHAR(10) NOT NULL CHECK (gender IN ('MALE', 'FEMALE')), -- Khusus Putra/Putri
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. Tabel KAMAR (Hujroh)
CREATE TABLE dorm_rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dorm_id UUID REFERENCES dorms(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,     -- Nomor/Nama Kamar (misal: A-101)
    capacity INT DEFAULT 10,       -- Kapasitas Santri
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. Profil SANTRI (Extend dari User)
CREATE TABLE santri_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE, -- Link ke tabel Users
    nis VARCHAR(50),               -- Nomor Induk Santri
    dorm_room_id UUID REFERENCES dorm_rooms(id), -- Link ke Kamar (Bisa NULL kalau Non-Mukim)
    parent_phone VARCHAR(50),      -- No HP Wali (untuk WA Gateway)
    tahfidz_target INT DEFAULT 0,  -- Target Hafalan (Juz)
    points INT DEFAULT 100,        -- Poin Kedisiplinan Awal
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_santri_nis ON santri_profiles(nis);
