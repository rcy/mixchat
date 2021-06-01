set("server.telnet",true)
set("server.telnet.port", 1234)
set("server.telnet.bind_addr", "0.0.0.0")

emergency = single("./emergency.ogg")

#library = playlist(reload_mode="watch", "/media")

def my_request_function() = 
  # Get the first line of my external process
  result = list.hd(default="", get_process_lines("sh ./request.sh"))
  # Create and return a request using this result
  [request.create(result)]
end

# Create the source
src = request.dynamic.list(my_request_function, id="dynlist")
#src = input.http("http://172.18.0.1:3010/next")

#requests = request.queue(id="request")

#tracks = fallback(track_sensitive=true, [requests, library])

radio = fallback(track_sensitive=false, [src, emergency])

output.icecast(%vorbis, radio, mount="emb.ogg", host="icecast", password=getenv("ICECAST_SOURCE_PASSWORD"))
output.icecast(%mp3, radio, mount="emb.mp3", host="icecast", password=getenv("ICECAST_SOURCE_PASSWORD"))


# Add a skip function to a source when it does not have one by default.
# @category Interaction
# @param s The source to attach the command to.
def add_skip_command(s) =
  # A command to skip
  def skip(_) =
    source.skip(s)
    "Done!"
  end
  # Register the command:
  server.register(namespace="#{source.id(s)}",
                  usage="skip",
                  description="Skip the current song.",
                  "skip",skip)
end

add_skip_command(src)
