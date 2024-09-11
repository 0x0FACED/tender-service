-- Удаляем записи из organization_responsible
DELETE FROM organization_responsible 
WHERE user_id IN (SELECT id FROM employee WHERE username IN ('john_doe', 'jane_smith', 'alice_brown', 'bob_jones', 'charlie_davis'));

-- Удаляем записи из organization
DELETE FROM organization 
WHERE name IN ('Tech Solutions', 'Global Logistics', 'BuildCo', 'Innovatech', 'EcoManufacture');

-- Удаляем записи из employee
DELETE FROM employee 
WHERE username IN ('john_doe', 'jane_smith', 'alice_brown', 'bob_jones', 'charlie_davis');

-- Удаляем запись о сотруднике, который не является ответственным за организацию
DELETE FROM employee 
WHERE username = 'michael_jordan';
