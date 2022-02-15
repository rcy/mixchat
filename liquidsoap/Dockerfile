FROM savonet/liquidsoap:v1.4.3

WORKDIR /work

COPY ./emergency.ogg /work/
COPY ./404.ogg /work/
COPY ./start.liq /work/
COPY ./request.sh /work/

ENTRYPOINT ["liquidsoap", "--force-start", "/work/start.liq"]
