import { Input } from "./input";

import { type IFormSectionProps, FormSection, Label } from "#components/forms";

export interface ITextProps extends IFormSectionProps {
  id: string;
  name: string;
}

export function Text({ id, name, children }: ITextProps) {
  return (
    <FormSection>
      <Label htmlFor={id}>{children}</Label>
      <Input id={id} type="text" name={name} />
    </FormSection>
  );
}
