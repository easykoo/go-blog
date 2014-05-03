DELETE FROM user;
INSERT INTO user (id, username, password, full_name, gender, qq, tel, postcode, address, email, role_id, dept_id, active, locked, create_user, create_date, update_user, update_date, version)
VALUES (1, 'admin', 'b0baee9d279d34fa1dfd71aadb908c3f', 'Admin', 1, 111111, '11122233344', '123456', '自由大道1号',
        'admin@admin.com', 1, 1, 1, 0, 'SYSTEM', now(), 'SYSTEM', now(), 1);

DELETE FROM role;
INSERT INTO role (id, description, create_user, create_date, update_user, update_date)
VALUES (1, '管理员', 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO role (id, description, create_user, create_date, update_user, update_date)
VALUES (2, '经理', 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO role (id, description, create_user, create_date, update_user, update_date)
VALUES (3, '员工', 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO role (id, description, create_user, create_date, update_user, update_date)
VALUES (4, '用户', 'SYSTEM', now(), 'SYSTEM', now());

DELETE FROM dept;
INSERT INTO dept (id, description, create_user, create_date, update_user, update_date)
VALUES (1, '默认', 'SYSTEM', now(), 'SYSTEM', now());

DELETE FROM module;
INSERT INTO module (id, description, create_user, create_date, update_user, update_date)
VALUES (1, 'Admin', 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO module (id, description, create_user, create_date, update_user, update_date)
VALUES (2, 'Account', 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO module (id, description, create_user, create_date, update_user, update_date)
VALUES (3, 'Feedback', 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO module (id, description, create_user, create_date, update_user, update_date)
VALUES (4, 'News', 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO module (id, description, create_user, create_date, update_user, update_date)
VALUES (5, 'Product', 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO module (id, description, create_user, create_date, update_user, update_date)
VALUES (6, 'Blog', 'SYSTEM', now(), 'SYSTEM', now());

DELETE FROM privilege;
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (1, 1, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (2, 1, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (2, 2, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (3, 1, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (3, 2, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (3, 3, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (4, 1, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (4, 2, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (4, 3, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (5, 1, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (5, 2, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (6, 1, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (6, 2, 1, 'SYSTEM', now(), 'SYSTEM', now());
INSERT INTO privilege (module_id, role_id, dept_id, create_user, create_date, update_user, update_date)
VALUES (6, 3, 1, 'SYSTEM', now(), 'SYSTEM', now());

DELETE FROM category;
INSERT INTO category (id, description, parent_id, create_user, create_date, update_user, update_date)
VALUES (1, '默认', 0, 'SYSTEM', now(), 'SYSTEM', now());

DELETE FROM settings;
INSERT INTO settings (id, app_name, owner_id, about, create_user, create_date, update_user, update_date)
VALUES (1, 'Easy Go', 1, null, 'SYSTEM', now(), 'SYSTEM', now());
