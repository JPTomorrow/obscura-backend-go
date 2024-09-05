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
}
