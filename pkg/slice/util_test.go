package slice

import "testing"

func TestReverse(t *testing.T) {
	type args struct {
		a []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Reverse the given slice",
			args: args{a: []int{5, 4, 3, 3, 2, 1}},
			want: []int{1, 2, 3, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reverse(tt.args.a)

			for index, value := range tt.args.a {
				if tt.want[index] != value {
					t.Errorf("Reverse() got = %d, want %d", value, tt.want[index])
				}
			}
		})
	}
}
