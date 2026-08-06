CREATE TABLE IF NOT EXISTS eg_sw_roadcuttinginfo (
    id                  VARCHAR(64) PRIMARY KEY,
    tenantid            VARCHAR(64) NOT NULL,
    swid                VARCHAR(64) REFERENCES eg_sw_connection (id) ON DELETE CASCADE,
    active              BOOLEAN DEFAULT TRUE,
    roadtype            VARCHAR(64),
    roadcuttingarea     DOUBLE PRECISION,
    createdby           VARCHAR(64),
    lastmodifiedby      VARCHAR(64),
    createdtime         BIGINT,
    lastmodifiedtime    BIGINT
);

CREATE INDEX IF NOT EXISTS idx_sw_roadcuttinginfo_swid ON eg_sw_roadcuttinginfo (swid);
