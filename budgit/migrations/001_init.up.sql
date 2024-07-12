CREATE TABLE
  accounts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cleared_balance BIGINT NOT NULL,
    effective_balance BIGINT NOT NULL,
    -- Fields concerning the linked external account. Optional.
    external_id TEXT,
    external_integration_id TEXT,
    external_last_sync_timestamp TIMESTAMPTZ,
    external_cleared_balance BIGINT,
    external_effective_balance BIGINT
  );

CREATE TABLE
  payees (id TEXT PRIMARY KEY, name TEXT NOT NULL);

CREATE INDEX payees_name_idx ON payees (name) INCLUDE (id);

CREATE TABLE
  transactions (
    id TEXT PRIMARY KEY,
    effective_date DATE NOT NULL,
    account_id TEXT NOT NULL,
    payee_id TEXT NOT NULL,
    is_payee_internal BOOLEAN NOT NULL,
    amount BIGINT NOT NULL,
    cleared BOOLEAN NOT NULL
  );

CREATE INDEX transactions_account_effective_date_amount_idx ON transactions (account_id, effective_date, amount);
CREATE INDEX transactions_effective_date_amount_idx ON transactions (effective_date, amount);