package middleware

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	. "github.com/easykoo/go-blog/common"

	"encoding/base32"
	"net/http"
	"strings"
	"time"
)

// MemoryStore ------------------------------------------------------------

func NewMemoryStore(age int) *MemoryStore {
	memoryStore := &MemoryStore{
		Codecs: securecookie.CodecsFromPairs([]byte("secret")),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: age, //default 30 minutes
		},
		Container: make(map[string]*SessionInfo),
	}
	go memoryStore.CheckMemorySessions()
	return memoryStore
}

type MemoryStore struct {
	Codecs    []securecookie.Codec
	Options   *sessions.Options // default configuration
	Container map[string]*SessionInfo
}

type SessionInfo struct {
	S *sessions.Session
	T time.Time
}

// Get returns a session for the given name after adding it to the registry.
//
func (s *MemoryStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	Log.Infol("MemoryStore Get()")
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
func (s *MemoryStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		Log.Infol("MemoryStore reading cookie")
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		PanicIf(err)
		if err == nil {
			Log.Infol("MemoryStore read cookie success")
			err = s.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (s *MemoryStore) Save(r *http.Request, w http.ResponseWriter,
	session *sessions.Session) error {
	Log.Infol("MemoryStore Save()", session.ID)
	if session.ID == "" {
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}
	if err := s.save(session); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID,
		s.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

// save writes encoded session.Values to a file.
func (s *MemoryStore) save(session *sessions.Session) error {
	Log.Infol("MemoryStore save()")
	s.Container[session.ID] = &SessionInfo{S: session, T: time.Now()}
	return nil
}

// load reads a file and decodes its content into session.Values.
func (s *MemoryStore) load(session *sessions.Session) error {
	Log.Infol("MemoryStore load()")
	Log.Infol("MemoryStore load session: ", session.ID)
	if _, ok := s.Container[session.ID]; ok {
		Log.Infol("MemoryStore load session OK ")
		session = s.Container[session.ID].S
		Log.Infol("MemoryStore load SignedUser: ", session.Values["SignedUser"])
	} else {
		Log.Infol("MemoryStore load session failed ")
	}
	return nil
}

func (s *MemoryStore) CheckMemorySessions() {
	timer := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-timer.C:
			go func() {
				s.removeMemorySessions()
			}()
		}
	}

}

func (s *MemoryStore) removeMemorySessions() {
	for sId, sessionInfo := range s.Container {
		if (time.Now().Unix() - sessionInfo.T.Unix()) >= int64(s.Options.MaxAge) {
			Log.Infol(time.Now().Unix() - sessionInfo.T.Unix())
			delete(s.Container, sId)
			Log.Infol("Removed: ", sId)
		}
	}
}
