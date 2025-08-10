# mod_browser_dl HTTP Server

This is a simple Go HTTP server/scraper for browsing and downloading module files from the Mod Archive.

## Features
- **Browse modules:** Search for songs/modules by keyword.
- **Download modules:** Download a module file by its ID.

## Endpoints
- `GET /browse?search=keyword` — Returns a JSON list of matching songs/modules.
- `POST /download?search={id}` — Downloads the module file with the given ID.

## How to Run
1. Make sure you have Go installed (Go 1.18 or newer recommended).
2. Navigate to the `src` directory:
   ```sh
   cd src
   ```
3. Run the server:
   ```sh
   go run httprequest.go httpretreiver.go
   ```
4. The server will start on port 3000 by default.

## Example Usage
- Browse: [http://localhost:3000/browse?search=thunder](http://localhost:3000/browse?search=thunder)
- Download (POST): `curl -X POST 'http://localhost:3000/download?search=194306'`

## Notes
- Downloads are saved to the `downloads/` directory.
- This project is for educational/demo purposes and may require further error handling for production use.

---

Feel free to modify or extend this project as needed!
