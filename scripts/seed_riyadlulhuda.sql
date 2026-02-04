-- Seed Script for Riyadlul Huda Premium Tenant
-- This script creates a Premium Hybrid tenant and an Admin user.

DO $$
DECLARE
    v_tenant_id UUID := gen_random_uuid();
    v_user_id UUID := gen_random_uuid();
    v_password_hash TEXT := '$2a$10$hZlIrOIPcPksVXT5rzRviugkt1v7.gZjtGNBZBSpg2YdX1QjTw4Ju'; -- 'admin123'
BEGIN
    -- 1. Create Tenant
    INSERT INTO tenants (
        id, 
        name, 
        school_name,
        pesantren_name,
        school_jenjangs,
        subdomain, 
        plan_type, 
        institution_type, 
        status, 
        subscription_tier,
        created_at, 
        updated_at
    ) VALUES (
        v_tenant_id,
        'Riyadlul Huda',
        'Sekolah Riyadlul Huda',
        'Pesantren Riyadlul Huda',
        ARRAY['TK', 'SD', 'MI', 'SMP', 'MTs', 'SMA', 'MA', 'SMK'],
        'riyadlulhuda',
        'hybrid',
        'sekolah_pesantren',
        'active',
        'premium',
        NOW(),
        NOW()
    ) ON CONFLICT (subdomain) DO UPDATE SET
        subscription_tier = 'premium',
        status = 'active',
        plan_type = 'hybrid',
        school_name = 'Sekolah Riyadlul Huda',
        pesantren_name = 'Pesantren Riyadlul Huda',
        school_jenjangs = ARRAY['TK', 'SD', 'MI', 'SMP', 'MTs', 'SMA', 'MA', 'SMK'];

    -- 2. Create Admin User
    INSERT INTO users (
        id,
        tenant_id,
        name,
        email,
        password_hash,
        role,
        is_active,
        email_verified_at,
        created_at,
        updated_at
    ) VALUES (
        v_user_id,
        v_tenant_id,
        'Admin Riyadlul Huda',
        'admin@riyadlulhuda.id',
        v_password_hash,
        'admin',
        true,
        NOW(),
        NOW(),
        NOW()
    ) ON CONFLICT (email) DO UPDATE SET
        tenant_id = EXCLUDED.tenant_id,
        is_active = true,
        password_hash = EXCLUDED.password_hash;

    -- 3. Create active subscription records to ensure feature gating works
    INSERT INTO subscriptions (
        id,
        tenant_id,
        plan_type,
        billing_cycle,
        status,
        subscription_tier,
        current_period_start,
        current_period_end,
        grace_period_end,
        created_at,
        updated_at
    ) VALUES (
        gen_random_uuid(),
        v_tenant_id,
        'hybrid',
        'annual',
        'active',
        'premium',
        NOW(),
        NOW() + INTERVAL '1 year',
        NOW() + INTERVAL '1 year 7 days',
        NOW(),
        NOW()
    ) ON CONFLICT (tenant_id) WHERE status IN ('active', 'grace_period') DO NOTHING;

    RAISE NOTICE 'Tenant Riyadlul Huda (riyadlulhuda) seeded successfully.';
    RAISE NOTICE 'Admin User: admin@riyadlulhuda.id / admin123';
END $$;
