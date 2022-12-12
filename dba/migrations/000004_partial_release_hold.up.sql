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

CREATE TYPE transaction_result AS
(
    hold_id             UUID,
    sender_entry_id     UUID,
    receiver_entry_id   UUID,
    sender_balance_id   UUID,
    receiver_balance_id UUID
);

CREATE OR REPLACE FUNCTION partial_release_hold(
    arg_transaction_id UUID,
    arg_request_id UUID,
    arg_sender_amount NUMERIC,
    arg_receiver_amount NUMERIC,
    arg_retail_fee_amount NUMERIC,
    arg_retail_fee_account_id UUID,
    arg_venue_fee_amount NUMERIC,
    arg_venue_fee_account_id UUID
) RETURNS transaction_result
    LANGUAGE plpgsql
AS
$$
DECLARE
    temp_transaction             transaction;
    temp_hold                    hold;
    sender_account               account;
    most_recent_sender_balance   account_balance;
    most_recent_receiver_balance account_balance;
    receiver_account             account;
    sender_entry_id              UUID;
    receiver_entry_id            UUID;
    sender_balance_id            UUID;
    receiver_balance_id          UUID;
    temp_balance_amount          NUMERIC;
    sender_hold_amount           NUMERIC;
    result                       transaction_result;
BEGIN
    SELECT * FROM transaction WHERE id = arg_transaction_id INTO temp_transaction;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'LGR404';
    END IF;

    --Locking
    LOCK TABLE account IN ROW EXCLUSIVE MODE;
    SELECT * FROM account WHERE id = temp_transaction.sender_id into sender_account FOR UPDATE;
    SELECT * FROM account WHERE id = temp_transaction.receiver_id into receiver_account FOR UPDATE;

    --idempotency
    PERFORM *
    FROM entry
    WHERE transaction_id = arg_transaction_id
      AND request_id = arg_request_id
      AND account_id = temp_transaction.sender_id;
    IF FOUND THEN
        SELECT id
        FROM entry
        WHERE transaction_id = arg_transaction_id
          AND request_id = arg_request_id
          AND account_id = temp_transaction.receiver_id
        LIMIT 1
        INTO receiver_entry_id;
        result.receiver_entry_id = receiver_entry_id;
        RETURN result;
    END IF;

    SELECT *
    FROM hold
    WHERE hold.transaction_id = temp_transaction.id
      AND hold.account_id = temp_transaction.sender_id
      AND NOT EXISTS(
            SELECT hold_id
            FROM released_hold
            WHERE released_hold.hold_id = hold.id)
    INTO temp_hold;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'LGR405';
    END IF;

    --Release the hold
    INSERT INTO released_hold (hold_id, request_id) VALUES (temp_hold.id, arg_request_id);

    --Get most recent sender balance
    SELECT * FROM get_latest_balance(temp_transaction.sender_id) INTO most_recent_sender_balance;

    --Insert Sender Entry for Amount Sent to Receiver
    INSERT INTO entry (account_id, transaction_id, amount, direction, request_id)
    VALUES (temp_transaction.sender_id, arg_transaction_id, arg_sender_amount, 'DEBIT', arg_request_id)
    RETURNING id INTO sender_entry_id;
    result.sender_entry_id = sender_entry_id;

    --Insert Sender Entry for Fee if Fee Paid
    IF arg_retail_fee_amount != 0 THEN
        INSERT INTO entry (account_id, transaction_id, amount, direction, request_id)
        VALUES (temp_transaction.sender_id, arg_transaction_id, arg_retail_fee_amount, 'DEBIT', arg_request_id);
        INSERT INTO entry (account_id, transaction_id, amount, direction, request_id)
        VALUES (arg_retail_fee_account_id, arg_transaction_id, arg_retail_fee_amount, 'CREDIT', arg_request_id);
    END IF;

    IF arg_venue_fee_amount != 0 THEN
        INSERT INTO entry (account_id, transaction_id, amount, direction, request_id)
        VALUES (temp_transaction.sender_id, arg_transaction_id, arg_venue_fee_amount, 'DEBIT', arg_request_id);
        INSERT INTO entry (account_id, transaction_id, amount, direction, request_id)
        VALUES (arg_venue_fee_account_id, arg_transaction_id, arg_venue_fee_amount, 'CREDIT', arg_request_id);
    END IF;

    --Insert Sender Balance
    temp_balance_amount =
                most_recent_sender_balance.balance - arg_sender_amount - arg_venue_fee_amount - arg_retail_fee_amount;
    sender_hold_amount = temp_hold.amount - arg_sender_amount - arg_retail_fee_amount - arg_venue_fee_amount;
    INSERT INTO account_balance(account_id, request_id, balance, hold, available, count)
    VALUES (temp_transaction.sender_id, arg_request_id, temp_balance_amount, sender_hold_amount,
            temp_balance_amount - sender_hold_amount, most_recent_sender_balance.count + 1)
    RETURNING id INTO sender_balance_id;
    result.sender_balance_id = sender_balance_id;

    --Insert Receiver Entry amd Update Account Balance
    INSERT INTO entry (account_id, transaction_id, amount, direction, request_id)
    VALUES (temp_transaction.receiver_id, arg_transaction_id, arg_receiver_amount, 'CREDIT', arg_request_id)
    RETURNING id INTO receiver_entry_id;
    result.receiver_entry_id = receiver_entry_id;

    SELECT * FROM get_latest_balance(temp_transaction.receiver_id) INTO most_recent_receiver_balance;

    temp_balance_amount = most_recent_receiver_balance.balance + arg_receiver_amount;
    INSERT INTO account_balance(account_id, request_id, balance, hold, available, count)
    VALUES (temp_transaction.receiver_id, arg_request_id, temp_balance_amount, most_recent_receiver_balance.hold,
            temp_balance_amount - most_recent_receiver_balance.hold, most_recent_receiver_balance.count + 1)
    RETURNING id INTO receiver_balance_id;
    result.receiver_balance_id = receiver_balance_id;

    --Insert New Hold if needed
    IF sender_hold_amount > 0 THEN
        INSERT INTO hold (account_id, transaction_id, amount, request_id)
        VALUES (temp_transaction.sender_id, temp_transaction.id, sender_hold_amount,
                arg_request_id)
        RETURNING * INTO temp_hold;
        result.hold_id = temp_hold.id;
    END IF;

    UPDATE account SET user_id = sender_account.user_id WHERE id = sender_account.id;
    UPDATE account SET user_id = receiver_account.user_id WHERE id = receiver_account.id;

    RETURN result;
END
$$;
