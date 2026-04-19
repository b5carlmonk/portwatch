package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"portwatch/profile"
)

func makeProfile(name string) profile.Profile {
	return profile.Profile{
		Name:     name,
		Hosts:    []string{"127.0.0.1"},
		Ports:    []int{80, 443},
		Protocol: "tcp",
	}
}

func TestAddAndGet(t *testing.T) {
	s := profile.New()
	p := makeProfile("web")
	if err := s.Add(p); err != nil {
		t.Fatal(err)
	}
	got, ok := s.Get("web")
	if !ok {
		t.Fatal("expected profile to exist")
	}
	if got.Name != "web" {
		t.Errorf("expected web, got %s", got.Name)
	}
}

func TestAddEmptyNameErrors(t *testing.T) {
	s := profile.New()
	err := s.Add(profile.Profile{})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRemove(t *testing.T) {
	s := profile.New()
	_ = s.Add(makeProfile("web"))
	s.Remove("web")
	_, ok := s.Get("web")
	if ok {
		t.Fatal("expected profile to be removed")
	}
}

func TestSaveAndLoad(t *testing.T) {
	s := profile.New()
	_ = s.Add(makeProfile("db"))
	path := filepath.Join(t.TempDir(), "profiles.json")
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := profile.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := loaded.Get("db"); !ok {
		t.Error("expected db profile after load")
	}
}

func TestLoadMissingFile(t *testing.T) {
	s, err := profile.Load(filepath.Join(t.TempDir(), "none.json"))
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestGetMissing(t *testing.T) {
	s := profile.New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected false for missing profile")
	}
	_ = os.Getenv("") // suppress unused import warning
}
