<script lang="tsx">
import type { PropType } from 'vue';

import type { DictFallback } from '#/components/dict/src/type';
import type { RelationOption } from '#/store/relation';

import { computed, defineComponent, h, isVNode } from 'vue';

import { Spin, Tag } from 'ant-design-vue';
import { isFunction, isString } from 'lodash-es';

import { tagTypes } from '#/components/dict/src/data';

/**
 * RelationTag 与 DictTag 行为类似，但用于 relation 数据（RelationOption）
 */
export default defineComponent({
  name: 'RelationTag',
  props: {
    relations: {
      required: false,
      type: Array as PropType<RelationOption[]>,
      default: () => [],
    },
    value: {
      required: true,
      type: [Number, String],
    },
    displayField: {
      required: false,
      type: String as PropType<string>,
      default: undefined,
    },
    valueField: {
      required: false,
      type: String as PropType<string>,
      default: undefined,
    },
    /** 直接传入Tag颜色或tagTypes key，优先使用 */
    color: {
      required: false,
      type: String as PropType<string>,
      default: undefined,
    },
    /** 直接传入额外的css类，优先使用 */
    cssClass: {
      required: false,
      type: String as PropType<string>,
      default: undefined,
    },
    fallback: {
      required: false,
      type: [String, Function] as PropType<DictFallback>,
      default: 'unknown',
    },
  },
  setup(props) {
    const current = computed(() => {
      return props.relations.find((item) => {
        const val = props.valueField
          ? (item as any)[props.valueField]
          : ((item as any).value ?? (item as any).dict_value);
        return String(val) === String(props.value);
      }) as RelationOption | undefined;
    });

    const tagColor = computed<string>(() => {
      const key = props.color;
      if (key) {
        if (Reflect.has(tagTypes, key)) return tagTypes[key]!.color;
        return key;
      }
      return '';
    });

    const tagCssClass = computed<string>(() => {
      return props.cssClass ?? '';
    });

    const label = computed<null | string>(() => {
      if (!current.value) return null;
      return (
        ((current.value as any)[props.displayField as string] as string) ??
        (current.value as any)?.label ??
        (current.value as any)?.dict_label ??
        null
      );
    });

    const loading = computed(() => {
      return props.relations?.length === 0;
    });

    return {
      tagColor,
      tagCssClass,
      label,
      loading,
    };
  },
  render() {
    const { color, cssClass, label, loading, fallback, value, $slots } =
      this as any;

    if (loading) {
      return (
        <div>
          <Spin size="small" spinning />
        </div>
      );
    }

    if (label === null) {
      if ($slots.fallback) {
        return $slots.fallback(value);
      }
      if (isFunction(fallback)) {
        const rValue = fallback(value);
        if (isVNode(rValue)) return h(rValue);
        return <div>{rValue}</div>;
      }
      if (isString(fallback)) return <div>{fallback}</div>;
    }

    if (color) {
      return (
        <div>
          <Tag class={cssClass} color={color}>
            {label}
          </Tag>
        </div>
      );
    }

    return (
      <div>
        <div class={cssClass}>{label}</div>
      </div>
    );
  },
});
</script>
