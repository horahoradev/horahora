import type { FormEvent } from "react";

/**
 * Convenience type for `HTMLFormElement.elements`
 * @link https://stackoverflow.com/questions/29907163/how-to-work-with-form-elements-in-typescript/70995964#70995964
 */
export type IFormElements<U extends string> = HTMLFormControlsCollection &
  Record<U, HTMLInputElement>;

/**
 * Convenience interface for submit event.
 */
export interface ISubmitEvent extends FormEvent<HTMLFormElement> {}
