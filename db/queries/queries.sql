-- name: InsertArticle :one
INSERT INTO articles (
	feed_id, 
	title, 
	link, 
	published, 
	date_found, 
	summary,
	scraped_html,
	article_content
) VALUES (
	 ?, 
	 ?, 
	 ?, 
	 ?, 
	 ?, 
	 ?, 
	 ?,
	 ?
 ) RETURNING * ;

-- name: InsertOrIgnoreArticle :exec
INSERT OR IGNORE INTO articles (
	feed_id, 
	title, 
	link, 
	published, 
	date_found,
	summary 
) VALUES (
	 ?, 
	 ?, 
	 ?,
	 ?, 
	 ?, 
	 ?
 );


-- name: InsertFeed :one
 INSERT INTO feeds (
	url, 
	title, 
	css_sel_container,
	css_sel_start,
	css_sel_stop,
	html_extraction_strategy,
	last_fetched
) VALUES (
	?, 
	?, 
	?, 
	?, 
	?, 
	?, 
	CURRENT_TIMESTAMP
) RETURNING *;

-- name: SelectAllFeeds :many
SELECT * from feeds;	