--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3 (Debian 13.3-1.pgdg100+1)
-- Dumped by pg_dump version 13.7

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
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: current_bucket(integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.current_bucket(_station_id integer) RETURNS integer
    LANGUAGE sql
    AS $$
  select coalesce(( select bucket from tracks where station_id = _station_id order by bucket asc limit 1), 0);
$$;


--
-- Name: set_track_due_played(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_track_due_played() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
  -- push the track forward 1 step
  perform update_track_bucket(new.track_id, 1);
  return new;
end
$$;


--
-- Name: set_track_due_skipped(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_track_due_skipped() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
  -- push the track forward as many times as its been skipped; every time it is skipped it gets pushed out one further
  perform update_track_bucket(
    new.track_id,
    (
      select count(1)
        from plays
      where
        track_id = new.track_id
      and
        action = 'skipped'
    )::integer
  );
  return new;
end
$$;


--
-- Name: set_track_due_yeeted(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_track_due_yeeted() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
  -- push the track forward a lot of steps, effectively taking it out of rotation
  perform update_track_bucket(new.track_id, 1000000);
  return new;
end
$$;


--
-- Name: trigger_job(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.trigger_job() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
  perform graphile_worker.add_job(tg_argv[0], json_build_object(
    'schema', tg_table_schema,
    'table', tg_table_name,
    'op', tg_op,
    'id', (case when tg_op = 'delete' then old.id else new.id end)
  ));
  return new;
end;
$$;


--
-- Name: trigger_notify(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.trigger_notify() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
  perform pg_notify(tg_argv[0], (case when tg_op = 'delete' then old.id else new.id end)::text);
  return new;
end;
$$;


--
-- Name: trigger_notify_insert_station_relation_row(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.trigger_notify_insert_station_relation_row() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
  perform pg_notify(
    'postgraphile:station:' || new.station_id || ':' || tg_argv[0],
    json_build_object(
      '__node__',
      json_build_array(tg_argv[0], new.id)
    )::text
  );
  return new;
end
$$;


--
-- Name: update_track_bucket(integer, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_track_bucket(track_id integer, count integer) RETURNS integer
    LANGUAGE plpgsql
    AS $$
  begin
    update tracks set bucket = bucket + count, fuzz = random() where id = track_id;
    return 1;
  end
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.events (
    id integer NOT NULL,
    name text,
    data jsonb,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    station_id integer NOT NULL
);


--
-- Name: events_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.events_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.events_id_seq OWNED BY public.events.id;


--
-- Name: irc_channels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.irc_channels (
    id integer NOT NULL,
    station_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    server text DEFAULT 'irc.libera.chat'::text NOT NULL,
    channel text NOT NULL
);


--
-- Name: irc_channels_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.irc_channels_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: irc_channels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.irc_channels_id_seq OWNED BY public.irc_channels.id;


--
-- Name: messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.messages (
    id integer NOT NULL,
    station_id integer NOT NULL,
    body text,
    nick text,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: messages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.messages_id_seq OWNED BY public.messages.id;


--
-- Name: track_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.track_events (
    id integer NOT NULL,
    track_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    action text NOT NULL,
    station_id integer NOT NULL
);


--
-- Name: plays; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.plays AS
 SELECT track_events.id,
    track_events.track_id,
    track_events.created_at,
    track_events.action
   FROM public.track_events;


--
-- Name: plays_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.plays_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: plays_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.plays_id_seq OWNED BY public.track_events.id;


--
-- Name: tracks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tracks (
    id integer NOT NULL,
    filename text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    bucket integer NOT NULL,
    fuzz real DEFAULT 0 NOT NULL,
    event_id integer NOT NULL,
    station_id integer NOT NULL,
    metadata jsonb
);


--
-- Name: recently_added; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.recently_added AS
 SELECT tracks.id,
    tracks.filename,
    tracks.created_at AS added_at,
    tracks.event_id,
    events.name,
    events.data
   FROM (public.tracks
     JOIN public.events ON ((tracks.event_id = events.id)))
  ORDER BY tracks.created_at DESC
 LIMIT 10;


--
-- Name: recently_played; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.recently_played AS
 SELECT track_events.id,
    track_events.track_id,
    tracks.created_at,
    track_events.created_at AS played_at,
    tracks.filename,
    tracks.event_id,
    track_events.action
   FROM (public.track_events
     JOIN public.tracks ON ((track_events.track_id = tracks.id)))
  WHERE ((track_events.action = 'played'::text) AND (track_events.created_at > (now() - '01:00:00'::interval hour)))
  ORDER BY track_events.created_at DESC
 LIMIT 100;


--
-- Name: results; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.results (
    id integer NOT NULL,
    name text,
    data jsonb,
    event_id integer NOT NULL,
    station_id integer NOT NULL
);


--
-- Name: results_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.results_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: results_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.results_id_seq OWNED BY public.results.id;


--
-- Name: skips; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.skips AS
 SELECT track_events.id,
    track_events.track_id,
    track_events.created_at AS ts
   FROM public.track_events
  WHERE (track_events.action = 'skipped'::text);


--
-- Name: stations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.stations (
    id integer NOT NULL,
    slug text NOT NULL,
    name text,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: stations_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.stations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: stations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.stations_id_seq OWNED BY public.stations.id;


--
-- Name: track_changes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.track_changes (
    id integer NOT NULL,
    track_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: track_changes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.track_changes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: track_changes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.track_changes_id_seq OWNED BY public.track_changes.id;


--
-- Name: tracks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tracks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tracks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tracks_id_seq OWNED BY public.tracks.id;


--
-- Name: events id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events ALTER COLUMN id SET DEFAULT nextval('public.events_id_seq'::regclass);


--
-- Name: irc_channels id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.irc_channels ALTER COLUMN id SET DEFAULT nextval('public.irc_channels_id_seq'::regclass);


--
-- Name: messages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages ALTER COLUMN id SET DEFAULT nextval('public.messages_id_seq'::regclass);


--
-- Name: results id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.results ALTER COLUMN id SET DEFAULT nextval('public.results_id_seq'::regclass);


--
-- Name: stations id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stations ALTER COLUMN id SET DEFAULT nextval('public.stations_id_seq'::regclass);


--
-- Name: track_changes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.track_changes ALTER COLUMN id SET DEFAULT nextval('public.track_changes_id_seq'::regclass);


--
-- Name: track_events id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.track_events ALTER COLUMN id SET DEFAULT nextval('public.plays_id_seq'::regclass);


--
-- Name: tracks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracks ALTER COLUMN id SET DEFAULT nextval('public.tracks_id_seq'::regclass);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: irc_channels irc_channels_channel_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.irc_channels
    ADD CONSTRAINT irc_channels_channel_key UNIQUE (channel);


--
-- Name: irc_channels irc_channels_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.irc_channels
    ADD CONSTRAINT irc_channels_pkey PRIMARY KEY (id);


--
-- Name: irc_channels irc_channels_station_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.irc_channels
    ADD CONSTRAINT irc_channels_station_id_key UNIQUE (station_id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: track_events plays_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.track_events
    ADD CONSTRAINT plays_pkey PRIMARY KEY (id);


--
-- Name: stations stations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stations
    ADD CONSTRAINT stations_pkey PRIMARY KEY (id);


--
-- Name: stations stations_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stations
    ADD CONSTRAINT stations_slug_key UNIQUE (slug);


--
-- Name: track_changes track_changes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.track_changes
    ADD CONSTRAINT track_changes_pkey PRIMARY KEY (id);


--
-- Name: tracks tracks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracks
    ADD CONSTRAINT tracks_pkey PRIMARY KEY (id);


--
-- Name: tracks tracks_unique_station_filename; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracks
    ADD CONSTRAINT tracks_unique_station_filename UNIQUE (station_id, filename);


--
-- Name: track_events broadcast_track_event; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER broadcast_track_event AFTER INSERT ON public.track_events FOR EACH ROW EXECUTE FUNCTION public.trigger_job('broadcast_track_event');


--
-- Name: events insert_event; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER insert_event AFTER INSERT ON public.events FOR EACH ROW EXECUTE FUNCTION public.trigger_job('event_created');


--
-- Name: stations insert_event; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER insert_event AFTER INSERT ON public.stations FOR EACH ROW EXECUTE FUNCTION public.trigger_job('station_created');


--
-- Name: messages insert_message_notify; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER insert_message_notify AFTER INSERT ON public.messages FOR EACH ROW EXECUTE FUNCTION public.trigger_notify_insert_station_relation_row('messages');


--
-- Name: results insert_result; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER insert_result AFTER INSERT ON public.results FOR EACH ROW EXECUTE FUNCTION public.trigger_notify('result');


--
-- Name: track_events queued_track_due_tg; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER queued_track_due_tg BEFORE INSERT ON public.track_events FOR EACH ROW WHEN ((new.action = 'queued'::text)) EXECUTE FUNCTION public.set_track_due_played();


--
-- Name: track_events skipped_track_due_tg; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER skipped_track_due_tg BEFORE INSERT ON public.track_events FOR EACH ROW WHEN ((new.action = 'skipped'::text)) EXECUTE FUNCTION public.set_track_due_skipped();


--
-- Name: tracks update_track_metadata; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_track_metadata AFTER INSERT ON public.tracks FOR EACH ROW EXECUTE FUNCTION public.trigger_job('update_track_metadata');


--
-- Name: track_events yeeted_track_due_tg; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER yeeted_track_due_tg BEFORE INSERT ON public.track_events FOR EACH ROW WHEN ((new.action = 'yeeted'::text)) EXECUTE FUNCTION public.set_track_due_yeeted();


--
-- Name: events events_station_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_station_id_fkey FOREIGN KEY (station_id) REFERENCES public.stations(id);


--
-- Name: results fk_results_event; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT fk_results_event FOREIGN KEY (event_id) REFERENCES public.events(id);


--
-- Name: irc_channels irc_channels_station_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.irc_channels
    ADD CONSTRAINT irc_channels_station_id_fkey FOREIGN KEY (station_id) REFERENCES public.stations(id);


--
-- Name: messages messages_station_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_station_id_fkey FOREIGN KEY (station_id) REFERENCES public.stations(id);


--
-- Name: track_events plays_track_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.track_events
    ADD CONSTRAINT plays_track_id_fkey FOREIGN KEY (track_id) REFERENCES public.tracks(id);


--
-- Name: results results_station_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT results_station_id_fkey FOREIGN KEY (station_id) REFERENCES public.stations(id);


--
-- Name: track_changes track_changes_track_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.track_changes
    ADD CONSTRAINT track_changes_track_id_fkey FOREIGN KEY (track_id) REFERENCES public.tracks(id);


--
-- Name: track_events track_events_station_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.track_events
    ADD CONSTRAINT track_events_station_id_fkey FOREIGN KEY (station_id) REFERENCES public.stations(id);


--
-- Name: tracks tracks_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracks
    ADD CONSTRAINT tracks_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id);


--
-- Name: tracks tracks_station_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracks
    ADD CONSTRAINT tracks_station_id_fkey FOREIGN KEY (station_id) REFERENCES public.stations(id);


--
-- PostgreSQL database dump complete
--

