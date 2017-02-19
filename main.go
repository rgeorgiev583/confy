// main.go
package main

import (
	"log"

	"net/http"

	"github.com/rgeorgiev583/confy/wrapper"
)

func main() {
	config, err := wrapper.New("/", "", 0)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer config.Close()

	handle := func(prefix string, method string, handler RequestHandler) {
		http.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
			config.HandleRequest(w, r, method, prefix, handler)
		})
	}

	handle("/api/list/", "GET", wrapper.WebList)
	handle("/api/match/", "GET", wrapper.WebMatch)
	handle("/api/get/", "GET", wrapper.WebGet)
	handle("/api/all/", "GET", wrapper.WebGetAll)
	handle("/api/label/", "GET", wrapper.WebGetLabel)
	handle("/api/set/", "PUT", wrapper.WebSet)
	handle("/api/multiset/", "PUT", wrapper.WebSetMultiple)
	handle("/api/clear/", "PUT", wrapper.WebClear)
	handle("/api/insert-before/", "POST", wrapper.WebInsertBefore)
	handle("/api/insert-after/", "POST", wrapper.WebInsertAfter)
	handle("/api/remove/", "DELETE", wrapper.WebRemove)
	handle("/api/move/", "PATCH", wrapper.WebMove)
	handle("/api/reload/", "PATCH", wrapper.WebReload)
	handle("/api/save/", "PATCH", wrapper.WebSave)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
