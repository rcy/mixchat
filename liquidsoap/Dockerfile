FROM phasecorex/liquidsoap

RUN apt-get update && apt-get install -y curl

WORKDIR /work

COPY ./emergency.ogg /work/
COPY ./start.liq /work/
COPY ./request.sh /work/

ENTRYPOINT ["user-entrypoint", "/app/liquidsoap", "/work/start.liq"]
