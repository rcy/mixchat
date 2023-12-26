FROM savonet/liquidsoap:v2.1.4

WORKDIR /work

COPY ./emergency.ogg /work/
COPY ./404.ogg /work/
COPY ./start.liq /work/

ENTRYPOINT ["liquidsoap", "--force-start", "/work/start.liq"]
