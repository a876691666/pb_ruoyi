import type { TenantPackage } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport } from '#/api/helper';
import { pb } from '#/api/request';

const collection = pb.collection<TenantPackage>('tenant_package');

// enum Api {
//   packageChangeStatus = '/system/tenant/package/changeStatus',
//   packageList = '/system/tenant/package/list',
//   packageSelectList = '/system/tenant/package/selectList',
//   root = '/system/tenant/package',
// }

/**
 * 租户套餐分页列表
 * @param params 请求参数
 * @returns 分页列表
 */
export function packageList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    requestKey: Math.random().toString(36).slice(2, 15),
    sort,
  });
}

/**
 * 租户套餐下拉框
 * @returns 下拉框
 */
export function packageSelectList() {
  return collection.getFullList({
    requestKey: Math.random().toString(36).slice(2, 15),
  });
}

/**
 * 租户套餐导出
 * @param params 查询参数
 * @returns blob
 */
export function packageExport(params?: PageQuery) {
  return commonExport('tenant_package', buildingQuery(params), {
    type: 'collection',
  });
}

/**
 * 租户套餐信息
 * @param id id
 * @returns 信息
 */
export function packageInfo(id: ID) {
  return collection.getOne(`${id}`);
}

/**
 * 租户套餐新增
 * @param data data
 * @returns void
 */
export function packageAdd(data: Partial<TenantPackage>) {
  return collection.create(data);
}

/**
 * 租户套餐更新
 * @param data data
 * @returns void
 */
export function packageUpdate(data: Partial<TenantPackage>) {
  return collection.update(`${data.id}`, data);
}

/**
 * 租户套餐状态变更
 * @param data data
 * @returns void
 */
export function packageChangeStatus(data: Partial<TenantPackage>) {
  return collection
    .update(`${data.id}`, { status: data.status })
    .then(() => undefined);
}

/**
 * 租户套餐移除
 * @param ids ids
 * @returns void
 */
export function packageRemove(ids: IDS) {
  return Promise.all(ids.map((id) => collection.delete(`${id}`)));
}
