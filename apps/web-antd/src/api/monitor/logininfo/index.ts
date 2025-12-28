import type { IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport } from '#/api/helper';
import { pb, requestClient } from '#/api/request';

const collection = pb.collection('logininfor');

/**
 * 登录日志列表
 * @param params 查询参数
 * @returns list[]
 */
export function loginInfoList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    sort,
  });
}

/**
 * 导出登录日志
 * @param params 查询参数
 * @returns excel
 */
export function loginInfoExport(params?: PageQuery) {
  return commonExport('logininfor', buildingQuery(params), {
    type: 'collection',
  });
}

/**
 * 移除登录日志
 * @param infoIds 登录日志id数组
 * @returns void
 */
export function loginInfoRemove(infoIds: IDS) {
  return Promise.all(infoIds.map((id) => collection.delete(`${id}`)));
}

/**
 * 账号解锁
 * @param username 用户名(账号)
 * @returns void
 */
export function userUnlock(username: string) {
  return requestClient.get<void>(`/monitor/logininfor/unlock/${username}`, {
    successMessageMode: 'message',
  });
}

/**
 * 清空全部登录日志
 * @returns void
 */
export function loginInfoClean() {
  return requestClient.deleteWithMsg<void>('/monitor/logininfor/clean');
}
