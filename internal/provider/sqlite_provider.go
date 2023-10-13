package provider

import (
	"database/sql"

	"github.com/callmemhz/alfred-cheatsheet/model"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteProviderFactory struct {
	Path string
}

func (factory *SqliteProviderFactory) NewProvider() (Provider, error) {
	db, err := sql.Open("sqlite3", factory.Path)
	if err != nil {
		return nil, err
	}
	p := &SqliteProvider{db: db}
	return p, nil
}

type SqliteProvider struct {
	db *sql.DB
}

func (p *SqliteProvider) Close() {
	p.db.Close()
}

func (p *SqliteProvider) Search(namespace, keyword string) ([]model.Entry, error) {
	_, err := p.db.Exec("create table if not exists documents (id integer not null primary key, title text, content text, type string, viewed integer);")
	if err != nil {
		return nil, err
	}

	stmt := "select title, content, type, viewed from documents"
	rows, err := p.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []model.Entry
	for rows.Next() {
		doc := new(Document)
		if err := rows.Scan(&doc.title, &doc.content, &doc.typ, &doc.viewed); err != nil {
			return nil, err
		}
		entries = append(entries, doc)
	}
	return entries, nil
}
