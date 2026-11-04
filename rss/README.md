# rss

Generates a simple RSS-style XML feed with `encoding/xml` and writes it to `feed.xml`.

## What the code does

`main.go` defines `Channel` and `Item` structs with `xml:"..."` tags. `generateRSSFeed` fills them in:

- a channel titled "Sample RSS Feed" (`http://example.com`)
- three hard-coded items, `Article 1`, `Article 2`, and `Article 3`, with `pubDate` set to now, now minus 1h, and now
- the result is marshaled with `xml.MarshalIndent` and the standard `<?xml ...?>` header is prepended

`main` writes the bytes to `feed.xml` (it overwrites any existing file) and prints `RSS feed generated successfully and saved to feed.xml`. The committed `feed.xml` is a sample output.

## Concepts you can learn

- **Struct tags for XML**: field tags control element names, and a slice field (`Items []Item`) becomes repeated `<item>` elements.
- **`MarshalIndent` + `xml.Header`**: produces pretty-printed XML with the declaration that most consumers expect.
- **`time.Time` in XML**: by default, `encoding/xml` writes times in RFC 3339 format.

## How to run

```bash
cd rss
go run .
```

You don't need any external services. Running it overwrites `rss/feed.xml`.

## Notes

- **The output isn't valid RSS 2.0.** The root element is `<Channel>` (from the Go type name, because there's no `XMLName` field) instead of `<rss version="2.0"><channel>...</channel></rss>`. Also, `pubDate` uses RFC 3339 instead of the RFC 822 date format that RSS requires. Feed readers may reject it.
