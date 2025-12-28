<!--
2025年03月08日重构为原生表单(反向重构??)
该文件作为例子 使用原生表单而非useVbenForm
-->
<script setup lang="ts">
import type { RuleObject } from 'ant-design-vue/es/form';

import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { DictEnum } from '@vben/constants';
import { $t } from '@vben/locales';
import { cloneDeep } from '@vben/utils';

import { Form, FormItem, Input, RadioGroup } from 'ant-design-vue';
import { pick } from 'lodash-es';

import { noticeAdd, noticeInfo, noticeUpdate } from '#/api/system/notice';
import { Tinymce } from '#/components/tinymce';
import { contentWithOssIdTransform } from '#/components/tinymce/src/helper';
import { getDictOptions } from '#/utils/dict';
import { useBeforeCloseDiff } from '#/utils/popup';

const emit = defineEmits<{ reload: [] }>();

const isUpdate = ref(false);
const title = computed(() => {
  return isUpdate.value ? $t('pages.common.edit') : $t('pages.common.add');
});

/**
 * 定义表单数据类型
 */
interface FormData {
  id?: number;
  notice_title?: string;
  status?: string;
  notice_type?: string;
  notice_content?: string;
}

/**
 * 定义默认值 用于reset
 */
const defaultValues: FormData = {
  id: undefined,
  notice_title: '',
  status: '0',
  notice_type: '1',
  notice_content: '',
};

/**
 * 表单数据ref
 */
const formData = ref(defaultValues);

type AntdFormRules<T> = Partial<Record<keyof T, RuleObject[]>> & {
  [key: string]: RuleObject[];
};
/**
 * 表单校验规则
 */
const formRules = ref<AntdFormRules<FormData>>({
  status: [{ required: true, message: $t('ui.formRules.selectRequired') }],
  notice_content: [{ required: true, message: $t('ui.formRules.required') }],
  notice_type: [{ required: true, message: $t('ui.formRules.selectRequired') }],
  notice_title: [{ required: true, message: $t('ui.formRules.required') }],
});

/**
 * useForm解构出表单方法
 */
const { validate, validateInfos, resetFields } = Form.useForm(
  formData,
  formRules,
);

function customFormValueGetter() {
  return JSON.stringify(formData.value);
}

const { onBeforeClose, markInitialized, resetInitialized } = useBeforeCloseDiff(
  {
    initializedGetter: customFormValueGetter,
    currentGetter: customFormValueGetter,
  },
);

const [BasicModal, modalApi] = useVbenModal({
  class: 'w-[800px]',
  fullscreenButton: true,
  onBeforeClose,
  onClosed: handleClosed,
  onConfirm: handleConfirm,
  onOpenChange: async (isOpen) => {
    if (!isOpen) {
      return null;
    }
    modalApi.modalLoading(true);

    const { id } = modalApi.getData() as { id?: number | string };
    isUpdate.value = !!id;
    if (isUpdate.value && id) {
      const record = await noticeInfo(id);
      // 只赋值存在的字段
      const filterRecord = pick(record, Object.keys(defaultValues));

      // 你可以调用这个方法来显示私有桶的图片（每次获取最新）
      // 如果你是公开桶 最好去掉这段代码 会造成不必要的查询
      filterRecord.notice_content =
        (await contentWithOssIdTransform(record.notice_content)) ?? '';

      formData.value = filterRecord;
    }
    await markInitialized();

    modalApi.modalLoading(false);
  },
});

async function handleConfirm() {
  try {
    modalApi.lock(true);
    await validate();
    // 可能会做数据处理 使用cloneDeep深拷贝
    const data = cloneDeep(formData.value);
    await (isUpdate.value ? noticeUpdate(data) : noticeAdd(data));
    resetInitialized();
    emit('reload');
    modalApi.close();
  } catch (error) {
    console.error(error);
  } finally {
    modalApi.lock(false);
  }
}

async function handleClosed() {
  formData.value = defaultValues;
  resetFields();
  resetInitialized();
}

const noticeStatusOptions = getDictOptions(DictEnum.SYS_NOTICE_STATUS);
const notice_typeOptions = getDictOptions(DictEnum.SYS_NOTICE_TYPE);
</script>

<template>
  <BasicModal :title="title">
    <Form layout="vertical">
      <FormItem label="公告标题" v-bind="validateInfos.notice_title">
        <Input
          :placeholder="$t('ui.formRules.required')"
          v-model:value="formData.notice_title"
        />
      </FormItem>
      <div class="grid sm:grid-cols-1 lg:grid-cols-2">
        <FormItem label="公告状态" v-bind="validateInfos.status">
          <RadioGroup
            button-style="solid"
            option-type="button"
            v-model:value="formData.status"
            :options="noticeStatusOptions"
          />
        </FormItem>
        <FormItem label="公告类型" v-bind="validateInfos.notice_type">
          <RadioGroup
            button-style="solid"
            option-type="button"
            v-model:value="formData.notice_type"
            :options="notice_typeOptions"
          />
        </FormItem>
      </div>
      <FormItem label="公告内容" v-bind="validateInfos.notice_content">
        <Tinymce v-model="formData.notice_content" />
      </FormItem>
    </Form>
  </BasicModal>
</template>
