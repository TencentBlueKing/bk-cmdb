package command

import (
	"configcenter/src/source_controller/api/metadata"
	"reflect"
	"testing"
)

func Test_getTopo(t *testing.T) {
	type args struct {
		root  string
		assts []metadata.ObjectAsst
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			want: []string{"biz", "yituo", "set", "module"},
			name: "",
			args: args{
				root: "biz",
				assts: []metadata.ObjectAsst{
					{
						ObjectID:    "set",
						ObjectAttID: "bk_childid",
						OwnerID:     "0",
						AsstObjID:   "yituo",
					},

					{
						ObjectID:    "module",
						ObjectAttID: "bk_childid",
						OwnerID:     "0",
						AsstObjID:   "set",
					},

					{
						ObjectID:    "host",
						ObjectAttID: "bk_childid",
						OwnerID:     "0",
						AsstObjID:   "module",
					},

					{
						ObjectID:    "host",
						ObjectAttID: "bk_cloud_id",
						OwnerID:     "0",
						AsstObjID:   "plat",
					},

					{
						ObjectID:    "yituo",
						ObjectAttID: "bk_childid",
						OwnerID:     "0",
						AsstObjID:   "biz",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTopo(tt.args.root, tt.args.assts)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTopo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTopo() = %v, want %v", got, tt.want)
			}
		})
	}
}
