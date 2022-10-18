--
-- PostgreSQL database dump
--

-- Dumped from database version 14.3 (Debian 14.3-1.pgdg110+1)
-- Dumped by pg_dump version 14.5 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- Name: account_result; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.account_result AS (
	id uuid,
	portfolio_id uuid,
	user_id uuid,
	currency text,
	created_at timestamp with time zone,
	balance numeric,
	hold numeric,
	available numeric
);


ALTER TYPE public.account_result OWNER TO postgres;

--
-- Name: direction; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.direction AS ENUM (
    'DEBIT',
    'CREDIT'
);


ALTER TYPE public.direction OWNER TO postgres;

--
-- Name: get_account_result; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.get_account_result AS (
	account_id uuid,
	currency text,
	balance numeric,
	hold numeric,
	available numeric,
	created_at timestamp with time zone
);


ALTER TYPE public.get_account_result OWNER TO postgres;

--
-- Name: transaction_result; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.transaction_result AS (
	hold_id uuid,
	sender_entry_id uuid,
	receiver_entry_id uuid,
	sender_balance_id uuid,
	receiver_balance_id uuid
);


ALTER TYPE public.transaction_result OWNER TO postgres;

--
-- Name: ttype; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.ttype AS ENUM (
    'ORDER',
    'TRANSFER',
    'CONVERT'
);


ALTER TYPE public.ttype OWNER TO postgres;

--
-- Name: cancel_transaction(uuid, uuid); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.cancel_transaction(arg_transaction_id uuid, arg_request_id uuid) RETURNS public.transaction_result
    LANGUAGE plpgsql
    AS $$
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

    temp_hold_amount = sender_most_recent_balance.hold - temp_hold.amount;

    INSERT INTO account_balance(account_id, request_id, balance, hold, available, count)
    VALUES (temp_transaction.sender_id, arg_request_id, sender_most_recent_balance.balance, temp_hold_amount,
            sender_most_recent_balance.balance - temp_hold_amount, sender_most_recent_balance.count + 1)
    RETURNING id INTO temp_sender_new_balance_id;

    result.sender_balance_id = temp_sender_new_balance_id;

    --Finalize Transaction
    INSERT INTO finalized_transaction (transaction_id, canceled_at, request_id)
    VALUES (arg_transaction_id, now(), arg_request_id);

    UPDATE account SET user_id = sender_account.user_id WHERE id = sender_account.id;

    return result;
END
$$;


ALTER FUNCTION public.cancel_transaction(arg_transaction_id uuid, arg_request_id uuid) OWNER TO postgres;

--
-- Name: complete_transaction(uuid, uuid); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.complete_transaction(arg_transaction_id uuid, arg_request_id uuid) RETURNS public.transaction_result
    LANGUAGE plpgsql
    AS $$
DECLARE
    temp_transaction      transaction;
    temp_hold             hold;
    sender_account        account;
    receiver_account      account;
    most_recent_sender_balance   account_balance;
    most_recent_receiver_balance account_balance;
    temp_balance_amount          NUMERIC;
    sender_hold_amount           NUMERIC;
    sender_entry_id              UUID;
    receiver_entry_id            UUID;
    sender_balance_id            UUID;
    receiver_balance_id          UUID;
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

    IF FOUND THEN
      --Release the hold
      INSERT INTO released_hold (hold_id, request_id) VALUES (temp_hold.id, arg_request_id);
      --Get most recent sender balance
      SELECT * FROM get_latest_balance(temp_transaction.sender_id) INTO most_recent_sender_balance;
      --Insert Sender Balance
      sender_hold_amount = most_recent_sender_balance.hold - temp_hold.amount;
      INSERT INTO account_balance(account_id, request_id, balance, hold, available, count)
      VALUES (temp_transaction.sender_id, arg_request_id,most_recent_sender_balance.balance, sender_hold_amount,
              most_recent_sender_balance.balance - sender_hold_amount, most_recent_sender_balance.count + 1)
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


ALTER FUNCTION public.complete_transaction(arg_transaction_id uuid, arg_request_id uuid) OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: transaction; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    sender_id uuid NOT NULL,
    receiver_id uuid NOT NULL,
    request_id uuid NOT NULL,
    transaction_type public.ttype NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now()
);


ALTER TABLE public.transaction OWNER TO postgres;

