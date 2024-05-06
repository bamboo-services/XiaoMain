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

package rss

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/os/glog"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
)

// GetAllLinkRssInfo 获取所有链接的Rss信息
// 用于获取所有链接的Rss信息
// 如果成功则返回 nil，否则返回错误
// 本接口会根据已有的链接信息对Rss信息进行获取，若获取失败返回失败信息，若成功返回成功信息
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
//
// 返回：
// err: 如果获取Rss信息成功，返回 nil；否则返回错误。
func (s *sRssLogic) GetAllLinkRssInfo(ctx context.Context) (err error) {
	glog.Noticef(ctx, "[LOGIC] 执行 RssLogic:GetAllLinkRssInfo 服务层")
	var getLink *[]entity.XfLinkList
	err = dao.XfLinkList.Ctx(ctx).
		Where(do.XfLinkList{DeletedAt: nil, Status: 1}).
		WhereNotIn("site_rss_url", nil).
		Scan(&getLink)
	if err != nil {
		glog.Errorf(ctx, "获取链接信息失败: %v", err)
		return errors.New("数据库查询失败<链接获取失败>")
	}
	if getLink != nil {
		return nil
	} else {
		return errors.New("没有数据")
	}
}