import {
    FormClient,
    type IFormElements,
    type ISubmitEvent,
  } from "#components/forms";
  import { Text } from "#components/inputs";
  import { addInferenceCategory, createNewArchivalRequest } from "#api/archives";

  const FIELD_NAMES = {
    NEW_TAG: "tag",
    NEW_CATEGORY: "category",
  } as const;
  type IFieldName = typeof FIELD_NAMES[keyof typeof FIELD_NAMES];

  export interface INewInferenceFormProps {
    onNewCategory: (category: string, tag: string) => Promise<void>;
  }

  export function NewInferenceForm({ onNewCategory }: INewInferenceFormProps) {
    async function handleSubmit(event: ISubmitEvent) {
      const fields = event.currentTarget.elements as IFormElements<IFieldName>;
      const categoryInput = fields[FIELD_NAMES.NEW_CATEGORY];
      const tagInput = fields[FIELD_NAMES.NEW_TAG];
      const formParams = new URLSearchParams([["tag", tagInput.value], ["category", categoryInput.value]]);

      await addInferenceCategory(formParams);
      await onNewCategory(categoryInput.value, tagInput.value);
    }

    return (
      <FormClient id="new-inference" onSubmit={handleSubmit}>
        <Text id="new-inference-category" name="category">
          Category
        </Text>
        <Text id="new-inference-tag" name="tag">
          Tag
        </Text>
      </FormClient>
    );
  }
