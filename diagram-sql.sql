CREATE TABLE employee (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM ('IE', 'LLC', 'JSC');

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id SERIAL PRIMARY KEY,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id INT REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TYPE tender_service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');
CREATE TYPE tender_status AS ENUM ('Created', 'Published', 'Closed');
CREATE TYPE bid_status AS ENUM ('Created', 'Published', 'Canceled', 'Approved', 'Rejected');
CREATE TYPE bid_author_type AS ENUM ('Organization', 'User');
CREATE TYPE bid_decision AS ENUM ('Approved', 'Rejected');

CREATE TABLE tenders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    service_type tender_service_type NOT NULL,
    status tender_status NOT NULL,
    organization_id UUID REFERENCES organization(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tender_versions (
    id SERIAL PRIMARY KEY, 
    tender_id UUID REFERENCES tenders(id) ON DELETE CASCADE NOT NULL,
    version_number INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    service_type tender_service_type NOT NULL,
    status tender_status NOT NULL,
    organization_id UUID REFERENCES organization(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_tender_versions_current ON tender_versions(tender_id, is_current);

CREATE TABLE bids (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    status bid_status NOT NULL,
    tender_id UUID REFERENCES tenders(id) ON DELETE CASCADE NOT NULL,
    author_type bid_author_type NOT NULL,
    author_id INT REFERENCES employee(id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bid_decisions (
    id SERIAL PRIMARY KEY,
    bid_id UUID REFERENCES bids(id) ON DELETE CASCADE NOT NULL,
    decision bid_decision NOT NULL
);

CREATE TABLE bid_versions (
    id SERIAL PRIMARY KEY,
    bid_id UUID REFERENCES bids(id) ON DELETE CASCADE NOT NULL,
    version_number INT NOT NULL,
    author_id INT REFERENCES employee(id) ON DELETE CASCADE,
    status bid_status NOT NULL,
    decision bid_decision,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_current BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_bid_versions_current ON bid_versions(bid_id, is_current);

CREATE TABLE bid_feedbacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bid_id UUID REFERENCES bids(id) ON DELETE CASCADE NOT NULL,
    author_id INT REFERENCES employee(id) ON DELETE CASCADE,
    description VARCHAR(1000) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
