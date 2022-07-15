import { Input } from "./input";
import { TextArea } from "./textarea";

import { type IFormSectionProps, FormSection, Label } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface ITextProps extends IFormSectionProps {
  id: string;
  name: string;
  maxLength?: number;
}

export const Text = blockComponent(undefined, Component);

function Component({
  id,
  name,
  maxLength,
  children,
  ...blockProps
}: ITextProps) {
  const isShort = maxLength && maxLength < 20;

  return (
    <FormSection {...blockProps}>
      <Label htmlFor={id}>{children}</Label>
      {isShort ? (
        <Input id={id} type="text" name={name} maxLength={maxLength} />
      ) : (
        <TextArea id={id} name={name} maxLength={maxLength} />
      )}
    </FormSection>
  );
}
