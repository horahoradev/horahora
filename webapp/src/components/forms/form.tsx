import { useState } from "react";

import { FormSection } from "./section";
import { type ISubmitEvent } from "./types";

import { type IBlockProps } from "#components/meta";
import { ButtonSubmit } from "#components/buttons";

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
      className="flex flex-col gap-2 max-w-md text-black dark:text-white bg-slate-50 dark:bg-stone-800 rounded-lg px-4 py-2 mx-auto"
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
        <ButtonSubmit>Submit</ButtonSubmit>
      </FormSection>
    </form>
  );
}
