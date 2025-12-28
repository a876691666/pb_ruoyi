import type { Tenant } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery, commonExport } from '#/api/helper';
import { pb, requestClient } from '#/api/request';

const collection = pb.collection<Tenant>('tenant');

enum Api {
  dictSync = '/system/tenant/syncTenantDict',
  root = '/system/tenant',
  tenantDynamic = '/system/tenant/dynamic',
  tenantDynamicClear = '/system/tenant/dynamic/clear',
  tenantExport = '/system/tenant/export',
  tenantList = '/system/tenant/list',
  tenantStatus = '/system/tenant/changeStatus',
  tenantSyncPackage = '/system/tenant/syncTenantPackage',
}

/**
 * 查询租户分页列表
 * @param params 参数
 * @returns 分页
 */
export function tenantList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    requestKey: Math.random().toString(36).slice(2, 15),
    sort,
  });
}

/**
 * 租户导出
 * @param params 查询参数
 * @returns void
 */
export function tenantExport(params?: PageQuery) {
  return commonExport('tenant', buildingQuery(params), { type: 'collection' });
}

/**
 * 查询租户信息
 * @param id id
 * @returns 租户信息
 */
export function tenantInfo(id: ID) {
  return collection.getOne(`${id}`);
}

/**
 * 新增租户 必须开启加密
 * @param data data
 * @returns void
 */
export function tenantAdd(data: Partial<Tenant>) {
  return collection.create(data);
}

/**
 * 租户更新
 * @param data data
 * @returns void
 */
export function tenantUpdate(data: Partial<Tenant>) {
  return collection.update(`${data.id}`, data);
}

/**
 * 租户状态更新
 * @param data data
 * @returns void
 */
export function tenantStatusChange(data: Partial<Tenant>) {
  return collection
    .update(`${data.id}`, { status: data.status })
    .then(() => undefined);
}

/**
 * 租户删除
 * @param ids ids
 * @returns void
 */
export function tenantRemove(ids: IDS) {
  return Promise.all(ids.map((id) => collection.delete(`${id}`)));
}

/**
 * 动态切换租户
 * @param tenantId 租户ID
 * @returns void
 */
export function tenantDynamicToggle(tenantId: string) {
  return requestClient.get<void>(`${Api.tenantDynamic}/${tenantId}`);
}

/**
 * 清除 动态切换租户
 * @returns void
 */
export function tenantDynamicClear() {
  return requestClient.get<void>(Api.tenantDynamicClear);
}

/**
 * 租户套餐同步
 * @param tenantId 租户id
 * @param packageId 套餐id
 * @returns void
 */
export function tenantSyncPackage(tenantId: string, packageId: string) {
  return requestClient.get<void>(Api.tenantSyncPackage, {
    params: { packageId, tenantId },
    successMessageMode: 'message',
  });
}

/**
 * 同步租户字典
 * @param tenantId 租户ID
 * @returns void
 */
export function dictSyncTenant(tenantId?: string) {
  return requestClient.get<void>(Api.dictSync, {
    params: { tenantId },
    successMessageMode: 'message',
  });
}
