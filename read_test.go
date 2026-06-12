package pdf

import "testing"

func TestOkayV4CryptFilterLength(t *testing.T) {
	// 构造一个最小可用的 encrypt dict，只变化 CF.<name>.Length
	makeEncrypt := func(length interface{}) dict {
		return dict{
			"CF": dict{
				"StdCF": dict{
					"Length": length,
					"CFM":    name("AESV2"),
				},
			},
			"StmF": name("StdCF"),
			"StrF": name("StdCF"),
		}
	}

	tests := []struct {
		name   string
		encrypt dict
		want   bool
	}{
		{
			name:    "Length 16 (byte) 应通过",
			encrypt: makeEncrypt(int64(16)),
			want:    true,
		},
		{
			name:    "Length 128 (bit) 应通过",
			encrypt: makeEncrypt(int64(128)),
			want:    true,
		},
		{
			name:    "Length 256 应拒绝",
			encrypt: makeEncrypt(int64(256)),
			want:    false,
		},
		{
			name:    "Length 缺失应通过",
			encrypt: makeEncrypt(nil),
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := okayV4(tt.encrypt)
			if got != tt.want {
				t.Errorf("okayV4() = %v, want %v", got, tt.want)
			}
		})
	}
}
