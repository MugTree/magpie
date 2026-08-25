-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS feeds
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    url          TEXT NOT NULL UNIQUE,
    title        TEXT NOT NULL DEFAULT '',
    folder_name        TEXT NOT NULL DEFAULT '',
    last_fetched DATETIME NOT NULL,
    css_sel_container TEXT NOT NULL DEFAULT '',
    css_sel_start  TEXT NOT NULL DEFAULT '',
    css_sel_stop  TEXT NOT NULL DEFAULT '',
    html_extraction_strategy  TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS articles
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id   INTEGER NOT NULL,
    title     TEXT NOT NULL DEFAULT '',
    link      TEXT NOT NULL DEFAULT '',
    published DATETIME NULL,
    date_found DATETIME NULL,
    article_content      TEXT NOT NULL DEFAULT '',
    scraped_html    TEXT NOT NULL DEFAULT '',
    summary   TEXT NOT NULL DEFAULT '',
    UNIQUE (feed_id, link),
    FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE
);
-- +goose StatementEnd