--
-- Name: create_transaction_and_place_hold(uuid, text, uuid, text, uuid, uuid, numeric, public.ttype); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.create_transaction_and_place_hold(arg_transaction_id uuid, arg_sender_currency text, arg_sender_user_id uuid, arg_receiver_currency text, arg_receiver_user_id uuid, arg_request_id uuid, arg_amount numeric, arg_type public.ttype) RETURNS public.transaction
    LANGUAGE plpgsql
    AS $$
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
    INTO temp_receiver_account FOR
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


ALTER FUNCTION public.create_transaction_and_place_hold(arg_transaction_id uuid, arg_sender_currency text, arg_sender_user_id uuid, arg_receiver_currency text, arg_receiver_user_id uuid, arg_request_id uuid, arg_amount numeric, arg_type public.ttype) OWNER TO postgres;

--
-- Name: fail_transaction(uuid, uuid); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.fail_transaction(arg_transaction_id uuid, arg_request_id uuid) RETURNS public.transaction_result
    LANGUAGE plpgsql
    AS $$
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

    temp_hold_amount = sender_most_recent_balance.hold - temp_hold.amount;

    INSERT INTO account_balance(account_id, request_id, balance, hold, available, count)
    VALUES (temp_transaction.sender_id, arg_request_id, sender_most_recent_balance.balance, temp_hold_amount,
            sender_most_recent_balance.balance - temp_hold_amount, sender_most_recent_balance.count + 1)
    RETURNING id INTO temp_sender_new_balance_id;

    result.sender_balance_id = temp_sender_new_balance_id;

    --Finalize Transaction
    INSERT INTO finalized_transaction (transaction_id, failed_at, request_id)
    VALUES (arg_transaction_id, now(), arg_request_id);

    UPDATE account SET user_id = sender_account.user_id WHERE id = sender_account.id;

    return result;
END
$$;


ALTER FUNCTION public.fail_transaction(arg_transaction_id uuid, arg_request_id uuid) OWNER TO postgres;

