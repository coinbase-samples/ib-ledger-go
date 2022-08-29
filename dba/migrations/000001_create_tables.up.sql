CREATE EXTENSION "uuid-ossp";

CREATE TYPE direction as ENUM ('DEBIT', 'CREDIT');
CREATE TYPE ttype as ENUM ('ORDER', 'TRANSFER', 'YIELD');

CREATE TABLE IF NOT EXISTS account
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    portfolio_id UUID                           NOT NULL,
    user_id      UUID                           NOT NULL,
    currency     TEXT                           NOT NULL,
    created_at   TIMESTAMPTZ(3)   DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS account_balance
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID                           NOT NULL,
    balance    NUMERIC          DEFAULT 0,
    hold       NUMERIC          DEFAULT 0,
    available  NUMERIC          DEFAULT 0,
    created_at TIMESTAMPTZ(3)   DEFAULT NOW() NOT NULL,
    request_id   UUID,
    count      NUMERIC          DEFAULT 0     NOT NULL
);

CREATE TABLE IF NOT EXISTS transaction
(
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id        UUID  NOT NULL,
    receiver_id      UUID  NOT NULL,
    request_id         UUID  NOT NULL,
    transaction_type TTYPE NOT NULL,
    created_at       TIMESTAMPTZ(3)   DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS finalized_transaction
(
    transaction_id UUID PRIMARY KEY,
    completed_at   TIMESTAMPTZ(3),
    canceled_at    TIMESTAMPTZ(3),
    failed_at      TIMESTAMPTZ(3),
    request_id       UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS entry
(
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id     UUID,
    transaction_id UUID                           NOT NULL,
    request_id       UUID                           NOT NULL,
    amount         NUMERIC,
    direction      DIRECTION                      NOT NULL,
    created_at     TIMESTAMPTZ(3)   DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS hold
(
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id     UUID                           NOT NULL,
    transaction_id UUID                           NOT NULL,
    request_id       UUID                           NOT NULL,
    amount         NUMERIC,
    created_at     TIMESTAMPTZ(3)   DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS released_hold
(
    hold_id     UUID PRIMARY KEY,
    released_at TIMESTAMPTZ(3) DEFAULT NOW() NOT NULL,
    request_id    UUID                         NOT NULL
);
