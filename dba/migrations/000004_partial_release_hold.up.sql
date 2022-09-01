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

CREATE OR REPLACE FUNCTION get_unreleased_hold(
    arg_transaction_id UUID,
    arg_account_id UUID
) RETURNS hold
    LANGUAGE plpgsql
AS
$$
DECLARE
    result_hold hold;
BEGIN
    SELECT *
    FROM hold
    WHERE hold.transaction_id = arg_transaction_id
      AND hold.account_id = arg_account_id
      AND NOT EXISTS(
            SELECT hold_id
            FROM released_hold
            WHERE released_hold.hold_id = hold.id)
    INTO result_hold;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'no unreleased hold found for transaction';
    end if;
    return result_hold;
END
$$;

CREATE TYPE transaction_result AS
(
    hold_id             UUID,
    sender_entry_id     UUID,
    receiver_entry_id   UUID,
    sender_balance_id   UUID,
    receiver_balance_id UUID
);

CREATE TYPE insert_entry_balance_result AS
(
    entry_id   UUID,
    balance_id UUID
);
CREATE OR REPLACE FUNCTION insert_entry_and_update_balance(
    arg_account_id UUID,
    arg_transaction_id UUID,
    arg_request_id UUID,
    arg_amount NUMERIC,
    arg_direction direction
) RETURNS insert_entry_balance_result
    LANGUAGE plpgsql
AS
$$
DECLARE
    temp_entry          entry;
    temp_balance        account_balance;
    temp_balance_amount NUMERIC;
    temp_hold_amount    NUMERIC;
    result              insert_entry_balance_result;
BEGIN
    INSERT INTO entry (account_id, transaction_id, amount, direction, request_id)
    VALUES (arg_account_id, arg_transaction_id, arg_amount, arg_direction, arg_request_id)
    RETURNING * INTO temp_entry;

    result.entry_id = temp_entry.id;

    SELECT * FROM get_latest_balance(arg_account_id) INTO temp_balance;

    temp_balance_amount = temp_balance.balance;
    temp_hold_amount = temp_balance.hold;
    if arg_direction = 'DEBIT' THEN
        temp_balance_amount = temp_balance_amount - arg_amount;
        temp_hold_amount = temp_hold_amount - arg_amount;
    end if;
    if arg_direction = 'CREDIT' THEN
        temp_balance_amount = temp_balance_amount + arg_amount;
    end if;
    INSERT INTO account_balance (account_id, request_id, balance, hold, available, count)
    VALUES (arg_account_id, arg_request_id, temp_balance_amount, temp_hold_amount, temp_balance_amount - temp_hold_amount,
            temp_balance.count + 1)
    RETURNING * INTO temp_balance;
    result.balance_id = temp_balance.id;
    return result;
END
$$;

CREATE OR REPLACE FUNCTION partial_release_hold(
    arg_transaction_id UUID,
    arg_request_id UUID,
    arg_sender_amount NUMERIC,
    arg_receiver_amount NUMERIC
) RETURNS transaction_result
    LANGUAGE plpgsql
AS
$$
DECLARE
    temp_transaction      transaction;
    temp_hold             hold;
    sender_account        account;
    receiver_account      account;
    sender_entry_result   insert_entry_balance_result;
    receiver_entry_result insert_entry_balance_result;
    result                transaction_result;
BEGIN
    SELECT * FROM transaction WHERE id = arg_transaction_id INTO temp_transaction;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'transaction not found';
    end if;
    --idempotency
    SELECT * FROM hold WHERE request_id = arg_request_id into temp_hold;
    if FOUND THEN
        result.hold_id = temp_hold.id;
        result.sender_entry_id = (SELECT id
                                  FROM entry
                                  where transaction_id = arg_transaction_id
                                    AND request_id = arg_request_id
                                    AND account_id = temp_transaction.sender_id);
        result.receiver_entry_id = (SELECT *
                                    FROM entry
                                    where transaction_id = arg_transaction_id
                                      AND request_id = arg_request_id
                                      AND account_id = temp_transaction.receiver_id);
        result.sender_balance_id = (SELECT id
                                    FROM account_balance
                                    WHERE request_id = arg_request_id
                                      AND account_id = temp_transaction.sender_id);
        result.receiver_balance_id = (SELECT id
                                      FROM account_balance
                                      WHERE request_id = arg_request_id
                                        AND account_id = temp_transaction.receiver_id);
        return result;
    END IF;

    --Locking
    LOCK TABLE account IN ROW EXCLUSIVE MODE;
    SELECT * FROM account WHERE id = temp_transaction.sender_id into sender_account FOR UPDATE;
    SELECT * FROM account WHERE id = temp_transaction.receiver_id into receiver_account FOR UPDATE;

    SELECT *
    FROM get_unreleased_hold(temp_transaction.id, temp_transaction.sender_id)
    INTO temp_hold;

    --Release the hold
    INSERT INTO released_hold (hold_id, request_id) VALUES (temp_hold.id, arg_request_id);

    --Insert Sender Entry and Update Account Balance
    SELECT *
    FROM insert_entry_and_update_balance(temp_transaction.sender_id, arg_transaction_id, arg_request_id, arg_sender_amount, 'DEBIT')
    INTO sender_entry_result;
    result.sender_entry_id = sender_entry_result.entry_id;
    result.sender_balance_id = sender_entry_result.balance_id;

    --Insert Receiver Entry amd Update Account Balance
    SELECT *
    FROM insert_entry_and_update_balance(temp_transaction.receiver_id, arg_transaction_id, arg_request_id, arg_receiver_amount, 'CREDIT')
    INTO receiver_entry_result;
    result.receiver_entry_id = receiver_entry_result.entry_id;
    result.receiver_balance_id = receiver_entry_result.balance_id;

    --Insert New Hold
    INSERT INTO hold (account_id, transaction_id, amount, request_id)
    VALUES (temp_transaction.sender_id, temp_transaction.id, temp_hold.amount - arg_sender_amount, arg_request_id)
    RETURNING * INTO temp_hold;
    result.hold_id = temp_hold.id;

    UPDATE account SET user_id = sender_account.user_id WHERE id = sender_account.id;
    UPDATE account SET user_id = receiver_account.user_id WHERE id = receiver_account.id;

    return result;
END
$$;