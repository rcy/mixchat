api_base = environment.get("API_BASE")
icecast_port = int_of_string(environment.get("ICECAST_PORT"))
icecast_host = environment.get("ICECAST_HOST")
icecast_source_password = environment.get("ICECAST_SOURCE_PASSWORD")
liquidsoap_broadcast_password = environment.get("LIQUIDSOAP_BROADCAST_PASSWORD")
station_slug = environment.get("STATION_SLUG")

print("API_BASE: " ^ api_base)
print("ICECAST_HOST: " ^ icecast_host)
print("ICECAST_PORT: " ^ string_of(icecast_port))
print("ICECAST_SOURCE_PASSWORD: " ^ icecast_source_password)
print("LIQUIDSOAP_BROADCAST_PASSWORD: " ^ liquidsoap_broadcast_password)
print("STATION_SLUG: " ^ station_slug)

set("init.allow_root",true)

set("server.telnet", true)
set("server.telnet.port", 1234)
set("server.telnet.bind_addr", "0.0.0.0")
set("harbor.bind_addr","0.0.0.0")

emergency = single("./emergency.ogg")

enable_replaygain_metadata()

def add_skip_command(s) =
  # A command to skip
  def skip(_) =
    source.skip(s)
    "Skipped!"
  end
  # Register the command:
  server.register(namespace="#{source.id(s)}",
                  usage="skip",
                  description="Skip the current song.",
                  "skip",skip)
end

def make_station(slug) =
  def my_request_function() =
    url = "#{api_base}/#{slug}/liq/pull"
    result = http.post(headers=[], url)

    # print("#{result}")
    print("my_request_function [#{url}] (#{slug}): #{result.status_code} #{result}")

    if result.status_code != 200 then
       []
    else
       [request.create("replaygain:" ^ result)]
    end
  end

  def track_change(m) =
    # print("track_change: " ^ m["filename"] ^ " -- " ^ m["replaygain_track_gain"] ^ " -- " ^ json_of(m))
    result = http.post(api_base ^ "/" ^ slug ^ "/liq/trackchange", data="filename=" ^ m["initial_uri"], headers=[("Content-Type","application/x-www-form-urlencoded")])
  end

  live = input.harbor("live/" ^ slug, port=8000, password=liquidsoap_broadcast_password)
  src = request.dynamic.list(my_request_function)
  src.on_track(track_change)
  radio = fallback(track_sensitive=false, [live, crossfade(blank.eat(replaygain(src))), emergency])

  #output.icecast(%vorbis, radio, mount=slug ^ ".ogg", port=icecast_port, host=icecast_host, password=icecast_source_password)
  output.icecast(%mp3,    radio, mount=slug ^ ".mp3", port=icecast_port, host=icecast_host, password=icecast_source_password)

  add_skip_command(src)
end

# define an output to keep liquidsoap from exiting
#output.icecast(%vorbis, emergency, mount="emergency.ogg", port=icecast_port, host=icecast_host, password=icecast_source_password)

make_station(station_slug)
