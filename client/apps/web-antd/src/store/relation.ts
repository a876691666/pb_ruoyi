/* eslint-disable @typescript-eslint/no-non-null-assertion */
import type { RecordModel } from 'pocketbase';

import { reactive } from 'vue';

import { defineStore } from 'pinia';

type RelationData = RecordModel;

/**
 * antd使用 select和radio通用
 * 本质上是对RelationData的拓展
 */
export interface RelationOption extends RelationData {
  disabled?: boolean;
  label: string;
  value: number | string;
}

/**
 * 将字典数据转为Options
 * @param data 字典数据
 * @param formatNumber 是否需要将value格式化为number类型
 * @returns options
 */
export function relationToOptions(
  data: RelationData[],
  formatNumber = false,
  displayField?: string,
  valueField?: string,
): RelationOption[] {
  return data.map((item) => {
    const label = displayField ? (item as any)[displayField] : item.dict_label;
    const valueRaw = valueField ? (item as any)[valueField] : item.dict_value;
    return {
      ...item,
      label,
      value: formatNumber ? Number(valueRaw) : valueRaw,
    };
  });
}

export const useRelationStore = defineStore('relation', () => {
  /**
   * select radio checkbox等使用 只能为固定格式{label, value}
   */
  const relationOptionsMap = reactive(new Map<string, RelationData[]>());
  // 记录每个 key 在写入时是否要求将 value 格式化为 number
  /**
   * 添加一个字典请求状态的缓存
   *
   * 主要解决多次请求重复api的问题(不能用abortController 会导致除了第一个其他的获取的全为空)
   * 比如在一个页面 index表单 modal drawer总共会请求三次 但是获取的都是一样的数据
   * 相当于加锁 保证只有第一次请求的结果能拿到
   */
  const relationRequestCache = reactive(
    new Map<string, Promise<RelationData[] | void>>(),
  );

  function makeKey(
    dictName: string,
    displayField?: string,
    valueField?: string,
  ) {
    return `${dictName}::${displayField || ''}::${valueField || ''}`;
  }

  function getRelationOptions(
    dictName: string,
    displayField?: string,
    valueField?: string,
    formatNumber = false,
  ): RelationOption[] {
    if (!dictName) return [];
    const key = makeKey(dictName, displayField, valueField);
    // 没有key 添加一个空数组
    if (!relationOptionsMap.has(key)) {
      relationOptionsMap.set(key, []);
    }
    // 这里拿到的是原始数据，根据写入时记录的 formatNumber 统一转换后返回
    const raw = relationOptionsMap.get(key)!;
    return relationToOptions(raw, formatNumber, displayField, valueField);
  }

  function resetRelationCache() {
    relationRequestCache.clear();
    relationOptionsMap.clear();
    /**
     * 不需要清空dictRequestCache 每次请求成功/失败都清空key
     */
  }

  /**
   * 核心逻辑
   *
   * 不能直接粗暴使用set 会导致之前return的空数组跟现在的数组指向不是同一个地址  数据也就为空了
   *
   * 判断是否已经存在key 并且数组长度为0 说明该次要处理的数据是return的空数组 直接push(不修改指向)
   * 否则 直接set
   *
   */
  function setRelationInfo(
    dictName: string,
    tableList: RelationData[],
    displayField?: string,
    valueField?: string,
  ) {
    const key = makeKey(dictName, displayField, valueField);
    if (
      relationOptionsMap.has(key) &&
      relationOptionsMap.get(key)?.length === 0
    ) {
      // 保持原引用，不改变地址（兼容外部持有引用的场景）
      relationOptionsMap.get(key)?.push(...tableList);
    } else {
      // 直接存原始数据，不在此处调用 relationToOptions
      relationOptionsMap.set(key, tableList);
    }
  }

  function $reset() {
    /**
     * doNothing
     */
  }

  return {
    $reset,
    relationOptionsMap,
    relationRequestCache,
    getRelationOptions,
    resetRelationCache,
    setRelationInfo,
  };
});
