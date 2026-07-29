# Postman Testing Guide — sw-services & sw-calculator

> **Base URLs**
> - sw-services → `http://localhost:8091`
> - sw-calculator → `http://localhost:8084`
>
> **All requests**: `POST` with `Content-Type: application/json` unless stated.

---

## 👥 Personas Reference

| Persona | `type` field | `roles[0].code` | Realistic Testing UUID |
|---------|-------------|-----------------|------------------------|
| Citizen | `CITIZEN` | `CITIZEN` | `e1a8a25c-cf6e-49b0-94df-6d7b42022b7c` |
| Counter Employee | `EMPLOYEE` | `SW_CEMP` | `5b6d9e03-7b40-42cf-90f7-11116de097ab` |
| Field Inspector | `EMPLOYEE` | `SW_FIELD_INSPECTOR` | `df90c678-bb34-45aa-b892-2222fa102b4d` |
| **Admin / Approver** | `EMPLOYEE` | `SW_APPROVER` | `8d407ead-21c2-4b41-8d52-e80055e0a74e` |
| Super Admin | `EMPLOYEE` | `SUPERUSER` | `8d407ead-21c2-4b41-8d52-e80055e0a74e` |
| Internal System | `SYSTEM` | `SYSTEM` | `00000000-0000-0000-0000-000000000000` |

---

## 🔄 The Workflow Story

```
CITIZEN creates application
        │
        ▼
   [INITIATED]  ← only visible to CITIZEN and employees
        │
   SW_CEMP reviews and forwards
        │
        ▼
[PENDING_FOR_FIELD_INSPECTION]
        │
   SW_FIELD_INSPECTOR does site visit and forwards
        │
        ▼
[PENDING_FOR_APPROVAL]
        │
   SW_APPROVER (ADMIN) gives final approval ← KEY STEP
        │
        ▼
   [APPROVED] + status = ACTIVE ← connection is live
        │
   SW_CEMP runs _calculate → billing demand created
```

> ⚠️ **Before admin approves**: CITIZEN cannot update, calculate, or modify anything.
> ✅ **After admin approves**: Calculation and billing become available.

---

---

# PART 1 — sw-services-go (Port 8091)

---

## STEP 1 — Health Check (Public, No Auth)

**GET** `http://localhost:8091/sw-services/actuator/health`
Body: *(none)*

**Expected Response:**
```json
{
  "status": "UP",
  "checkedAt": "2026-07-28T..."
}
```

---

## STEP 2 — CITIZEN Creates a Connection Application

**POST** `http://localhost:8091/sw-services/swc/_create`

**What happens:** Citizen submits a new sewerage connection application. Status starts as `INITIATED`. Connection is inactive until admin approves.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "citizen-create-001",
    "userInfo": {
      "uuid": "e1a8a25c-cf6e-49b0-94df-6d7b42022b7c",
      "userName": "john_citizen",
      "type": "CITIZEN",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "CITIZEN", "name": "Citizen", "tenantId": "pb" }
      ]
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
        "name": "John Citizen",
        "mobileNumber": "9876543210",
        "gender": "MALE",
        "ownerType": "NONE",
        "isPrimaryOwner": true
      }
    ]
  }
}
```

**Expected Response (200 OK):**
```json
{
  "ResponseInfo": { "status": "successful" },
  "SewerageConnections": [
    {
      "id": "<<SAVE THIS ID>>",
      "applicationNo": "<<SAVE THIS e.g. SW-APP-ca219524>>",
      "applicationStatus": "INITIATED",
      "status": "INACTIVE",
      "tenantId": "pb.amritsar",
      "propertyId": "PB-PT-2026-27-21"
    }
  ]
}
```

> 📝 **Save**: `id` and `applicationNo` — you need them for all following steps.

---

## STEP 3 — CITIZEN Tries to Update → ❌ 403 ERROR (Expected)

**POST** `http://localhost:8091/sw-services/swc/_update`

**What happens:** Citizen tries to advance their own application — **BLOCKED**. Only employees can update.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "citizen-update-attempt",
    "userInfo": {
      "uuid": "e1a8a25c-cf6e-49b0-94df-6d7b42022b7c",
      "userName": "john_citizen",
      "type": "CITIZEN",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "CITIZEN", "name": "Citizen", "tenantId": "pb" }
      ]
    }
  },
  "sewerageConnection": {
    "id": "<<CONNECTION ID FROM STEP 2>>",
    "tenantId": "pb.amritsar"
  }
}
```

**Expected Response (403 Forbidden):**
```json
{
  "Errors": [
    {
      "code": "EG_SW_ACCESS_DENIED",
      "message": "Access denied: insufficient role for this operation"
    }
  ]
}
```

✅ **This 403 is CORRECT** — auth is working.

---

## STEP 4 — SW_CEMP Forwards to Field Inspector

**POST** `http://localhost:8091/sw-services/swc/_update`

