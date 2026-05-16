#!/usr/bin/env python3
"""Regenerate docs/swagger.json from the live API route map."""
import json
from pathlib import Path

ERR = {"$ref": "#/definitions/http.ErrorResponse"}
SEC = [{"BearerAuth": []}]

def op(tags, summary, method_responses, security=None, consumes=None, produces=None, parameters=None):
    o = {"tags": tags, "summary": summary, "responses": method_responses}
    if security:
        o["security"] = security
    if consumes:
        o["consumes"] = consumes
    if produces:
        o["produces"] = produces
    if parameters:
        o["parameters"] = parameters
    return o

def body_param(ref, required=True):
    return {"in": "body", "name": "body", "required": required, "schema": {"$ref": f"#/definitions/{ref}"}}

def path_param(name, desc="MongoDB ObjectID (hex)"):
    return {"name": name, "in": "path", "required": True, "type": "string", "description": desc}

def ok(ref):
    return {"200": {"description": "OK", "schema": {"$ref": f"#/definitions/{ref}"}}}

def created(ref):
    return {"201": {"description": "Created", "schema": {"$ref": f"#/definitions/{ref}"}}}

def std_errors(*codes):
    r = {}
    for c in codes:
        r[str(c)] = {"description": {400: "Bad Request", 401: "Unauthorized", 403: "Forbidden", 404: "Not Found"}.get(c, "Error"), "schema": ERR}
    return r

spec = {
    "swagger": "2.0",
    "info": {
        "title": "Driving Authority API",
        "description": "Full REST API: auth, identity, licenses, theory exams, practical tests, vehicles, inspections, monitoring, payments, analytics.",
        "version": "1.0",
        "contact": {"name": "API Support", "email": "support@example.com"},
    },
    "host": "api-production-5e10.up.railway.app",
    "basePath": "/api/v1",
    "schemes": ["https", "http"],
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Bearer {JWT access_token}",
        }
    },
    "tags": [
        {"name": "system"},
        {"name": "auth"},
        {"name": "identity"},
        {"name": "licenses"},
        {"name": "exams"},
        {"name": "practical"},
        {"name": "vehicles"},
        {"name": "inspections"},
        {"name": "monitoring"},
        {"name": "payments"},
        {"name": "analytics"},
        {"name": "admin"},
    ],
    "paths": {},
    "definitions": {},
}

P = spec["paths"]
D = spec["definitions"]

