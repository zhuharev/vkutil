package vkutil

import "testing"

func TestParse(t *testing.T) {
	input := []byte(`{"response":[{"id":780206,"views":1511,"owner":-57466174},{"id":780202,"views":1856,"owner":-57466174},{"id":780195,"views":2657,"owner":-57466174},{"id":780186,"views":2592,"owner":-57466174},{"id":780176,"views":2102,"owner":-57466174},{"id":780170,"views":2127,"owner":-57466174},{"id":779914,"views":6416,"owner":-57466174},{"id":780144,"views":2807,"owner":-57466174},{"id":780143,"views":2605,"owner":-57466174},{"id":780136,"views":2668,"owner":-57466174},{"id":780131,"views":2454,"owner":-57466174},{"id":780123,"views":2705,"owner":-57466174},{"id":780120,"views":3011,"owner":-57466174},{"id":780094,"views":6403,"owner":-57466174},{"id":780034,"views":7564,"owner":-57466174},{"id":780023,"views":5911,"owner":-57466174},{"id":779998,"views":5730,"owner":-57466174},{"id":779993,"views":10174,"owner":-57466174},{"id":779984,"views":5690,"owner":-57466174},{"id":742778,"views":10241,"owner":-57466174},{"id":779955,"views":7639,"owner":-57466174},{"id":779952,"views":5569,"owner":-57466174},{"id":779947,"views":5744,"owner":-57466174},{"id":779945,"views":6277,"owner":-57466174},{"id":779938,"views":5293,"owner":-57466174},{"id":779935,"views":5511,"owner":-57466174},{"id":779930,"views":5622,"owner":-57466174},{"id":779199,"views":null,"owner":-57466174},{"id":779120,"views":null,"owner":-57466174},{"id":779122,"views":null,"owner":-57466174},{"id":779118,"views":null,"owner":-57466174},{"id":779116,"views":null,"owner":-57466174},{"id":779115,"views":null,"owner":-57466174},{"id":779081,"views":null,"owner":-57466174},{"id":778786,"views":null,"owner":-57466174},{"id":778782,"views":null,"owner":-57466174},{"id":778779,"views":null,"owner":-57466174},{"id":778770,"views":null,"owner":-57466174},{"id":778763,"views":null,"owner":-57466174},{"id":778752,"views":null,"owner":-57466174},{"id":778737,"views":null,"owner":-57466174},{"id":778725,"views":null,"owner":-57466174},{"id":778723,"views":null,"owner":-57466174},{"id":778794,"views":null,"owner":-57466174},{"id":778796,"views":null,"owner":-57466174}]}`)
	res, err := ParseIDViewsResponse(input)
	if err != nil {
		t.Errorf("parse unexpected error: %s", err)
	}
	// if len(res) != 50 {
	// 	t.Errorf("parse res len: %d", len(res))
	// }
	if res[0].Views != 1511 {
		t.Errorf("unexpected views: %d", res[0].Views)
	}
	if res[0].ID != 780206 {
		t.Errorf("unexpected id: %d", res[0].ID)
	}
	if res[0].Owner != -57466174 {
		t.Errorf("unexpected owner: %d", res[0].Owner)
	}
}
