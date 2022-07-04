import { useState } from "react";

import { FormSection } from "./section";
import { type ISubmitEvent } from "./types";

import { type IBlockProps } from "#components/meta";

export interface IFormProps extends IBlockProps<"form"> {
  onSubmit?: (event: ISubmitEvent) => Promise<void> | void;
}

export function Form({ onSubmit, children, ...blockProps }: IFormProps) {
  const [isSubmitting, switchSubmit] = useState(false);
  const [errors, changeErrors] = useState<string[]>([]);

  async function handleSubmit(event: ISubmitEvent) {
    // do not resubmit while the current one is pending
    if (isSubmitting) {
      event.preventDefault();
      return;
    }

    try {
      switchSubmit(true);

      if (onSubmit) {
        await onSubmit(event);
      }
      changeErrors([]);
    } catch (error) {
      event.preventDefault();
      changeErrors([String(error)]);
    } finally {
      // enable submit again regardless of outcome of the current submit
      switchSubmit(false);
    }
  }

  return (
    <form
      className="text-black dark:text-white"
      onSubmit={handleSubmit}
      {...blockProps}
    >
      {children}
      <FormSection>
        <ul>
          {isSubmitting ? (
            <li>Submit is in progress...</li>
          ) : (
            errors.map((message, index) => <li key={index}>{message}</li>)
          )}
        </ul>
      </FormSection>
      <FormSection>
        <button type="submit">Submit</button>
      </FormSection>
    </form>
  );
}