# --- definitions ---
D.update({
    "http.ErrorResponse": {"type": "object", "properties": {"error": {"type": "string"}}},
    "http.HealthResponse": {"type": "object", "properties": {"ok": {"type": "boolean"}}},
    "http.MeResponse": {"type": "object", "properties": {
        "id": {"type": "string"}, "email": {"type": "string"},
        "role": {"type": "string", "enum": ["citizen", "admin", "examiner", "officer"]},
    }},
    "http.AdminPingResponse": {"type": "object", "properties": {"admin": {"type": "boolean"}}},
    "auth.RegisterInput": {"type": "object", "required": ["first_name", "last_name", "email", "password"], "properties": {
        "first_name": {"type": "string"}, "last_name": {"type": "string"}, "email": {"type": "string", "format": "email"},
        "password": {"type": "string", "minLength": 8}, "phone": {"type": "string"},
    }},
    "auth.LoginInput": {"type": "object", "required": ["email", "password"], "properties": {
        "email": {"type": "string"}, "password": {"type": "string"},
    }},
    "auth.AuthOutput": {"type": "object", "properties": {
        "access_token": {"type": "string"}, "user_id": {"type": "string"}, "email": {"type": "string"}, "role": {"type": "string"},
        "verification_token": {"type": "string"},
    }},
    "auth.VerifyEmailInput": {"type": "object", "required": ["token"], "properties": {"token": {"type": "string"}}},
    "auth.ForgotPasswordInput": {"type": "object", "required": ["email"], "properties": {"email": {"type": "string"}}},
    "auth.ForgotPasswordOutput": {"type": "object", "properties": {"message": {"type": "string"}, "reset_token": {"type": "string"}}},
    "auth.ResetPasswordInput": {"type": "object", "required": ["token", "new_password"], "properties": {
        "token": {"type": "string"}, "new_password": {"type": "string", "minLength": 8},
    }},
    "auth.BootstrapAdminInput": {"type": "object", "required": ["secret", "first_name", "last_name", "email", "password"], "properties": {
        "secret": {"type": "string"}, "first_name": {"type": "string"}, "last_name": {"type": "string"},
        "email": {"type": "string"}, "password": {"type": "string"}, "phone": {"type": "string"},
    }},
    "identity.SubmitInput": {"type": "object", "required": ["national_id_number", "document_front_path", "document_back_path", "selfie_path"], "properties": {
        "national_id_number": {"type": "string"}, "document_front_path": {"type": "string"},
        "document_back_path": {"type": "string"}, "selfie_path": {"type": "string"},
    }},
    "identity.DecisionInput": {"type": "object", "properties": {"comment": {"type": "string"}}},
    "identity.IdentityVerification": {"type": "object", "properties": {
        "id": {"type": "string"}, "user_id": {"type": "string"}, "national_id_number": {"type": "string"},
        "document_front_path": {"type": "string"}, "document_back_path": {"type": "string"}, "selfie_path": {"type": "string"},
        "status": {"type": "string", "enum": ["", "pending", "approved", "rejected"]},
        "reviewed_by": {"type": "string"}, "review_comment": {"type": "string"},
        "submitted_at": {"type": "string", "format": "date-time"}, "reviewed_at": {"type": "string", "format": "date-time"},
    }},
    "licenses.CreateInput": {"type": "object", "required": ["type"], "properties": {
        "type": {"type": "string", "enum": ["motorcycle", "car", "truck", "bus"]},
    }},
    "licenses.License": {"type": "object", "properties": {
        "id": {"type": "string"}, "user_id": {"type": "string"}, "license_number": {"type": "string"},
        "type": {"type": "string"}, "status": {"type": "string"},
        "issue_date": {"type": "string", "format": "date-time"}, "expiry_date": {"type": "string", "format": "date-time"},
    }},
    "exams.StartInput": {"type": "object", "required": ["license_type"], "properties": {
        "license_type": {"type": "string", "enum": ["motorcycle", "car", "truck", "bus"]},
    }},
    "exams.StartOutput": {"type": "object", "properties": {"attempt": {"type": "object"}, "questions": {"type": "array", "items": {"type": "object"}}}},
    "exams.SubmitInput": {"type": "object", "required": ["answers"], "properties": {
        "answers": {"type": "array", "items": {"$ref": "#/definitions/exams.AnswerInput"}},
    }},
    "exams.AnswerInput": {"type": "object", "properties": {
        "question_id": {"type": "string"}, "selected_option_id": {"type": "string"},
    }},
    "exams.SubmitOutput": {"type": "object", "properties": {"attempt": {"type": "object"}, "passed": {"type": "boolean"}}},
    "practical.BookInput": {"type": "object", "required": ["slot_id"], "properties": {"slot_id": {"type": "string"}}},
    "practical.ResultInput": {"type": "object", "properties": {
        "result": {"type": "string", "enum": ["pass", "fail"]}, "comments": {"type": "string"}, "dangerous_action": {"type": "boolean"},
    }},
    "vehicles.CreateInput": {"type": "object", "required": ["vin", "plate_number", "make", "model", "year"], "properties": {
        "vin": {"type": "string"}, "plate_number": {"type": "string"}, "make": {"type": "string"}, "model": {"type": "string"}, "year": {"type": "integer"},
    }},
    "vehicles.TransferInput": {"type": "object", "required": ["buyer_id"], "properties": {"buyer_id": {"type": "string"}}},
    "inspections.ScheduleInput": {"type": "object", "required": ["vehicle_id"], "properties": {"vehicle_id": {"type": "string"}}},
    "inspections.UploadReportInput": {"type": "object", "required": ["report_path"], "properties": {"report_path": {"type": "string"}}},
    "monitoring.DeviceDataInput": {"type": "object", "properties": {
        "device_serial": {"type": "string"}, "vehicle_id": {"type": "string"}, "speed": {"type": "number"},
        "events": {"type": "array", "items": {"type": "object"}},
    }},
    "payments.InitiateInput": {"type": "object", "required": ["service_type"], "properties": {
        "service_type": {"type": "string", "enum": ["theory_exam", "license", "inspection", "transfer"]},
    }},
    "payments.Payment": {"type": "object", "properties": {
        "id": {"type": "string"}, "service_type": {"type": "string"}, "amount": {"type": "number"},
        "status": {"type": "string"}, "transaction_id": {"type": "string"},
    }},
    "analytics.Overview": {"type": "object", "properties": {
        "total_users": {"type": "integer"}, "active_licenses": {"type": "integer"},
        "pending_identity": {"type": "integer"}, "total_vehicles": {"type": "integer"},
    }},
    "analytics.Revenue": {"type": "object", "properties": {"total_revenue": {"type": "number"}, "paid_count": {"type": "integer"}}},
    "analytics.Exams": {"type": "object", "properties": {
        "total_attempts": {"type": "integer"}, "passed": {"type": "integer"}, "failed": {"type": "integer"},
    }},
})

