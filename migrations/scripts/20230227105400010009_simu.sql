CREATE TABLE client (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    legal_name VARCHAR(100) NOT NULL,
    api_key VARCHAR(100) NOT NULL,
    status BOOLEAN DEFAULT true,
    initial VARCHAR(10),
    website VARCHAR(200)
);

-- Add sample account for testing the SDK
insert into client (id, name, api_key, status, legal_name) values (gen_random_uuid(), 'Sample testing account', 'a701551d-70fe-4d0f-b100-3a274fe4f225', true, '');
