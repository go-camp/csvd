package csvd

import "testing"

func TestCanonicalHeaderKey(t *testing.T) {
	tcs := []struct {
		Name string

		Key       string
		ExpectKey string
	}{
		{
			Name: "left space",

			Key:       "   \t\v\ruser id",
			ExpectKey: "user id",
		},
		{
			Name: "middle space",

			Key:       "user\r \v\tid",
			ExpectKey: "user id",
		},
		{
			Name: "right space",

			Key:       "user id\r \v \t",
			ExpectKey: "user id",
		},
		{
			Name: "all spaces",

			Key:       " \r\v \v\r \t\t    ",
			ExpectKey: "",
		},
		{
			Name: "cn",

			Key:       " 用户名 ",
			ExpectKey: "用户名",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			key := CanonicalHeaderKey(tc.Key)
			if key != tc.ExpectKey {
				t.Errorf("CanonicalHeaderKey(%q), expect return %q, got %q", tc.Key, tc.ExpectKey, key)
			}
		})
	}
}

func BenchmarkCanonicalHeaderKey_leftSpaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CanonicalHeaderKey("     \t\t\v user id")
	}
}

func BenchmarkCanonicalHeaderKey_rightSpaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CanonicalHeaderKey("user id     \t\t\v ")
	}
}

func BenchmarkCanonicalHeaderKey_middleSpaces(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CanonicalHeaderKey("user     \t\t\v id")
	}
}

func BenchmarkCanonicalHeaderKey_alreadyCanonical(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CanonicalHeaderKey("user id")
	}
}
