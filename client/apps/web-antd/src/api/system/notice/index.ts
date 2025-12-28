import type { Notice } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery } from '#/api/helper';
import { pb } from '#/api/request';

const collection = pb.collection<Notice>('notice');

/**
 * 通知公告分页
 * @param params 分页参数
 * @returns 分页结果
 */
export function noticeList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return collection.getList(currentPage, pageSize, {
    filter,
    sort,
  });
}

/**
 * 通知公告详情
 * @param noticeId id
 * @returns 详情
 */
export function noticeInfo(noticeId: ID) {
  return collection.getOne(`${noticeId}`);
}

/**
 * 通知公告新增
 * @param data 参数
 */
export function noticeAdd(data: Partial<Notice>) {
  return collection.create(data);
}

/**
 * 通知公告更新
 * @param data 参数
 */
export function noticeUpdate(data: any) {
  return collection.update(`${data.id}`, data);
}

/**
 * 通知公告删除
 * @param noticeIds ids
 */
export function noticeRemove(noticeIds: IDS) {
  return Promise.all(noticeIds.map((id) => collection.delete(`${id}`)));
}
