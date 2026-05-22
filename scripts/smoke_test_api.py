#!/usr/bin/env python3
"""Full API smoke test — all /api/v1 routes."""
import json
import sys
import time
import urllib.error
import urllib.request
from dataclasses import dataclass, asdict
from datetime import datetime, timedelta, timezone

BASE = sys.argv[1] if len(sys.argv) > 1 else "https://api-production-5e10.up.railway.app/api/v1"
OUT = sys.argv[2] if len(sys.argv) > 2 else ""


@dataclass
class Result:
    name: str
    method: str
    path: str
    status: int
    ok: bool
    note: str = ""


results: list[Result] = []
tokens: dict[str, str] = {}
ids: dict[str, str] = {}


def log(msg: str):
    print(msg, flush=True)
    if OUT:
        with open(OUT, "a", encoding="utf-8") as f:
            f.write(msg + "\n")


def req(method: str, path: str, token: str | None = None, body: dict | None = None):
    url = BASE + path
    data = json.dumps(body).encode() if body is not None else None
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    r = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(r, timeout=90) as resp:
            raw = resp.read().decode()
            try:
                return resp.status, json.loads(raw) if raw else {}
            except json.JSONDecodeError:
                return resp.status, raw
    except urllib.error.HTTPError as e:
        raw = e.read().decode()
        try:
            return e.code, json.loads(raw) if raw else {}
        except json.JSONDecodeError:
            return e.code, raw


def record(name: str, method: str, path: str, status: int, ok: bool | None = None, note: str = ""):
    if ok is None:
        ok = 200 <= status < 300
    r = Result(name, method, path, status, ok, note)
    results.append(r)
    mark = "PASS" if ok else "FAIL"
    log(f"{mark} {method} {path} -> {status} {note}")


def login(role: str, email: str):
    status, data = req("POST", "/auth/login", body={"email": email, "password": "Password123!"})
    record(f"POST /auth/login ({role})", "POST", "/auth/login", status)
    if status == 200 and isinstance(data, dict):
        tokens[role] = data["access_token"]
        ids[f"{role}_user"] = data.get("user_id", "")


