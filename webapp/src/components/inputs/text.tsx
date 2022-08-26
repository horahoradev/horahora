import { IInputProps, Input } from "./input";
import { TextArea } from "./textarea";

import { type IFormSectionProps, FormSection, Label } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface ITextProps extends IFormSectionProps {
  id: string;
  name: string;
  minLength?: IInputProps["minLength"];
  maxLength?: IInputProps["maxLength"];
  /**
   * @TODO typed values
   */
  autoComplete?: IInputProps["autoComplete"];
}

export const Text = blockComponent(undefined, Component);

function Component({
  id,
  name,
  minLength,
  maxLength,
  autoComplete,
  children,
  ...blockProps
}: ITextProps) {
  const isShort = maxLength && maxLength < 20;

  return (
    <FormSection {...blockProps}>
      <Label htmlFor={id}>{children}</Label>
      {isShort ? (
        <Input
          id={id}
          type="text"
          name={name}
          minLength={minLength}
          maxLength={maxLength}
          autoComplete={autoComplete}
        />
      ) : (
        <TextArea
          id={id}
          name={name}
          minLength={minLength}
          maxLength={maxLength}
          autoComplete={autoComplete}
        />
      )}
    </FormSection>
  );
}
