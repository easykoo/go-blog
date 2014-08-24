package middleware

import (
	"github.com/gorilla/securecookie"
	. "github.com/gorilla/sessions"

	. "github.com/easykoo/go-blog/common"

	"encoding/base32"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// FileStore ------------------------------------------------------------

var fileMutex sync.RWMutex

// NewSessionStore returns a new FileStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewFileStore(age int) *FileStore {
	path := "."
	if path[len(path)-1] != '/' {
		path += "/"
	}
	sessionStore := &FileStore{
		Codecs: securecookie.CodecsFromPairs([]byte("secret")),
		Options: &Options{
			Path: "/",
			//			MaxAge: 86400 * 30,
			MaxAge: age, //30 minutes
		},
		path: path,
	}
	go sessionStore.CheckFileSessions()
	return sessionStore
}

// FileStore stores sessions in the filesystem.
//
// It also serves as a referece for custom stores.
//
// This store is still experimental and not well tested. Feedback is welcome.
type FileStore struct {
	Codecs  []securecookie.Codec
	Options *Options // default configuration
	path    string
}

// MaxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new FileStore is 4096.
func (s *FileStore) MaxLength(l int) {
	for _, c := range s.Codecs {
		if codec, ok := c.(*securecookie.SecureCookie); ok {
			codec.MaxLength(l)
		}
	}
}

// Get returns a session for the given name after adding it to the registry.
//
// See CookieStore.Get().
func (s *FileStore) Get(r *http.Request, name string) (*Session, error) {
	return GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See CookieStore.New().
func (s *FileStore) New(r *http.Request, name string) (*Session, error) {
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
func (s *FileStore) Save(r *http.Request, w http.ResponseWriter,
	session *Session) error {
	if session.ID == "" {
		// Because the ID is used in the filename, encode it to
		// use alphanumeric characters only.
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

// save writes encoded session.Values to a file.
func (s *FileStore) save(session *Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		s.Codecs...)
	if err != nil {
		return err
	}
	filename := s.path + "session_" + session.ID
	fileMutex.Lock()
	defer fileMutex.Unlock()
	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	if _, err = fp.Write([]byte(encoded)); err != nil {
		return err
	}
	fp.Close()
	return nil
}

// load reads a file and decodes its content into session.Values.
func (s *FileStore) load(session *Session) error {
	filename := s.path + "session_" + session.ID
	fp, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer fp.Close()
	var fdata []byte
	buf := make([]byte, 128)
	for {
		var n int
		n, err = fp.Read(buf[0:])
		fdata = append(fdata, buf[0:n]...)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	if err = securecookie.DecodeMulti(session.Name(), string(fdata),
		&session.Values, s.Codecs...); err != nil {
		return err
	}
	return nil
}

func (s *FileStore) CheckFileSessions() {
	s.removeFileSessions(true)
	timer := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-timer.C:
			go func() {
				s.removeFileSessions(false)
			}()
		}
	}

}

func (s *FileStore) removeFileSessions(first bool) {
	dir, err := os.Open(".")
	PanicIf(err)
	files, err := dir.Readdir(0)
	PanicIf(err)
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "session_") {
			if first {
				Log.Info("Removed: ", f.Name())
			} else {
				if time.Now().Unix()-f.ModTime().Unix() >= 60*30 {
					Log.Info(time.Now().Unix() - f.ModTime().Unix())
					os.Remove(f.Name())
					Log.Info("Removed: ", f.Name())
				}
			}
		}
	}
}
