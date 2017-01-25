package mal

import (
	"os"
	"testing"
)

func testMAL() *MAL {
	// Enter your username and password to test.
	mal := NewMAL(os.Getenv("MAL_USERNAME"), os.Getenv("MAL_PASSWORD"))
	return mal
}

func TestSearch(t *testing.T) {
	mal := testMAL()
	results, err := mal.Search("naruto")
	if err != nil {
		t.Errorf("Error searching: %v", err)
	}

	if results == nil || len(results.Anime) == 0 {
		t.Errorf("Bad results: %v", results)
	}

	t.Log(*results)
}

func TestVerify(t *testing.T) {
	mal := testMAL()
	user, err := mal.Verify()
	if err != nil {
		t.Error(err.Error())
	}

	if user == nil {
		t.Error("Unknown error occurred")
	}

	if user.Username != mal.Username {
		t.Error("Username mismatch")
	}

	t.Log(*user)
}

func TestMyAnime(t *testing.T) {
	mal := testMAL()

	anime, err := mal.MyAnime()
	if err != nil {
		t.Error(err.Error())
	}

	if anime == nil {
		t.Error("Anime list is nil expected empty")
	}

	if len(anime.Anime) == 0 {
		t.Error("No anime found")
	}

	if anime.Username != mal.Username {
		t.Errorf("Expected username %s got %s", mal.Username, anime.Username)
	}

	t.Log(anime)
}

func TestAdd(t *testing.T) {
	mal := testMAL()

	ap := AnimePayload{
		Status:  "watching",
		Episode: 1,
		Tags:    []string{"Cool", "rad", "awesome"},
	}

	err := mal.Add(21, ap)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	mal := testMAL()

	ap := AnimePayload{
		Status: "completed",
		Tags:   []string{"long"},
	}

	err := mal.Update(21, ap)
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	mal := testMAL()
	err := mal.Delete(21)
	if err != nil {
		t.Error(err)
	}
}
