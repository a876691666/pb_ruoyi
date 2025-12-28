import type { PBQuery } from './common';

import { $t } from '@vben/locales';

import { message, Modal } from 'ant-design-vue';
import dayjs from 'dayjs';

import { useAuthStore } from '#/store';

import { Ands, requestClient } from './request';

/**
 * @description:  contentType
 */
export const ContentTypeEnum = {
  // form-data  upload
  FORM_DATA: 'multipart/form-data;charset=UTF-8',
  // form-data qs
  FORM_URLENCODED: 'application/x-www-form-urlencoded;charset=UTF-8',
  // json
  JSON: 'application/json;charset=UTF-8',
} as const;

/**
 * 通用下载接口 封装一层
 * @param url 请求地址
 * @param data  请求参数
 * @returns blob二进制
 */
/**
 * 通用导出函数
 * 兼容两种模式:
 * 1. 普通业务导出(old): POST form-url-encoded 到具体url (保持向后兼容)
 * 2. 集合导出(collection): GET /collections/{collection}/export?filter&sort
 * @param target url 或 collection 名称
 * @param data 参数; collection模式下需包含 { filter, sort }
 * @param options 配置对象
 * @param options.type 当值为 'collection' 时走集合导出
 */
export function commonExport(
  target: string,
  data: Record<string, any>,
  options?: { type?: 'api' | 'collection' },
) {
  const type = options?.type ?? 'api';
  if (type === 'collection') {
    return requestClient.get<Blob>(`/collections/${target}/export`, {
      params: data,
      responseType: 'blob',
      isTransformResponse: false,
    });
  }
  return requestClient.post<Blob>(target, data, {
    data,
    headers: { 'Content-Type': ContentTypeEnum.FORM_URLENCODED },
    isTransformResponse: false,
    responseType: 'blob',
  });
}

/**
 * 定义一个401专用异常 用于可能会用到的区分场景?
 */
export class UnauthorizedException extends Error {}

/**
 * logout这种接口都返回401 抛出这个异常
 */
export class ImpossibleReturn401Exception extends Error {}

/**
 * 是否已经处在登出过程中了 一个标志位
 * 主要是防止一个页面会请求多个api 都401 会导致登出执行多次
 */
let isLogoutProcessing = false;
/**
 * 防止 调用logout接口 logout又返回401 然后又走到Logout逻辑死循环
 */
let lockLogoutRequest = false;

/**
 * 登出逻辑 两个地方用到 提取出来
 * @throws UnauthorizedException 抛出特定的异常
 */
export function handleUnauthorizedLogout() {
  const timeoutMsg = $t('http.loginTimeout');
  /**
   * lock 不再请求logout接口
   * 这里已经算异常情况了
   */
  if (lockLogoutRequest) {
    throw new UnauthorizedException(timeoutMsg);
  }
  // 已经在登出过程中 不再执行
  if (isLogoutProcessing) {
    throw new UnauthorizedException(timeoutMsg);
  }
  isLogoutProcessing = true;
  const userStore = useAuthStore();
  userStore
    .logout()
    .catch((error) => {
      /**
       * logout接口返回了401
       * 做Lock处理 且 该标志位不会复位(因为这种场景出现 系统已经算故障了)
       * 因为这已经不符合正常的逻辑了
       */
      if (error instanceof ImpossibleReturn401Exception) {
        lockLogoutRequest = true;
        if (import.meta.env.DEV) {
          Modal.error({
            title: '提示',
            centered: true,
            content:
              '检测到你的logout接口返回了401, 去检查你的后端配置 这已经不符合正常逻辑(该提示不会在非dev环境弹出)',
          });
        }
      }
    })
    .finally(() => {
      message.error(timeoutMsg);
      isLogoutProcessing = false;
    });
  // 不再执行下面逻辑
  throw new UnauthorizedException(timeoutMsg);
}

export function buildingQuery(pageQuery?: PBQuery) {
  if (!pageQuery) {
    return {
      filter: '',
      sort: '',
      currentPage: undefined,
      pageSize: undefined,
    };
  }
  // 构建 filter: 根据 params 和 queryType 生成 PocketBase SQL
  const filters: string[] = [];

  if (pageQuery?.params && pageQuery?.queryType) {
    Object.entries(pageQuery.queryType).forEach(([key, queryType]) => {
      const value = pageQuery.params?.[key];
      if (value === undefined || value === null || value === '') return;

      switch (queryType) {
        case 'AEQ': {
          filters.push(`${key} ?= "${value}"`);
          break;
        }
        case 'AGE': {
          filters.push(`${key} ?>= "${value}"`);
          break;
        }
        case 'AGT': {
          filters.push(`${key} ?> "${value}"`);
          break;
        }
        case 'AIN': {
          if (Array.isArray(value)) {
            filters.push(value.map((v) => `${key} = "${v}"`).join(' || '));
          }
          break;
        }
        case 'ALE': {
          filters.push(`${key} ?>= "${value}"`);
          break;
        }
        case 'ALIKE': {
          filters.push(`${key} ?~ "${value}"`);
          break;
        }
        case 'ALT': {
          filters.push(`${key} ?< "${value}"`);
          break;
        }
        case 'ANE': {
          filters.push(`${key} ?!= "${value}"`);
          break;
        }
        case 'BETWEEN': {
          if (Array.isArray(value) && value.length === 2) {
            filters.push(
              [
                `${key} >= "${dayjs(value[0]).format('YYYY-MM-DD 00:00:00')}"`,
                `${key} <= "${dayjs(value[1]).format('YYYY-MM-DD 23:59:59')}"`,
              ].join('&&'),
            );
          }
          break;
        }
        case 'EQ': {
          filters.push(`${key} = "${value}"`);
          break;
        }
        case 'GE': {
          filters.push(`${key} >= "${value}"`);
          break;
        }
        case 'GT': {
          filters.push(`${key} > "${value}"`);
          break;
        }
        case 'IN': {
          if (Array.isArray(value)) {
            filters.push(value.map((v) => `${key} = "${v}"`).join(' || '));
          }
          break;
        }
        case 'LE': {
          filters.push(`${key} <= "${value}"`);
          break;
        }
        case 'LIKE': {
          filters.push(`${key} ~ "${value}"`);
          break;
        }
        case 'LT': {
          filters.push(`${key} < "${value}"`);
          break;
        }
        case 'NE': {
          filters.push(`${key} != "${value}"`);
          break;
        }
      }
    });
  }

  const filter = Ands(filters);

  // 构建 sort: 根据 sorts 数组生成排序字符串
  let sort = '';
  if (pageQuery?.sorts && pageQuery.sorts.length > 0) {
    sort = pageQuery.sorts
      .sort((a, b) => (a.sortTime || 0) - (b.sortTime || 0))
      .map((s) => {
        const prefix = s.order?.toLowerCase() === 'desc' ? '-' : '+';
        return `${prefix}${s.field}`;
      })
      .join(',');
  }

  return {
    filter,
    sort,
    currentPage: pageQuery?.page?.currentPage,
    pageSize: pageQuery?.page?.pageSize,
  };
}
