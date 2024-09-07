/*
These are structures that can be easily be swapped between JSON, to return to the client side, and SQL, to
store them in a databasae.
*/
package db

type YoutubeVideo struct {
	Id          int    `json:"id" sql_name:"id" sql_props:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	Title       string `json:"title" sql_name:"title" sql_props:"TEXT UNIQUE NOT NULL"`
	Description string `json:"description" sql_name:"description" sql_props:"TEXT NOT NULL"`
	VideoTag    string `json:"video_tag" sql_name:"video_tag" sql_props:"TEXT NOT NULL"`
	Upvotes     int    `json:"upvotes" sql_name:"upvotes" sql_props:"INTEGER NOT NULL DEFAULT 0"`
	Downvotes   int    `json:"downvotes" sql_name:"downvotes" sql_props:"INTEGER NOT NULL DEFAULT 0"`
}
