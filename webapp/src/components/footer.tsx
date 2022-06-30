import { PrivacyPolicy } from "./privacy-policy";
import { TermsOfService } from "./terms-of-service";

interface IFooterProps extends Record<string, unknown> {}

export function Footer(props: IFooterProps) {
  const { userData, dataless } = props;

  return (
    <nav className="h-8 w-full bg-white dark:bg-gray-900 shadow">
      <div className="flex justify-around p-auto">
        <PrivacyPolicy />
        <TermsOfService />
      </div>
    </nav>
  );
}
