package controllers

import (
	"github.com/fyuan1316/fundpool/api/v1alpha1"
	"github.com/fyuan1316/fundpool/controllers/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

var pools = []v1alpha1.FundPool{
	{
		ObjectMeta: metav1.ObjectMeta{Name: "fundpool-1"},
		Spec:       v1alpha1.FundPoolSpec{Balance: 120},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "fundpool-2"},
		Spec:       v1alpha1.FundPoolSpec{Balance: 100},
	},
}

func TestCalc(t *testing.T) {
	type args struct {
		fetch int64
		pools []v1alpha1.FundPool
	}
	tests := []struct {
		name string
		args args
		want []*types.Pool
	}{
		{
			name: "test-1",
			args: args{
				fetch: 40,
				pools: pools,
			},
			want: []*types.Pool{
				{"fundpool-1", 120, 30},
				{"fundpool-2", 100, 10},
			},
		},
		{
			name: "test-2",
			args: args{
				fetch: 300,
				pools: pools,
			},
			//want: wantedPools,
			want: []*types.Pool{
				{"fundpool-1", 120, 20},
				{"fundpool-2", 100, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Calc(tt.args.fetch, tt.args.pools)
			if len(got) != len(tt.want) {
				t.Errorf("get.size = %v, want %v", len(got), len(tt.want))
			} else {
				for i := range got {
					if !reflect.DeepEqual(got[i], tt.want[i]) {
						t.Errorf("getMaxPool() = %v, want %v", got[i], tt.want[i])
					}
				}
			}
		})
	}
}
