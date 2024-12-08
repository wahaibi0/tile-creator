#!/bin/bash

HIGHEST_ZOOM=${HIGHEST_ZOOM:-5}
MIN_LAT=${MIN_LAT:-16.5}
MAX_LAT=${MAX_LAT:-26.5}
MIN_LON=${MIN_LON:-51.8}
MAX_LON=${MAX_LON:-60.0}
OUTPUT_DIR=${OUTPUT_DIR:-"./tiles"}
OUTPUT_MBTILES=${OUTPUT_MBTILES:-"./output.mbtiles"}

echo "Downloading tiles..."
./tile_creator -zoom "$HIGHEST_ZOOM" -minlat "$MIN_LAT" -maxlat "$MAX_LAT" -minlon "$MIN_LON" -maxlon "$MAX_LON" -output "$OUTPUT_DIR"

echo "Converting tiles to MBTiles..."
mb-util --image_format=png "$OUTPUT_DIR" "$OUTPUT_MBTILES"

echo "Inserting metadata into MBTiles..."

sqlite3 "$OUTPUT_MBTILES" <<EOF
CREATE TABLE IF NOT EXISTS metadata (name TEXT, value TEXT);
DELETE FROM metadata;

INSERT INTO metadata (name, value) VALUES ('name', 'Downloaded Tiles');
INSERT INTO metadata (name, value) VALUES ('type', 'overlay');
INSERT INTO metadata (name, value) VALUES ('version', '1.0');
INSERT INTO metadata (name, value) VALUES ('description', 'Satellite imagery tiles downloaded using Go');
INSERT INTO metadata (name, value) VALUES ('format', 'png');
INSERT INTO metadata (name, value) VALUES ('minzoom', '0');
INSERT INTO metadata (name, value) VALUES ('maxzoom', '$HIGHEST_ZOOM');
INSERT INTO metadata (name, value) VALUES ('bounds', '$MIN_LON,$MIN_LAT,$MAX_LON,$MAX_LAT');
EOF

echo "MBTiles file created at: $OUTPUT_MBTILES"
