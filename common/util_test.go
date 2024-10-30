package common

import "testing"

func TestValidateEmail(t *testing.T) {
	// 使用表测试
	tests := []struct {
		email string
		want  bool
	}{
		{
			email: "t@gmail.com",
			want:  true,
		},
		{
			email: "@qq.com",
			want:  true,
		},
		{
			email: "1@gmail",
			want:  true,
		},
	}

	for _, test := range tests {
		if result := ValidateEmail(test.email); result != test.want {
			t.Errorf("Validate Email %s: got result %v, want %v", test.email, result, test.want)
		}
	}
}
