import { type IInputProps } from "./input";

import { type IFormSectionProps } from "#components/forms";

/**
 * A helper interface for input sections.
 * Components have to manually assign these values
 * to the underlying input component.
 */
export interface IInputSectionProps extends IFormSectionProps {
  id: IInputProps["id"];
  name: IInputProps["name"];
  disabled?: IInputProps["disabled"];
  form?: IInputProps["form"];
  required?: IInputProps["required"];
  defaultValue?: IInputProps["defaultValue"];
}

/**
 * https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/password#allowing_autocomplete
 */
export const PASSWORD_AUTOCOMPLETE = {
  /**
   * Allow the browser or a password manager to automatically fill out the password field. This isn't as informative as using either `current-password` or `new-password`.
   */
  ON: "on",
  /**
   * Don't allow the browser or password manager to automatically fill out the password field. Note that some software ignores this value, since it's typically harmful to users' ability to maintain safe password practices.
   */
  OFF: "off",
  /**
   * Allow the browser or password manager to enter the current password for the site. This provides more information than `on` does, since it lets the browser or password manager automatically enter currently-known password for the site in the field, but not to suggest a new one.
   */
  CURRENT_PASSWORD: "current-password",
  /**
   * Allow the browser or password manager to automatically enter a new password for the site; this is used on "change your password" and "new user" forms, on the field asking the user for a new password. The new password may be generated in a variety of ways, depending on the password manager in use. It may fill in a new suggested password, or it might show the user an interface for creating one.
   */
  NEW_PASSWORD: "new-password",
} as const;

export type IPasswordAutoComplete =
  typeof PASSWORD_AUTOCOMPLETE[keyof typeof PASSWORD_AUTOCOMPLETE];
