package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	createSlotQuery = "SELECT * FROM pg_create_logical_replication_slot('%s', '%s');"
)

func createSlot(db *sqlx.DB, slot, plugin string) (created bool, err error) {
	var result []struct {
		Slot     string `db:"slot_name"`
		Position string `db:"xlog_position"`
	}
	query := fmt.Sprintf(createSlotQuery, slot, plugin)
	err = db.Select(&result, query)
	return result[0].Slot == slot, errors.Wrapf(err, "Could not create slot, query: %s", query)
}
