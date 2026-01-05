import type { ShallowRef } from 'vue';

import type { RelationOption } from '#/store/relation';

import { shallowRef, watchEffect } from 'vue';

import { UnauthorizedException } from '#/api/helper';
import { pb } from '#/api/request';
import { useRelationStore } from '#/store/relation';

/**
 * 一般是Select, Radio, Checkbox等组件使用
 * @warning 注意内部为异步实现 所以不要写这种`getDictOptions()[0]`的代码 会获取不到
 * @warning 需要保持`formatNumber`统一 在所有调用地方需要一致 不能出现A处为string B处为number
 *
 * @param tableName 表名称
 * @param formatNumber 是否格式化字典value为number类型
 * @returns 返回 shallowRef，且相同 tableName 返回同一个 shallowRef
 */
// 缓存每个 tableName+fields 对应的 shallowRef，确保同名返回同一引用
const relationOptionsRefMap = new Map<string, ShallowRef<RelationOption[]>>();

// 缓存每个 id 校验的 shallowRef
const relationValidationRefMap = new Map<string, ShallowRef<boolean>>();

/**
 * 校验 relation 缓存中是否存在指定 id
 * @param tableName 表名称
 * @param id 要校验的 id
 * @param valueField value 字段名，默认为 'id'
 * @returns 返回 shallowRef<boolean>，当 relation 数据更新时自动更新
 */
export function validateRelationId(
  tableName: string,
  id: number | string,
  valueField: string = 'id',
): ShallowRef<boolean> {
  const key = `${tableName}::${valueField}::${id}`;
  let validationRef = relationValidationRefMap.get(key);

  if (!validationRef) {
    validationRef = shallowRef(false);
    relationValidationRefMap.set(key, validationRef);

    const optionsRef = getRelationOptions(tableName, 'name', valueField);
    watchEffect(() => {
      if (!validationRef) return;
      validationRef.value = optionsRef.value.some((opt) => opt.value === id);
    });
  }

  return validationRef;
}

export function getRelationOptions(
  tableName: string,
  displayField: string = 'name',
  valueField: string = 'id',
  formatNumber = false,
): ShallowRef<RelationOption[]> {
  const {
    relationRequestCache,
    setRelationInfo,
    getRelationOptions: getStoreRelationOptions,
  } = useRelationStore();
  const dataList = getStoreRelationOptions(tableName, displayField, valueField);
  const key = `${tableName}::${displayField || ''}::${valueField || ''}`;
  // 初始化或复用 shallowRef
  let ref = relationOptionsRefMap.get(key);
  if (!ref) {
    const created = shallowRef<RelationOption[]>(dataList);
    relationOptionsRefMap.set(key, created);
    ref = created;
  }

  // 检查请求状态缓存
  if (dataList.length === 0 && !relationRequestCache.has(key)) {
    relationRequestCache.set(
      key,
      pb
        .collection(tableName)
        .getFullList()
        .then((resp) => {
          // 缓存到store 这样就不用重复获取了
          // 内部处理了push的逻辑 这里不用push
          setRelationInfo(tableName, resp, displayField, valueField);
          ref.value = getStoreRelationOptions(
            tableName,
            displayField,
            valueField,
            formatNumber,
          );
        })
        .catch((error) => {
          /**
           * 需要判断是否为401抛出的特定异常 401清除缓存
           * 其他error清除缓存会导致无限循环调用字典接口 则不做处理
           */
          if (error instanceof UnauthorizedException) {
            // 401时 移除字典缓存 下次登录重新获取
            relationRequestCache.delete(key);
          }
          // 其他不做处理
        })
        .finally(() => {
          // 移除请求状态缓存
          /**
           * 这里主要判断字典item为空的情况(无奈兼容 不给字典item本来就是错误用法)
           * 会导致if一直进入逻辑导致接口无限刷新
           * 在这里dictList为空时 不删除缓存
           */
          if (dataList.length > 0) {
            relationRequestCache.delete(key);
          }
        }),
    );
  }

  return ref as ShallowRef<RelationOption[]>;
}