**What happens:** Counter Employee reviews and sends application for field inspection. Status → `PENDING_FOR_FIELD_INSPECTION`.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "cemp-update-001",
    "userInfo": {
      "uuid": "5b6d9e03-7b40-42cf-90f7-11116de097ab",
      "userName": "counter_employee",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "SW_CEMP", "name": "SW Counter Employee", "tenantId": "pb" }
      ]
    }
  },
  "sewerageConnection": {
    "id": "<<CONNECTION ID FROM STEP 2>>",
    "tenantId": "pb.amritsar"
  }
}
```

**Expected Response (200 OK):**
```json
{
  "SewerageConnections": [
    {
      "applicationStatus": "PENDING_FOR_FIELD_INSPECTION",
      "status": "INACTIVE"
    }
  ]
}
```

---

## STEP 5 — SW_FIELD_INSPECTOR Sends for Approval

**POST** `http://localhost:8091/sw-services/swc/_update`

**What happens:** Field Inspector visits site and sends application for final approval. Status → `PENDING_FOR_APPROVAL`.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "inspector-update-001",
    "userInfo": {
      "uuid": "df90c678-bb34-45aa-b892-2222fa102b4d",
      "userName": "field_inspector",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "SW_FIELD_INSPECTOR", "name": "SW Field Inspector", "tenantId": "pb" }
      ]
    }
  },
  "sewerageConnection": {
    "id": "<<CONNECTION ID FROM STEP 2>>",
    "tenantId": "pb.amritsar"
  }
}
```

**Expected Response (200 OK):**
```json
{
  "SewerageConnections": [
    {
      "applicationStatus": "PENDING_FOR_APPROVAL",
      "status": "INACTIVE"
    }
  ]
}
```

---

## STEP 6 — SW_APPROVER (ADMIN) Searches for Pending Applications

**POST** `http://localhost:8091/sw-services/swc/_search?tenantId=pb.amritsar&applicationStatus=PENDING_FOR_APPROVAL`

**What happens:** The Admin/Approver searches for a list of all applications that have been processed by inspectors and are now waiting for final approval.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "admin-search",
    "userInfo": {
      "uuid": "8d407ead-21c2-4b41-8d52-e80055e0a74e",
      "userName": "sw_admin",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "SW_APPROVER", "name": "SW Approver", "tenantId": "pb" }
      ]
    }
  }
}
```

**Expected Response (200 OK):**
Returns the list of connections ready to be approved. Copy the `id` from one of the records in `"SewerageConnections"` to use in the approval step below.

---

## STEP 7 — ⭐ SW_APPROVER (ADMIN) Gives Final Approval

**POST** `http://localhost:8091/sw-services/swc/_update`

**What happens:** Admin reviews and **approves** the application. This is the key admin action. Status → `APPROVED`, connection status → `ACTIVE`.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "approver-approve-001",
    "userInfo": {
      "uuid": "8d407ead-21c2-4b41-8d52-e80055e0a74e",
      "userName": "sw_admin",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "SW_APPROVER", "name": "SW Approver", "tenantId": "pb" }
      ]
    }
  },
  "sewerageConnection": {
    "id": "<<CONNECTION ID COPPIED FROM SEARCH OR STEP 2>>",
    "tenantId": "pb.amritsar"
  }
}
```

**Expected Response (200 OK):**
```json
{
  "SewerageConnections": [
    {
      "applicationStatus": "APPROVED",
      "status": "ACTIVE"
    }
  ]
}
```

✅ `status: ACTIVE` means the sewerage connection is now live!

---

## STEP 8 — CITIZEN Searches Their Connection

**POST** `http://localhost:8091/sw-services/swc/_search?tenantId=pb.amritsar&applicationNumber=<<APP_NO>>`

**What happens:** Citizen can view their own connection — now shows APPROVED/ACTIVE status.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "citizen-search-001",
    "userInfo": {
      "uuid": "e1a8a25c-cf6e-49b0-94df-6d7b42022b7c",
      "userName": "john_citizen",
      "type": "CITIZEN",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "CITIZEN", "name": "Citizen", "tenantId": "pb" }
      ]
    }
  }
}
```

**Expected Response (200 OK):** Returns the connection with `applicationStatus: APPROVED`, `status: ACTIVE`.

---

## STEP 9 — CITIZEN Tries Plain Search → ❌ 403 ERROR (Expected)

**POST** `http://localhost:8091/sw-services/swc/_plainsearch?tenantId=pb.amritsar`

