import type { GrantType } from '@vben/common-ui';
import type { UserInfo } from '@vben/types';

import { useAppConfig } from '@vben/hooks';

import { pb, requestClient } from '#/api/request';

const { sseEnable } = useAppConfig(import.meta.env, import.meta.env.PROD);

export namespace AuthApi {
  /**
   * @description: 所有登录类型都需要用到的
   * @param clientId 客户端ID 这里为必填项 但是在loginApi内部处理了 所以为可选
   * @param grantType 授权/登录类型
   * @param tenantId 租户id
   */
  export interface BaseLoginParams {
    clientId?: string;
    grantType: GrantType;
    tenantId: string;
  }

  /**
   * @description: oauth登录需要用到的参数
   * @param socialCode 第三方参数
   * @param socialState 第三方参数
   * @param source 与后端的 justauth.type.xxx的回调地址的source对应
   */
  export interface OAuthLoginParams extends BaseLoginParams {
    socialCode: string;
    socialState: string;
    source: string;
  }

  /**
   * @description: 验证码登录需要用到的参数
   * @param code 验证码 可选(未开启验证码情况)
   * @param uuid 验证码ID 可选(未开启验证码情况)
   * @param username 用户名
   * @param password 密码
   */
  export interface SimpleLoginParams extends BaseLoginParams {
    code?: string;
    uuid?: string;
    username: string;
    password: string;
  }

  export type LoginParams = OAuthLoginParams | SimpleLoginParams;

  // /** 登录接口参数 */
  // export interface LoginParams {
  //   code?: string;
  //   grantType: string;
  //   password: string;
  //   tenantId: string;
  //   username: string;
  //   uuid?: string;
  // }

  /** 登录接口返回值 */
  export interface LoginResult {
    access_token: string;
    client_id: string;
    expire_in: number;
  }

  export interface RefreshTokenResult {
    data: string;
    status: number;
  }
}

/**
 * 登录
 */
export async function loginApi(data: AuthApi.LoginParams) {
  if ('username' in data) {
    const { username, password, ...arg } = data;
    return pb
      .collection('users')
      .authWithPassword<UserInfo>(username, password, {
        body: {
          identity: username,
          password,
          ...arg,
        },
      });
  }
}

/**
 * 用户登出
 * @returns void
 */
export function doLogout() {
  return pb.authStore.clear();
}

/**
 * 关闭sse连接
 * @returns void
 */
export function seeConnectionClose() {
  /**
   * 未开启sse 不需要处理
   */
  if (!sseEnable) {
    return;
  }
  return requestClient.get<void>('/resource/sse/close');
}

/**
 * @param companyName 租户/公司名称
 * @param domain 绑定域名(不带http(s)://) 可选
 * @param tenantId 租户id
 */
export interface TenantOption {
  company_name: string;
  domain?: string;
  id: string;
}

/**
 * @param tenantEnabled 是否启用租户
 * @param list 租户列表
 */
export interface TenantResp {
  tenantEnabled: boolean;
  list: TenantOption[];
}

/**
 * 获取租户列表 下拉框使用
 */
export async function tenantList() {
  const tenantEnabled = await pb
    .collection('global_config')
    .getFirstListItem("key='tenantEnabled'")
    .then((item) => item.value as boolean);

  const list = await pb
    .collection('tenant')
    .getFullList<TenantOption>({ fields: 'company_name,domain,id' });

  return {
    tenantEnabled,
    list,
  };
}

/**
 * vben的 先不删除
 * @returns string[]
 */
export async function getAccessCodesApi() {
  return requestClient.get<string[]>('/auth/codes');
}

/**
 * 绑定第三方账号
 * @param source 绑定的来源
 * @returns 跳转url
 */
export function authBinding(source: string, tenantId: string) {
  return requestClient.get<string>(`/auth/binding/${source}`, {
    params: {
      domain: window.location.host,
      tenantId,
    },
  });
}

/**
 * 取消绑定
 * @param id id
 */
export function authUnbinding(id: string) {
  return requestClient.deleteWithMsg<void>(`/auth/unlock/${id}`);
}

/**
 * oauth授权回调
 * @param data oauth授权
 * @returns void
 */
export function authCallback(data: AuthApi.OAuthLoginParams) {
  return requestClient.post<void>('/auth/social/callback', data);
}
