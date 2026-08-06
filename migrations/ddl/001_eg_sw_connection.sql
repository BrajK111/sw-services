CREATE TABLE IF NOT EXISTS eg_sw_connection (
    id                          VARCHAR(64) PRIMARY KEY,
    tenantid                    VARCHAR(64) NOT NULL,
    property_id                 VARCHAR(64),
    applicationno                VARCHAR(64),
    applicationstatus           VARCHAR(64),
    status                      VARCHAR(64),
    connectionno                VARCHAR(64),
    oldconnectionno             VARCHAR(64),
    roadtype                    VARCHAR(64),
    roadcuttingarea             DOUBLE PRECISION,
    adhocrebate                 NUMERIC(12, 2),
    adhocpenalty                NUMERIC(12, 2),
    adhocpenaltyreason          VARCHAR(256),
    adhocpenaltycomment         VARCHAR(1024),
    adhocrebatereason           VARCHAR(256),
    adhocrebatecomment          VARCHAR(1024),
    applicationtype             VARCHAR(64),
    dateeffectivefrom           BIGINT,
    locality                    VARCHAR(256),
    isoldapplication            BOOLEAN DEFAULT FALSE,
    additionaldetails           JSONB,
    channel                     VARCHAR(64),
    isdisconnectiontemporary    BOOLEAN DEFAULT FALSE,
    disconnectionreason         VARCHAR(1024),
    createdby                   VARCHAR(64),
    lastmodifiedby              VARCHAR(64),
    createdtime                 BIGINT,
    lastmodifiedtime            BIGINT
);

CREATE INDEX IF NOT EXISTS idx_sw_connection_tenantid ON eg_sw_connection (tenantid);
CREATE INDEX IF NOT EXISTS idx_sw_connection_applicationno ON eg_sw_connection (applicationno);
CREATE INDEX IF NOT EXISTS idx_sw_connection_connectionno ON eg_sw_connection (connectionno);
CREATE INDEX IF NOT EXISTS idx_sw_connection_oldconnectionno ON eg_sw_connection (oldconnectionno);
CREATE INDEX IF NOT EXISTS idx_sw_connection_property_id ON eg_sw_connection (property_id);
CREATE INDEX IF NOT EXISTS idx_sw_connection_applicationstatus ON eg_sw_connection (applicationstatus);
