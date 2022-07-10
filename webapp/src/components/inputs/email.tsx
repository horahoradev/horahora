import { Input } from "./input";

import { FormSection, type IFormSectionProps, Label } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface IEmailProps extends IFormSectionProps {
  id: string;
  name: string;
}

export const Email = blockComponent(undefined, Component);

function Component({ id, name, children, ...blockProps }: IEmailProps) {
  return (
    <FormSection {...blockProps}>
      <Label htmlFor={id}>{children}</Label>
      <Input id={id} name={name} />
    </FormSection>
  );
}
