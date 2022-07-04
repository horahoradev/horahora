import { useState } from "react";

import { Page } from "#components/page";
import {
  FormClient,
  FormSection,
  type IFormElements,
  type ISubmitEvent,
} from "#components/forms";
import { toJSON } from "#lib/json";
import { uploadFile } from "#api/upload";

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

            if (!fileElement.files) {
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
      <FormClient id="file-upload" onSubmit={handleFileUpload}>
        <FormSection>
          <label htmlFor="file-upload-title">Title</label>
          <input id="file-upload-title" type="text" name="title" />
        </FormSection>
        <FormSection>
          <label htmlFor="file-upload-description">Description</label>
          <input id="file-upload-description" type="text" name="description" />
        </FormSection>
        {/* TODO: standalone component */}
        <FormSection>
          <label htmlFor="file-upload-tags">Tags</label>
          <textarea id="file-upload-tags" name="tags"></textarea>
          <p>Space-separated list of tag names.</p>
        </FormSection>
        <FormSection>
          <label htmlFor="file-upload-thumb">Thumbnail</label>
          <input id="file-upload-thumb" type="file" name="file[0]" />
        </FormSection>
        <FormSection>
          <label htmlFor="file-upload-file">Video</label>
          <input id="file-upload-file" type="file" name="file[1]" />
        </FormSection>
      </FormClient>
    </Page>
  );
}

export default UploadPage;
