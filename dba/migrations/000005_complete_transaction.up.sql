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

CREATE OR REPLACE FUNCTION complete_transaction(
    arg_transaction_id UUID,
    arg_request_id UUID,
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

    --Release the hold
    INSERT INTO released_hold (hold_id, request_id) VALUES (temp_hold.id, arg_request_id);

    --Insert Sender Entry and Update Account Balance
    SELECT *
    FROM insert_entry_and_update_balance(
            temp_transaction.sender_id,
            arg_transaction_id,
            arg_request_id,
            temp_hold.amount,
            'DEBIT')
    INTO sender_entry_result;
    result.sender_entry_id = sender_entry_result.entry_id;
    result.sender_balance_id = sender_entry_result.balance_id;

    --Insert Receiver Entry amd Update Account Balance
    SELECT *
    FROM insert_entry_and_update_balance(temp_transaction.receiver_id, arg_transaction_id, arg_request_id, arg_receiver_amount, 'CREDIT')
    INTO receiver_entry_result;
    result.receiver_entry_id = receiver_entry_result.entry_id;
    result.receiver_balance_id = receiver_entry_result.balance_id;

    --Finalize Transaction
    INSERT INTO finalized_transaction (transaction_id, completed_at, request_id)
    VALUES (arg_transaction_id, now(), arg_request_id);

    UPDATE account SET user_id = sender_account.user_id WHERE id = sender_account.id;
    UPDATE account SET user_id = receiver_account.user_id WHERE id = receiver_account.id;

    return result;
END
$$;