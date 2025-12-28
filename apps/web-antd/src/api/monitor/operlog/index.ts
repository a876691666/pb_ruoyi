import type { OperationLog } from './model';

import type { IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport } from '#/api/helper';
import { pb, requestClient } from '#/api/request';

const collection = pb.collection<OperationLog>('oper_log');

/**
 * 操作日志分页
 * @param params 查询参数
 * @returns 分页结果
 */
export function operLogList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    sort,
  });
}

/**
 * 删除操作日志
 * @param operIds id/ids
 */
export function operLogDelete(operIds: IDS) {
  return Promise.all(operIds.map((id) => collection.delete(`${id}`)));
}

/**
 * 清空全部分页日志
 */
export function operLogClean() {
  return requestClient.deleteWithMsg<void>('/monitor/operlog/clean');
}

/**
 * 导出操作日志
 * @param params 查询参数
 */
export function operLogExport(params?: PageQuery) {
  return commonExport('oper_log', buildingQuery(params), {
    type: 'collection',
  });
}
