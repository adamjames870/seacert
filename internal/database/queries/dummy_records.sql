-- name: CreateDummyCertTypes :exec

INSERT INTO certificate_types (
    id,
    name,
    short_name,
    stcw_reference,
    normal_validity_months,
    status,
    created_at,
    updated_at
) VALUES
-- CoC
('10000000-0000-0000-0000-000000000001', 'Certificate of Competency - Deck - Chief Mate Unlimited', 'CoC-CM', 'A-II / 2', 60, 'approved', NOW(), NOW()),
('10000000-0000-0000-0000-000000000002', 'Certificate of Competency - Deck - OOW Unlimited', 'CoC-OOW', 'A-II / 1', 60, 'approved', NOW(), NOW()),
('10000000-0000-0000-0000-000000000003', 'Certificate of Competency - Deck - Master Unlimited', 'CoC-MS', 'A-II / 2', 60, 'approved', NOW(), NOW()),

-- Medical
(gen_random_uuid(), 'Seafarer Medical - UK - ENG1', 'ENG1', 'A-I / 9', 24, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Seafarer Medical', 'MED', 'A-I / 9', 24, 'approved', NOW(), NOW()),

-- Radio
(gen_random_uuid(), 'GMDSS General Operator', 'GMDSS', 'A-IV / 2', 60, 'approved', NOW(), NOW()),

-- Basic Training
(gen_random_uuid(), 'Basic Training STCW', 'BT', 'A-VI / 1', 60, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Personal Survival Techniques', 'PST', 'A-VI / 1-1', 60, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Fire Fighting & Prevention', 'FFP', 'A-VI / 1-2', 60, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Elementary First Aid', 'EFA', 'A-VI / 1-3', NULL, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Personal Safety & Social Responsibility', 'PSSR', 'A-VI / 1-4', NULL, 'approved', NOW(), NOW()),

-- Advanced Fire
(gen_random_uuid(), 'Advanced Firefighting Update', 'AFF-U', 'A-VI / 3', 60, 'approved', NOW(), NOW()),

-- Tanker
(gen_random_uuid(), 'Tanker Firefighting', 'TFF', 'A-V / 1', 60, 'approved', NOW(), NOW()),

-- Survival Craft
(gen_random_uuid(), 'Proficiency in Survival Craft and Rescue Boats (Except Fast Rescue Boats)', 'PSCRB', 'A-VI / 2-1', 60, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Proficiency in Survival Craft and Rescue Boats (Except Fast Rescue Boats) Update', 'PSCRB-U', 'A-VI / 2-1', 60, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Proficiency in Fast Rescue Boats Update', 'FRC-U', 'A-VI / 2-1', 60, 'approved', NOW(), NOW()),

-- Medical Proficiency
(gen_random_uuid(), 'Proficiency in Medical First Aid', 'MFA', 'A-VI / 4-1', NULL, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Proficiency in Medical Care', 'PMC', 'A-VI / 4-2', NULL, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Proficiency in Medical Care Update', 'PMC-U', 'A-VI / 4-2', NULL, 'approved', NOW(), NOW()),

-- Security
(gen_random_uuid(), 'Proficiency in Security Awareness', 'PSA', 'A-VI / 6-1', NULL, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Proficiency in Designated Security Duties', 'DSD', 'A-VI / 6-2', NULL, 'approved', NOW(), NOW()),

-- Onboard Medical
(gen_random_uuid(), 'Proficiency in Medical Care Onboard', 'MCOB', 'A-VI / 4-2', 60, 'approved', NOW(), NOW()),
(gen_random_uuid(), 'Proficiency in Medical Care Onboard Update', 'MCOBU', 'A-VI / 4-2', 60, 'approved', NOW(), NOW());

-- name: CreateDummyIssuers :exec

-- one-time setup (safe to run repeatedly)

INSERT INTO issuers (
    id,
    name,
    country,
    website,
    created_at,
    updated_at
) VALUES
      ('10000000-0000-0000-0000-000000000001', 'Maritime & Coastguard Agency (MCA)', 'GB', 'https://www.gov.uk/government/organisations/maritime-and-coastguard-agency', NOW(), NOW()),
      (gen_random_uuid(), 'Warsash Maritime Academy (WMA)', 'GB', 'https://maritime.solent.ac.uk/', NOW(), NOW()),
      (gen_random_uuid(), 'Maritime Skills Academy (MSA)', 'GB', 'https://www.maritimeskillsacademy.com/', NOW(), NOW()),
      (gen_random_uuid(), 'Stream Marine Training (SMT)', 'GB', 'https://streammarinetraining.com/', NOW(), NOW()),
      (gen_random_uuid(), 'Medical Support Offshore (MSOS)', 'GB', 'https://www.msos.org.uk/', NOW(), NOW()),
      (gen_random_uuid(), 'Fire-Aid Academy', 'GB', 'https://fireaid.com/', NOW(), NOW()),
      (gen_random_uuid(), 'Marina', 'PH', NULL, NOW(), NOW());

-- name: CreateDummyUsers :exec

-- Notification test users
--
-- A: >7 days old, no certificates
--    EXPECTED: eligible for no_certificates_7d
--
-- B: <7 days old, no certificates
--    EXPECTED: not eligible
--
-- C: >7 days old, will be given a certificate
--    EXPECTED: not eligible
--
-- D: >7 days old, no certificates, but already has a 7d notification
--    EXPECTED: not eligible

INSERT INTO users (
    id,
    created_at,
    updated_at,
    forename,
    surname,
    email,
    nationality,
    email_consent
)
VALUES
    (
        '10000000-0000-0000-0000-000000000001',
        NOW() - INTERVAL '10 days',
        NOW(),
        'Test',
        'Eligible',
        'hello@seacert.app',
        'GB',
        TRUE
    ),
    (
        '10000000-0000-0000-0000-000000000002',
        NOW() - INTERVAL '3 days',
        NOW(),
        'Test',
        'TooYoung',
        'hello@seacert.app',
        'GB',
        TRUE
    ),
    (
        '10000000-0000-0000-0000-000000000003',
        NOW() - INTERVAL '10 days',
        NOW(),
        'Test',
        'HasCertificate',
        'hello@seacert.app',
        'GB',
        TRUE
    ),
    (
        '10000000-0000-0000-0000-000000000004',
        NOW() - INTERVAL '10 days',
        NOW(),
        'Test',
        'AlreadyNotified',
        'hello@seacert.app',
        'GB',
        TRUE
    );

-- name: CreateDummyCertificates :exec

INSERT INTO certificates (id, created_at, updated_at, user_id, cert_type_id, cert_number, issuer_id, issued_date)
VALUES (
        '10000000-0000-0000-0000-000000000001',
        NOW(),
        NOW(),
        '10000000-0000-0000-0000-000000000003',
        '10000000-0000-0000-0000-000000000001',
        '123456789',
        '10000000-0000-0000-0000-000000000001',
        NOW()
       );

-- name: CreateDummyNotifications :exec

INSERT INTO notifications (
    id,
    user_id,
    notification_type,
    notification_key,
    status,
    payload,
    scheduled_at,
    created_at,
    updated_at
)
VALUES (
           gen_random_uuid(),
           '10000000-0000-0000-0000-000000000004',
           'no_certificates_7d',
           'no-certificates:10000000-0000-0000-0000-000000000004:7d',
           'completed',
           '{}'::jsonb,
           NOW() - INTERVAL '1 day',
           NOW() - INTERVAL '1 day',
           NOW()
       );