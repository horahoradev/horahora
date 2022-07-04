import { type IFormSectionProps, FormSection } from "#components/forms";

export interface IFileProps extends IFormSectionProps {
  id: string;
  name: string;
}

export function File({ id, name, children }: IFileProps) {
  return (
    <FormSection>
      <label htmlFor={id}>{children}</label>
      <input id={id} type="file" name={name} />
    </FormSection>
  );
}
