import {
  FormClient,
  type IFormElements,
  type ISubmitEvent,
} from "#components/forms";
import { Text } from "#components/inputs";
import { createNewArchivalRequest } from "#api/archives";

const FIELD_NAMES = {
  NEW_URL: "url",
} as const;
type IFieldName = typeof FIELD_NAMES[keyof typeof FIELD_NAMES];

export interface INewVideoFormProps {
  onNewURL: (url: string) => Promise<void>;
}

export function NewVideoForm({ onNewURL }: INewVideoFormProps) {
  async function handleSubmit(event: ISubmitEvent) {
    const fields = event.currentTarget.elements as IFormElements<IFieldName>;
    const urlInput = fields[FIELD_NAMES.NEW_URL];
    const newURL = urlInput.value;
    const formParams = new URLSearchParams([["url", newURL]]);

    await createNewArchivalRequest(formParams);
    await onNewURL(newURL);
  }

  return (
    <FormClient id="new-video" onSubmit={handleSubmit}>
      <Text id="new-video-url" name="url">
        New video URL
      </Text>
    </FormClient>
  );
}
