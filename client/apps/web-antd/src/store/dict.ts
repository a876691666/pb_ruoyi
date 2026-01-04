/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { useRelationStore, relationToOptions } from '#/store/relation';
import type { RelationOption } from '#/store/relation';
import type { DictData } from '#/api/system/dict/dict-data-model';

/**
 * 向后兼容的精简包装：保留旧 API（`useDictStore` / `dictToOptions` / `DictOption`），
 * 实际委托到 `useRelationStore` 实现，避免重复逻辑并保持代码简洁。
 */
export type DictOption = RelationOption;

export function dictToOptions(data: DictData[], formatNumber = false): DictOption[] {
  return relationToOptions(data, formatNumber);
}

export const useDictStore = () => {
  const s = useRelationStore();
  return {
    $reset: s.$reset,
    dictOptionsMap: s.relationOptionsMap,
    dictRequestCache: s.relationRequestCache,
    getDictOptions: s.getRelationOptions,
    resetCache: s.resetRelationCache,
    setDictInfo: s.setRelationInfo,
  } as const;
};
