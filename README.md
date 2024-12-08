# Tile Creator and MBTiles Generator

This project automates the process of **downloading satellite imagery tiles** for a specified region and **converting them into an MBTiles file**. The workflow uses **Go**, **Bash**, and **Docker**, making it easy to configure, run, and package tile downloads.

---

## Project Structure

```
tile_creator/
├── Dockerfile             # Dockerfile to build the environment
├── download_tiles.go      # Go script to download tiles
└── entrypoint.sh          # Bash script to convert tiles and add metadata
```

---

## Features

- **Concurrent Downloads**: Speeds up tile downloads using Go's concurrency.
- **Configurable Regions and Zoom Levels**: Specify latitude, longitude, and zoom levels.
- **MBTiles Generation**: Converts downloaded tiles to MBTiles using `mb-util`.
- **Metadata Injection**: Adds essential metadata to the MBTiles file.
- **Dockerized Workflow**: Encapsulates dependencies in a Docker container for easy setup.

---

## Installation

1. **Clone the Repository**:

   ```bash
   git clone <repo-url>
   cd tile_creator
   ```

2. **Build the Docker Image**:

   ```bash
   docker build -t tile_creator .
   ```

---

## Usage

Run the following command to download tiles and generate an MBTiles file:

```bash
docker run --rm -it \
  -e HIGHEST_ZOOM=10 \
  -e MIN_LAT=16.5 \
  -e MAX_LAT=26.5 \
  -e MIN_LON=51.8 \
  -e MAX_LON=60.0 \
  -e OUTPUT_DIR="/app/tiles" \
  -e OUTPUT_MBTILES="/app/output/output.mbtiles" \
  -v $(pwd)/output:/app/output \
  -v $(pwd)/tiles:/app/tiles \
  tile_creator

```

### Explanation of Parameters

| **Parameter**    | **Description**                     | **Example**        |
| ---------------- | ----------------------------------- | ------------------ |
| `HIGHEST_ZOOM`   | Maximum zoom level for the download | `10`               |
| `MIN_LAT`        | Minimum latitude of the region      | `16.5`             |
| `MAX_LAT`        | Maximum latitude of the region      | `26.5`             |
| `MIN_LON`        | Minimum longitude of the region     | `51.8`             |
| `MAX_LON`        | Maximum longitude of the region     | `60.0`             |
| `OUTPUT_DIR`     | Directory to store downloaded tiles | `./tiles`          |
| `OUTPUT_MBTILES` | Name of the generated MBTiles file  | `./output.mbtiles` |

---

## Understanding Zoom Levels

Zoom levels determine the scale and detail of the map tiles. Here’s a breakdown of typical zoom levels:

| **Zoom Level** | **Description**                  | **Approximate Scale** | **Tile Grid Size**    |
| -------------- | -------------------------------- | --------------------- | --------------------- |
| **0**          | Entire world in one tile         | 1:500 million         | 1 x 1                 |
| **1**          | World divided into 4 tiles       | 1:250 million         | 2 x 2                 |
| **5**          | Large countries visible          | 1:4 million           | 32 x 32               |
| **10**         | Cities and towns visible         | 1:70,000              | 1024 x 1024           |
| **15**         | Streets and buildings visible    | 1:2,000               | 32,768 x 32,768       |
| **20**         | Individual buildings and details | 1:500                 | 1,048,576 x 1,048,576 |

- **Zoom Level 0**: The entire world fits in a single tile.
- Each increase in zoom level doubles the number of tiles along each axis.

---

## Dependencies

All dependencies are installed within the Docker container:

- **Go** (1.23)
- **Python 3** with `mb-util`
- **SQLite3**
- **Docker**

---

## Generated MBTiles

After running the script, you'll get an MBTiles file (e.g., `output.mbtiles`) in the specified directory. This file can be used with:

- **Map Servers**: TileServer-GL, GeoServer
- **GIS Tools**: QGIS, Mapbox
- **Custom Applications**: Web or mobile mapping apps

---

## Example Command for Oman at Zoom Level 10

To download tiles for **Oman** with a maximum zoom level of **10**, run:

```bash
docker run --rm -it \
  -e HIGHEST_ZOOM=10 \
  -e MIN_LAT=16.5 \
  -e MAX_LAT=26.5 \
  -e MIN_LON=51.8 \
  -e MAX_LON=60.0 \
  -e OUTPUT_DIR="./tiles" \
  -e OUTPUT_MBTILES="./oman.mbtiles" \
  tile_creator
```

This will generate an MBTiles file named `oman.mbtiles`.