**What happens:** Plain search is internal/system use only. Citizen gets blocked.

*(Use same CITIZEN RequestInfo from Step 8)*

**Expected Response (403 Forbidden):**
```json
{
  "Errors": [
    { "code": "EG_SW_ACCESS_DENIED", "message": "Access denied: insufficient role for this operation" }
  ]
}
```

---

---

# PART 2 — sw-calculator-go (Port 8084)

> **Important:** Run these AFTER the connection is `APPROVED` + `ACTIVE` (Step 7 above). Use `applicationNo` from Step 2.

---

## STEP 10 — CITIZEN Gets Fee Estimate (Allowed)

**POST** `http://localhost:8084/sw-calculator/sewerageCalculator/_estimate`

**What happens:** Citizens can preview fees before payment. **Read-only, no billing demand created.**

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "citizen-estimate-001",
    "userInfo": {
      "uuid": "e1a8a25c-cf6e-49b0-94df-6d7b42022b7c",
      "userName": "john_citizen",
      "type": "CITIZEN",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "CITIZEN", "name": "Citizen", "tenantId": "pb" }
      ]
    }
  },
  "CalculationCriteria": [
    {
      "applicationNo": "<<APP_NO FROM STEP 2>>",
      "tenantId": "pb.amritsar"
    }
  ]
}
```

**Expected Response (200 OK):**
```json
{
  "Calculation": [
    {
      "applicationNo": "SW-APP-ca219524",
      "totalAmount": 176,
      "fee": 175,
      "charge": 1,
      "taxHeadEstimates": [
        { "taxHeadCode": "SW_FORM_FEE",       "estimateAmount": 100, "category": "FEE" },
        { "taxHeadCode": "SW_SCRUTINY_FEE",   "estimateAmount": 50,  "category": "FEE" },
        { "taxHeadCode": "SW_OTHER_CHARGE",   "estimateAmount": 25,  "category": "FEE" },
        { "taxHeadCode": "SW_SECURITY_CHARGE","estimateAmount": 1,   "category": "CHARGES" }
      ]
    }
  ]
}
```

---

## STEP 11 — CITIZEN Tries _calculate → ❌ 403 ERROR (Expected)

**POST** `http://localhost:8084/sw-calculator/sewerageCalculator/_calculate`

**What happens:** Only employees/admins can create billing demands. Citizen is blocked.

*(Use same CITIZEN RequestInfo + CalculationCriteria from Step 10)*

**Expected Response (403 Forbidden):**
```json
{
  "Errors": [
    {
      "code": "EG_SW_CALC_ACCESS_DENIED",
      "message": "Access denied: insufficient role for this operation"
    }
  ]
}
```

✅ **This 403 is CORRECT.**

---

## STEP 12 — SW_CEMP Runs Calculate (Creates Billing Demand)

**POST** `http://localhost:8084/sw-calculator/sewerageCalculator/_calculate`

**What happens:** Counter Employee triggers billing demand creation. Demand is saved in billing-service.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "cemp-calculate-001",
    "userInfo": {
      "uuid": "5b6d9e03-7b40-42cf-90f7-11116de097ab",
      "userName": "counter_employee",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "SW_CEMP", "name": "SW Counter Employee", "tenantId": "pb" }
      ]
    }
  },
  "CalculationCriteria": [
    {
      "applicationNo": "<<APP_NO FROM STEP 2>>",
      "tenantId": "pb.amritsar"
    }
  ]
}
```

**Expected Response (200 OK):**
```json
{
  "Calculation": [
    {
      "applicationNo": "SW-APP-ca219524",
      "totalAmount": 176,
      "fee": 175,
      "charge": 1
    }
  ]
}
```

---

## STEP 13 — CITIZEN Tries _applyAdhocTax → ❌ 403 ERROR (Expected)

**POST** `http://localhost:8084/sw-calculator/sewerageCalculator/_applyAdhocTax`

