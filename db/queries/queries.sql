

-- name: InsertOrIgnoreArticle :exec
INSERT OR IGNORE INTO articles (
	feed_id, 
	title, 
	link, 
	published, 
	date_found,
	summary,
	scraped_html,
	formatted_html
) VALUES (
	 ?, 
	 ?, 
	 ?, 
	 ?,
	 ?, 
	 ?, 
	 ?, 
	 ? 
 );

-- name: UpdateArticleByID :exec
UPDATE articles SET markdown = ? WHERE id = ?;

-- name: SelectArticleByFeedIDAndLink :one
SELECT * FROM articles WHERE feed_id = ? AND link = ?;

-- name: InsertFeed :one
 INSERT INTO feeds (
	url, 
	title, 
	css_sel_container,
	css_sel_start,
	css_sel_stop,
	html_extraction_strategy,
	last_fetched,
	folder_name
) VALUES (
	?, 
	?, 
	?, 
	?, 
	?, 
	?, 
	CURRENT_TIMESTAMP,
	?
) RETURNING *;

-- name: SelectAllFeeds :many
SELECT * from feeds;

-- name: InsertScraperRan :one
INSERT INTO log (
	articles_created, 
	run_type,
	time_ran
) VALUES (
	?,
	?, 
	CURRENT_TIMESTAMP
) RETURNING *;