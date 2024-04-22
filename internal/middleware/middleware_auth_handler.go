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

package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
	"xiaoMain/utility"
	"xiaoMain/utility/result"
)

// MiddleAuthHandler 是用于处理用户授权的中间件。
// 它检查用户的授权信息是否有效。
//
// 参数:
// r: 请求的上下文，用于管理请求的信息。
//
// 返回:
// 无
func MiddleAuthHandler(r *ghttp.Request) {
	getAuthorize, err := utility.GetAuthorizationFromHeader(r)
	if err != nil {
		glog.Warning(r.Context(), "[MIDDLE] 用户授权异常|无授权头|用户未登录")
		result.NotLoggedIn.SetErrorMessage(err.Error()).Response(r)
		return
	}
	getUUID, err := utility.GetUUIDFromHeader(r)
	if err != nil {
		glog.Warning(r.Context(), "[MIDDLE] 用户授权异常|无UUID标记|用户未登录")
		result.NotLoggedIn.SetErrorMessage(err.Error()).Response(r)
		return
	}
	getVerify, err := utility.GetVerifyFromHeader(r)
	if err != nil {
		glog.Warning(r.Context(), "[MIDDLE] 用户授权异常|无校验标记|用户未登录")
		result.NotLoggedIn.SetErrorMessage(err.Error()).Response(r)
		return
	}
	// 数据库检查
	if getAuthorize != nil && getUUID != nil && getVerify != nil {
		// 数据库检查
		var tokenInfo *entity.XfToken
		err = dao.XfToken.Ctx(r.Context()).Where(do.XfToken{
			UserToken:    getAuthorize,
			UserUuid:     getUUID,
			Verification: getVerify,
		}).Limit(1).OrderDesc("expired_at").Scan(&tokenInfo)
		if err != nil {
			glog.Error(r.Context(), "[MIDDLE] 数据库查询错误")
			result.ServerInternalError.SetErrorMessage("数据库查询错误").Response(r)
			return
		}
		// 对数据库进行有效性检查
		if tokenInfo != nil {
			if gtime.Now().Before(tokenInfo.ExpiredAt) {
				glog.Infof(r.Context(), "[MIDDLE] 用户授权有效|用户UUID[%s]", tokenInfo.UserUuid)
				r.Middleware.Next()
			} else {
				glog.Warning(r.Context(), "[MIDDLE] 用户授权异常|授权已过期|用户未登录")
				result.NotLoggedIn.SetErrorMessage("授权已过期").Response(r)
				// 删除数据库中的授权信息
				_, _ = dao.XfToken.Ctx(r.Context()).Where(do.XfToken{Id: tokenInfo.Id}).Delete()
			}
		} else {
			glog.Warning(r.Context(), "[MIDDLE] 用户授权异常|授权不存在|用户未登录")
		}
	} else {
		glog.Error(r.Context(), "[MIDDLE] 服务器策略错误 X00001")
		result.ServerInternalError.Response(r)
	}
}