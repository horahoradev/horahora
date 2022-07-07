import { Input } from "./input";

import { type IFormSectionProps, FormSection, Label } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface IFileProps extends IFormSectionProps {
  id: string;
  name: string;
}

export const File = blockComponent(undefined, Component);

function Component({ id, name, children, ...blockProps }: IFileProps) {
  return (
    <FormSection {...blockProps}>
      <Label htmlFor={id}>{children}</Label>
      <Input id={id} type="file" name={name} />
    </FormSection>
  );
}
