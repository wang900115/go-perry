# multipartload

File upload over HTTP with only the standard library: a raw-body upload endpoint and a `multipart/form-data` endpoint, plus an HTML page with one form for each.

## What the code does

Module `multipart` (no dependencies).

- **`server/main.go`**: `http.ServeMux` on `:8080` with:
  - `/` -> `Hello, World` (also catches every unmatched path).
  - `/upload` -> `controllers.FileUpload`.
  - `/upload_multipart` -> `controllers.FileUploadMultipart`.
- **`server/controllers/file.go`**:
  - `FileUpload`: limits the body to 100 MB, calls `ParseForm`, reads the target name from the form/query value `file`, creates `./uploads/<name>` and copies the **raw request body** into it.
  - `FileUploadMultipart`: limits the body to 100 MB, `ParseMultipartForm(10 MB)` (larger parts spill to temp files), takes the `file` part, prints its filename, size and MIME header on the server, and saves it as `./uploads/<original filename>`.
  - Both reply `File upload successfully: <name>`.
- **`client/client.html`**: two forms posting to `http://localhost:8080/upload` (no `enctype`) and `/upload_multipart` (`multipart/form-data`).

## Concepts you can learn

- **`multipart/form-data`**: the browser's way of sending files with boundaries and per-part headers; `FormFile` gives you the part and its `FileHeader`.
- **Memory vs. disk in `ParseMultipartForm`**: the argument caps in-memory size, the rest goes to temporary files.
- **`http.MaxBytesReader`**: rejects oversized bodies instead of reading unlimited data.
- **Streaming with `io.Copy`**: the file is copied to disk without loading it fully into memory.

## How to run

```bash
cd multipartload
go run ./server
```

Uploads land in `uploads/` under the **current working directory** of the server process. Test with curl (verified):

```bash
curl -F "file=@photo.png" http://localhost:8080/upload_multipart
curl -H "Content-Type: application/octet-stream" --data-binary @photo.png \
  "http://localhost:8080/upload?file=photo.png"
```

Or open `client/client.html` in a browser.

## Notes

- The first form in `client.html` does not work as a file upload: without `enctype` the browser sends only `file=<filename>` as urlencoded data, `ParseForm` consumes that body, and `/upload` creates a **0-byte** file. Send a raw body with a non-form `Content-Type` instead (as in the curl example).
- Neither handler checks the HTTP method.
- `/upload` takes the file name from the client's `file` query/form value and joins it into the path unchecked, so `?file=../x` writes outside `uploads/` (verified). `/upload_multipart` is not affected: `multipart.Part.FileName()` already applies `filepath.Base`. Existing files are overwritten in both.
