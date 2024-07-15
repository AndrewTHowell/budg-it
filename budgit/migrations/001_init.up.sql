CREATE TABLE
  accounts (
    request_id TEXT PRIMARY KEY,
    valid_from_timestamp TIMESTAMPTZ,
    valid_to_timestamp TIMESTAMPTZ,

    id TEXT NOT NULL,
    name TEXT NOT NULL,
    cleared_balance BIGINT NOT NULL,
    effective_balance BIGINT NOT NULL,
    -- Fields concerning the linked external account. Optional.
    external_id TEXT,
    external_name TEXT,
    external_integration_id TEXT,
    external_last_sync_timestamp TIMESTAMPTZ,
    external_cleared_balance BIGINT,
    external_effective_balance BIGINT
  );

CREATE INDEX accounts_request_id_idx ON accounts (request_id);
CREATE INDEX accounts_id_idx ON accounts (id) WHERE valid_to_timestamp = 'infinity';

CREATE TABLE
  payees (
    request_id TEXT PRIMARY KEY,
    valid_from_timestamp TIMESTAMPTZ,
    valid_to_timestamp TIMESTAMPTZ,

    id TEXT NOT NULL,
    name TEXT NOT NULL
  );

CREATE INDEX payees_request_id_idx ON payees (request_id);
CREATE INDEX payees_id_idx ON payees (id) WHERE valid_to_timestamp = 'infinity';
CREATE INDEX payees_name_idx ON payees (name) WHERE valid_to_timestamp = 'infinity';

CREATE TABLE
  transactions (
    request_id TEXT PRIMARY KEY,
    valid_from_timestamp TIMESTAMPTZ,
    valid_to_timestamp TIMESTAMPTZ,

    id TEXT NOT NULL,
    effective_date DATE NOT NULL,
    account_id TEXT NOT NULL,
    payee_id TEXT NOT NULL,
    is_payee_internal BOOLEAN NOT NULL,
    amount BIGINT NOT NULL,
    cleared BOOLEAN NOT NULL
  );

CREATE INDEX transactions_request_id_idx ON transactions (request_id);
CREATE INDEX transactions_id_idx ON transactions (id) WHERE valid_to_timestamp = 'infinity';
CREATE INDEX transactions_account_effective_date_amount_idx ON transactions (account_id, effective_date, amount) WHERE valid_to_timestamp = 'infinity';
CREATE INDEX transactions_effective_date_amount_idx ON transactions (effective_date, amount) WHERE valid_to_timestamp = 'infinity';