CREATE OR REPLACE FUNCTION fail_transaction(
    arg_transaction_id UUID,
    arg_request_id UUID
) RETURNS transaction_result
    LANGUAGE plpgsql
AS
$$
DECLARE
    temp_transaction           transaction;
    temp_hold                  hold;
    sender_account             account;
    sender_most_recent_balance account_balance;
    temp_hold_amount           NUMERIC;
    temp_sender_new_balance_id uuid;
    result                     transaction_result;
BEGIN
    SELECT *
    FROM transaction
    WHERE id = arg_transaction_id
      AND EXISTS(
            SELECT * FROM finalized_transaction WHERE transaction_id = arg_transaction_id)
    INTO temp_transaction;
    IF FOUND THEN
        result.sender_balance_id = (SELECT id
                                    FROM account_balance
                                    WHERE request_id = arg_request_id
                                      AND account_id = temp_transaction.sender_id);
        result.receiver_balance_id = (SELECT id
                                      FROM account_balance
                                      WHERE request_id = arg_request_id
                                        AND account_id = temp_transaction.receiver_id);
        return result;
    end if;

    --Locking
    LOCK TABLE account IN ROW EXCLUSIVE MODE;
    SELECT * FROM account WHERE id = temp_transaction.sender_id into sender_account FOR UPDATE;

    SELECT * FROM transaction WHERE id = arg_transaction_id INTO temp_transaction;

    SELECT *
    FROM get_unreleased_hold(temp_transaction.id, temp_transaction.sender_id)
    INTO temp_hold;

    --Release the hold
    INSERT INTO released_hold (hold_id, request_id) VALUES (temp_hold.id, arg_request_id);

    result.hold_id = temp_hold.id;

    --Update Sender Balance
    SELECT *
    FROM get_latest_balance(
            temp_transaction.sender_id
        )
    INTO sender_most_recent_balance;

    temp_hold_amount = sender_most_recent_balance - temp_hold.amount;

    INSERT INTO account_balance(account_id, request_id, balance, hold, available)
    VALUES (temp_transaction.sender_id, arg_request_id, sender_most_recent_balance.balance, temp_hold_amount,
            sender_most_recent_balance - temp_hold_amount)
    RETURNING id INTO temp_sender_new_balance_id;

    result.sender_balance_id = temp_sender_new_balance_id;

    --Finalize Transaction
    INSERT INTO finalized_transaction (transaction_id, failed_at, request_id)
    VALUES (arg_transaction_id, now(), arg_request_id);

    UPDATE account SET user_id = sender_account.user_id WHERE id = sender_account.id;

    return result;
END
$$;

CREATE OR REPLACE FUNCTION cancel_transaction(
    arg_transaction_id UUID,
    arg_request_id UUID
) RETURNS transaction_result
    LANGUAGE plpgsql
AS
$$
DECLARE
    temp_transaction           transaction;
    temp_hold                  hold;
    sender_account             account;
    sender_most_recent_balance account_balance;
    temp_hold_amount           NUMERIC;
    temp_sender_new_balance_id uuid;
    result                     transaction_result;
BEGIN
    SELECT *
    FROM transaction
    WHERE id = arg_transaction_id
      AND EXISTS(
            SELECT * FROM finalized_transaction WHERE transaction_id = arg_transaction_id)
    INTO temp_transaction;
    IF FOUND THEN
        result.sender_balance_id = (SELECT id
                                    FROM account_balance
                                    WHERE request_id = arg_request_id
                                      AND account_id = temp_transaction.sender_id);
        result.receiver_balance_id = (SELECT id
                                      FROM account_balance
                                      WHERE request_id = arg_request_id
                                        AND account_id = temp_transaction.receiver_id);
        return result;
    end if;

    --Locking
    LOCK TABLE account IN ROW EXCLUSIVE MODE;
    SELECT * FROM account WHERE id = temp_transaction.sender_id into sender_account FOR UPDATE;

    SELECT * FROM transaction WHERE id = arg_transaction_id INTO temp_transaction;

    SELECT *
    FROM get_unreleased_hold(temp_transaction.id, temp_transaction.sender_id)
    INTO temp_hold;

    --Release the hold
    INSERT INTO released_hold (hold_id, request_id) VALUES (temp_hold.id, arg_request_id);

    result.hold_id = temp_hold.id;

    --Update Sender Balance
    SELECT *
    FROM get_latest_balance(
            temp_transaction.sender_id
        )
    INTO sender_most_recent_balance;

    temp_hold_amount = sender_most_recent_balance - temp_hold.amount;

    INSERT INTO account_balance(account_id, request_id, balance, hold, available)
    VALUES (temp_transaction.sender_id, arg_request_id, sender_most_recent_balance.balance, temp_hold_amount,
            sender_most_recent_balance - temp_hold_amount)
    RETURNING id INTO temp_sender_new_balance_id;

    result.sender_balance_id = temp_sender_new_balance_id;

    --Finalize Transaction
    INSERT INTO finalized_transaction (transaction_id, canceled_at, request_id)
    VALUES (arg_transaction_id, now(), arg_request_id);

    UPDATE account SET user_id = sender_account.user_id WHERE id = sender_account.id;

    return result;
END
$$;