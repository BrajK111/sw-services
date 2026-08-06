CREATE TABLE IF NOT EXISTS eg_sw_applicationdocument (
    id              VARCHAR(64) PRIMARY KEY,
    tenantid        VARCHAR(64) NOT NULL,
    documenttype    VARCHAR(64),
    filestoreid     VARCHAR(256),
    documentuid     VARCHAR(256),
    swid            VARCHAR(64) REFERENCES eg_sw_connection (id) ON DELETE CASCADE,
    active          BOOLEAN DEFAULT TRUE,
    createdby       VARCHAR(64),
    lastmodifiedby  VARCHAR(64),
    createdtime     BIGINT,
    lastmodifiedtime BIGINT
);

CREATE INDEX IF NOT EXISTS idx_sw_applicationdocument_swid ON eg_sw_applicationdocument (swid);