**What happens:** Only Approvers/Admins can apply penalties/rebates. Citizen blocked.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "citizen-adhoc-attempt",
    "userInfo": {
      "uuid": "e1a8a25c-cf6e-49b0-94df-6d7b42022b7c",
      "userName": "john_citizen",
      "type": "CITIZEN",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "CITIZEN", "name": "Citizen", "tenantId": "pb" }
      ]
    }
  },
  "businessService": "SW.ONE_TIME_FEE",
  "consumerCode": "<<APP_NO FROM STEP 2>>",
  "tenantId": "pb.amritsar",
  "adhocPenalty": 500
}
```

**Expected Response (403 Forbidden):**
```json
{
  "Errors": [
    { "code": "EG_SW_CALC_ACCESS_DENIED", "message": "Access denied: insufficient role for this operation" }
  ]
}
```

---

## STEP 14 — ⭐ SW_APPROVER (ADMIN) Applies Adhoc Tax/Rebate

**POST** `http://localhost:8084/sw-calculator/sewerageCalculator/_applyAdhocTax`

**What happens:** Admin applies a penalty (e.g. late payment) and/or a rebate to the billing demand.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "approver-adhoc-001",
    "userInfo": {
      "uuid": "8d407ead-21c2-4b41-8d52-e80055e0a74e",
      "userName": "sw_admin",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "SW_APPROVER", "name": "SW Approver", "tenantId": "pb" }
      ]
    }
  },
  "businessService": "SW.ONE_TIME_FEE",
  "consumerCode": "<<APP_NO FROM STEP 2>>",
  "tenantId": "pb.amritsar",
  "adhocPenalty": 500,
  "adhocRebate": 100
}
```

**Expected Response (200 OK):**
```json
{
  "Calculation": [
    {
      "totalAmount": 400,
      "taxHeadEstimates": [
        { "taxHeadCode": "SW_ADHOC_PENALTY", "estimateAmount": 500,  "category": "PENALTY" },
        { "taxHeadCode": "SW_ADHOC_REBATE",  "estimateAmount": -100, "category": "REBATE" }
      ]
    }
  ]
}
```

---

# 📋 Full Test Checklist for Postman

| # | Service | Role | Endpoint | Expect |
|---|---------|------|----------|--------|
| 1 | sw-services | Public | `GET /actuator/health` | ✅ 200 `UP` |
| 2 | sw-calculator | Public | `GET /actuator/health` | ✅ 200 `UP` |
| 3 | sw-services | CITIZEN | `POST /swc/_create` | ✅ 200 `INITIATED` |
| 4 | sw-services | CITIZEN | `POST /swc/_update` | 🔒 **403** (blocked) |
| 5 | sw-services | SW_CEMP | `POST /swc/_update` | ✅ 200 `PENDING_FOR_FIELD_INSPECTION` |
| 6 | sw-services | SW_FIELD_INSPECTOR | `POST /swc/_update` | ✅ 200 `PENDING_FOR_APPROVAL` |
| 7 | sw-services | **SW_APPROVER** | `POST /swc/_search` | ✅ 200 Admin gets pending list ⭐ |
| 8 | sw-services | **SW_APPROVER** | `POST /swc/_update` | ✅ 200 `APPROVED` + `ACTIVE` ⭐ |
| 9 | sw-services | CITIZEN | `POST /swc/_search` | ✅ 200 sees APPROVED status |
| 10 | sw-services | CITIZEN | `POST /swc/_plainsearch` | 🔒 **403** (blocked) |
| 11 | sw-calculator | CITIZEN | `POST /_estimate` | ✅ 200 totalAmount=176 |
| 12 | sw-calculator | CITIZEN | `POST /_calculate` | 🔒 **403** (blocked) |
| 13 | sw-calculator | SW_CEMP | `POST /_calculate` | ✅ 200 demand created |
| 14 | sw-calculator | CITIZEN | `POST /_applyAdhocTax` | 🔒 **403** (blocked) |
| 15 | sw-calculator | **SW_APPROVER** | `POST /_applyAdhocTax` | ✅ 200 penalty+rebate applied ⭐ |
| 16 | sw-services | SW_CEMP | `POST /swc/_update` | ✅ 200 Modify connection details |
| 17 | sw-services | **SW_APPROVER** | `POST /swc/_update` | ✅ 200 Disconnection request |
| 18 | sw-services | **SW_APPROVER** | `POST /swc/_update` | ✅ 200 Reject application |
| 19 | sw-services | CITIZEN | `POST /swc/_search` | ✅ 200 sees REJECTED status |

---

## 💡 Postman Tips

1. **Set up a collection variable** `conn_id` — after Step 3, copy the `id` from the response and save it.
2. **Set up a collection variable** `app_no` — copy `applicationNo` from Step 3.
3. **Use {{conn_id}}** and **{{app_no}}** in subsequent requests.
4. Steps 3–10 must be run **in order** — each step depends on the previous status.
5. Steps 11–15 (calculator) require Step 8 (admin approval update) to be done first.

---

## ✏️ Workflow: Modify Connection Request

This workflow is used to make updates to an existing connection (e.g., changing the number of toilets, plumbers, or holder metadata).

### STEP 16 — SW_CEMP / EMPLOYEE Modifies Connection Details

**POST** `http://localhost:8091/sw-services/swc/_update`

