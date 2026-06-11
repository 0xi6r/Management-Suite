-- Seed the default roles used by the application
INSERT INTO auth.roles (name, description) VALUES
    ('patient', 'Patient role – can manage own profile, appointments, and health data'),
    ('doctor', 'Doctor / health practitioner – can access and manage patient records'),
    ('admin', 'Institutional administrator – manages staff, branches, and system config'),
    ('insurance', 'Third‑party insurance provider – limited access to claims and documents')
ON CONFLICT (name) DO NOTHING;
