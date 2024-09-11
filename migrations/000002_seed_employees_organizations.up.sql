-- Заполняем таблицу employee, если записи не существуют
INSERT INTO employee (username, first_name, last_name)
SELECT 'john_doe', 'John', 'Doe'
WHERE NOT EXISTS (SELECT 1 FROM employee WHERE username = 'john_doe');

INSERT INTO employee (username, first_name, last_name)
SELECT 'jane_smith', 'Jane', 'Smith'
WHERE NOT EXISTS (SELECT 1 FROM employee WHERE username = 'jane_smith');

INSERT INTO employee (username, first_name, last_name)
SELECT 'alice_brown', 'Alice', 'Brown'
WHERE NOT EXISTS (SELECT 1 FROM employee WHERE username = 'alice_brown');

INSERT INTO employee (username, first_name, last_name)
SELECT 'bob_jones', 'Bob', 'Jones'
WHERE NOT EXISTS (SELECT 1 FROM employee WHERE username = 'bob_jones');

INSERT INTO employee (username, first_name, last_name)
SELECT 'charlie_davis', 'Charlie', 'Davis'
WHERE NOT EXISTS (SELECT 1 FROM employee WHERE username = 'charlie_davis');

-- Добавляем сотрудника, который не является ответственным за организацию
INSERT INTO employee (username, first_name, last_name)
SELECT 'michael_jordan', 'Michael', 'Jordan'
WHERE NOT EXISTS (SELECT 1 FROM employee WHERE username = 'michael_jordan');

-- Заполняем таблицу organization, если записи не существуют
INSERT INTO organization (name, description, type)
SELECT 'Tech Solutions', 'IT Consulting Company', 'LLC'
WHERE NOT EXISTS (SELECT 1 FROM organization WHERE name = 'Tech Solutions');

INSERT INTO organization (name, description, type)
SELECT 'Global Logistics', 'Logistics and Delivery', 'LLC'
WHERE NOT EXISTS (SELECT 1 FROM organization WHERE name = 'Global Logistics');

INSERT INTO organization (name, description, type)
SELECT 'BuildCo', 'Construction Company', 'JSC'
WHERE NOT EXISTS (SELECT 1 FROM organization WHERE name = 'BuildCo');

INSERT INTO organization (name, description, type)
SELECT 'Innovatech', 'Research and Development', 'LLC'
WHERE NOT EXISTS (SELECT 1 FROM organization WHERE name = 'Innovatech');

INSERT INTO organization (name, description, type)
SELECT 'EcoManufacture', 'Eco-friendly Manufacturing', 'IE'
WHERE NOT EXISTS (SELECT 1 FROM organization WHERE name = 'EcoManufacture');

-- Заполняем таблицу organization_responsible, если записи не существуют
-- Связываем первых 5 сотрудников с организациями
INSERT INTO organization_responsible (organization_id, user_id)
SELECT (SELECT id FROM organization WHERE name = 'Tech Solutions'), (SELECT id FROM employee WHERE username = 'john_doe')
WHERE NOT EXISTS (SELECT 1 FROM organization_responsible WHERE organization_id = (SELECT id FROM organization WHERE name = 'Tech Solutions') AND user_id = (SELECT id FROM employee WHERE username = 'john_doe'));

INSERT INTO organization_responsible (organization_id, user_id)
SELECT (SELECT id FROM organization WHERE name = 'Global Logistics'), (SELECT id FROM employee WHERE username = 'jane_smith')
WHERE NOT EXISTS (SELECT 1 FROM organization_responsible WHERE organization_id = (SELECT id FROM organization WHERE name = 'Global Logistics') AND user_id = (SELECT id FROM employee WHERE username = 'jane_smith'));

INSERT INTO organization_responsible (organization_id, user_id)
SELECT (SELECT id FROM organization WHERE name = 'BuildCo'), (SELECT id FROM employee WHERE username = 'alice_brown')
WHERE NOT EXISTS (SELECT 1 FROM organization_responsible WHERE organization_id = (SELECT id FROM organization WHERE name = 'BuildCo') AND user_id = (SELECT id FROM employee WHERE username = 'alice_brown'));

INSERT INTO organization_responsible (organization_id, user_id)
SELECT (SELECT id FROM organization WHERE name = 'Innovatech'), (SELECT id FROM employee WHERE username = 'bob_jones')
WHERE NOT EXISTS (SELECT 1 FROM organization_responsible WHERE organization_id = (SELECT id FROM organization WHERE name = 'Innovatech') AND user_id = (SELECT id FROM employee WHERE username = 'bob_jones'));

INSERT INTO organization_responsible (organization_id, user_id)
SELECT (SELECT id FROM organization WHERE name = 'EcoManufacture'), (SELECT id FROM employee WHERE username = 'charlie_davis')
WHERE NOT EXISTS (SELECT 1 FROM organization_responsible WHERE organization_id = (SELECT id FROM organization WHERE name = 'EcoManufacture') AND user_id = (SELECT id FROM employee WHERE username = 'charlie_davis'));
