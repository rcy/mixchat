#FROM savonet/liquidsoap:v2.1.4
#FROM savonet/liquidsoap:ef2a2db-minimal
FROM savonet/liquidsoap:v2.3.0
#FROM phasecorex/liquidsoap:latest

WORKDIR /work

COPY ./emergency.ogg /work/
COPY ./404.ogg /work/
COPY ./start.liq /work/

ENTRYPOINT ["liquidsoap", "--force-start", "/work/start.liq"]
