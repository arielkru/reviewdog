package cookieman

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeCipher struct {
	Cipher
	fakeEncrypt func(plaintext []byte) ([]byte, error)
	fakeDecrypt func(ciphertext []byte) ([]byte, error)
}

func (f *fakeCipher) Encrypt(plaintext []byte) ([]byte, error) {
	if f.fakeEncrypt != nil {
		return f.fakeEncrypt(plaintext)
	}
	return plaintext, nil
}

func (f *fakeCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	if f.fakeDecrypt != nil {
		return f.fakeDecrypt(ciphertext)
	}
	return ciphertext, nil
}

func TestCookieStore_Set_Get(t *testing.T) {
	opt := CookieOption{}
	cookieman := New(&fakeCipher{}, opt)
	name := "vim"
	value := "vim vim vim"
	w := httptest.NewRecorder()

	vimStore := cookieman.NewCookieStore(name, nil)

	if vimStore.Name() != name {
		t.Errorf("CookieStore.Name() = %q, want %q", vimStore.Name(), name)
	}

	if err := vimStore.Set(w, []byte(value)); err != nil {
		t.Error(err)
	}

	response := w.Result()

	gotSetCookie := response.Header.Get("Set-Cookie")
	wantSetCookie := fmt.Sprintf("%s=%s", name, base64.URLEncoding.EncodeToString([]byte(value)))
	if gotSetCookie != wantSetCookie {
		t.Errorf("CookieStore.Get: Set-Cookie value: got %q, want %q", gotSetCookie, wantSetCookie)
	}

	req, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range response.Cookies() {
		req.AddCookie(c)
	}

	b, err := vimStore.Get(req)
	if err != nil {
		t.Fatal(err)
	}

	if got := string(b); got != value {
		t.Errorf("CookieStore.Get: got %q, want %q", got, value)
	}
}

func TestCookieStore_Clear(t *testing.T) {
	opt := CookieOption{}
	cookieman := New(&fakeCipher{}, opt)
	name := "vim"
	w := httptest.NewRecorder()

	vimStore := cookieman.NewCookieStore(name, nil)
	vimStore.Clear(w)

	if cookieLen := len(w.Result().Cookies()); cookieLen != 1 {
		t.Fatalf("got %d cookies, want 1 cookie", cookieLen)
	}
	cookie := w.Result().Cookies()[0]

	if cookie.Name != name {
		t.Errorf("Cookie.Name = %q, want %q", cookie.Name, name)
	}
	if cookie.MaxAge != -1 {
		t.Errorf("Cookie.MaxAge = %d, want -1", cookie.MaxAge)
	}
}

func TestCookieOption(t *testing.T) {
	defaultOpt := CookieOption{
		http.Cookie{
			Secure: false,
			MaxAge: 30,
		},
	}
	cookieman := New(&fakeCipher{}, defaultOpt)

	w := httptest.NewRecorder()
	opt := &CookieOption{
		http.Cookie{
			Domain:   "domain",
			Expires:  time.Now(),
			HttpOnly: true,
			MaxAge:   14,
			Path:     "/",
			Secure:   true,
		},
	}
	if err := cookieman.Set(w, "n", []byte(""), opt); err != nil {
		t.Fatal(err)
	}

	if cookieLen := len(w.Result().Cookies()); cookieLen != 1 {
		t.Fatalf("got %d cookies, want 1 cookie", cookieLen)
	}
	cookie := w.Result().Cookies()[0]

	if cookie.Domain != opt.Domain {
		t.Errorf("Cookie.Domain = %q, want %q", cookie.Domain, opt.Domain)
	}
	if cookie.Expires.Second() != opt.Expires.Second() {
		t.Errorf("Cookie.Expires = %q, want %q", cookie.Expires, opt.Expires)
	}
	if cookie.HttpOnly != opt.HttpOnly {
		t.Errorf("Cookie.HttpOnly = %v, want %v", cookie.HttpOnly, opt.HttpOnly)
	}
	if cookie.MaxAge != opt.MaxAge {
		t.Errorf("Cookie.MaxAge = %v, want %v", cookie.MaxAge, opt.MaxAge)
	}
	if cookie.Path != opt.Path {
		t.Errorf("Cookie.Path = %q, want %q", cookie.Path, opt.Path)
	}
	if cookie.Secure != opt.Secure {
		t.Errorf("Cookie.Secure = %v, want %v", cookie.Secure, opt.Secure)
	}
}
