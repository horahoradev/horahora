import { Input } from "./input";

import { type IFormSectionProps, FormSection, Label } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface ISearchProps extends IFormSectionProps {
  id: string;
  name: string;
}

export const Search = blockComponent(undefined, Component);

function Component({ id, name, children, ...blockProps }: ISearchProps) {
  return (
    <FormSection {...blockProps}>
      <Label htmlFor={id}>{children}</Label>
      <Input id={id} name={name} type="search" />
    </FormSection>
  );
}
