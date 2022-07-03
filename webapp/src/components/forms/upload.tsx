import { useState } from "react";

import type { IFormElements, ISubmitEvent } from "./types";

import { toJSON } from "#lib/json";

const FIELD_NAMES = {
  TITLE: "title",
  DESCRIPTION: "description",
  TAGS: "tags",
  THUMBNAIL: "file[1]",
  FILE: "file[0]",
} as const;
type IFieldName = typeof FIELD_NAMES[keyof typeof FIELD_NAMES];

export function UploadForm() {
  const [isSubmitting, switchSubmit] = useState(false);
  const [errors, setErrors] = useState<string[]>([]);

  async function handleFileUpload(event: ISubmitEvent) {
    // do not resubmit while the current one is pending
    if (isSubmitting) {
      return;
    }

    try {
      switchSubmit(true);
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

              formData.set(fieldName, fileElement.files[0])

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

      

    } catch (error) {
      setErrors([String(error)]);
    } finally {
      // enable submit again regardless of outcome of the current submit
      switchSubmit(false);
    }
  }

  return (
    <form id="file-upload" onSubmit={handleFileUpload}>
      <div>
        <label htmlFor="file-upload-title">Title</label>
        <input id="file-upload-title" type="text" name="title" />
      </div>
      <div>
        <label htmlFor="file-upload-description">Description</label>
        <input id="file-upload-description" type="text" name="description" />
      </div>
      {/* TODO: standalone component */}
      <div>
        <label htmlFor="file-upload-tags">Tags</label>
        <textarea id="file-upload-tags" name="tags"></textarea>
        <p>Space-separated list of tag names.</p>
      </div>
      <div>
        <label htmlFor="file-upload-thumb">Thumbnail</label>
        <input id="file-upload-thumb" type="file" name="file[0]" />
      </div>
      <div>
        <label htmlFor="file-upload-file">Video</label>
        <input id="file-upload-file" type="file" name="file[1]" />
      </div>
      <div>
        <ul>
          {errors.map((message, index) => (
            <li key={index}>{message}</li>
          ))}
        </ul>
      </div>
      <div>
        <button type="submit">Submit</button>
      </div>
    </form>
  );
}
