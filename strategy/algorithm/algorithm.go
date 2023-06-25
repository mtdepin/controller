package algorithm

import (
	"controller/strategy/dict"
)

const (
	WEIGHT        = 1.5 //1.5,  1:供测试版本使用
	MAX_REPLICATE = 9
)

//预估算法
type Estimate struct {
}

//
/*fid   orderid  region  minRep maxRep  realMinRep realMaxRep  Weight  Expire CreateTime  status                   used
1111    order1    cd      1      2       1           2           1                        init success fail
1111    order2    cd      3      4       4           6           1.5
1111    order3    cd      5      6       7           9           1.5
1111    order4    cd      2      3       2           3           1.5                       status 如果备份成功了，相同备份次数的订单不必重复执行。
1111    order5    cd      6      7       6           7           1.5                       如果此集群备份失败，删除此记录，备份迁移到其它集群(能满足备份需求的集群).
*/
// 对于备份成功的数据，动态调整备份数据，如果文件很热(时序数据,最近一周下载频次太高), 提升文件当前备份数，如果数据下载频次低,降低备份数。 f1*w1 + f2*w2 + ... + fn*wn = weight, 调整后不进行计费, 调整后，变更文件对应订单的备份详情。
//动态负载均衡: 如果当前备份集群 cpu负载(时序数据, 最近一个小时平均负载)，如果负载极高，迁移到其它集群进行备份， 如果负载中高，本集群进行部分备份, 另外一部分备份到其它集群。
//f1*w1 + f2*w2 + ... + fn*wn = weight
//weight * max_replication

//1.计算当前订单实际最小，最大备份数.
//2.计算当前订单应该进行的备份数， 遍历region 中最大的一组realMinRep， realMaxRep 作为此订单备份策略的备份数 4, 6; 6,9
func (p *Estimate) CalculateRep(orderId string, fidReps map[string]*dict.Rep, rep *dict.RepInfo) (initMinRep, initMaxRep, realMinRep, realMaxRep int) {
	max := true
	initMinRep = rep.MinRep
	initMaxRep = rep.MaxRep

	initRealMinRep := rep.MinRep
	initRealMaxRep := rep.MaxRep

	bExist := false

	for key, fidRep := range fidReps { //-条记录
		fidRep.Used = 0     //重置used = 0
		if key == orderId { //跳过自己， 存在，用自己第一次初始化值, 不存在，默认就是初始值,  minRep*1.5, maxRep*1.5
			initMinRep = fidRep.MinRep
			initMaxRep = fidRep.MaxRep
			initRealMinRep = fidRep.RealMinRep
			initRealMaxRep = fidRep.RealMaxRep
			bExist = true
			continue
		}

		if rep.MinRep <= fidRep.RealMinRep { //查找系统最大备份数，作为当前倍数数。
			rep.MinRep = fidRep.RealMinRep
			rep.MaxRep = fidRep.RealMaxRep
			max = false
		}
	}

	if max && (len(fidReps) > 1 || len(fidReps) == 1 && bExist == false) { //只有一条记录，且是最大的，不做处理，如果是多条记录，是最大的，调整最大备份数。
		realMinRep = int(float64(initMinRep) * WEIGHT)
		realMaxRep = int(float64(initMaxRep) * WEIGHT)
		if realMinRep > MAX_REPLICATE || realMaxRep > MAX_REPLICATE { //最大备份限制
			realMinRep = MAX_REPLICATE
			realMaxRep = MAX_REPLICATE
		}

		rep.MinRep = realMinRep
		rep.MaxRep = realMaxRep
	} else { //实际备份数维持不变,
		realMinRep = initRealMinRep
		realMaxRep = initRealMaxRep
	}

	return
}

//
/*fid   orderid  region  minRep maxRep  realMinRep realMaxRep  Weight  Expire CreateTime  status                   used
1111    order1    cd      1      2       1           2           1                        init success fail
1111    order2    cd      3      4       4           6           1.5
1111    order3    cd      5      6       7           9           1.5
1111    order4    cd      2      3       2           3           1.5                       status 如果备份成功了，相同备份次数的订单不必重复执行。
1111    order5    cd      6      7       6           7           1.5                       如果此集群备份失败，删除此记录，备份迁移到其它集群(能满足备份需求的集群).
*/
func (p *Estimate) CalculateDeleteRep(orderId string, fidReps map[string]*dict.Rep, rep *dict.RepInfo) error {
	realMinRep := rep.MinRep
	secondMinRep := rep.MinRep
	secondMaxRep := rep.MaxRep
	max := true

	used := 0
	for _, rep := range fidReps {
		if rep.Used == 1 {
			used = rep.Used
			break
		}
	}

	if _, ok := fidReps[orderId]; !ok { //不存在。
		rep.MinRep = 0
		rep.MaxRep = 0
		rep.Status = 0
		if len(fidReps) > 0 || used == 1 { //如果还有其它订单，则直接返回删除成功， 否则没有其它订单了，返回可以删除。
			rep.Status = dict.TASK_DEL_SUC //不做处理, 只有一条记录，记录数大于0，则不进行删除
		}
		return nil
	}

	if len(fidReps) == 1 { //1条记录 and used == 0
		rep.MinRep = 0
		rep.MaxRep = 0
		rep.Status = 0
		if used > 0 { //被占用了.
			rep.Status = dict.TASK_DEL_SUC //删除成功
		} else { //
			delete(fidReps, orderId) //delete region, 此region 没有订单，delete region
		}

		return nil
	}

	//处理多条记录.
	delete(fidReps, orderId) //过滤掉本身记录。
	//余下的记录选择次大的。
	idx := 0
	for _, fidRep := range fidReps {
		if realMinRep > fidRep.RealMinRep { //delete map: 0, 1, 2, 3: sort
			idx++
			if idx == 1 {
				secondMinRep = fidRep.RealMinRep
				secondMaxRep = fidRep.RealMaxRep
			} else { //比较
				if secondMinRep < fidRep.RealMinRep {
					secondMinRep = fidRep.RealMinRep
					secondMaxRep = fidRep.RealMaxRep
				}
			}
		} else {
			max = false
			break
		}
	}

	if max { //当前实际删除备份数是最大的，返回次大的一组.
		rep.MinRep = secondMinRep
		rep.MaxRep = secondMaxRep
		rep.Status = 0
	} else { //次大的，不用删除.
		rep.Status = dict.TASK_DEL_SUC
	}

	//map 没记录了，delete fid
	return nil
}
