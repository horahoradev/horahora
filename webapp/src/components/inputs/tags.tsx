import { TextArea } from "./textarea";

import { type IFormSectionProps, FormSection, Label } from "#components/forms";

export interface ITagsProps extends IFormSectionProps {
  id: string;
  name: string;
}

export function Tags({ id, name, children }: ITagsProps) {
  return (
    <FormSection>
      <Label htmlFor={id}>{children}</Label>
      <TextArea id={id} name={name}></TextArea>
      <p>Space-separated list of tag names.</p>
    </FormSection>
  );
}
