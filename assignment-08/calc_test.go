package main
import "testing"
func TestAddTableDriven(t *testing.T){
	tests := []struct{
		name string
		a, b int
		want int
	}{
		{"both positive", 2,3,5},
		{"positive + zero", 5, 0,5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -2, -4},
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T){
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}


func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
	}{
		{"normal", 10, 2, 5, false},
		{"zero division", 10, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 5, 3, 2},
		{"positive minus zero", 5, 0, 5},
		{"negative minus zero", -5, 0, -5},
		{"negative minus positive", -2, 3, -5},
		{"positive minus negative", -2, -3, 1}, /// i added some prosto tak 
		{"both negative", -5, -2, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}