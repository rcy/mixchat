--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3 (Debian 13.3-1.pgdg100+1)
-- Dumped by pg_dump version 13.3 (Debian 13.3-1.pgdg100+1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: events; Type: TABLE; Schema: public; Owner: app
--

CREATE TABLE public.events (
    event_id text NOT NULL,
    event_type text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    payload jsonb NOT NULL
);


ALTER TABLE public.events OWNER TO app;

--
-- Name: schema_version; Type: TABLE; Schema: public; Owner: app
--

CREATE TABLE public.schema_version (
    version integer NOT NULL
);


ALTER TABLE public.schema_version OWNER TO app;

--
-- Name: station_messages; Type: TABLE; Schema: public; Owner: app
--

CREATE TABLE public.station_messages (
    station_message_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    type text NOT NULL,
    station_id text NOT NULL,
    parent_id text NOT NULL,
    nick text NOT NULL,
    body text NOT NULL
);


ALTER TABLE public.station_messages OWNER TO app;

--
-- Name: stations; Type: TABLE; Schema: public; Owner: app
--

CREATE TABLE public.stations (
    station_id text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    slug text NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    active boolean NOT NULL
);


ALTER TABLE public.stations OWNER TO app;

--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (event_id);


--
-- Name: station_messages station_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.station_messages
    ADD CONSTRAINT station_messages_pkey PRIMARY KEY (station_message_id);


--
-- Name: stations stations_pkey; Type: CONSTRAINT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.stations
    ADD CONSTRAINT stations_pkey PRIMARY KEY (station_id);


--
-- Name: stations stations_slug_key; Type: CONSTRAINT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.stations
    ADD CONSTRAINT stations_slug_key UNIQUE (slug);


--
-- Name: events events_after_insert; Type: TRIGGER; Schema: public; Owner: app
--

CREATE TRIGGER events_after_insert AFTER INSERT ON public.events FOR EACH ROW EXECUTE FUNCTION public.notify_event_insert();


--
-- PostgreSQL database dump complete
--
