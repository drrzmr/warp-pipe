package sql

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// SlotInfo object
type SlotInfo struct {
	SlotName   string
	Plugin     string
	SlotType   string
	Database   string
	RestartLSN string
	ActivePID  int64
	Active     bool
}

// ReplicationEvent object
type ReplicationEvent struct {
	Location string `db:"location"`
	XID      string `db:"xid"`
	Data     string `db:"data"`
}

const (
	pluginArgs    = "'include-xids', '1', 'include-timestamp', '1', 'skip-empty-xacts', '0', 'only-local', '0'"
	listSlotsCols = "slot_name, plugin, slot_type, database, active, active_pid, restart_lsn"

	createSlotQuery  = "SELECT * FROM pg_create_logical_replication_slot('%s', '%s');"
	listSlotsQuery   = "SELECT " + listSlotsCols + " FROM pg_replication_slots;"
	getChangesQuery  = "SELECT * FROM pg_logical_slot_get_changes('%s', NULL, NULL, " + pluginArgs + ");"
	peekChangesQuery = "SELECT * FROM pg_logical_slot_peek_changes('%s', NULL, NULL, " + pluginArgs + ");"
	dropSlotQuery    = "SELECT pg_drop_replication_slot('%s');"
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

func listSlots(db *sqlx.DB) (result []SlotInfo, err error) {

	var list []struct {
		SlotName   string        `db:"slot_name"`
		Plugin     string        `db:"plugin"`
		SlotType   string        `db:"slot_type"`
		Database   string        `db:"database"`
		RestartLSN string        `db:"restart_lsn"`
		ActivePID  sql.NullInt64 `db:"active_pid"`
		Active     bool          `db:"active"`
	}

	err = db.Select(&list, listSlotsQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not read slot list, query: %s", listSlotsQuery)
	}

	var activePID int64

	for _, info := range list {

		activePID = -1
		if info.ActivePID.Valid {
			activePID = info.ActivePID.Int64
		}

		result = append(result, SlotInfo{
			SlotName:   info.SlotName,
			Plugin:     info.Plugin,
			SlotType:   info.SlotType,
			Database:   info.Database,
			RestartLSN: info.RestartLSN,
			ActivePID:  activePID,
			Active:     info.Active,
		})
	}

	return result, nil
}

func getAllChanges(db *sqlx.DB, slot string) (result []ReplicationEvent, err error) {
	query := fmt.Sprintf(getChangesQuery, slot)
	err = db.Select(&result, query)
	return result, errors.Wrapf(err, "Could not get all changes, query: %s", query)
}

func peekAllChanges(db *sqlx.DB, slot string) (result []ReplicationEvent, err error) {
	query := fmt.Sprintf(peekChangesQuery, slot)
	err = db.Select(&result, query)
	return result, errors.Wrapf(err, "Could not peek all changes, query: %s", query)
}

func dropSlot(db *sqlx.DB, slot string) (err error) {
	query := fmt.Sprintf(dropSlotQuery, slot)
	_, err = db.Exec(query)
	return errors.Wrapf(err, "Could not drop slot, query: %s", query)
}
