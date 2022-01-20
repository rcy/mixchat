--! Previous: sha1:c692b7f918e3cc7575bcf3887ef4cdf0d1771f68
--! Hash: sha1:7399968c0d242af7c14e1c6d2a8c04023848f356

-- Enter migration here

-- trigger job to wire up station to radio after it is created
drop trigger if exists insert_event on stations;

create trigger insert_event
  after insert on stations
  for each row
  execute procedure trigger_job('station_created');
