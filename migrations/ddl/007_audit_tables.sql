-- Audit-mirror tables, created but not populated in Phase 1 — Java doesn't
-- populate these from the sw-services code path either (it happens in the
-- persister/indexer pipeline, out of scope until Phase 3's Kafka work).
CREATE TABLE IF NOT EXISTS eg_sw_connection_audit (
    LIKE eg_sw_connection INCLUDING DEFAULTS
);

CREATE TABLE IF NOT EXISTS eg_sw_service_audit (
    LIKE eg_sw_service INCLUDING DEFAULTS
);
