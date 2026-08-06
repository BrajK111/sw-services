-- NOTE: Java's eg_sw_connectionholder has UNIQUE(connectionid), which silently
-- limits every connection to a single holder despite the model supporting a
-- list (see Resource/sw-services issues.pdf, finding B). This DDL intentionally
-- does NOT repeat that constraint, and adds a proper surrogate id PK (Java's
-- table has none) plus a numeric holdershippercentage (Java stores it as
-- varchar despite the model field being a Double).
CREATE TABLE IF NOT EXISTS eg_sw_connectionholder (
    id                      VARCHAR(64) PRIMARY KEY,
    tenantid                VARCHAR(64) NOT NULL,
    connectionid            VARCHAR(64) REFERENCES eg_sw_connection (id) ON DELETE CASCADE,
    status                  VARCHAR(64),
    userid                  VARCHAR(64),
    isprimaryholder         BOOLEAN DEFAULT FALSE,
    connectionholdertype    VARCHAR(64),
    holdershippercentage    NUMERIC(5, 2),
    relationship            VARCHAR(32),
    createdby               VARCHAR(64),
    createdtime             BIGINT,
    lastmodifiedby          VARCHAR(64),
    lastmodifiedtime        BIGINT
);

CREATE INDEX IF NOT EXISTS idx_sw_connectionholder_userid ON eg_sw_connectionholder (userid);
CREATE INDEX IF NOT EXISTS idx_sw_connectionholder_connectionid ON eg_sw_connectionholder (connectionid);