def main():
    if OUT:
        open(OUT, "w", encoding="utf-8").write(f"API smoke test {datetime.now().isoformat()}\nBase: {BASE}\n\n")

    log(f"Testing {BASE}\n")

    # --- Public / auth ---
    s, _ = req("GET", "/health")
    record("GET /health", "GET", "/health", s)

    s, _ = req("POST", "/auth/seed-demo")
    record("POST /auth/seed-demo", "POST", "/auth/seed-demo", s)

    login("citizen", "citizen@example.com")
    login("admin", "admin@example.com")
    login("examiner", "examiner@example.com")
    login("officer", "officer@example.com")

    c, a, e, o = tokens["citizen"], tokens["admin"], tokens["examiner"], tokens["officer"]

    ts = int(time.time())
    s, reg = req("POST", "/auth/register", body={
        "full_name": "Smoke User",
        "email": f"smoke{ts}@test.local",
        "phone": "+10000000099",
        "password": "Password123!",
    })
    record("POST /auth/register", "POST", "/auth/register", s)
    if s in (200, 201) and isinstance(reg, dict):
        ids["verify_token"] = reg.get("verification_token", "")
        ids["reset_token"] = reg.get("reset_token", "")

    s, forgot = req("POST", "/auth/forgot-password", body={"email": "citizen@example.com"})
    record("POST /auth/forgot-password", "POST", "/auth/forgot-password", s)
    if s == 200 and isinstance(forgot, dict) and forgot.get("reset_token"):
        ids["reset_token"] = forgot["reset_token"]
        s2, _ = req("POST", "/auth/reset-password", body={
            "token": ids["reset_token"],
            "new_password": "Password123!",
        })
        record("POST /auth/reset-password", "POST", "/auth/reset-password", s2, ok=(s2 in (200, 400)))

    if ids.get("verify_token"):
        s, _ = req("POST", "/auth/verify-email", body={"token": ids["verify_token"]})
        record("POST /auth/verify-email", "POST", "/auth/verify-email", s, ok=(s in (200, 400)))
    else:
        s, _ = req("POST", "/auth/verify-email", body={"token": "invalid"})
        record("POST /auth/verify-email", "POST", "/auth/verify-email", s, ok=(s in (200, 400)))

    s, _ = req("POST", "/auth/bootstrap-admin", body={
        "secret": "wrong", "full_name": "X", "last_name": "Y",
        "email": "x@test.com", "password": "Password123!",
    })
    record("POST /auth/bootstrap-admin", "POST", "/auth/bootstrap-admin", s, ok=(s in (403, 400)))

    s, _ = req("GET", "/admin/ping", c)
    record("GET /admin/ping (citizen expect 403)", "GET", "/admin/ping", s, ok=(s == 403))

    s, _ = req("GET", "/me", c)
    record("GET /me", "GET", "/me", s)

    # --- Citizen GET ---
    for path, name in [
        ("/notifications", "GET /notifications"),
        ("/activity", "GET /activity"),
        ("/licenses/me", "GET /licenses/me"),
        ("/vehicles/me", "GET /vehicles/me"),
        ("/payments/history", "GET /payments/history"),
        ("/exam/history", "GET /exam/history"),
        ("/inspection", "GET /inspection"),
        ("/violations", "GET /violations (citizen)"),
        ("/identity/status", "GET /identity/status"),
        ("/centers", "GET /centers"),
    ]:
        s, data = req("GET", path, c)
        record(name, "GET", path, s)
        if path == "/notifications" and s == 200 and isinstance(data, list) and data:
            ids["notification"] = data[0]["id"]
        if path == "/licenses/me" and s == 200 and isinstance(data, list):
            for lic in data:
                if lic.get("status") == "submitted" and lic.get("id"):
                    ids["license_pending"] = lic["id"]
                if lic.get("status") == "issued" and lic.get("id"):
                    ids["license_issued"] = lic["id"]
        if path == "/vehicles/me" and s == 200 and isinstance(data, list) and data:
            if data[0].get("id"):
                ids["vehicle"] = data[0]["id"]
        if path == "/payments/history" and s == 200 and isinstance(data, list) and data:
            if data[0].get("id"):
                ids["payment"] = data[0]["id"]
        if path == "/inspection" and s == 200 and isinstance(data, list) and data:
            if data[0].get("id"):
                ids["inspection"] = data[0]["id"]

    # --- Citizen POST ---
    s, data = req("POST", "/identity/submit", c, {
        "national_id_number": f"ID-{ts}",
        "document_front_path": "/uploads/front.pdf",
        "document_back_path": "/uploads/back.pdf",
        "selfie_path": "/uploads/selfie.jpg",
    })
    record("POST /identity/submit", "POST", "/identity/submit", s)
    if s == 200 and isinstance(data, dict) and data.get("id"):
        ids["identity"] = data["id"]

    s, data = req("POST", "/licenses", c, {
        "name": "Smoke Test",
        "dob": "1998-06-15",
        "gender": "Male",
        "nationality": "Testland",
        "address": "1 Main St",
        "city": "Capital",
        "postal": "10001",
        "license_type": "Car license",
    })
    record("POST /licenses", "POST", "/licenses", s)
    if s in (200, 201) and isinstance(data, dict):
        lic = data.get("data") or data
        if isinstance(lic, dict) and lic.get("id"):
            ids["license_new"] = lic["id"]

    vin = f"SMK{ts % 100000000:08d}"
    s, data = req("POST", "/vehicles", c, {
        "vin": vin, "plate": f"P-{vin[-4:]}",
        "make": "Toyota", "model": "Corolla", "year": 2022, "color": "White",
    })
    record("POST /vehicles", "POST", "/vehicles", s)
    if s in (200, 201) and isinstance(data, dict) and data.get("id"):
        ids["vehicle"] = data["id"]
    if not ids.get("vehicle"):
        s, data = req("GET", "/admin/vehicles", a)
        if s == 200 and isinstance(data, list):
            for v in data:
                if v.get("id"):
                    ids["vehicle"] = v["id"]
                    break

    s, _ = req("POST", "/payments/initiate", c, {"service_type": "license"})
    record("POST /payments/initiate", "POST", "/payments/initiate", s)

    s, data = req("POST", "/exam/start", c, {"license_type": "car"})
    record("POST /exam/start", "POST", "/exam/start", s, ok=(s in (200, 201, 400)))
    if s in (200, 201) and isinstance(data, dict):
        attempt = data.get("attempt") or {}
        attempt_id = attempt.get("id")
        questions = data.get("questions") or []
        if attempt_id and questions:
            answers = []
            for item in questions:
                qobj = item.get("question") or item
                opts = item.get("options") or []
                qid = qobj.get("id")
                if qid and opts:
                    correct = next((o for o in opts if o.get("is_correct")), opts[0])
                    answers.append({
                        "question_id": qid,
                        "selected_option_id": correct.get("id"),
                    })
            s2, _ = req("POST", f"/exam/{attempt_id}/submit", c, {"answers": answers})
            record("POST /exam/{id}/submit", "POST", f"/exam/{attempt_id}/submit", s2, ok=(s2 in (200, 201, 400)))

    s, data = req("GET", "/centers", c)
    if s == 200 and isinstance(data, list) and data:
        cid = data[0]["id"]
        s2, slots = req("GET", f"/centers/{cid}/slots", c)
        record("GET /centers/{id}/slots", "GET", f"/centers/{cid}/slots", s2)
        if s2 == 200 and isinstance(slots, list):
            for sl in slots:
                if sl.get("booked", 0) < sl.get("capacity", 1):
                    s3, book = req("POST", "/practical/book", c, {"slot_id": sl["id"]})
                    record("POST /practical/book", "POST", "/practical/book", s3)
                    if s3 in (200, 201) and isinstance(book, dict):
                        ids["booking"] = book.get("id", "")
                    break

    if ids.get("vehicle"):
        s, insp = req("POST", "/inspection/schedule", c, {
            "vehicle_id": ids["vehicle"],
            "inspection_date": "2026-07-01",
        })
        record("POST /inspection/schedule", "POST", "/inspection/schedule", s)
        if s in (200, 201) and isinstance(insp, dict) and insp.get("id"):
            ids["inspection"] = insp["id"]
    else:
        record("POST /inspection/schedule", "POST", "/inspection/schedule", 0, ok=False, note="no vehicle id")

    if ids.get("inspection"):
        s, _ = req("POST", f"/inspection/{ids['inspection']}/upload-report", c, {
            "report_path": "/reports/smoke.pdf",
            "status": "passed",
        })
        record("POST /inspection/{id}/upload-report", "POST", f"/inspection/{ids['inspection']}/upload-report", s)

    if ids.get("vehicle"):
        now = datetime.now(timezone.utc)
        s, _ = req("POST", "/devices/data", c, {
            "device_serial": f"DEV-{ts}",
            "vehicle_id": ids["vehicle"],
            "user_id": ids.get("citizen_user", ""),
            "trip": {
                "start_time": (now - timedelta(hours=1)).isoformat().replace("+00:00", "Z"),
                "end_time": now.isoformat().replace("+00:00", "Z"),
                "distance": 15.0,
                "average_speed": 50.0,
                "safety_score": 85.0,
            },
            "events": [{
                "event_type": "harsh_brake",
                "severity": "low",
                "timestamp": now.isoformat().replace("+00:00", "Z"),
            }],
        })
        record("POST /devices/data", "POST", "/devices/data", s)
        s, _ = req("GET", f"/monitoring/trips/{ids['vehicle']}", c)
        record("GET /monitoring/trips/{vehicleId}", "GET", f"/monitoring/trips/{ids['vehicle']}", s)
    else:
        record("POST /devices/data", "POST", "/devices/data", 0, ok=False, note="skip")

    if ids.get("citizen_user"):
        s, _ = req("GET", f"/monitoring/score/{ids['citizen_user']}", c)
        record("GET /monitoring/score/{userId}", "GET", f"/monitoring/score/{ids['citizen_user']}", s)

    if ids.get("notification"):
        s, _ = req("PATCH", f"/notifications/{ids['notification']}/read", c)
        record("PATCH /notifications/{id}/read", "PATCH", f"/notifications/{ids['notification']}/read", s, ok=(s in (200, 400)))
    else:
        s, _ = req("PATCH", "/notifications/000000000000000000000099/read", c)
        record("PATCH /notifications/{id}/read", "PATCH", "/notifications/.../read", s, ok=(s in (400, 404)))

    if ids.get("vehicle"):
        s, tr = req("POST", f"/vehicles/{ids['vehicle']}/transfer", c, {"buyer_email": "admin@example.com"})
        record("POST /vehicles/{id}/transfer", "POST", f"/vehicles/{ids['vehicle']}/transfer", s, ok=(s in (200, 201, 400)))
        if s in (200, 201) and isinstance(tr, dict) and tr.get("id"):
            ids["transfer"] = tr["id"]

    # --- Admin ---
    for path, name in [
        ("/admin/users", "GET /admin/users"),
        ("/admin/applications", "GET /admin/applications"),
        ("/admin/audit-logs", "GET /admin/audit-logs"),
        ("/admin/vehicles", "GET /admin/vehicles"),
        ("/admin/analytics/overview", "GET /admin/analytics/overview"),
        ("/admin/analytics/revenue", "GET /admin/analytics/revenue"),
        ("/admin/analytics/exams", "GET /admin/analytics/exams"),
        ("/admin/analytics/trends", "GET /admin/analytics/trends"),
        ("/exam/questions", "GET /exam/questions"),
        ("/admin/ping", "GET /admin/ping"),
    ]:
        s, _ = req("GET", path, a)
        record(name, "GET", path, s)

    if ids.get("identity"):
        s, _ = req("PUT", f"/admin/identity/{ids['identity']}/approve", a, {"comment": "approved"})
        record("PUT /admin/identity/{id}/approve", "PUT", f"/admin/identity/{ids['identity']}/approve", s, ok=(s in (200, 400)))

    lid_reject = ids.get("license_new") or ids.get("license_pending")
    if lid_reject:
        s, _ = req("PUT", f"/admin/licenses/{lid_reject}/reject", a)
        record("PUT /admin/licenses/{id}/reject", "PUT", f"/admin/licenses/{lid_reject}/reject", s, ok=(s in (200, 400)))

    s, lics = req("GET", "/licenses/me", c)
    if s == 200 and isinstance(lics, list):
        for lic in lics:
            if lic.get("status") == "submitted" and lic.get("id"):
                s2, _ = req("PUT", f"/admin/licenses/{lic['id']}/approve", a)
                record("PUT /admin/licenses/{id}/approve", "PUT", f"/admin/licenses/{lic['id']}/approve", s2)
                ids["license_issued"] = lic["id"]
                break

    if ids.get("license_issued"):
        s, _ = req("PUT", f"/licenses/{ids['license_issued']}/renew", c)
        record("PUT /licenses/{id}/renew", "PUT", f"/licenses/{ids['license_issued']}/renew", s, ok=(s in (200, 400)))

    if ids.get("payment"):
        s, _ = req("PUT", f"/admin/payments/{ids['payment']}/mark-paid", a)
        record("PUT /admin/payments/{id}/mark-paid", "PUT", f"/admin/payments/{ids['payment']}/mark-paid", s, ok=(s in (200, 400)))

    if ids.get("transfer"):
        s, _ = req("PUT", f"/admin/transfer/{ids['transfer']}/approve", a)
        record("PUT /admin/transfer/{id}/approve", "PUT", f"/admin/transfer/{ids['transfer']}/approve", s, ok=(s in (200, 400)))

    # --- Officer ---
    s, _ = req("GET", "/violations", o)
    record("GET /violations (officer)", "GET", "/violations", s)

    s, vdata = req("POST", "/violations", o, {
        "driver": "Test Driver",
        "driver_id": ids.get("citizen_user", ""),
        "type": "Speeding",
        "severity": "Minor",
    })
    record("POST /violations", "POST", "/violations", s)
    if s in (200, 201) and isinstance(vdata, dict) and vdata.get("id") and len(vdata["id"]) == 24:
        s2, _ = req("PUT", f"/violations/{vdata['id']}/status", o, {"status": "paid"})
        record("PUT /violations/{id}/status", "PUT", f"/violations/{vdata['id']}/status", s2, ok=(s2 in (200, 400)))

    # --- Examiner ---
    if ids.get("booking"):
        s, _ = req("PUT", f"/examiner/practical/{ids['booking']}/result", e, {
            "result": "pass",
            "comments": "smoke test pass",
        })
        record("PUT /examiner/practical/{id}/result", "PUT", f"/examiner/practical/{ids['booking']}/result", s, ok=(s in (200, 400)))

    # --- Summary ---
    passed = sum(1 for r in results if r.ok)
    failed = [r for r in results if not r.ok]
    log(f"\n=== SUMMARY: {passed}/{len(results)} passed, {len(failed)} failed ===")
    for r in failed:
        log(f"  FAIL {r.method} {r.path} -> {r.status} {r.note}")

    if OUT:
        with open(OUT.replace(".txt", "-summary.json") if OUT.endswith(".txt") else OUT + ".json", "w", encoding="utf-8") as f:
            json.dump([asdict(r) for r in results], f, indent=2)

    return 0 if not failed else 1


if __name__ == "__main__":
    sys.exit(main())
