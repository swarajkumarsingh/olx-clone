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

-- begin; drop table client cascade; drop table kyc_session cascade; drop table kyc_session_redirect; drop table kyc_session_digilocker; drop table kyc_session_media; drop table kyc_session_partner_data; drop table session_kyc_evaluations; delete from migrations_metadata  where script_name in ('scripts/20221212114411011119_client.sql','scripts/20221212114420011209_kyc_session.sql','scripts/20221212114431011319_kyc_evaluations.sql','20221223160447004479_vault.sql'); end;