package vkutil

import "testing"

func TestParseCallbackURL(t *testing.T) {
	var inputs = []struct {
		URL   string
		Token string
	}{
		{"https://oauth.vk.com/blank.html#access_token=340277088251f4206985357919e28fd8879271bff9e509d205b8d36b7352cdd83557a8c13c08e026f2cef&expires_in=0&user_id=154766939",
			"340277088251f4206985357919e28fd8879271bff9e509d205b8d36b7352cdd83557a8c13c08e026f2cef"},
	}

	for _, inp := range inputs {
		token, _ := ParseCallbackURL(inp.URL)
		if token != inp.Token {
			t.Fatalf("error parsing callback url expected token %s != got token %s", inp.Token, token)
		}
	}

}

func TestParseDomain(t *testing.T) {
	var inputs = []struct {
		URL    string
		Domain string
	}{
		{"https://vk.com/zhuharev",
			"zhuharev"},
	}

	for _, inp := range inputs {
		token, err := ParseDomain(inp.URL)
		if err != nil {
			t.Fatalf("err get fomain %s", err)
		}
		if token != inp.Domain {
			t.Fatalf("error parsing expected domain %s != got domain %s", inp.Domain, token)
		}
	}
}

func TestParseBoardURL(t *testing.T) {
	var inputs = []struct {
		URL     string
		GroupID int
		TopicID int
		PostID  int
	}{
		{"https://vk.com/topic-57466174_31977455?offset=31400",
			57466174, 31977455, 31400},
		{"https://vk.com/topic-57466174_31977455",
			57466174, 31977455, 0},
	}

	for _, inp := range inputs {
		gID, tID, pID, err := ParseBoardURL(inp.URL)
		if err != nil {
			t.Fatalf("err get fomain %s", err)
		}
		if gID != inp.GroupID {
			t.Fatalf("error parsing callback url expected token %d != got token %d", inp.GroupID, gID)
		}
		if tID != inp.TopicID {
			t.Fatalf("error parsing callback url expected token %d != got token %d", inp.TopicID, tID)
		}
		if pID != inp.PostID {
			t.Fatalf("error parsing callback url expected token %d != got token %d", inp.PostID, pID)
		}
	}

}
