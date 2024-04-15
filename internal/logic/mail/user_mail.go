/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

package mail

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
	"xiaoMain/internal/service"
	"xiaoMain/utility"
)

type sMailUserLogic struct {
}

func init() {
	service.RegisterMailUserLogic(New())
}

func New() *sMailUserLogic {
	return &sMailUserLogic{}
}

// VerificationCodeHasCorrect
// 验证验证码是否正确，若验证码正确将会返回 true，否则返回 false；
// 若返回错误的内容将会返回具体的错误原因，不会抛出 Error
func (s *sMailUserLogic) VerificationCodeHasCorrect(
	ctx context.Context,
	email string,
	code string,
	scenes string,
) (isCorrect bool, info string) {
	glog.Info(ctx, "[LOGIC] 执行 MailUserLogic:VerificationCodeHasCorrect 服务层")
	// 获取邮箱以及验证码
	var getCode []entity.XfVerificationCode
	if dao.XfVerificationCode.Ctx(ctx).Where(do.XfVerificationCode{
		Type:    true,
		Contact: email,
		Scenes:  scenes,
	}).Scan(&getCode) != nil {
		glog.Info(ctx, "[LOGIC] 用户的验证码不存在")
		return false, "验证码不存在"
	}
	// 对获取的验证码进行查询
	for _, verificationCode := range getCode {
		// 检查是否有匹配项
		if verificationCode.Code == code {
			// 检查是否过期
			if verificationCode.ExpiredAt.After(gtime.Now()) {
				return true, "验证码正确"
			}
		}
	}
	// 验证码已过期
	glog.Info(ctx, "[LOGIC] 用户的验证码已过期")
	return false, "验证码已过期"
}

// SendEmailVerificationCode
// 根据输入的场景进行邮箱的发送，需要保证场景的合法性，场景的合法性参考 consts.Scenes 的参考值
// 若邮件发送的过程中出现错误将会终止发件并且返回 error 信息，发件成功返回 nil
func (s *sMailUserLogic) SendEmailVerificationCode(ctx context.Context, mail string, scenes string) (err error) {
	glog.Info(ctx, "[LOGIC] 执行 MailUserLogic:SendEmailVerificationCode 服务层")
	// 场景检查
	if !utility.CheckScenesInScope(scenes) {
		glog.Warningf(ctx, "[LOGIC] 场景内容不正确，输入是 %s 场景", scenes)
		return errors.New("场景内容不正确")
	}
	// 验证码存入数据库
	err = dao.XfVerificationCode.Ctx(ctx).Transaction(ctx, func(_ context.Context, tx gdb.TX) error {
		// 创建验证码并发送
		getCode, err := s.sendVerifyCodeMail(ctx, mail)
		if err != nil {
			glog.Warning(ctx, "[LOGIC] 发送邮件时发生错误，不进行数据库插入")
			return err
		}
		// 存入验证码
		_, err = tx.Insert(dao.XfVerificationCode.Table(), do.XfVerificationCode{
			Type:      true,
			Contact:   mail,
			Code:      getCode,
			ExpiredAt: gtime.NewFromTimeStamp(gtime.Timestamp() + 300),
		})
		if err != nil {
			glog.Warningf(ctx, "[LOGIC] 发送邮件时发生错误，回滚 %s 数据表中发送给 %s 的数据",
				dao.XfVerificationCode.Table(),
				mail,
			)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 发送成功
	return nil
}