--
-- Name: get_account_and_latest_balance(uuid, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.get_account_and_latest_balance(arg_user_id uuid, arg_currency text) RETURNS SETOF public.get_account_result
    LANGUAGE plpgsql
    AS $$
BEGIN
    return QUERY
        SELECT acc.id, acc.currency, ab.balance, ab.hold, ab.available, ab.created_at
        FROM (select id, currency FROM account WHERE user_id = arg_user_id AND currency = arg_currency) acc
                 INNER JOIN
             (SELECT account_id, balance, hold, available, created_at FROM account_balance HAVING count = MAX(count)) ab
        ON acc.id = ab.account_id;
END
$$;


ALTER FUNCTION public.get_account_and_latest_balance(arg_user_id uuid, arg_currency text) OWNER TO postgres;

--
-- Name: get_balances_for_users(uuid); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.get_balances_for_users(arg_user_id uuid) RETURNS SETOF public.get_account_result
    LANGUAGE plpgsql
    AS $$
BEGIN
    return QUERY
        SELECT acc.id, acc.currency, ab.balance, ab.hold, ab.available, ab.created_at
        FROM (select id, currency FROM account WHERE user_id = arg_user_id) acc
                 INNER JOIN
             (SELECT account_id, MAX(count) as max FROM account_balance WHERE account_id IN (
                select id FROM account WHERE user_id = arg_user_id
                ) GROUP BY account_id ) recent_balance
             ON acc.id = recent_balance.account_id
                 INNER JOIN
             account_balance ab
             ON recent_balance.account_id = ab.account_id and recent_balance.max = ab.count
        GROUP BY ab.count, acc.id, acc.currency, ab.balance, ab.hold, ab.available, ab.created_at
        HAVING recent_balance.count = acc.count;
END
$$;


ALTER FUNCTION public.get_balances_for_users(arg_user_id uuid) OWNER TO postgres;

--
-- Name: account_balance; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.account_balance (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    account_id uuid NOT NULL,
    balance numeric DEFAULT 0 NOT NULL,
    hold numeric DEFAULT 0 NOT NULL,
    available numeric DEFAULT 0 NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    request_id uuid,
    count numeric DEFAULT 1 NOT NULL,
    CONSTRAINT account_balance_available_check CHECK ((available >= (0)::numeric)),
    CONSTRAINT account_balance_balance_check CHECK ((balance >= (0)::numeric)),
    CONSTRAINT account_balance_hold_check CHECK ((hold >= (0)::numeric))
);


ALTER TABLE public.account_balance OWNER TO postgres;

--
-- Name: get_latest_balance(uuid); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.get_latest_balance(arg_account_id uuid) RETURNS public.account_balance
    LANGUAGE plpgsql
    AS $$
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


ALTER FUNCTION public.get_latest_balance(arg_account_id uuid) OWNER TO postgres;

--
-- Name: hold; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hold (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    account_id uuid NOT NULL,
    transaction_id uuid NOT NULL,
    request_id uuid NOT NULL,
    amount numeric NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    CONSTRAINT hold_amount_check CHECK ((amount > (0)::numeric))
);


ALTER TABLE public.hold OWNER TO postgres;

--
-- Name: get_unreleased_hold(uuid, uuid); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.get_unreleased_hold(arg_transaction_id uuid, arg_account_id uuid) RETURNS public.hold
    LANGUAGE plpgsql
    AS $$
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


ALTER FUNCTION public.get_unreleased_hold(arg_transaction_id uuid, arg_account_id uuid) OWNER TO postgres;

--
-- Name: initialize_account(uuid, uuid, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.initialize_account(arg_portfolio_id uuid, arg_user_id uuid, arg_currency text) RETURNS public.account_result
    LANGUAGE plpgsql
    AS $$
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


ALTER FUNCTION public.initialize_account(arg_portfolio_id uuid, arg_user_id uuid, arg_currency text) OWNER TO postgres;

--
-- Name: partial_release_hold(uuid, uuid, numeric, numeric, numeric, uuid, numeric, uuid); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.partial_release_hold(arg_transaction_id uuid, arg_request_id uuid, arg_sender_amount numeric, arg_receiver_amount numeric, arg_retail_fee_amount numeric, arg_retail_fee_account_id uuid, arg_venue_fee_amount numeric, arg_venue_fee_account_id uuid) RETURNS public.transaction_result
    LANGUAGE plpgsql
    AS $$
DECLARE
    temp_transaction             transaction;
    temp_sender_entry            entry;
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
        RAISE EXCEPTION 'transaction not found';
    end if;
    --idempotency
    SELECT * FROM entry WHERE transaction_id = arg_transaction_id AND request_id = arg_request_id AND account_id = temp_transaction.sender_id into temp_sender_entry;
    if FOUND THEN
        result.sender_entry_id = temp_sender_entry.id;
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
    temp_balance_amount = most_recent_sender_balance.balance - arg_sender_amount - arg_venue_fee_amount - arg_retail_fee_amount;
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

    return result;
END
$$;


ALTER FUNCTION public.partial_release_hold(arg_transaction_id uuid, arg_request_id uuid, arg_sender_amount numeric, arg_receiver_amount numeric, arg_retail_fee_amount numeric, arg_retail_fee_account_id uuid, arg_venue_fee_amount numeric, arg_venue_fee_account_id uuid) OWNER TO postgres;

--
-- Name: account; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.account (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    portfolio_id uuid NOT NULL,
    user_id uuid NOT NULL,
    currency text NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.account OWNER TO postgres;

--
-- Name: entry; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.entry (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    account_id uuid NOT NULL,
    transaction_id uuid NOT NULL,
    request_id uuid NOT NULL,
    amount numeric NOT NULL,
    direction public.direction NOT NULL,
    created_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    CONSTRAINT entry_amount_check CHECK ((amount > (0)::numeric))
);


ALTER TABLE public.entry OWNER TO postgres;

--
-- Name: finalized_transaction; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.finalized_transaction (
    transaction_id uuid NOT NULL,
    completed_at timestamp(3) with time zone,
    canceled_at timestamp(3) with time zone,
    failed_at timestamp(3) with time zone,
    request_id uuid NOT NULL
);


ALTER TABLE public.finalized_transaction OWNER TO postgres;

--
-- Name: released_hold; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.released_hold (
    hold_id uuid NOT NULL,
    released_at timestamp(3) with time zone DEFAULT now() NOT NULL,
    request_id uuid NOT NULL
);


ALTER TABLE public.released_hold OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Data for Name: account; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.account (id, portfolio_id, user_id, currency, created_at) FROM stdin;
b72d0e55-f53a-4db0-897e-2ce4a73cb94b	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	116bde43-7733-43a1-a85a-fc8627e6da8e	USD	2022-10-18 14:54:43.491+00
c4d0e14e-1b2b-4023-afa6-8891ad1960c9	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	433c0c15-0a44-49c4-a207-4501bb11f48c	USD	2022-10-18 14:54:43.491+00
0adbb104-fc18-46ca-a4eb-beee7775eb69	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	c5af3271-7185-4a52-9d0c-1c4b418317d8	USD	2022-10-18 14:54:43.491+00
af271fce-c63c-4d79-af50-dcded9f806a4	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	c5af3271-7185-4a52-9d0c-1c4b418317d8	ETH	2022-10-18 14:54:43.491+00
b6fc22c2-84b6-4018-8619-93db67d5bb1f	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	c5af3271-7185-4a52-9d0c-1c4b418317d8	SOL	2022-10-18 14:54:43.491+00
67655fcb-c375-408d-975d-b5ab346de3bb	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	c5af3271-7185-4a52-9d0c-1c4b418317d8	BTC	2022-10-18 14:54:43.491+00
a76fbcd5-d436-4e00-9e4f-40cd4cf68c2a	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	c5af3271-7185-4a52-9d0c-1c4b418317d8	ADA	2022-10-18 14:54:43.491+00
cfc27886-e5ac-4183-949a-cfcf183b8258	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	c5af3271-7185-4a52-9d0c-1c4b418317d8	MATIC	2022-10-18 14:54:43.491+00
de2d4e35-8018-43d1-a29a-8fcfd42ddaf2	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	c5af3271-7185-4a52-9d0c-1c4b418317d8	ATOM	2022-10-18 14:54:43.491+00
8e2ee8eb-2057-4996-8d32-9eef6a6ab824	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	36ae23e7-d79e-4901-89d5-3fa87cca1abf	USD	2022-10-18 14:54:43.491+00
63b1016e-e62b-4f71-94bb-94d399576ac7	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	36ae23e7-d79e-4901-89d5-3fa87cca1abf	ETH	2022-10-18 14:54:43.491+00
4ee7c378-345d-4451-8b4b-ecd89eb929be	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	36ae23e7-d79e-4901-89d5-3fa87cca1abf	SOL	2022-10-18 14:54:43.491+00
ae4a10aa-9f30-4b43-92bb-a5788d40e565	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	36ae23e7-d79e-4901-89d5-3fa87cca1abf	BTC	2022-10-18 14:54:43.491+00
744d2936-2b30-4f50-a555-f19c71a3264f	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	36ae23e7-d79e-4901-89d5-3fa87cca1abf	ADA	2022-10-18 14:54:43.491+00
71838256-a3df-4838-93f3-134eb94382dc	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	36ae23e7-d79e-4901-89d5-3fa87cca1abf	MATIC	2022-10-18 14:54:43.491+00
96cb37f3-69ce-44ad-aeec-c38a770e79ff	d263e7e3-24d7-4c04-8d67-ea3a0be7907e	36ae23e7-d79e-4901-89d5-3fa87cca1abf	ATOM	2022-10-18 14:54:43.491+00
\.


--
-- Data for Name: account_balance; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.account_balance (id, account_id, balance, hold, available, created_at, request_id, count) FROM stdin;
d170ce3e-cae7-4e48-87df-5b51bda2474f	0adbb104-fc18-46ca-a4eb-beee7775eb69	10000	0	10000	2022-10-18 14:54:43.491+00	\N	1
c81ade5d-6e5e-44d0-80c9-87ae403644ac	af271fce-c63c-4d79-af50-dcded9f806a4	0	0	0	2022-10-18 14:54:43.491+00	\N	1
23d65705-3af2-4a07-8543-b8249d5715b6	b6fc22c2-84b6-4018-8619-93db67d5bb1f	0	0	0	2022-10-18 14:54:43.491+00	\N	1
683c2cc0-f829-4605-a033-327ceb042ac5	67655fcb-c375-408d-975d-b5ab346de3bb	0	0	0	2022-10-18 14:54:43.491+00	\N	1
7daecb6c-5995-43f0-9d3d-07a6b4ec6b72	a76fbcd5-d436-4e00-9e4f-40cd4cf68c2a	0	0	0	2022-10-18 14:54:43.491+00	\N	1
bb4478a4-673a-4d73-aace-bfa0876db33a	cfc27886-e5ac-4183-949a-cfcf183b8258	0	0	0	2022-10-18 14:54:43.491+00	\N	1
06b6846c-6579-41f5-b6e2-3caf85a28b35	de2d4e35-8018-43d1-a29a-8fcfd42ddaf2	0	0	0	2022-10-18 14:54:43.491+00	\N	1
0d63ab7f-ecf0-4604-a7a6-82dfc5458afd	8e2ee8eb-2057-4996-8d32-9eef6a6ab824	10000	0	10000	2022-10-18 14:54:43.491+00	\N	0
fe740d41-44ae-4e05-b80a-55bd1affc0dd	63b1016e-e62b-4f71-94bb-94d399576ac7	0	0	0	2022-10-18 14:54:43.491+00	\N	1
933f5347-968c-4d64-95ff-e0d83a1b3847	4ee7c378-345d-4451-8b4b-ecd89eb929be	0	0	0	2022-10-18 14:54:43.491+00	\N	1
cd12d369-aab8-4ed0-aa97-3a3d6c3555cf	ae4a10aa-9f30-4b43-92bb-a5788d40e565	0	0	0	2022-10-18 14:54:43.491+00	\N	1
21fc0067-73b1-4534-8164-747007ebe8f4	744d2936-2b30-4f50-a555-f19c71a3264f	0	0	0	2022-10-18 14:54:43.491+00	\N	1
c91f2188-f96e-4683-913f-e12412015dac	71838256-a3df-4838-93f3-134eb94382dc	0	0	0	2022-10-18 14:54:43.491+00	\N	1
0e879393-ed82-4200-80b8-ef45a4cd3c52	96cb37f3-69ce-44ad-aeec-c38a770e79ff	0	0	0	2022-10-18 14:54:43.491+00	\N	1
\.


--
-- Data for Name: entry; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.entry (id, account_id, transaction_id, request_id, amount, direction, created_at) FROM stdin;
\.


--
-- Data for Name: finalized_transaction; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.finalized_transaction (transaction_id, completed_at, canceled_at, failed_at, request_id) FROM stdin;
\.


--
-- Data for Name: hold; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.hold (id, account_id, transaction_id, request_id, amount, created_at) FROM stdin;
\.


--
-- Data for Name: released_hold; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.released_hold (hold_id, released_at, request_id) FROM stdin;
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
14	f
\.


--
-- Data for Name: transaction; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transaction (id, sender_id, receiver_id, request_id, transaction_type, created_at) FROM stdin;
\.


--
-- Name: account_balance account_balance_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_balance
    ADD CONSTRAINT account_balance_pkey PRIMARY KEY (id);


--
-- Name: account account_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_pkey PRIMARY KEY (id);


--
-- Name: entry entry_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entry
    ADD CONSTRAINT entry_pkey PRIMARY KEY (id);


--
-- Name: finalized_transaction finalized_transaction_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.finalized_transaction
    ADD CONSTRAINT finalized_transaction_pkey PRIMARY KEY (transaction_id);


--
-- Name: hold hold_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hold
    ADD CONSTRAINT hold_pkey PRIMARY KEY (id);


--
-- Name: released_hold released_hold_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.released_hold
    ADD CONSTRAINT released_hold_pkey PRIMARY KEY (hold_id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: transaction transaction_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction
    ADD CONSTRAINT transaction_pkey PRIMARY KEY (id);


--
-- Name: account_balance account_balance_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_balance
    ADD CONSTRAINT account_balance_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: entry entry_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entry
    ADD CONSTRAINT entry_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: entry entry_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entry
    ADD CONSTRAINT entry_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.transaction(id);


--
-- Name: finalized_transaction finalized_transaction_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.finalized_transaction
    ADD CONSTRAINT finalized_transaction_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.transaction(id);


--
-- Name: hold hold_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hold
    ADD CONSTRAINT hold_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: hold hold_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hold
    ADD CONSTRAINT hold_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.transaction(id);


--
-- Name: released_hold released_hold_hold_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.released_hold
    ADD CONSTRAINT released_hold_hold_id_fkey FOREIGN KEY (hold_id) REFERENCES public.hold(id);


--
-- Name: transaction transaction_receiver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction
    ADD CONSTRAINT transaction_receiver_id_fkey FOREIGN KEY (receiver_id) REFERENCES public.account(id);


--
-- Name: transaction transaction_sender_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction
    ADD CONSTRAINT transaction_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES public.account(id);


--
-- PostgreSQL database dump complete
--

