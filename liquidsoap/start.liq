set("server.telnet",true)
set("server.telnet.port", 1234)
set("server.telnet.bind_addr", "0.0.0.0")
set("harbor.bind_addr","0.0.0.0")

emergency = single("./emergency.ogg")

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
    # Get the first line of my external process
    result = list.hd(default="", get_process_lines("sh ./request.sh " ^ slug))
    # Create and return a request using this result
    [request.create(result)]
  end

  def track_change(metadata) =
    print("track_change")
    filename = metadata["filename"]
    duration = request.duration(filename)
    metadata = list.append(metadata, [("duration", string_of(int_of_float(duration)))])
    print(metadata)
    result = http.post(getenv("API_BASE") ^ "/now/" ^ slug, data=json_of(metadata), headers=[("Content-Type","application/json")])
  end

  live = input.harbor("live/" ^ slug, port=8000, password=getenv("LIQUIDSOAP_BROADCAST_PASSWORD"))
  src = request.dynamic.list(my_request_function, id="dynlist_" ^ slug)
  radio = fallback(track_sensitive=false, [live, src, emergency])
  radio = on_track(track_change, radio)

  output.icecast(%vorbis, radio, mount=slug ^ ".ogg", host="icecast", password=getenv("ICECAST_SOURCE_PASSWORD"))
  output.icecast(%mp3, radio, mount=slug ^ ".mp3", host="icecast", password=getenv("ICECAST_SOURCE_PASSWORD"))

  add_skip_command(src)
end

# define an output to keep liquidsoap from exiting
output.icecast(%vorbis, emergency, mount="emergency.ogg", host="icecast", password=getenv("ICECAST_SOURCE_PASSWORD"))

def make_station_wrapper(slug) =
  make_station(slug)
  "Station created:" ^ slug
end
server.register(namespace="meta",usage="make_station slug",description="Make a station","make_station",make_station_wrapper)

make_station("emb")
