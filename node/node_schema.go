package node

import (
	"database/sql"
)

func CreateTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS template (
			node_type VARCHAR(32) PRIMARY KEY NOT NULL,
			template_name VARCHAR(120)
		)`)
	checkerr(err)
	db.Exec(`INSERT INTO template
		(node_type, template_name)
		VALUES
		(?, ?)`, "webpage", "Web page")
	db.Exec(`INSERT INTO template
		(node_type, template_name)
		VALUES
		(?, ?)`, "folder", "Folder")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS field (
			field_id INTEGER PRIMARY KEY,
			node_type VARCHAR(32),
			field_name VARCHAR(120),
			field_type VARCHAR(32),
			field_size INTEGER,
			field_default VARCHAR(120),
			field_description VARCHAR(120),
			field_tooltip VARCHAR(360),
			FOREIGN KEY (node_type) REFERENCES template(node_type) ON DELETE CASCADE
		)`)
	checkerr(err)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_field_name ON field(field_name)`)
	checkerr(err)
	db.Exec(`INSERT INTO field
		(node_type, field_name, field_type, field_size, field_default, field_description, field_tooltip)
		VALUES
		(?, ?, ?, ?, ?, ?, ?)`, "webpage", "URL", "String", 180, "", "URL", "URL")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS node (
            node_id CHAR(36) PRIMARY KEY NOT NULL,
            node_type VARCHAR(32),
            node_name VARCHAR(120),
			parent_id CHAR(36),
			node_values BLOB,
			FOREIGN KEY (node_type) REFERENCES template(node_type) ON DELETE CASCADE
        )`)
	checkerr(err)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_node_name ON node(node_name)`)
	checkerr(err)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_parent_id ON node(parent_id)`)
	checkerr(err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tag (
			tag_name VARCHAR(120) PRIMARY KEY NOT NULL,
			tag_count INTEGER DEFAULT 0,
			tag_timestamp INTEGER DEFAULT CURRENT_TIMESTAMP
		)`)
	checkerr(err)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_tag_count ON tag(tag_count)`)
	checkerr(err)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_tag_timestamp on tag(tag_timestamp)`)
	checkerr(err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS node_tags (
			node_tags_id INTEGER PRIMARY KEY,
			node_id CHAR(36),
			tag_name VARCHAR(120),
			FOREIGN KEY (node_id) REFERENCES node(node_id) ON DELETE CASCADE,
			FOREIGN KEY (tag_name) REFERENCES tag(tag_name) ON DELETE CASCADE
		)`)
	checkerr(err)
	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	checkerr(err)
}
