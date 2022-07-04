import { type IFormSectionProps, FormSection, Label } from "#components/forms";

export interface ITagsProps extends IFormSectionProps {
  id: string;
  name: string;
}

export function Tags({ id, name, children }: ITagsProps) {
  return (
    <FormSection>
      <Label htmlFor={id}>{children}</Label>
      <textarea id={id} name={name}></textarea>
      <p>Space-separated list of tag names.</p>
    </FormSection>
  );
}
