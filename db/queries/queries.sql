

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

-- name: UpdateArticleByID :one
UPDATE articles SET markdown = ? WHERE id = ? RETURNING *;

-- name: SelectArticleByFeedIDAndLink :one
SELECT * FROM articles WHERE feed_id = ? AND link = ?;

-- name: SelectArticleByID :one
SELECT * FROM articles WHERE id = ?;

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

-- name: UpdateArticleSetStarredValue :one
UPDATE articles SET starred = ? WHERE id = ? RETURNING *;

-- name: SelectArticlesByFeedIDWithLimit :many
SELECT
   a.id as article_id,
	a.link as article_link,
	a.title as article_title,
	a.starred as article_stars,
	a.published as article_published,
	a.read as article_read,
	a.feed_id as article_feed_id,
    f.title as feed_title
FROM articles a
    INNER JOIN feeds f ON f.id = a.feed_id
WHERE a.feed_id  = ?
ORDER BY published DESC
LIMIT ? OFFSET ?;

-- name: SelectFeedByID :one
SELECT * FROM feeds where id = ?;

-- name: SelectArticleCountByFeedID :one
SELECT COUNT(*) FROM articles WHERE feed_id = ?;

-- name: SelectArticlesWithFeedName :many
SELECT a.*, f.title AS feed_title FROM articles a INNER JOIN feeds f ON f.id = a.feed_id;
