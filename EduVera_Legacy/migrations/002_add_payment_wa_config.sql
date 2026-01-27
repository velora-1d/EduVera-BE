-- 002_add_payment_wa_config.sql

-- Menambahkan konfigurasi Finansial & WA ke tabel TENANTS
ALTER TABLE tenants
    -- Info Bank (Untuk Pencairan Dana)
    ADD COLUMN bank_name VARCHAR(100),          -- Nama Bank (BCA, Mandiri, BS)
    ADD COLUMN bank_account_number VARCHAR(50), -- No Rekening
    ADD COLUMN bank_account_holder VARCHAR(150),-- Atas Nama Rekening

    -- Konfigurasi Payment Gateway (Xendit/Midtrans) - Disimpan Encrypted nanti di aplikasi
    ADD COLUMN payment_gateway_key VARCHAR(255),
    ADD COLUMN payment_gateway_secret VARCHAR(255),

    -- Konfigurasi WA Gateway
    ADD COLUMN wa_gateway_token VARCHAR(255),   -- Token/API Key WA
    ADD COLUMN wa_gateway_device_id VARCHAR(100); -- ID Perangkat WA
