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
		log.Fatalln(err.Error())
	}

	listMethod := "/api/list/"
	http.HandleFunc(listMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "GET", listMethod, wrapper.WebList)
	})
	matchMethod := "/api/match/"
	http.HandleFunc(matchMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "GET", matchMethod, wrapper.WebMatch)
	})
	getMethod := "/api/get/"
	http.HandleFunc(getMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "GET", getMethod, wrapper.WebGet)
	})
	allMethod := "/api/all/"
	http.HandleFunc(allMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "GET", allMethod, wrapper.WebGetAll)
	})
	labelMethod := "/api/label/"
	http.HandleFunc(labelMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "GET", labelMethod, wrapper.WebGetLabel)
	})
	setMethod := "/api/set/"
	http.HandleFunc(setMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "PUT", setMethod, wrapper.WebSet)
	})
	multisetMethod := "/api/multiset/"
	http.HandleFunc(multisetMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "PUT", multisetMethod, wrapper.WebSetMultiple)
	})
	clearMethod := "/api/clear/"
	http.HandleFunc(clearMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "DELETE", clearMethod, wrapper.WebClear)
	})
	insertBeforeMethod := "/api/insert-before/"
	http.HandleFunc(insertBeforeMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "POST", insertBeforeMethod, wrapper.WebInsertBefore)
	})
	insertAfterMethod := "/api/insert-after/"
	http.HandleFunc(insertAfterMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "POST", insertAfterMethod, wrapper.WebInsertAfter)
	})
	removeMethod := "/api/remove/"
	http.HandleFunc(removeMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "DELETE", removeMethod, wrapper.WebRemove)
	})
	moveMethod := "/api/move/"
	http.HandleFunc(moveMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "PATCH", moveMethod, wrapper.WebMove)
	})
	reloadMethod := "/api/reload/"
	http.HandleFunc(reloadMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "PATCH", reloadMethod, wrapper.WebReload)
	})
	saveMethod := "/api/save/"
	http.HandleFunc(saveMethod, func(w http.ResponseWriter, r *http.Request) {
		config.HandleRequest(w, r, "PATCH", saveMethod, wrapper.WebSave)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

	config.Close()
}
