import type { AxiosRequestConfig } from '@vben/request';

import type { OssFile } from './model';

import type { ID, IDS, PageQuery } from '#/api/common';

import { buildingQuery } from '#/api/helper';
import { pb, requestClient } from '#/api/request';

const ossCollection = pb.collection<OssFile>('oss');
const usersCollection = pb.collection('users');

enum Api {
  ossDownload = '/resource/oss/download',
  ossInfo = '/resource/oss/listByIds',
  ossList = '/resource/oss/list',
  root = '/resource/oss',
}

/**
 * 文件list
 * @param params 参数
 * @returns 分页
 */
export function ossList(params?: PageQuery) {
  const { filter, sort, pageSize, currentPage } = buildingQuery(params);
  return ossCollection.getList(currentPage, pageSize, {
    filter,
    sort,
  });
}

/**
 * 查询文件信息 返回为数组
 * @param ossIds id数组
 * @returns 信息数组
 */
export function ossInfo(ossIds: ID | IDS) {
  return requestClient.get<OssFile[]>(`${Api.ossInfo}/${ossIds}`);
}

/**
 * @deprecated 使用apps/web-antd/src/api/core/upload.ts uploadApi方法
 * @param file 文件
 * @returns void
 */
export function ossUpload(file: Blob | File) {
  const formData = new FormData();
  formData.append('file', file);
  return ossCollection.create(formData);
}

/**
 * 下载文件  返回为二进制
 * @param ossId ossId
 * @param onDownloadProgress 下载进度(可选)
 * @returns blob
 */
export function ossDownload(
  ossId: ID,
  onDownloadProgress?: AxiosRequestConfig['onDownloadProgress'],
) {
  return requestClient.get<Blob>(`${Api.ossDownload}/${ossId}`, {
    responseType: 'blob',
    timeout: 30 * 1000,
    isTransformResponse: false,
    onDownloadProgress,
  });
}

/**
 * 在使用浏览器原生下载前检测是否登录
 * 这里的方案为请求一次接口 如果登录超时会走到response的401逻辑
 * 如果没有listByIds的权限 也不会弹出无权限提示
 * 仅仅是为了检测token是否有效使用
 *
 * @returns void
 */
export function checkLoginBeforeDownload() {
  return usersCollection.authRefresh();
}

/**
 * 删除文件
 * @param ossIds id数组
 * @returns void
 */
export function ossRemove(ossIds: IDS) {
  return Promise.all(ossIds.map((id) => ossCollection.delete(`${id}`)));
}
