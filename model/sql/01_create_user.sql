-- psql -U mylladmin -d myll -h 127.0.0.1 -p 5432
SET ON_ERROR STOP;
CREATE user mylluser password 'example';
CREATE schema myll;
GRANT ALL ON schema myll to mylluser;
CREATE schema extensions;
CREATE extension IF NOT EXISTS pgcrypto schema extensions;
GRANT usage ON schema extensions to mylluser;

