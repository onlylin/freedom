package aggregate

import (
	"errors"

	"github.com/8treenet/freedom/example/fshop/adapter/repository"
	"github.com/8treenet/freedom/example/fshop/domain/entity"
	"github.com/8treenet/freedom/infra/transaction"
)

// 订单发货聚合根
type DeliveryCmd struct {
	entity.Order
	adminEntity  *entity.Admin
	orderRepo    repository.OrderRepo
	deliveryRepo repository.DeliveryRepo
	tx           transaction.Transaction
}

// Run .
func (cmd *DeliveryCmd) Run(trackingNumber string) error {
	//调用订单父类 判断是否支付
	if !cmd.IsPay() {
		return errors.New("该订单未支付")
	}

	deliveryEntity, err := cmd.deliveryRepo.New()
	if err != nil {
		return err
	}

	//设置发货数据
	deliveryEntity.SetOrderNo(cmd.OrderNo)
	deliveryEntity.SetAdminId(cmd.adminEntity.Id)
	deliveryEntity.SetTrackingNumber(trackingNumber)

	//调用订单父类发货
	cmd.Shipment()

	return cmd.tx.Execute(func() error {
		if e := cmd.orderRepo.Save(&cmd.Order); e != nil {
			return e
		}

		return cmd.deliveryRepo.Save(deliveryEntity)
	})
}
