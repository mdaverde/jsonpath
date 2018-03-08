# jsonpath

Used to get and set data using JSONpath

## Example

```go
sample := `{ owner: { name: "john doe", contact: { phone: "555-555-5555" } } }`

var payload interface{}

err := json.Unmarshal([]byte(sample), &payload)
if err != nil {
    panic(err)
}

err = jsonpath.Set(&payload, "owner.contact.phone", "333-333-3333")
if err != nil {
    panic(err)
}
```

## Install

```bash
$ go get github.com/mdaverde/jsonpath
```
