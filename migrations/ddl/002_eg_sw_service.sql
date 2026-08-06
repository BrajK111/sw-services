CREATE TABLE IF NOT EXISTS eg_sw_service (
    connection_id               VARCHAR(64) PRIMARY KEY REFERENCES eg_sw_connection (id) ON DELETE CASCADE,
    connectionexecutiondate     BIGINT,
    disconnectionexecutiondate  BIGINT,
    noofwaterclosets            INTEGER,
    nooftoilets                 INTEGER,
    connectiontype              VARCHAR(64),
    proposedwaterclosets        INTEGER,
    proposedtoilets             INTEGER,
    appcreateddate              BIGINT,
    detailsprovidedby           VARCHAR(64),
    estimationfilestoreid       VARCHAR(256),
    sanctionfilestoreid         VARCHAR(256),
    estimationletterdate        BIGINT,
    createdby                   VARCHAR(64),
    lastmodifiedby              VARCHAR(64),
    createdtime                 BIGINT,
    lastmodifiedtime            BIGINT
);

CREATE INDEX IF NOT EXISTS idx_sw_service_appcreateddate ON eg_sw_service (appcreateddate);
