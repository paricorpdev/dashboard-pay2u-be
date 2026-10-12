DELETE FROM role_permissions WHERE role_id IN (SELECT id FROM roles WHERE code IN ('ROLE_SUPER_ADMIN', 'ROLE_ADMIN', 'ROLE_USER'));
DELETE FROM permissions WHERE code IN ('users.read', 'users.create', 'users.update', 'users.delete', 'users.activate', 'users.deactivate');
DELETE FROM roles WHERE code IN ('ROLE_SUPER_ADMIN', 'ROLE_ADMIN', 'ROLE_USER');
