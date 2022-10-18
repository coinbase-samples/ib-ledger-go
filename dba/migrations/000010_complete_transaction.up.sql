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
DROP FUNCTION IF EXISTS complete_transaction cascade;
CREATE OR REPLACE FUNCTION complete_transaction(
    arg_transaction_id UUID,
    arg_request_id UUID
) RETURNS transaction_result
    LANGUAGE plpgsql
AS
$$
DECLARE
    temp_transaction             transaction;
    temp_hold                    hold;
    sender_account               account;
    receiver_account             account;
    most_recent_sender_balance   account_balance;
    most_recent_receiver_balance account_balance;
    temp_balance_amount          NUMERIC;
    sender_hold_amount           NUMERIC;
    sender_entry_id              UUID;
    receiver_entry_id            UUID;
    sender_balance_id            UUID;
    receiver_balance_id          UUID;
    result                       transaction_result;
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
    SELECT * FROM account WHERE id = temp_transaction.receiver_id into receiver_account FOR UPDATE;

    SELECT * FROM transaction WHERE id = arg_transaction_id INTO temp_transaction;

    SELECT *
    FROM get_unreleased_hold(temp_transaction.id, temp_transaction.sender_id)
    INTO temp_hold;

    IF FOUND THEN
        --Release the hold
        INSERT INTO released_hold (hold_id, request_id) VALUES (temp_hold.id, arg_request_id);
        --Insert Sender Balance
        temp_balance_amount = most_recent_sender_balance.balance - temp_hold.amount;
        sender_hold_amount = most_recent_sender_balance.hold - temp_hold.amount;
        INSERT INTO account_balance(account_id, request_id, balance, hold, available, count)
        VALUES (temp_transaction.sender_id, arg_request_id, temp_balance_amount, sender_hold_amount,
                temp_balance_amount - sender_hold_amount, most_recent_sender_balance.count + 1)
        RETURNING id INTO sender_balance_id;
        result.sender_balance_id = sender_balance_id;
    end if;

    --Finalize Transaction
    INSERT INTO finalized_transaction (transaction_id, completed_at, request_id)
    VALUES (arg_transaction_id, now(), arg_request_id);

    UPDATE account SET user_id = sender_account.user_id WHERE id = sender_account.id;
    UPDATE account SET user_id = receiver_account.user_id WHERE id = receiver_account.id;

    return result;
END
$$;