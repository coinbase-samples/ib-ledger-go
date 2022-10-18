/**
 * Copyright 2022 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

CREATE OR REPLACE FUNCTION create_transaction_and_place_hold(
    arg_transaction_id UUID,
    arg_sender_currency TEXT,
    arg_sender_user_id UUID,
    arg_receiver_currency TEXT,
    arg_receiver_user_id UUID,
    arg_request_id UUID,
    arg_amount NUMERIC,
    arg_type TTYPE
) RETURNS transaction
    LANGUAGE plpgsql
AS
$$
DECLARE
    result_transaction    transaction;
    temp_sender_account   account;
    temp_receiver_account account;
    most_recent_balance   account_balance;
BEGIN
    --Idempotency
    SELECT id, sender_id, receiver_id, request_id, transaction_type, created_at
    FROM transaction
    WHERE id = arg_transaction_id
    INTO result_transaction;
    IF FOUND THEN
        RETURN result_transaction;
    END IF;

    LOCK TABLE account IN ROW EXCLUSIVE MODE;
    --validate accounts exist and lock account rows if they do
    SELECT *
    FROM account
    WHERE currency = arg_sender_currency
      AND user_id = arg_sender_user_id
    INTO temp_sender_account FOR
    UPDATE;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'sender account missing';
    END IF;
    SELECT *
    FROM account
    WHERE currency = arg_receiver_currency
      AND user_id = arg_receiver_user_id
    INTO temp_receiver_account
    FOR
        UPDATE;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'receiver account missing';
    END IF;

    -- validate that we have sufficient balance to complete transaction
    SELECT * FROM get_latest_balance(temp_sender_account.id) INTO most_recent_balance;
    IF most_recent_balance.available < arg_amount THEN
        RAISE EXCEPTION 'insufficient available balance in sender account to process transaction';
    end if;

    -- initialize the transaction
    INSERT INTO transaction (id, sender_id, receiver_id, transaction_type, request_id)
    VALUES (arg_transaction_id, temp_sender_account.id, temp_receiver_account.id, arg_type, arg_request_id)
    RETURNING id, sender_id, receiver_id, request_id, transaction_type, created_at INTO result_transaction;

    -- initialize the hold
    INSERT INTO hold (account_id, transaction_id, amount, request_id)
    VALUES (temp_sender_account.id, result_transaction.id, arg_amount, arg_request_id);

    -- get the most recent balance
    SELECT * FROM get_latest_balance(temp_sender_account.id) INTO most_recent_balance;

    -- create a new balance entry
    INSERT INTO account_balance(account_id, balance, hold, available, request_id, count)
    VALUES (temp_sender_account.id, most_recent_balance.balance, most_recent_balance.hold + arg_amount,
            most_recent_balance.balance - most_recent_balance.hold - arg_amount, arg_request_id,
            most_recent_balance.count + 1);

    UPDATE account SET user_id = temp_sender_account.user_id WHERE id = temp_sender_account.id;
    UPDATE account SET user_id = temp_receiver_account.user_id WHERE id = temp_receiver_account.id;

    return result_transaction;

END
$$;