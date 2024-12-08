FROM golang:1.23

RUN apt-get update && apt-get install -y \
    python3 \
    python3-venv \
    sqlite3 \
    curl

RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

RUN pip install git+https://github.com/mapbox/mbutil.git

WORKDIR /app

COPY . .

RUN go build -o tile_creator tile_creator.go

ENTRYPOINT ["/bin/bash", "entrypoint.sh"]
