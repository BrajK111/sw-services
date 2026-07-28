# SW-Services-Go — Live API Testing Guide
## CDPI Internship Project | DIGIT Platform (Go Migration)

> Service  : sw-services-go — Sewerage Connection Service (Go rewrite of Java DIGIT microservice)
> Base URL : http://localhost:8091
> Status   : All 5 endpoints live and returning 200 OK

---

## Infrastructure Running (Docker)

| Container              | Role                  | Port  |
|------------------------|-----------------------|-------|
| sw-postgres            | PostgreSQL DB         | 35432 |
| sw-kafka               | Kafka broker          | 39092 |
| sw-redis               | Redis cache           | 36379 |
| sw-egov-mdms-service   | MDMS (master data)    | 3456  |
| sw-egov-idgen          | ID generation         | 3457  |
| sw-property-services   | Property validation   | 3466  |

---

## ENDPOINT 1 — Health Check

Method : GET
URL    : http://localhost:8091/sw-services/actuator/health
Body   : (none)

What we do     : Ping the service to confirm it is alive and running.
What we get    : 200 OK — status "UP" with a timestamp.

---

## ENDPOINT 2 — Search Sewerage Connections

Method       : POST
URL          : http://localhost:8091/sw-services/swc/_search?tenantId=pb.amritsar
Content-Type : application/json

Request Body:
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "test-msg-1",
    "userInfo": {
      "uuid": "test-uuid",
      "userName": "testuser",
      "type": "SYSTEM",
      "tenantId": "pb.amritsar"
    }
  }
}

Optional Query Params (add to URL):
  tenantId=pb.amritsar        (required)
  connectionNumber=...
  applicationNumber=...
  status=ACTIVE
  limit=10
  offset=0

What we do     : Search all sewerage connections for a given tenant.
What we get    : 200 OK — list of SewerageConnection objects with their IDs,
                 application numbers, status, property IDs, and connection details.

---

## ENDPOINT 3 — Plain Search

Method       : POST
URL          : http://localhost:8091/sw-services/swc/_plainsearch?tenantId=pb.amritsar
Content-Type : application/json

Request Body:
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "test-msg-2",
    "userInfo": {
      "uuid": "test-uuid",
      "userName": "testuser",
      "type": "SYSTEM",
      "tenantId": "pb.amritsar"
    }
  }
}

What we do     : Same as search but used internally by sw-calculator to
                 fetch connection data before computing bills and demands.
What we get    : 200 OK — same list of connections as _search.

---

## ENDPOINT 4 — Create New Sewerage Connection

Method       : POST
URL          : http://localhost:8091/sw-services/swc/_create
Content-Type : application/json

Request Body:
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "test-create-1",
    "userInfo": {
      "uuid": "test-uuid",
      "userName": "testuser",
      "type": "CITIZEN",
      "tenantId": "pb.amritsar"
    }
  },
  "sewerageConnection": {
    "tenantId": "pb.amritsar",
    "propertyId": "PB-PT-2026-27-21",
    "applicationType": "NEW_CONNECTION",
    "connectionType": "Metered",
    "noOfToilets": 2,
    "noOfWaterClosets": 1,
    "proposedToilets": 2,
    "proposedWaterClosets": 1,
    "channel": "COUNTER",
    "connectionHolders": [
      {
        "name": "Test User",
        "mobileNumber": "9999999999",
        "gender": "MALE",
        "ownerType": "NONE",
        "isPrimaryOwner": true
      }
    ]
  }
}

What we do     : Creates a brand new sewerage connection application.
                 Internally calls property-services to validate the property,
                 MDMS to validate connection type, and idgen to auto-generate
                 an application number. Then saves to PostgreSQL.
What we get    : 200 OK — the newly created connection with a system-generated
                 UUID and application number (e.g. SW-APP-xxxxx),
                 status set to INITIATED.

---

## ENDPOINT 5 — Update Sewerage Connection

Method       : POST
URL          : http://localhost:8091/sw-services/swc/_update
Content-Type : application/json

Real IDs from live database:
  id          : a5726086-8a12-4483-9abd-8073dfd091fa
  applicationNo : SW-APP-49b1525e
  propertyId  : PB-PT-2026-27-21

Request Body:
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "test-update-1",
    "userInfo": {
      "uuid": "test-uuid",
      "userName": "testuser",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar"
    }
  },
  "sewerageConnection": {
    "id": "a5726086-8a12-4483-9abd-8073dfd091fa",
    "applicationNo": "SW-APP-49b1525e",
    "tenantId": "pb.amritsar",
    "propertyId": "PB-PT-2026-27-21",
    "connectionType": "Metered",
    "noOfToilets": 3,
    "noOfWaterClosets": 2
  }
}

What we do     : Updates an existing sewerage connection record by its ID.
                 Merges the new field values into the existing DB record.
What we get    : 200 OK — the updated connection record reflecting the
                 new noOfToilets and noOfWaterClosets values.

---

## QUICK REFERENCE — All Endpoints

  #  | Method | URL
-----|--------|---------------------------------------------------
  1  | GET    | http://localhost:8091/sw-services/actuator/health
  2  | POST   | http://localhost:8091/sw-services/swc/_search?tenantId=pb.amritsar
  3  | POST   | http://localhost:8091/sw-services/swc/_plainsearch?tenantId=pb.amritsar
  4  | POST   | http://localhost:8091/sw-services/swc/_create
  5  | POST   | http://localhost:8091/sw-services/swc/_update

---

## Notes for Mentors

- Full Go rewrite of the Java DIGIT sw-services microservice
- Framework  : Gin (HTTP router) + GORM (ORM)
- Database   : PostgreSQL (rainmaker DB, same schema as Java original)
- API parity : 100% request/response contract match with Java service
- Peer calls : property-services, MDMS, idgen — all live and integrated
- Port       : 8091 (same as Java original)
- All 5 endpoints tested and verified — 200 OK

---
Generated: 2026-07-15 | CDPI Internship — DIGIT Platform Go Migration
