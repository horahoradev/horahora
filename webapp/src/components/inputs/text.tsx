import { Input } from "./input";
import { TextArea } from "./textarea";

import { type IFormSectionProps, FormSection, Label } from "#components/forms";

export interface ITextProps extends IFormSectionProps {
  id: string;
  name: string;
  maxLength?: number;
}

export function Text({ id, name, maxLength, children }: ITextProps) {
  return (
    <FormSection>
      <Label htmlFor={id}>{children}</Label>
      {maxLength && maxLength < 20 ? (
        <Input id={id} type="text" name={name} maxLength={maxLength} />
      ) : (
        <TextArea id={id} name={name} maxLength={maxLength} />
      )}
    </FormSection>
  );
}
