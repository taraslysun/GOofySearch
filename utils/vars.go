package utils

import (
	"bytes"

	"github.com/elastic/go-elasticsearch/v8"
)

var CFG = elasticsearch.Config{
	CloudID: "4bcfd7625572429ab9bbf56b399bf60b:dXMtY2VudHJhbDEuZ2NwLmNsb3VkLmVzLmlvOjQ0MyQ1ZDQ2YzhlNjZjZDY0Y2YwOGQ1MDc4ODU5NzVmZDRjNyQ3M2I0ODg0YzM0MjA0ZjA5YTMwZThjZWQ2ZjM4MGQwMw==",
	APIKey:  "TzROSDNJMEIwUHd2b1l2S0VOblE6TG1XZFF3Y2xRSVdwLXRuUFk4UGtuZw==",
}

var BUF = bytes.NewBufferString(`
{"index":{"_id":"9780553351927"}}
{"name":"Snow Crash","author":"Neal Stephenson","release_date":"1992-06-01","page_count": 470}
{ "index": { "_id": "9780441017225"}}
{"name": "Revelation Space", "author": "Alastair Reynolds", "release_date": "2000-03-15", "page_count": 585}
{ "index": { "_id": "9780451524935"}}
{"name": "1984", "author": "George Orwell", "release_date": "1985-06-01", "page_count": 328}
{ "index": { "_id": "9781451673319"}}
{"name": "Fahrenheit 451", "author": "Ray Bradbury", "release_date": "1953-10-15", "page_count": 227}
{ "index": { "_id": "9780060850524"}}
{"name": "Brave New World", "author": "Aldous Huxley", "release_date": "1932-06-01", "page_count": 268}
{ "index": { "_id": "9780385490818"}}
{"name": "The Handmaid's Tale", "author": "Margaret Atwood", "release_date": "1985-06-01", "page_count": 311}
`)