json_obj = {"type": "object"}
arr_lic = {"type": "array", "items": {"$ref": "#/definitions/licenses.License"}}
arr_pay = {"type": "array", "items": {"$ref": "#/definitions/payments.Payment"}}

# system
P["/health"] = {"get": op(["system"], "Health check", {**ok("http.HealthResponse")}, produces=["application/json"])}

# auth (public)
for path, ref_in, ref_out, code in [
    ("/auth/register", "auth.RegisterInput", "auth.AuthOutput", 201),
    ("/auth/login", "auth.LoginInput", "auth.AuthOutput", 200),
    ("/auth/verify-email", "auth.VerifyEmailInput", None, 200),
    ("/auth/forgot-password", "auth.ForgotPasswordInput", "auth.ForgotPasswordOutput", 200),
    ("/auth/reset-password", "auth.ResetPasswordInput", None, 200),
    ("/auth/bootstrap-admin", "auth.BootstrapAdminInput", "auth.AuthOutput", 201),
]:
    resp = {str(code): {"description": "OK" if code == 200 else "Created"}}
    if ref_out:
        resp[str(code)]["schema"] = {"$ref": f"#/definitions/{ref_out}"}
    P[path] = {"post": op(["auth"], path.split("/")[-1], {**resp, **std_errors(400, 401, 403)},
               consumes=["application/json"], produces=["application/json"],
               parameters=[body_param(ref_in)])}

P["/me"] = {"get": op(["auth"], "Current user", {**ok("http.MeResponse"), **std_errors(401)}, security=SEC, produces=["application/json"])}
P["/admin/ping"] = {"get": op(["admin"], "Admin ping", {**ok("http.AdminPingResponse"), **std_errors(401, 403)}, security=SEC, produces=["application/json"])}

