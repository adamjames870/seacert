-- +goose Up

-- Reference table for ship types (e.g., Tanker, Bulk Carrier, Container)
CREATE TABLE ship_types (
                            id UUID PRIMARY KEY,
                            name TEXT NOT NULL UNIQUE,
                            description TEXT
);

INSERT INTO ship_types (id, name, description) VALUES
                                                   ('00000000-0000-0000-0000-000000000010', 'oil_tanker', 'Oil Tanker'),
                                                   ('00000000-0000-0000-0000-000000000011', 'chemical_tanker', 'Chemical Tanker'),
                                                   ('00000000-0000-0000-0000-000000000012', 'gas_carrier', 'Liquefied Gas Carrier'),
                                                   ('00000000-0000-0000-0000-000000000013', 'bulk_carrier', 'Bulk Carrier'),
                                                   ('00000000-0000-0000-0000-000000000014', 'container_ship', 'Container Ship'),
                                                   ('00000000-0000-0000-0000-000000000015', 'general_cargo', 'General Cargo Ship'),
                                                   ('00000000-0000-0000-0000-000000000016', 'ro_ro', 'Roll-on/Roll-off Vessel'),
                                                   ('00000000-0000-0000-0000-000000000017', 'passenger_ship', 'Passenger/Cruise Ship'),
                                                   ('00000000-0000-0000-0000-000000000018', 'offshore_supply', 'Offshore Support/Supply Vessel'),
                                                   ('00000000-0000-0000-0000-000000000019', 'tug', 'Tug/Towing Vessel'),
                                                   ('00000000-0000-0000-0000-000000000020', 'other', 'Other vessel type');

-- Reference table for voyage types
CREATE TABLE voyage_types (
                              id UUID PRIMARY KEY,
                              name TEXT NOT NULL UNIQUE,
                              description TEXT
);

INSERT INTO voyage_types (id, name, description) VALUES
                                                     ('00000000-0000-0000-0000-000000000001', 'near_coastal', 'NCV (Near Coastal)'),
                                                     ('00000000-0000-0000-0000-000000000002', 'international', 'FGN (International)'),
                                                     ('00000000-0000-0000-0000-000000000003', 'offshore', 'OFFSH (Offshore)');

-- Reference table for specialized period types
CREATE TABLE seatime_period_types (
                                      id UUID PRIMARY KEY,
                                      name TEXT NOT NULL UNIQUE,
                                      description TEXT
);

INSERT INTO seatime_period_types (id, name, description) VALUES
                                                             ('00000000-0000-0000-0000-000000000003', 'polar', 'Polar waters'),
                                                             ('00000000-0000-0000-0000-000000000004', 'dp', 'Dynamic Positioning'),
                                                             ('00000000-0000-0000-0000-000000000005', 'tanker', 'Service on Tankers');

-- Ships table referencing ship_types
CREATE TABLE ships (
                       id UUID PRIMARY KEY,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       name TEXT NOT NULL,
                       ship_type_id UUID NOT NULL REFERENCES ship_types(id),
                       imo_number TEXT NOT NULL UNIQUE,
                       gt INT NOT NULL,
                       flag TEXT NOT NULL,
                       propulsion_power INT, -- Optional
                       status VARCHAR(20) NOT NULL DEFAULT 'approved',
                       created_by UUID,

                       CONSTRAINT fk_cert_types_created_by
                        FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Main seatime record table
CREATE TABLE seatime (
                         id UUID PRIMARY KEY,
                         user_id UUID NOT NULL REFERENCES users(id),
                         ship_id UUID NOT NULL REFERENCES ships(id),
                         voyage_type_id UUID NOT NULL REFERENCES voyage_types(id),
                         created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

                         start_date DATE NOT NULL,
                         start_location TEXT NOT NULL,
                         end_date DATE NOT NULL,
                         end_location TEXT NOT NULL,
                         total_days INT NOT NULL,

                         company TEXT NOT NULL,
                         capacity TEXT NOT NULL,
                         is_watchkeeping BOOLEAN NOT NULL DEFAULT FALSE
);

-- Specialized service periods within a voyage
CREATE TABLE seatime_periods (
                                 id UUID PRIMARY KEY,
                                 seatime_id UUID NOT NULL REFERENCES seatime(id) ON DELETE CASCADE,
                                 period_type_id UUID NOT NULL REFERENCES seatime_period_types(id),
                                 start_date DATE NOT NULL,
                                 end_date DATE NOT NULL,
                                 days INT NOT NULL,
                                 remarks TEXT
);

-- +goose Down
DROP TABLE seatime_periods;
DROP TABLE seatime;
DROP TABLE ships;
DROP TABLE seatime_period_types;
DROP TABLE voyage_types;
DROP TABLE ship_types;