**JSON Body:**
```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "modify-connection-001",
    "userInfo": {
      "uuid": "5b6d9e03-7b40-42cf-90f7-11116de097ab",
      "userName": "counter_employee",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "SW_CEMP", "name": "SW Counter Employee", "tenantId": "pb" }
      ]
    }
  },
  "sewerageConnection": {
    "id": "<<CONNECTION ID>>",
    "tenantId": "pb.amritsar",
    "noOfToilets": 5,
    "noOfWaterClosets": 3
  }
}
```

**Expected Response (200 OK):**
Returns the updated connection payload showing `"noOfToilets": 5` and `"noOfWaterClosets": 3`.

---

## 🔌 Workflow: Disconnection Request

This workflow is used to request the disconnection of an existing active connection.

### STEP 17 — SW_APPROVER (ADMIN) Requests Disconnection

**POST** `http://localhost:8091/sw-services/swc/_update`

**JSON Body:**
```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "disconnect-request-001",
    "userInfo": {
      "uuid": "8d407ead-21c2-4b41-8d52-e80055e0a74e",
      "userName": "sw_admin",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "SW_APPROVER", "name": "SW Approver", "tenantId": "pb" }
      ]
    }
  },
  "disconnectRequest": true,
  "sewerageConnection": {
    "id": "<<CONNECTION ID>>",
    "tenantId": "pb.amritsar"
  }
}
```

**Expected Response (200 OK):**
```json
{
  "SewerageConnections": [
    {
      "applicationStatus": "PENDING_APPROVAL_FOR_DISCONNECTION",
      "status": "INACTIVE"
    }
  ]
}
```
*(Connection goes INACTIVE and enters disconnection workflow)*

---

## 🗑️ Action: Reject or Cancel an Application

In the DIGIT system, there is no physical `DELETE` endpoint. Instead, the Admin rejects or cancels an application by updating its status.

### STEP 18 — SW_APPROVER (ADMIN) Rejects/Cancels Application

**POST** `http://localhost:8091/sw-services/swc/_update`

**JSON Body:**
```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "admin-reject-001",
    "userInfo": {
      "uuid": "8d407ead-21c2-4b41-8d52-e80055e0a74e",
      "userName": "sw_admin",
      "type": "EMPLOYEE",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "SW_APPROVER", "name": "SW Approver", "tenantId": "pb" }
      ]
    }
  },
  "sewerageConnection": {
    "id": "<<CONNECTION ID>>",
    "tenantId": "pb.amritsar",
    "applicationStatus": "REJECTED"
  }
}
```
*(You can use either `"REJECTED"` or `"CANCELLED"` as the `applicationStatus`)*

**Expected Response (200 OK):**
```json
{
  "SewerageConnections": [
    {
      "applicationStatus": "REJECTED",
      "status": "INACTIVE"
    }
  ]
}
```

### 💡 Database Hard Delete (For Testing Cleanup Only)
If you need to wipe out the connection from your database completely to reuse property ids or clear space:

```sql
-- Connect to Postgres database (Port 35432, Database: rainmaker)
DELETE FROM eg_sw_connection WHERE id = '<<CONNECTION ID>>';
```

---

### STEP 19 — CITIZEN Searches Connection (Sees REJECTED Status)

**POST** `http://localhost:8091/sw-services/swc/_search?tenantId=pb.amritsar&applicationNumber=<<APP_NO>>`

**What happens:** When the Citizen searches for their connection application status after the Admin rejected it, they see it reflected as `REJECTED` and `INACTIVE`.

```json
{
  "RequestInfo": {
    "apiId": "Rainmaker",
    "ver": "01",
    "ts": 1720000000000,
    "msgId": "citizen-search-rejected",
    "userInfo": {
      "uuid": "e1a8a25c-cf6e-49b0-94df-6d7b42022b7c",
      "userName": "john_citizen",
      "type": "CITIZEN",
      "tenantId": "pb.amritsar",
      "roles": [
        { "code": "CITIZEN", "name": "Citizen", "tenantId": "pb" }
      ]
    }
  }
}
```

**Expected Response (200 OK):**
The connection payload returns with `"applicationStatus": "REJECTED"` and `"status": "INACTIVE"`.
