set("server.telnet",true)
set("server.telnet.port", 1234)
set("server.telnet.bind_addr", "0.0.0.0")

emergency = single("./emergency.ogg")

def my_request_function() = 
  # Get the first line of my external process
  result = list.hd(default="", get_process_lines("sh ./request.sh"))
  # Create and return a request using this result
  [request.create(result)]
end

src = request.dynamic.list(my_request_function, id="dynlist")
radio = fallback(track_sensitive=false, [src, emergency])

output.icecast(%vorbis, radio, mount="emb.ogg", host="icecast", password=getenv("ICECAST_SOURCE_PASSWORD"))
output.icecast(%mp3, radio, mount="emb.mp3", host="icecast", password=getenv("ICECAST_SOURCE_PASSWORD"))

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

add_skip_command(src)
