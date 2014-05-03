package middleware

import (
	"github.com/gorilla/securecookie"
	. "github.com/gorilla/sessions"

	. "github.com/easykoo/go-blog/common"
	"github.com/easykoo/go-blog/model"

	"encoding/base32"
	"net/http"
	"strings"
	"time"
)

// DbStore ------------------------------------------------------------

// NewSessionStore returns a new DbStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewDbStore(age int) *DbStore {
	sessionStore := &DbStore{
		Codecs: securecookie.CodecsFromPairs([]byte("secret")),
		Options: &Options{
			Path:   "/",
			MaxAge: age, //seconds
		},
	}
	go sessionStore.CheckDbSessions()
	return sessionStore
}

// DbStore stores sessions in the filesystem.
//
// It also serves as a referece for custom stores.
//
// This store is still experimental and not well tested. Feedback is welcome.
type DbStore struct {
	Codecs  []securecookie.Codec
	Options *Options // default configuration
}

// MaxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new DbStore is 4096.
func (s *DbStore) MaxLength(l int) {
	for _, c := range s.Codecs {
		if codec, ok := c.(*securecookie.SecureCookie); ok {
			codec.MaxLength(l)
		}
	}
}

// Get returns a session for the given name after adding it to the registry.
//
// See CookieStore.Get().
func (s *DbStore) Get(r *http.Request, name string) (*Session, error) {
	return GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See CookieStore.New().
func (s *DbStore) New(r *http.Request, name string) (*Session, error) {
	session := NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		PanicIf(err)
		if err == nil {
			err = s.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (s *DbStore) Save(r *http.Request, w http.ResponseWriter,
	session *Session) error {
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
	http.SetCookie(w, NewCookie(session.Name(), encoded, session.Options))
	return nil
}

// save writes encoded session.Values to db.
func (s *DbStore) save(session *Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		s.Codecs...)
	if err != nil {
		return err
	}
	sessionInfo := &model.SessionInfo{Id: session.ID, Content: encoded, Age: session.Options.MaxAge}
	if (&model.SessionInfo{Id: session.ID}).Exist() {
		sessionInfo.Update()
	} else {
		sessionInfo.Insert()
	}
	return nil
}

// load db and decodes its content into session.Values.
func (s *DbStore) load(session *Session) error {
	sessionInfo := (&model.SessionInfo{Id: session.ID}).GetSessionInfo()
	if err := securecookie.DecodeMulti(session.Name(), sessionInfo.Content,
		&session.Values, s.Codecs...); err != nil {
		return err
	}
	return nil
}

func (s *DbStore) CheckDbSessions() {
	timer := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-timer.C:
			go func() {
				s.removeDbSessions()
			}()
		}
	}

}

func (s *DbStore) removeDbSessions() {
	err := (&model.SessionInfo{}).RemoveExpiredSession()
	PanicIf(err)
}
