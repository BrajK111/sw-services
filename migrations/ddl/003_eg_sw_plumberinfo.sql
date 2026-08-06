CREATE TABLE IF NOT EXISTS eg_sw_plumberinfo (
    id                      VARCHAR(64) PRIMARY KEY,
    tenantid                VARCHAR(64) NOT NULL,
    name                    VARCHAR(256),
    licenseno               VARCHAR(64),
    mobilenumber            VARCHAR(16),
    gender                  VARCHAR(16),
    fatherorhusbandname     VARCHAR(256),
    correspondenceaddress   VARCHAR(1024),
    relationship            VARCHAR(32),
    swid                    VARCHAR(64) REFERENCES eg_sw_connection (id) ON DELETE CASCADE,
    createdby               VARCHAR(64),
    lastmodifiedby          VARCHAR(64),
    createdtime             BIGINT,
    lastmodifiedtime        BIGINT
);

CREATE INDEX IF NOT EXISTS idx_sw_plumberinfo_swid ON eg_sw_plumberinfo (swid);
