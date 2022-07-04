import { type IFormSectionProps, FormSection, Label } from "#components/forms";

export interface IFileProps extends IFormSectionProps {
  id: string;
  name: string;
}

export function File({ id, name, children }: IFileProps) {
  return (
    <FormSection>
      <Label htmlFor={id}>{children}</Label>
      <input id={id} type="file" name={name} />
    </FormSection>
  );
}
