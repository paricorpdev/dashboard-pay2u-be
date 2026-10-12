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
    ('users.deactivate', 'users', 'deactivate', 'Deactivate users')
ON CONFLICT (code) DO NOTHING;

-- Link all permissions to ROLE_SUPER_ADMIN
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'ROLE_SUPER_ADMIN'
ON CONFLICT DO NOTHING;
