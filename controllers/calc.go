package controllers

import (
	"fmt"
	"github.com/fyuan1316/fundpool/api/v1alpha1"
	"github.com/fyuan1316/fundpool/controllers/types"
	"sort"
)

func parse(fundPools []v1alpha1.FundPool) []*types.Pool {
	var P []*types.Pool
	for _, fundPool := range fundPools {
		var pool = &types.Pool{Name: fundPool.Name, Balance: fundPool.Spec.Balance}
		P = append(P, pool)
	}
	return P
}
func getMaxPool(pools *[]*types.Pool) *types.Pool {
	var pMax = (*pools)[0]
	for i := range *pools {
		cur := (*pools)[i]
		if cur.Max() > pMax.Max() {
			pMax = cur
		}
	}
	return pMax
}

var step = int64(1)

func Calc(fetch int64, pools []v1alpha1.FundPool) []*types.Pool {
	wrapPool := parse(pools)
	sort.Slice(wrapPool, func(i, j int) bool {
		return wrapPool[i].Balance > wrapPool[j].Balance
	})
	var accumulator = int64(0)
	var maxPool *types.Pool
	for accumulator < fetch {
		maxPool = getMaxPool(&wrapPool)
		maxPool.StepFetch(step)
		accumulator = accumulator + step
	}
	for _, p := range wrapPool {
		fmt.Println(p.Info())
	}
	return wrapPool
}