# identity
P["/identity/submit"] = {"post": op(["identity"], "Submit identity", {**ok("identity.IdentityVerification"), **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"], parameters=[body_param("identity.SubmitInput")])}
P["/identity/status"] = {"get": op(["identity"], "My identity status", {**ok("identity.IdentityVerification"), **std_errors(401, 403)},
    security=SEC, produces=["application/json"])}

for action in ["approve", "reject"]:
    P[f"/admin/identity/{{id}}/{action}"] = {"put": op(["admin", "identity"], f"Identity {action}", {**ok("identity.IdentityVerification"), **std_errors(400, 401, 403)},
        security=SEC, consumes=["application/json"], produces=["application/json"],
        parameters=[path_param("id"), body_param("identity.DecisionInput", False)])}

# licenses
P["/licenses"] = {"post": op(["licenses"], "Apply for license", {**created("licenses.License"), **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"], parameters=[body_param("licenses.CreateInput")])}
P["/licenses/me"] = {"get": op(["licenses"], "My licenses", {"200": {"description": "OK", "schema": arr_lic}, **std_errors(401, 403)},
    security=SEC, produces=["application/json"])}
P["/licenses/{id}/renew"] = {"put": op(["licenses"], "Renew license", {**ok("licenses.License"), **std_errors(400, 401, 403)},
    security=SEC, produces=["application/json"], parameters=[path_param("id")])}
P["/admin/licenses/{id}/approve"] = {"put": op(["admin", "licenses"], "Approve license", {**ok("licenses.License"), **std_errors(400, 401, 403)},
    security=SEC, produces=["application/json"], parameters=[path_param("id")])}

# exams
P["/exam/questions"] = {"get": op(["exams", "admin"], "List questions (admin)", {"200": {"description": "OK", "schema": {"type": "array", "items": {"type": "object"}}}, **std_errors(401, 403)},
    security=SEC, produces=["application/json"])}
P["/exam/start"] = {"post": op(["exams"], "Start theory exam", {**ok("exams.StartOutput"), **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"], parameters=[body_param("exams.StartInput")])}
P["/exam/{attemptId}/submit"] = {"post": op(["exams"], "Submit exam", {**ok("exams.SubmitOutput"), **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"],
    parameters=[path_param("attemptId", "Exam attempt ID"), body_param("exams.SubmitInput")])}
P["/exam/history"] = {"get": op(["exams"], "Exam history", {"200": {"description": "OK", "schema": {"type": "array", "items": {"type": "object"}}}, **std_errors(401, 403)},
    security=SEC, produces=["application/json"])}

# practical
P["/centers"] = {"get": op(["practical"], "List test centers", {"200": {"description": "OK", "schema": {"type": "array", "items": {"type": "object"}}}, **std_errors(401)},
    security=SEC, produces=["application/json"])}
P["/centers/{id}/slots"] = {"get": op(["practical"], "List center slots", {"200": {"description": "OK", "schema": {"type": "array", "items": {"type": "object"}}}, **std_errors(401, 404)},
    security=SEC, produces=["application/json"], parameters=[path_param("id", "Center ID")])}
P["/practical/book"] = {"post": op(["practical"], "Book practical test", {"200": {"description": "OK", "schema": {"type": "object"}}, **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"], parameters=[body_param("practical.BookInput")])}
P["/examiner/practical/{id}/result"] = {"put": op(["practical"], "Record practical result (examiner)", {"200": {"description": "OK", "schema": {"type": "object"}}, **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"],
    parameters=[path_param("id", "Booking ID"), body_param("practical.ResultInput")])}

# vehicles
P["/vehicles"] = {"post": op(["vehicles"], "Register vehicle", {"200": {"description": "OK", "schema": {"type": "object"}}, **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"], parameters=[body_param("vehicles.CreateInput")])}
P["/vehicles/me"] = {"get": op(["vehicles"], "My vehicles", {"200": {"description": "OK", "schema": {"type": "array", "items": {"type": "object"}}}, **std_errors(401, 403)},
    security=SEC, produces=["application/json"])}
P["/vehicles/{id}/transfer"] = {"post": op(["vehicles"], "Request transfer", {"200": {"description": "OK", "schema": {"type": "object"}}, **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"],
    parameters=[path_param("id", "Vehicle ID"), body_param("vehicles.TransferInput")])}
P["/admin/transfer/{id}/approve"] = {"put": op(["admin", "vehicles"], "Approve transfer", {"200": {"description": "OK", "schema": {"type": "object"}}, **std_errors(400, 401, 403)},
    security=SEC, produces=["application/json"], parameters=[path_param("id", "Transfer request ID")])}

# inspections
P["/inspection/schedule"] = {"post": op(["inspections"], "Schedule inspection", {"200": {"description": "OK", "schema": {"type": "object"}}, **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"], parameters=[body_param("inspections.ScheduleInput")])}
P["/inspection/{id}/upload-report"] = {"post": op(["inspections"], "Upload inspection report path", {"200": {"description": "OK", "schema": {"type": "object"}}, **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"],
    parameters=[path_param("id", "Inspection ID"), body_param("inspections.UploadReportInput")])}

# monitoring
P["/devices/data"] = {"post": op(["monitoring"], "Ingest device telemetry", {"200": {"description": "OK", "schema": {"type": "object"}}, **std_errors(400, 401)},
    security=SEC, consumes=["application/json"], produces=["application/json"], parameters=[body_param("monitoring.DeviceDataInput")])}
P["/monitoring/trips/{vehicleId}"] = {"get": op(["monitoring"], "Trips by vehicle", {"200": {"description": "OK", "schema": {"type": "array", "items": {"type": "object"}}}, **std_errors(401)},
    security=SEC, produces=["application/json"], parameters=[path_param("vehicleId", "Vehicle ID")])}
P["/monitoring/score/{userId}"] = {"get": op(["monitoring"], "Safety score by user", {"200": {"description": "OK", "schema": {"type": "object"}}, **std_errors(401)},
    security=SEC, produces=["application/json"], parameters=[path_param("userId", "User ID")])}

# payments
P["/payments/initiate"] = {"post": op(["payments"], "Initiate payment", {**ok("payments.Payment"), **std_errors(400, 401, 403)},
    security=SEC, consumes=["application/json"], produces=["application/json"], parameters=[body_param("payments.InitiateInput")])}
P["/payments/history"] = {"get": op(["payments"], "Payment history", {"200": {"description": "OK", "schema": arr_pay}, **std_errors(401, 403)},
    security=SEC, produces=["application/json"])}
P["/admin/payments/{id}/mark-paid"] = {"put": op(["admin", "payments"], "Mark payment paid", {**ok("payments.Payment"), **std_errors(400, 401, 403)},
    security=SEC, produces=["application/json"], parameters=[path_param("id", "Payment ID")])}

# analytics
for ep, ref in [("/admin/analytics/overview", "analytics.Overview"), ("/admin/analytics/revenue", "analytics.Revenue"), ("/admin/analytics/exams", "analytics.Exams")]:
    P[ep] = {"get": op(["analytics", "admin"], ep.split("/")[-1], {**ok(ref), **std_errors(401, 403)}, security=SEC, produces=["application/json"])}

out = Path(__file__).resolve().parent.parent / "docs" / "swagger.json"
out.write_text(json.dumps(spec, indent=2), encoding="utf-8")
print(f"Wrote {len(P)} paths to {out}")
