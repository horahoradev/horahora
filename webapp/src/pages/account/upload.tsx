import { useState } from "react";

import { Page } from "#components/page";
import {
  FormClient,
  type IFormElements,
  type ISubmitEvent,
} from "#components/forms";
import { toJSON } from "#lib/json";
import { uploadFile } from "#api/upload";
import { Text, File, Tags } from "#components/inputs";
import { LinkInternal } from "#components/links";

const FIELD_NAMES = {
  TITLE: "title",
  DESCRIPTION: "description",
  TAGS: "tags",
  THUMBNAIL: "file[1]",
  FILE: "file[0]",
} as const;
type IFieldName = typeof FIELD_NAMES[keyof typeof FIELD_NAMES];

function UploadPage() {
  const [newVideoID, changeVideoID] = useState<number>();

  async function handleFileUpload(event: ISubmitEvent) {
    const fields = event.currentTarget.elements as IFormElements<IFieldName>;

    const formData = Object.values(FIELD_NAMES).reduce(
      (formData, fieldName) => {
        switch (fieldName) {
          case FIELD_NAMES.TITLE:
          case FIELD_NAMES.DESCRIPTION: {
            const fieldElement = fields[fieldName];
            formData.set(fieldName, fieldElement.value);
            break;
          }

          case FIELD_NAMES.TAGS: {
            const tagsElement = fields[fieldName];
            // the endpoint requires a json array string
            const tagsValue = toJSON(tagsElement.value.split(" "));
            formData.set(fieldName, tagsValue);
            break;
          }

          case FIELD_NAMES.THUMBNAIL:
          case FIELD_NAMES.FILE: {
            const fileElement = fields[fieldName];

            if (!(fileElement.files && fileElement.files[0])) {
              const message = "No thumbnail or file was provided.";
              throw new Error(message);
            }

            formData.set(fieldName, fileElement.files[0]);

            break;
          }

          default:
            throw new Error(
              `The field "${fieldName}" is missing from the form.`
            );
        }

        return formData;
      },
      new FormData()
    );

    const newID = await uploadFile(formData);
    changeVideoID(newID);
  }

  return (
    <Page>
      {newVideoID ? (
        <p>
          New video is available at{" "}
          <LinkInternal href={`/videos/${newVideoID}`} />
        </p>
      ) : (
        <FormClient id="file-upload" onSubmit={handleFileUpload}>
          <Text id="file-upload-title" name={FIELD_NAMES.TITLE}>
            Title
          </Text>
          <Text id="file-upload-description" name={FIELD_NAMES.DESCRIPTION}>
            Description
          </Text>
          <Tags id="file-upload-tags" name={FIELD_NAMES.TAGS}>
            Tags
          </Tags>
          <File id="file-upload-thumb" name={FIELD_NAMES.THUMBNAIL}>
            Thumbnail
          </File>
          <File id="file-upload-file" name={FIELD_NAMES.FILE}>
            Video
          </File>
        </FormClient>
      )}
    </Page>
  );
}

export default UploadPage;
