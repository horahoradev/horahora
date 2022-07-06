import { Form, type IFormProps } from "./form";
import { type ISubmitEvent } from "./types";

import { ClientComponent } from "#components/meta";

export interface IFormClientProps extends IFormProps {
  onSubmit: (event: ISubmitEvent) => Promise<void>;
}

/**
 * Form which only renders upon hydration.
 *
 * Submit is default prevented and `onSubmit()` is always async and is required.
 */
export function FormClient({
  onSubmit,
  children,
  ...blockProps
}: IFormClientProps) {
  async function handleSubmit(
    ...args: Parameters<IFormClientProps["onSubmit"]>
  ) {
    const [event] = args;
    event.preventDefault();
    await onSubmit(event);
  }

  return (
    <ClientComponent>
      <Form onSubmit={handleSubmit} {...blockProps}>
        {children}
      </Form>
    </ClientComponent>
  );
}
