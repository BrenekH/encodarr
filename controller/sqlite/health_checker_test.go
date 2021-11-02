package sqlite_test

import (
	"reflect"
	"testing"

	"github.com/BrenekH/encodarr/controller"
	"github.com/BrenekH/encodarr/controller/sqlite"
)

// TODO: Write tests

func Test_dbDispatchedJob_ToController(t *testing.T) {
	tests := []struct {
		name  string
		input sqlite.DispatchedJob
		want  controller.DispatchedJob
	}{
		{
			name: "ok",
			input: sqlite.DispatchedJob{
				UUID:   "98765",
				Runner: "runner",
				Job:    []byte(`{"uuid": "12345"}`),
				Status: []byte(`{"stage": "one"}`),
			},
			want: controller.DispatchedJob{
				UUID:   controller.UUID("98765"),
				Runner: "runner",
				Job: controller.Job{
					UUID: controller.UUID("12345"),
				},
				Status: controller.JobStatus{
					Stage: "one",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.input.ToController()
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DispatchedJob.ToController() = %v, want %v", got, tt.want)
			}
		})
	}
}
