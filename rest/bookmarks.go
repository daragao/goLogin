package rest

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/daragao/goLogin/db"
	//"github.com/daragao/goLogin/logger"
	"github.com/daragao/goLogin/session"
	"net/http"
	"net/url"
	"strconv"
)

type Bookmarks struct {
	Id           int
	UserId       int
	BookmarkType int
	NewsItemId   string
}

func (self *Bookmarks) GetBookmarkByID(w rest.ResponseWriter, r *rest.Request) {
	idStr := r.PathParam("id")
	id, err := strconv.Atoi(idStr)
	bookmarkRow, err := db.GetBookmarkByID(id)
	if err == nil {
		w.WriteJson(bookmarkRow)
	} else {
		rest.Error(w, "Bookmark not found: "+err.Error(), http.StatusNotFound)
	}
}

func (self *Bookmarks) DeleteBookmarkByID(w rest.ResponseWriter, r *rest.Request) {
	idStr := r.PathParam("id")
	err := db.DeleteBookmark(idStr)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		rest.Error(w, "Bookmark not found: "+err.Error(), http.StatusInternalServerError)
	}
}

func (self *Bookmarks) GetAllBookmarks(w rest.ResponseWriter, r *rest.Request) {
	// TODO: set limit and offset
	currSession, _ := session.Get(r.Request)
	userId := currSession.Values["userId"]
	if userId == nil {
		rest.Error(w, "Invalid login", http.StatusUnauthorized)
		return
	}
	req := *r.Request
	urlQuery, err := url.ParseQuery(req.URL.RawQuery)
	bookmarkType := 0
	bookmarkTypeStr, errUrlQueryMap := urlQuery["bookmarkType"]
	if errUrlQueryMap {
		bookmarkType, err = strconv.Atoi(bookmarkTypeStr[0])
	}
	if err != nil {
		rest.Error(w, "Bookmark not found: "+err.Error(), http.StatusInternalServerError)
		return
	}
	bookmarkRows, err := db.GetAllBookmarks(100, 0, userId.(int), bookmarkType)
	if err == nil {
		w.WriteJson(bookmarkRows)
	} else {
		rest.Error(w, "Bookmark not found: "+err.Error(), http.StatusNotFound)
	}
}

func (self *Bookmarks) InsertBookmark(w rest.ResponseWriter, r *rest.Request) {
	bookmarkStruct := Bookmarks{}
	err := r.DecodeJsonPayload(&bookmarkStruct)
	if err != nil {
		rest.Error(w, "Could not decode JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bookmarkRow, err := db.InsertBookmark(bookmarkStruct.UserId,
		bookmarkStruct.NewsItemId, bookmarkStruct.BookmarkType)
	if err == nil {
		w.WriteJson(bookmarkRow)
	} else {
		rest.Error(w, "Could not insert bookmark: "+err.Error(), http.StatusInternalServerError)
	}
}
