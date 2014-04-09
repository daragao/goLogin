package db

type Bookmark struct {
	// contains filtered or unexported fields
	id           int
	UserId       int
	NewsItemId   string
	BookmarkType int
}

func GetBookmarkByID(idArg int) (bookmarkRow *Bookmark, err error) {
	row, err := GetRowBy("bookmarks", "id,userId,newsItemId,bookmarkType", "id", idArg)
	bookmarkRow = &Bookmark{}
	err = row.Scan(&bookmarkRow.id, &bookmarkRow.UserId,
		&bookmarkRow.NewsItemId, &bookmarkRow.BookmarkType)
	return
}

func GetAllBookmarks(limit int, offset int, userId int,
	bookmarkType int) (bookmarks []Bookmark, err error) {
	// THIS IS VERY VERY BAD!!!!!!!
	// TODO: limit to type and user
	// TODO: GENERIC FUNCTION TO GET THE ROWS WITH CONSTRAINS!!!
	con, err := connectDB()
	defer con.Close()
	sqlStatement := "SELECT id,userId,newsItemId,bookmarkType FROM bookmarks" +
		" WHERE userid = $3 AND bookmarktype = $4 LIMIT $1 OFFSET $2"
	rows, err := con.Query(sqlStatement, limit, offset, userId, bookmarkType)
	if err != nil {
		//logger.ERRO.Println("Failed to get row: " + err.Error() + "\n\t" + sqlStatement)
	}
	// TODO: GENERIC FUNCTION TO GET THE ROWS WITH CONSTRAINS!!!
	bookmarks = make([]Bookmark, 0)
	for rows.Next() {
		bookmarkRow := Bookmark{}
		err = rows.Scan(&bookmarkRow.id, &bookmarkRow.UserId,
			&bookmarkRow.NewsItemId, &bookmarkRow.BookmarkType)
		bookmarks = append(bookmarks, bookmarkRow)
	}
	return
}

func InsertBookmark(userId int, newsItemId string,
	bookmarkType int) (bookmarkRow *Bookmark, err error) {
	con, err := connectDB()
	defer con.Close()
	sqlStatement := insertRowStr(bookmarks_table_name, "userid", "newsitemid", "bookmarktype")
	_, err = con.Exec(sqlStatement, userId, newsItemId, bookmarkType)
	if err == nil {
		//userRow, err = GetUserByUsername(username)
	} else {
		//userRow = nil
	}
	return
}

func DeleteBookmark(newsItemId string) (err error) {
	con, err := connectDB()
	defer con.Close()
	sqlStatement := "DELETE FROM bookmarks WHERE newsitemid = $1"
	_, err = con.Exec(sqlStatement, newsItemId)
	return
}
