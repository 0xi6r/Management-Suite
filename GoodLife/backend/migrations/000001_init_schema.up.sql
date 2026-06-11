-- GoodLife initial schema setup
-- This migration creates the domain-specific schemas for future modules.

CREATE SCHEMA IF NOT EXISTS auth;        -- users, roles, sessions
CREATE SCHEMA IF NOT EXISTS ehr;         -- patient records, diagnoses, prescriptions
CREATE SCHEMA IF NOT EXISTS scheduling;  -- appointments, schedules, waitlists
CREATE SCHEMA IF NOT EXISTS messaging;   -- internal messages
CREATE SCHEMA IF NOT EXISTS billing;     -- invoices, payments, insurance reconciliation
CREATE SCHEMA IF NOT EXISTS audit;       -- access and activity logs
