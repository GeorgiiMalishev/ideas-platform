DO $$
DECLARE
    admin_user_id UUID := 'f47ac10b-58cc-4372-a567-0e02b2c3d479';
    system_coffee_shop_id UUID := 'a47ac10b-58cc-4372-a567-0e02b2c3d479';
    admin_role_id UUID := 'b47ac10b-58cc-4372-a567-0e02b2c3d479';
    worker_coffee_shop_id UUID := 'c47ac10b-58cc-4372-a567-0e02b2c3d479';
BEGIN
    INSERT INTO "users" (id, name, phone, is_deleted, updated_at, created_at)
    VALUES (admin_user_id, 'Super Admin', '999', false, NOW(), NOW())
    ON CONFLICT (phone) DO NOTHING;

    INSERT INTO "coffee_shop" (id, creator_id, name, address, updated_at, created_at)
    VALUES (system_coffee_shop_id, admin_user_id, 'system', 'system', NOW(), NOW())
    ON CONFLICT (name) DO NOTHING;

    INSERT INTO "role" (id, name, is_deleted, updated_at, created_at)
    VALUES (admin_role_id, 'admin', false, NOW(), NOW())
    ON CONFLICT (name) DO NOTHING;

    INSERT INTO "worker_coffee_shop" (id, worker_id, coffee_shop_id, role_id, is_deleted, created_at)
    VALUES (worker_coffee_shop_id, admin_user_id, system_coffee_shop_id, admin_role_id, false, NOW())
    ON CONFLICT (worker_id, coffee_shop_id) DO NOTHING;
END $$;
