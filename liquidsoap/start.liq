set("init.allow_root",true)

set("server.telnet",true)
set("server.telnet.port", 1234)
set("server.telnet.bind_addr", "0.0.0.0")
set("harbor.bind_addr","0.0.0.0")

emergency = single("./emergency.ogg")

api_base = getenv("API_BASE")
icecast_port = int_of_string(getenv("ICECAST_PORT"))
icecast_host = getenv("ICECAST_HOST")
icecast_source_password = getenv("ICECAST_SOURCE_PASSWORD")
liquidsoap_broadcast_password = getenv("LIQUIDSOAP_BROADCAST_PASSWORD")

print("API_BASE: " ^ api_base)
print("ICECAST_HOST: " ^ icecast_host)
print("ICECAST_PORT: " ^ string_of(icecast_port))
print("ICECAST_SOURCE_PASSWORD: " ^ icecast_source_password)
print("LIQUIDSOAP_BROADCAST_PASSWORD: " ^ liquidsoap_broadcast_password)

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
    result = http.get(api_base ^ "/next/" ^ slug)

    print("my_request_function (#{slug}): #{result}")

    if result == "500" then
      []
    elsif result == "404" then
      []
    else
      # Create and return a request using this result
      [request.create("replaygain:" ^ result)]
    end
  end

  def track_change(m) =
    print("track_change: " ^ m["filename"] ^ " -- " ^ m["replaygain_track_gain"])
    result = http.post(api_base ^ "/now/" ^ slug, data="filename=" ^ m["filename"], headers=[("Content-Type","application/x-www-form-urlencoded")])
  end

  live = input.harbor("live/" ^ slug, port=8000, password=liquidsoap_broadcast_password)
  src = request.dynamic.list(my_request_function, id="dynlist_" ^ slug)
  src.on_track(track_change)
  radio = fallback(track_sensitive=false, [live, crossfade(blank.eat(replaygain(src))), emergency])

  #output.icecast(%vorbis, radio, mount=slug ^ ".ogg", port=icecast_port, host=icecast_host, password=icecast_source_password)
  output.icecast(%mp3,    radio, mount=slug ^ ".mp3", port=icecast_port, host=icecast_host, password=icecast_source_password)

  add_skip_command(src)
end

# define an output to keep liquidsoap from exiting
#output.icecast(%vorbis, emergency, mount="emergency.ogg", port=icecast_port, host=icecast_host, password=icecast_source_password)

def make_station_wrapper(slug) =
  make_station(slug)
  "Station created: " ^ slug
end
server.register(namespace="meta",usage="make_station slug",description="Make a station","make_station",make_station_wrapper)

make_station("emb")
make_station("tlm")
