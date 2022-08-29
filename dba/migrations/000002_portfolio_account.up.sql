CREATE OR REPLACE FUNCTION get_latest_balance(
    arg_account_id UUID
) RETURNS account_balance
    LANGUAGE plpgsql
AS
$$
DECLARE
    result_balance account_balance;
BEGIN
    SELECT *
    FROM account_balance
    WHERE account_id = arg_account_id
    ORDER BY count DESC
    LIMIT 1
    INTO result_balance;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'no balance found';
    END IF;
    return result_balance;
END
$$;

CREATE TYPE account_result AS
(
    id           UUID,
    portfolio_id UUID,
    user_id      UUID,
    currency     TEXT,
    created_at   TIMESTAMPTZ,
    balance      NUMERIC,
    hold         NUMERIC,
    available    NUMERIC
);

CREATE OR REPLACE FUNCTION initialize_account(
    arg_portfolio_id UUID,
    arg_user_id UUID,
    arg_currency TEXT
) RETURNS account_result
    LANGUAGE plpgsql
AS
$$
DECLARE
    temp_balance   account_balance;
    result_account account_result;
BEGIN
    --idempotency
    SELECT *
    FROM account
    WHERE account.portfolio_id = arg_portfolio_id
      AND account.currency = arg_currency
      AND account.user_id = arg_user_id
    INTO result_account;
    IF FOUND THEN
        SELECT * FROM get_latest_balance(result_account.id) INTO temp_balance;
        result_account.balance = temp_balance.balance;
        result_account.hold = temp_balance.hold;
        result_account.available = temp_balance.available;
        return result_account;
    end if;
    INSERT INTO account (portfolio_id, user_id, currency)
    VALUES (arg_portfolio_id, arg_user_id, arg_currency)
    RETURNING id, portfolio_id, user_id, currency, created_at INTO result_account;

    INSERT INTO account_balance (account_id)
    VALUES (result_account.id)
    RETURNING id, account_id, balance, hold, available INTO temp_balance;

    result_account.balance = temp_balance.balance;
    result_account.hold = temp_balance.hold;
    result_account.available = temp_balance.available;

    return result_account;
END
$$;
