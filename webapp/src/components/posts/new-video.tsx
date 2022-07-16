import {
  FormClient,
  type IFormElements,
  type ISubmitEvent,
} from "#components/forms";
import { Text } from "#components/inputs";

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

    await onNewURL(urlInput.value);
  }

  return (
    <FormClient id="new-video" onSubmit={handleSubmit}>
      <Text id="new-video-url" name="url">
        New video URL
      </Text>
    </FormClient>
  );
}
