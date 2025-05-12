package session

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockStore_PutGetUpdate(t *testing.T) {
	store := &MockStore{Sessions: make(map[string]*Session)}
	sess := &Session{SessionID: "abc", ClientID: "cid", AccountID: "aid", RoleName: "role", State: "state", Status: "pending"}
	// Put
	err := store.Put(context.Background(), sess)
	assert.NoError(t, err)
	// Get
	got, err := store.Get(context.Background(), "abc")
	assert.NoError(t, err)
	assert.Equal(t, sess, got)
	// Update
	sess.Status = "complete"
	err = store.Update(context.Background(), sess)
	assert.NoError(t, err)
	got, _ = store.Get(context.Background(), "abc")
	assert.Equal(t, "complete", got.Status)
}
