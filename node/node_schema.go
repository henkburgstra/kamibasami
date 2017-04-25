package node

import (
	"database/sql"
)

func CreateTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS node (
            node_id CHAR(36) PRIMARY KEY NOT NULL,
            node_type VARCHAR(32),
            node_name VARCHAR(120),
			parent_id CHAR(36),
			node_values BLOB
        )`)
	checkerr(err)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_node_type ON node(node_type)`)
	checkerr(err)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_node_name ON node(node_name)`)
	checkerr(err)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_parent_id ON node(parent_id)`)
	checkerr(err)
}
