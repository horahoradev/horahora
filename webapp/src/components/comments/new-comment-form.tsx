import {
  FormClient,
  type IFormElements,
  type ISubmitEvent,
} from "#components/forms";
import { Hidden, Text } from "#components/inputs";

export interface ICommentInit {
  video_id: string;
  content: string;
  parent: string;
}

const FIELD_NAMES = {
  MESSAGE: "content",
  VIDEO_ID: "video_id",
  PARENT: "parent",
} as const;
type IFieldName = typeof FIELD_NAMES[keyof typeof FIELD_NAMES];

export interface INewCommentFormProps {
  videoID: number;
  parentID?: string;

  onNewComment: (commentInit: URLSearchParams) => Promise<void>;
}

export function NewCommentForm({
  videoID,
  parentID,
  onNewComment,
}: INewCommentFormProps) {
  async function handleSubmit(event: ISubmitEvent) {
    const fields = event.currentTarget.elements as IFormElements<IFieldName>;
    const commentInit = Object.values(FIELD_NAMES).reduce(
      (formParams, fieldName) => {
        switch (fieldName) {
          case FIELD_NAMES.MESSAGE:
          case FIELD_NAMES.PARENT:
          case FIELD_NAMES.VIDEO_ID: {
            const fieldElement = fields[fieldName];
            formParams.set(fieldName, fieldElement.value);
            break;
          }

          default:
            throw new Error(
              `The field "${fieldName}" is missing from the form.`
            );
        }
        return formParams;
      },
      new URLSearchParams()
    );

    await onNewComment(commentInit);
    event.currentTarget.reset();
  }

  return (
    <FormClient onSubmit={handleSubmit} id="new-comment">
      <Hidden
        id="new-comment-parent-id"
        name={FIELD_NAMES.PARENT}
        defaultValue={parentID}
      />
      <Hidden
        id="new-comment-video-id"
        name={FIELD_NAMES.VIDEO_ID}
        defaultValue={videoID}
      />
      <Text id="new-comment-content" name={FIELD_NAMES.MESSAGE}>
        New comment
      </Text>
    </FormClient>
  );
}
