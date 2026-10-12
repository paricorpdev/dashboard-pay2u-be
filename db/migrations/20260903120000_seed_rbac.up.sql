-- Seed default roles
INSERT INTO roles (code, name, description) VALUES
    ('ROLE_SUPER_ADMIN', 'Super Admin', 'Full access to all resources'),
    ('ROLE_ADMIN', 'Admin', 'Administrative access'),
    ('ROLE_USER', 'User', 'Standard user access')
ON CONFLICT (code) DO NOTHING;

-- Seed default permissions
INSERT INTO permissions (code, resource, action, description) VALUES
    ('users.read', 'users', 'read', 'View users'),
    ('users.create', 'users', 'create', 'Create users'),
    ('users.update', 'users', 'update', 'Update users'),
    ('users.delete', 'users', 'delete', 'Delete users'),
    ('users.activate', 'users', 'activate', 'Activate users'),
    ('users.deactivate', 'users', 'deactivate', 'Deactivate users'),
    ('merchants.read', 'merchants', 'read', 'View merchants'),
    ('merchants.create', 'merchants', 'create', 'Create merchants'),
    ('merchants.update', 'merchants', 'update', 'Update merchants'),
    ('merchants.delete', 'merchants', 'delete', 'Delete merchants'),
    ('transactions.read', 'transactions', 'read', 'View transactions'),
    ('analytics.read', 'analytics', 'read', 'View analytics'),
    ('roles.read', 'roles', 'read', 'View roles and permissions'),
    ('permissions.read', 'permissions', 'read', 'View permissions')
ON CONFLICT (code) DO NOTHING;

-- Link all permissions to ROLE_SUPER_ADMIN
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'ROLE_SUPER_ADMIN'
ON CONFLICT DO NOTHING;

-- Seed default super admin user (password: password123)
INSERT INTO users (username, email, password_hash, full_name, status) VALUES
    ('admin', 'admin@pay2u.com', '$2a$10$XJCTzhQW2r84M1RE4bQ.Y.44XU//5U0/lBYftjtoxLyB1ywY7fww.', 'Super Administrator', 'ACTIVE')
ON CONFLICT (username) DO NOTHING;

-- Link admin to ROLE_SUPER_ADMIN
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'admin' AND r.code = 'ROLE_SUPER_ADMIN'
ON CONFLICT DO NOTHING;

-- Link permissions to ROLE_ADMIN
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'ROLE_ADMIN'
  AND p.code IN (
      'merchants.read', 'merchants.create', 'merchants.update',
      'transactions.read', 'analytics.read',
      'roles.read', 'permissions.read'
  )
ON CONFLICT DO NOTHING;

-- Link permissions to ROLE_USER
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'ROLE_USER'
  AND p.code IN ('transactions.read')
ON CONFLICT DO NOTHING;

-- Seed default admin & staff users (password: password123)
INSERT INTO users (username, email, password_hash, full_name, status) VALUES
    ('admin_ops', 'ops@pay2u.com', '$2a$10$XJCTzhQW2r84M1RE4bQ.Y.44XU//5U0/lBYftjtoxLyB1ywY7fww.', 'Operations Admin', 'ACTIVE'),
    ('staff', 'staff@pay2u.com', '$2a$10$XJCTzhQW2r84M1RE4bQ.Y.44XU//5U0/lBYftjtoxLyB1ywY7fww.', 'Staff Member', 'ACTIVE')
ON CONFLICT (username) DO NOTHING;

-- Link admin_ops to ROLE_ADMIN
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'admin_ops' AND r.code = 'ROLE_ADMIN'
ON CONFLICT DO NOTHING;

-- Link staff to ROLE_USER
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'staff' AND r.code = 'ROLE_USER'
ON CONFLICT DO NOTHING;
