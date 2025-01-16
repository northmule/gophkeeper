package storage

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	session := NewSession()
	assert.NotNil(t, session)
	assert.IsType(t, &Session{}, session)
}

func TestAdd(t *testing.T) {
	session := NewSession().(*Session)

	token := "test-token"
	expire := time.Now().Add(1 * time.Hour)
	session.Add(token, expire)

	session.mx.RLock()
	defer session.mx.RUnlock()
	assert.Equal(t, expire, session.values[token])
}

func TestIsValid(t *testing.T) {
	session := NewSession().(*Session)

	token := "test-token"
	expire := time.Now().Add(1 * time.Hour)
	session.Add(token, expire)
	assert.True(t, session.IsValid(token))

	expiredToken := "expired-token"
	expiredExpire := time.Now().Add(-1 * time.Hour)
	session.Add(expiredToken, expiredExpire)
	assert.False(t, session.IsValid(expiredToken))

	nonExistentToken := "non-existent-token"
	assert.False(t, session.IsValid(nonExistentToken))
}

func TestConcurrency(t *testing.T) {

	session := NewSession().(*Session)

	var wg sync.WaitGroup
	numRoutines := 1000

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			token := fmt.Sprintf("token-%d", i)
			expire := time.Now().Add(1 * time.Hour)
			session.Add(token, expire)
		}(i)
	}
	wg.Wait()

	for i := 0; i < numRoutines; i++ {
		token := fmt.Sprintf("token-%d", i)
		assert.True(t, session.IsValid(token))
	}

}
