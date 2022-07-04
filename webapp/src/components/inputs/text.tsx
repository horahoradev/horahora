import { type IFormSectionProps, FormSection } from "#components/forms";

export interface ITextProps extends IFormSectionProps {
  id: string;
  name: string;
}

export function Text({ id, name, children }: ITextProps) {
  return (
    <FormSection>
      <label htmlFor={id}>{children}</label>
      <input id={id} type="text" name={name} />
    </FormSection>
  );
}
