package core

import (
	"database/sql"
	"encoding/json"
)

// GetSegments returns the saved segments as a raw JSON array. Segments are a
// fork feature: named, reusable subscriber query expressions stored under the
// 'segments' key in the settings table (no dedicated table / migration needed).
// Returns "[]" when none have been saved yet.
func (c *Core) GetSegments() (json.RawMessage, error) {
	var b []byte
	if err := c.q.GetSegments.Get(&b); err != nil {
		if err == sql.ErrNoRows {
			return json.RawMessage("[]"), nil
		}
		c.log.Printf("error fetching segments: %v", err)
		return nil, err
	}
	if len(b) == 0 {
		return json.RawMessage("[]"), nil
	}
	return json.RawMessage(b), nil
}

// SaveSegments upserts the full segments array (raw JSON) into settings.
func (c *Core) SaveSegments(raw json.RawMessage) error {
	if len(raw) == 0 {
		raw = json.RawMessage("[]")
	}
	if _, err := c.q.SaveSegments.Exec(string(raw)); err != nil {
		c.log.Printf("error saving segments: %v", err)
		return err
	}
	return nil
}
